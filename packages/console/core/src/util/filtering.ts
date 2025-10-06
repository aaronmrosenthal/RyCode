import { z } from "zod"
import { SQL, sql, and, or, eq, ne, gt, gte, lt, lte, like, inArray } from "drizzle-orm"

export namespace Filtering {
  /**
   * Filter operator types
   */
  export const operators = [
    "eq", // equals
    "ne", // not equals
    "gt", // greater than
    "gte", // greater than or equal
    "lt", // less than
    "lte", // less than or equal
    "like", // LIKE pattern
    "in", // IN array
    "between", // BETWEEN two values
  ] as const

  export type Operator = (typeof operators)[number]

  /**
   * Filter condition schema
   */
  export const filterSchema = z.object({
    field: z.string().max(100), // Prevent excessively long field names
    operator: z.enum(operators),
    value: z.union([
      z.string().max(1000), // Limit string length to prevent DoS
      z.number(),
      z.boolean(),
      z.array(z.any()).max(100), // Limit array size
    ]),
  })

  export type Filter = z.infer<typeof filterSchema>

  /**
   * Apply a single filter condition
   */
  export function applyFilter(column: any, filter: Filter): SQL | undefined {
    switch (filter.operator) {
      case "eq":
        return eq(column, filter.value)
      case "ne":
        return ne(column, filter.value)
      case "gt":
        return gt(column, filter.value as any)
      case "gte":
        return gte(column, filter.value as any)
      case "lt":
        return lt(column, filter.value as any)
      case "lte":
        return lte(column, filter.value as any)
      case "like":
        return like(column, `%${filter.value}%`)
      case "in":
        return inArray(column, filter.value as any[])
      case "between":
        if (Array.isArray(filter.value) && filter.value.length === 2) {
          return and(gte(column, filter.value[0]), lte(column, filter.value[1]))
        }
        return undefined
      default:
        return undefined
    }
  }

  /**
   * Apply multiple filters (AND logic)
   */
  export function applyFilters(table: any, filters: Filter[]): SQL | undefined {
    // Prevent DoS with excessive filters
    if (filters.length > 20) {
      throw new Error("Too many filters. Maximum 20 filters allowed.")
    }

    const conditions = filters
      .map((filter) => {
        const column = table[filter.field]
        // Security: Only apply filter if column exists in table
        // This prevents SQL injection via unknown field names
        if (!column) {
          console.warn(`Attempted to filter on unknown field: ${filter.field}`)
          return undefined
        }
        return applyFilter(column, filter)
      })
      .filter((c): c is SQL => c !== undefined)

    return conditions.length > 0 ? and(...conditions) : undefined
  }

  /**
   * Sorting configuration schema
   */
  export const sortSchema = z.object({
    field: z.string(),
    order: z.enum(["asc", "desc"]).default("desc"),
  })

  export type Sort = z.infer<typeof sortSchema>

  /**
   * Search configuration
   */
  export const searchSchema = z.object({
    query: z.string().min(1).max(500), // Prevent DoS with long queries
    fields: z.array(z.string()).min(1).max(10), // Limit search fields
  })

  export type Search = z.infer<typeof searchSchema>

  /**
   * Apply search across multiple fields (OR logic)
   */
  export function applySearch(table: any, search: Search): SQL | undefined {
    // Sanitize search query to prevent SQL injection
    const sanitizedQuery = search.query.replace(/[%_]/g, "\\$&") // Escape LIKE wildcards

    const conditions = search.fields
      .map((field) => {
        const column = table[field]
        // Security: Only search on fields that exist in table
        if (!column) {
          console.warn(`Attempted to search on unknown field: ${field}`)
          return undefined
        }
        return like(column, `%${sanitizedQuery}%`)
      })
      .filter((c): c is SQL => c !== undefined)

    return conditions.length > 0 ? or(...conditions) : undefined
  }

  /**
   * Combined query params schema
   */
  export const querySchema = z.object({
    filters: z.array(filterSchema).optional(),
    sort: sortSchema.optional(),
    search: searchSchema.optional(),
    page: z.number().min(1).default(1),
    pageSize: z.number().min(1).max(100).default(20),
  })

  export type QueryParams = z.infer<typeof querySchema>

  /**
   * Build WHERE clause from query params
   */
  export function buildWhereClause(table: any, params: QueryParams): SQL | undefined {
    const conditions: (SQL | undefined)[] = []

    // Add filters
    if (params.filters && params.filters.length > 0) {
      conditions.push(applyFilters(table, params.filters))
    }

    // Add search
    if (params.search) {
      conditions.push(applySearch(table, params.search))
    }

    const validConditions = conditions.filter((c): c is SQL => c !== undefined)
    return validConditions.length > 0 ? and(...validConditions) : undefined
  }
}
