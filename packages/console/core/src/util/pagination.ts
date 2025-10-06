import { z } from "zod"
import { SQL, sql, asc, desc } from "drizzle-orm"

export namespace Pagination {
  /**
   * Standard pagination parameters schema
   */
  export const schema = z.object({
    page: z.number().min(1).default(1),
    pageSize: z.number().min(1).max(100).default(20),
    sortBy: z.string().optional(),
    sortOrder: z.enum(["asc", "desc"]).default("desc"),
  })

  export type Params = z.infer<typeof schema>

  /**
   * Calculate offset from page number
   */
  export function getOffset(params: Params): number {
    return (params.page - 1) * params.pageSize
  }

  /**
   * Build pagination response
   */
  export function buildResponse<T>(
    items: T[],
    totalCount: number,
    params: Params,
  ): PaginationResponse<T> {
    const totalPages = Math.ceil(totalCount / params.pageSize)
    const hasNextPage = params.page < totalPages
    const hasPreviousPage = params.page > 1

    return {
      items,
      pagination: {
        page: params.page,
        pageSize: params.pageSize,
        totalCount,
        totalPages,
        hasNextPage,
        hasPreviousPage,
      },
    }
  }

  /**
   * Pagination response type
   */
  export type PaginationResponse<T> = {
    items: T[]
    pagination: {
      page: number
      pageSize: number
      totalCount: number
      totalPages: number
      hasNextPage: boolean
      hasPreviousPage: boolean
    }
  }

  /**
   * Apply sorting to query
   */
  export function applySorting(column: any, order: "asc" | "desc") {
    return order === "asc" ? asc(column) : desc(column)
  }

  /**
   * Cursor-based pagination schema (for infinite scroll)
   */
  export const cursorSchema = z.object({
    cursor: z.string().optional(),
    limit: z.number().min(1).max(100).default(20),
    sortOrder: z.enum(["asc", "desc"]).default("desc"),
  })

  export type CursorParams = z.infer<typeof cursorSchema>

  /**
   * Build cursor-based pagination response
   */
  export function buildCursorResponse<T extends { id: string }>(
    items: T[],
    params: CursorParams,
  ): CursorResponse<T> {
    const hasMore = items.length === params.limit
    const nextCursor = hasMore ? items[items.length - 1]?.id : null

    return {
      items,
      cursor: {
        nextCursor,
        hasMore,
      },
    }
  }

  /**
   * Cursor-based pagination response type
   */
  export type CursorResponse<T> = {
    items: T[]
    cursor: {
      nextCursor: string | null
      hasMore: boolean
    }
  }
}
