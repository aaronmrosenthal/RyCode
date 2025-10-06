import { z } from "zod"
import { and, eq, getTableColumns, isNull, sql, desc, asc } from "drizzle-orm"
import { fn } from "./util/fn"
import { Database } from "./drizzle"
import { TestTable, TestStatus } from "./schema/test.sql"
import { Actor } from "./actor"
import { Identifier } from "./identifier"
import { Pagination } from "./util/pagination"
import { Filtering } from "./util/filtering"

export namespace Test {
  /**
   * CREATE - Create a new test entity
   */
  export const create = fn(
    z.object({
      name: z.string().min(1).max(255),
      description: z.string().optional(),
      status: z.enum(TestStatus).default("active"),
    }),
    async ({ name, description, status }) => {
      const workspaceID = Actor.workspace()

      const id = Identifier.create("test")

      await Database.use((tx) =>
        tx.insert(TestTable).values({
          id,
          workspaceID,
          name,
          description: description ?? null,
          status,
        }),
      )

      return { id, name, description, status }
    },
  )

  /**
   * READ - Get a single test entity by ID
   */
  export const fromID = fn(z.string(), (id) =>
    Database.use((tx) =>
      tx
        .select()
        .from(TestTable)
        .where(
          and(
            eq(TestTable.workspaceID, Actor.workspace()),
            eq(TestTable.id, id),
            isNull(TestTable.timeDeleted),
          ),
        )
        .then((rows) => rows[0]),
    ),
  )

  /**
   * READ - List all test entities in workspace
   */
  export const list = fn(
    z.object({
      status: z.enum(TestStatus).optional(),
      limit: z.number().min(1).max(100).default(50),
      offset: z.number().min(0).default(0),
    }),
    ({ status, limit, offset }) => {
      const workspaceID = Actor.workspace()

      return Database.use((tx) => {
        const conditions = [
          eq(TestTable.workspaceID, workspaceID),
          isNull(TestTable.timeDeleted),
        ]

        // Filter by status if provided
        if (status) {
          conditions.push(eq(TestTable.status, status))
        }

        return tx
          .select()
          .from(TestTable)
          .where(and(...conditions))
          .limit(limit)
          .offset(offset)
      })
    },
  )

  /**
   * UPDATE - Update an existing test entity
   */
  export const update = fn(
    z.object({
      id: z.string(),
      name: z.string().min(1).max(255).optional(),
      description: z.string().optional(),
      status: z.enum(TestStatus).optional(),
    }),
    async ({ id, name, description, status }) => {
      const workspaceID = Actor.workspace()

      // Build update object with only provided fields
      const updateData: any = {
        timeUpdated: sql`now()`,
      }

      if (name !== undefined) updateData.name = name
      if (description !== undefined) updateData.description = description
      if (status !== undefined) updateData.status = status

      const result = await Database.use((tx) =>
        tx
          .update(TestTable)
          .set(updateData)
          .where(
            and(
              eq(TestTable.id, id),
              eq(TestTable.workspaceID, workspaceID),
              isNull(TestTable.timeDeleted),
            ),
          ),
      )

      // Return updated entity
      return await fromID(id)
    },
  )

  /**
   * DELETE - Soft delete a test entity
   */
  export const remove = fn(z.string(), async (id) => {
    const workspaceID = Actor.workspace()

    await Database.use((tx) =>
      tx
        .update(TestTable)
        .set({
          timeDeleted: sql`now()`,
        })
        .where(and(eq(TestTable.id, id), eq(TestTable.workspaceID, workspaceID))),
    )

    return { success: true, id }
  })

  /**
   * DELETE - Permanently delete a test entity (hard delete)
   */
  export const destroy = fn(z.string(), async (id) => {
    const workspaceID = Actor.workspace()

    await Database.use((tx) =>
      tx
        .delete(TestTable)
        .where(and(eq(TestTable.id, id), eq(TestTable.workspaceID, workspaceID))),
    )

    return { success: true, id }
  })

  /**
   * COUNT - Count test entities in workspace
   */
  export const count = fn(
    z.object({
      status: z.enum(TestStatus).optional(),
    }),
    ({ status }) => {
      const workspaceID = Actor.workspace()

      return Database.use(async (tx) => {
        const conditions = [
          eq(TestTable.workspaceID, workspaceID),
          isNull(TestTable.timeDeleted),
        ]

        if (status) {
          conditions.push(eq(TestTable.status, status))
        }

        const result = await tx
          .select({ count: sql<number>`count(*)` })
          .from(TestTable)
          .where(and(...conditions))

        return result[0]?.count ?? 0
      })
    },
  )

  /**
   * SEARCH - Search test entities by name
   */
  export const search = fn(
    z.object({
      query: z.string().min(1),
      limit: z.number().min(1).max(100).default(20),
    }),
    ({ query, limit }) => {
      const workspaceID = Actor.workspace()

      return Database.use((tx) =>
        tx
          .select()
          .from(TestTable)
          .where(
            and(
              eq(TestTable.workspaceID, workspaceID),
              sql`${TestTable.name} LIKE ${`%${query}%`}`,
              isNull(TestTable.timeDeleted),
            ),
          )
          .limit(limit),
      )
    },
  )

  /**
   * LIST with advanced pagination, filtering, and sorting
   */
  export const listAdvanced = fn(
    z.object({
      page: z.number().min(1).default(1),
      pageSize: z.number().min(1).max(100).default(20),
      sortBy: z.enum(["name", "status", "timeCreated", "timeUpdated"]).optional(),
      sortOrder: z.enum(["asc", "desc"]).default("desc"),
      filters: z
        .array(
          z.object({
            field: z.enum(["status", "name"]),
            operator: z.enum(["eq", "like"]),
            value: z.union([z.string(), z.array(z.string())]),
          }),
        )
        .optional(),
      search: z.string().optional(),
    }),
    async ({ page, pageSize, sortBy, sortOrder, filters, search }) => {
      const workspaceID = Actor.workspace()

      return Database.use(async (tx) => {
        // Base conditions
        const baseConditions = [
          eq(TestTable.workspaceID, workspaceID),
          isNull(TestTable.timeDeleted),
        ]

        // Apply filters
        if (filters && filters.length > 0) {
          const filterCondition = Filtering.applyFilters(TestTable, filters as any)
          if (filterCondition) {
            baseConditions.push(filterCondition)
          }
        }

        // Apply search
        if (search) {
          const searchCondition = Filtering.applySearch(TestTable, {
            query: search,
            fields: ["name", "description"],
          })
          if (searchCondition) {
            baseConditions.push(searchCondition)
          }
        }

        const whereClause = and(...baseConditions)

        // Get total count
        const countResult = await tx
          .select({ count: sql<number>`count(*)` })
          .from(TestTable)
          .where(whereClause)

        const totalCount = countResult[0]?.count ?? 0

        // Build sort column
        const sortColumn = sortBy ? TestTable[sortBy] : TestTable.timeCreated
        const orderBy = sortOrder === "asc" ? asc(sortColumn) : desc(sortColumn)

        // Get paginated results
        const offset = (page - 1) * pageSize
        const items = await tx
          .select()
          .from(TestTable)
          .where(whereClause)
          .orderBy(orderBy)
          .limit(pageSize)
          .offset(offset)

        return Pagination.buildResponse(items, totalCount, {
          page,
          pageSize,
          sortBy,
          sortOrder,
        })
      })
    },
  )

  /**
   * LIST with cursor-based pagination (for infinite scroll)
   */
  export const listCursor = fn(
    z.object({
      cursor: z.string().optional(),
      limit: z.number().min(1).max(100).default(20),
      sortOrder: z.enum(["asc", "desc"]).default("desc"),
    }),
    async ({ cursor, limit, sortOrder }) => {
      const workspaceID = Actor.workspace()

      return Database.use(async (tx) => {
        const conditions = [
          eq(TestTable.workspaceID, workspaceID),
          isNull(TestTable.timeDeleted),
        ]

        // Apply cursor if provided
        if (cursor) {
          if (sortOrder === "desc") {
            conditions.push(sql`${TestTable.id} < ${cursor}`)
          } else {
            conditions.push(sql`${TestTable.id} > ${cursor}`)
          }
        }

        const items = await tx
          .select()
          .from(TestTable)
          .where(and(...conditions))
          .orderBy(sortOrder === "asc" ? asc(TestTable.id) : desc(TestTable.id))
          .limit(limit)

        return Pagination.buildCursorResponse(items, { cursor, limit, sortOrder })
      })
    },
  )
}
