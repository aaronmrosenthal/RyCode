import { describe, test, expect } from "bun:test"
import { Pagination } from "../../src/util/pagination"

describe("Pagination Utility", () => {
  describe("schema", () => {
    test("should accept valid pagination params", () => {
      const params = {
        page: 1,
        pageSize: 20,
        sortBy: "name",
        sortOrder: "desc" as const,
      }

      expect(() => Pagination.schema.parse(params)).not.toThrow()
    })

    test("should use default values", () => {
      const params = Pagination.schema.parse({})

      expect(params.page).toBe(1)
      expect(params.pageSize).toBe(20)
      expect(params.sortOrder).toBe("desc")
    })

    test("should reject invalid page number", () => {
      expect(() => Pagination.schema.parse({ page: 0 })).toThrow()
      expect(() => Pagination.schema.parse({ page: -1 })).toThrow()
    })

    test("should reject invalid page size", () => {
      expect(() => Pagination.schema.parse({ pageSize: 0 })).toThrow()
      expect(() => Pagination.schema.parse({ pageSize: 101 })).toThrow()
    })

    test("should reject invalid sort order", () => {
      expect(() => Pagination.schema.parse({ sortOrder: "invalid" })).toThrow()
    })
  })

  describe("getOffset", () => {
    test("should calculate correct offset for page 1", () => {
      const offset = Pagination.getOffset({ page: 1, pageSize: 20, sortOrder: "desc" })

      expect(offset).toBe(0)
    })

    test("should calculate correct offset for page 2", () => {
      const offset = Pagination.getOffset({ page: 2, pageSize: 20, sortOrder: "desc" })

      expect(offset).toBe(20)
    })

    test("should calculate correct offset for page 5", () => {
      const offset = Pagination.getOffset({ page: 5, pageSize: 10, sortOrder: "desc" })

      expect(offset).toBe(40)
    })

    test("should handle custom page sizes", () => {
      const offset = Pagination.getOffset({ page: 3, pageSize: 50, sortOrder: "desc" })

      expect(offset).toBe(100)
    })
  })

  describe("buildResponse", () => {
    test("should build correct response for first page", () => {
      const items = Array.from({ length: 20 }, (_, i) => ({ id: i }))
      const response = Pagination.buildResponse(items, 100, {
        page: 1,
        pageSize: 20,
        sortOrder: "desc",
      })

      expect(response.items).toHaveLength(20)
      expect(response.pagination.page).toBe(1)
      expect(response.pagination.pageSize).toBe(20)
      expect(response.pagination.totalCount).toBe(100)
      expect(response.pagination.totalPages).toBe(5)
      expect(response.pagination.hasNextPage).toBe(true)
      expect(response.pagination.hasPreviousPage).toBe(false)
    })

    test("should build correct response for middle page", () => {
      const items = Array.from({ length: 20 }, (_, i) => ({ id: i }))
      const response = Pagination.buildResponse(items, 100, {
        page: 3,
        pageSize: 20,
        sortOrder: "desc",
      })

      expect(response.pagination.page).toBe(3)
      expect(response.pagination.hasNextPage).toBe(true)
      expect(response.pagination.hasPreviousPage).toBe(true)
    })

    test("should build correct response for last page", () => {
      const items = Array.from({ length: 20 }, (_, i) => ({ id: i }))
      const response = Pagination.buildResponse(items, 100, {
        page: 5,
        pageSize: 20,
        sortOrder: "desc",
      })

      expect(response.pagination.page).toBe(5)
      expect(response.pagination.hasNextPage).toBe(false)
      expect(response.pagination.hasPreviousPage).toBe(true)
    })

    test("should handle partial last page", () => {
      const items = Array.from({ length: 10 }, (_, i) => ({ id: i }))
      const response = Pagination.buildResponse(items, 95, {
        page: 5,
        pageSize: 20,
        sortOrder: "desc",
      })

      expect(response.items).toHaveLength(10)
      expect(response.pagination.totalPages).toBe(5)
      expect(response.pagination.hasNextPage).toBe(false)
    })

    test("should handle single page result", () => {
      const items = Array.from({ length: 15 }, (_, i) => ({ id: i }))
      const response = Pagination.buildResponse(items, 15, {
        page: 1,
        pageSize: 20,
        sortOrder: "desc",
      })

      expect(response.pagination.totalPages).toBe(1)
      expect(response.pagination.hasNextPage).toBe(false)
      expect(response.pagination.hasPreviousPage).toBe(false)
    })

    test("should handle empty result", () => {
      const response = Pagination.buildResponse([], 0, {
        page: 1,
        pageSize: 20,
        sortOrder: "desc",
      })

      expect(response.items).toHaveLength(0)
      expect(response.pagination.totalCount).toBe(0)
      expect(response.pagination.totalPages).toBe(0)
      expect(response.pagination.hasNextPage).toBe(false)
      expect(response.pagination.hasPreviousPage).toBe(false)
    })
  })

  describe("cursorSchema", () => {
    test("should accept valid cursor params", () => {
      const params = {
        cursor: "abc123",
        limit: 20,
        sortOrder: "desc" as const,
      }

      expect(() => Pagination.cursorSchema.parse(params)).not.toThrow()
    })

    test("should use default values", () => {
      const params = Pagination.cursorSchema.parse({})

      expect(params.limit).toBe(20)
      expect(params.sortOrder).toBe("desc")
      expect(params.cursor).toBeUndefined()
    })

    test("should reject invalid limit", () => {
      expect(() => Pagination.cursorSchema.parse({ limit: 0 })).toThrow()
      expect(() => Pagination.cursorSchema.parse({ limit: 101 })).toThrow()
    })
  })

  describe("buildCursorResponse", () => {
    test("should build response with next cursor when hasMore", () => {
      const items = Array.from({ length: 20 }, (_, i) => ({ id: `id_${i}` }))
      const response = Pagination.buildCursorResponse(items, {
        limit: 20,
        sortOrder: "desc",
      })

      expect(response.items).toHaveLength(20)
      expect(response.cursor.hasMore).toBe(true)
      expect(response.cursor.nextCursor).toBe("id_19")
    })

    test("should build response with no next cursor when at end", () => {
      const items = Array.from({ length: 15 }, (_, i) => ({ id: `id_${i}` }))
      const response = Pagination.buildCursorResponse(items, {
        limit: 20,
        sortOrder: "desc",
      })

      expect(response.items).toHaveLength(15)
      expect(response.cursor.hasMore).toBe(false)
      expect(response.cursor.nextCursor).toBeNull()
    })

    test("should handle empty result", () => {
      const response = Pagination.buildCursorResponse([], {
        limit: 20,
        sortOrder: "desc",
      })

      expect(response.items).toHaveLength(0)
      expect(response.cursor.hasMore).toBe(false)
      expect(response.cursor.nextCursor).toBeNull()
    })

    test("should handle single item result", () => {
      const items = [{ id: "id_0" }]
      const response = Pagination.buildCursorResponse(items, {
        limit: 20,
        sortOrder: "desc",
      })

      expect(response.items).toHaveLength(1)
      expect(response.cursor.hasMore).toBe(false)
      expect(response.cursor.nextCursor).toBeNull()
    })

    test("should include cursor in params", () => {
      const items = Array.from({ length: 20 }, (_, i) => ({ id: `id_${i}` }))
      const response = Pagination.buildCursorResponse(items, {
        cursor: "previous_cursor",
        limit: 20,
        sortOrder: "desc",
      })

      expect(response.cursor.hasMore).toBe(true)
    })
  })

  describe("edge cases", () => {
    test("should handle large page numbers", () => {
      const offset = Pagination.getOffset({ page: 1000, pageSize: 50, sortOrder: "desc" })

      expect(offset).toBe(49950)
    })

    test("should handle page size of 1", () => {
      const items = [{ id: 1 }]
      const response = Pagination.buildResponse(items, 100, {
        page: 50,
        pageSize: 1,
        sortOrder: "desc",
      })

      expect(response.pagination.totalPages).toBe(100)
      expect(response.pagination.page).toBe(50)
    })

    test("should calculate total pages correctly with exact division", () => {
      const items = Array.from({ length: 20 }, (_, i) => ({ id: i }))
      const response = Pagination.buildResponse(items, 100, {
        page: 1,
        pageSize: 20,
        sortOrder: "desc",
      })

      expect(response.pagination.totalPages).toBe(5)
    })

    test("should calculate total pages correctly with remainder", () => {
      const items = Array.from({ length: 20 }, (_, i) => ({ id: i }))
      const response = Pagination.buildResponse(items, 95, {
        page: 1,
        pageSize: 20,
        sortOrder: "desc",
      })

      expect(response.pagination.totalPages).toBe(5)
    })
  })
})
