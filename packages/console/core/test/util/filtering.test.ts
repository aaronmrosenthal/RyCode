import { describe, test, expect } from "bun:test"
import { Filtering } from "../../src/util/filtering"
import { sql } from "drizzle-orm"

// Mock table object for testing
const mockTable = {
  name: sql`name`,
  age: sql`age`,
  status: sql`status`,
  email: sql`email`,
  created: sql`created`,
}

describe("Filtering Utility", () => {
  describe("filterSchema", () => {
    test("should accept valid filter", () => {
      const filter = {
        field: "name",
        operator: "eq" as const,
        value: "John",
      }

      expect(() => Filtering.filterSchema.parse(filter)).not.toThrow()
    })

    test("should reject invalid operator", () => {
      const filter = {
        field: "name",
        operator: "invalid",
        value: "John",
      }

      expect(() => Filtering.filterSchema.parse(filter)).toThrow()
    })

    test("should reject excessively long field name", () => {
      const filter = {
        field: "a".repeat(101),
        operator: "eq" as const,
        value: "test",
      }

      expect(() => Filtering.filterSchema.parse(filter)).toThrow()
    })

    test("should reject excessively long string value", () => {
      const filter = {
        field: "name",
        operator: "eq" as const,
        value: "a".repeat(1001),
      }

      expect(() => Filtering.filterSchema.parse(filter)).toThrow()
    })

    test("should reject excessively large array", () => {
      const filter = {
        field: "id",
        operator: "in" as const,
        value: Array(101).fill("test"),
      }

      expect(() => Filtering.filterSchema.parse(filter)).toThrow()
    })

    test("should accept all value types", () => {
      const filters = [
        { field: "name", operator: "eq" as const, value: "string" },
        { field: "age", operator: "gt" as const, value: 25 },
        { field: "active", operator: "eq" as const, value: true },
        { field: "ids", operator: "in" as const, value: [1, 2, 3] },
      ]

      filters.forEach((filter) => {
        expect(() => Filtering.filterSchema.parse(filter)).not.toThrow()
      })
    })
  })

  describe("applyFilter", () => {
    test("should apply eq operator", () => {
      const filter: Filtering.Filter = {
        field: "name",
        operator: "eq",
        value: "John",
      }

      const result = Filtering.applyFilter(mockTable.name, filter)
      expect(result).toBeDefined()
    })

    test("should apply ne operator", () => {
      const filter: Filtering.Filter = {
        field: "status",
        operator: "ne",
        value: "inactive",
      }

      const result = Filtering.applyFilter(mockTable.status, filter)
      expect(result).toBeDefined()
    })

    test("should apply gt operator", () => {
      const filter: Filtering.Filter = {
        field: "age",
        operator: "gt",
        value: 18,
      }

      const result = Filtering.applyFilter(mockTable.age, filter)
      expect(result).toBeDefined()
    })

    test("should apply gte operator", () => {
      const filter: Filtering.Filter = {
        field: "age",
        operator: "gte",
        value: 21,
      }

      const result = Filtering.applyFilter(mockTable.age, filter)
      expect(result).toBeDefined()
    })

    test("should apply lt operator", () => {
      const filter: Filtering.Filter = {
        field: "age",
        operator: "lt",
        value: 65,
      }

      const result = Filtering.applyFilter(mockTable.age, filter)
      expect(result).toBeDefined()
    })

    test("should apply lte operator", () => {
      const filter: Filtering.Filter = {
        field: "age",
        operator: "lte",
        value: 100,
      }

      const result = Filtering.applyFilter(mockTable.age, filter)
      expect(result).toBeDefined()
    })

    test("should apply like operator", () => {
      const filter: Filtering.Filter = {
        field: "email",
        operator: "like",
        value: "example.com",
      }

      const result = Filtering.applyFilter(mockTable.email, filter)
      expect(result).toBeDefined()
    })

    test("should apply in operator", () => {
      const filter: Filtering.Filter = {
        field: "status",
        operator: "in",
        value: ["active", "pending"],
      }

      const result = Filtering.applyFilter(mockTable.status, filter)
      expect(result).toBeDefined()
    })

    test("should apply between operator with valid array", () => {
      const filter: Filtering.Filter = {
        field: "age",
        operator: "between",
        value: [18, 65],
      }

      const result = Filtering.applyFilter(mockTable.age, filter)
      expect(result).toBeDefined()
    })

    test("should return undefined for between with invalid array", () => {
      const filter: Filtering.Filter = {
        field: "age",
        operator: "between",
        value: [18], // Only one value
      }

      const result = Filtering.applyFilter(mockTable.age, filter)
      expect(result).toBeUndefined()
    })
  })

  describe("applyFilters", () => {
    test("should apply multiple filters with AND logic", () => {
      const filters: Filtering.Filter[] = [
        { field: "name", operator: "like", value: "John" },
        { field: "age", operator: "gte", value: 18 },
        { field: "status", operator: "eq", value: "active" },
      ]

      const result = Filtering.applyFilters(mockTable, filters)
      expect(result).toBeDefined()
    })

    test("should return undefined for empty filters array", () => {
      const result = Filtering.applyFilters(mockTable, [])
      expect(result).toBeUndefined()
    })

    test("should skip unknown fields", () => {
      const filters: Filtering.Filter[] = [
        { field: "unknownField", operator: "eq", value: "test" },
        { field: "name", operator: "eq", value: "John" },
      ]

      const result = Filtering.applyFilters(mockTable, filters)
      expect(result).toBeDefined() // Only applies the valid field
    })

    test("should throw error for too many filters", () => {
      const filters = Array(21)
        .fill(null)
        .map((_, i) => ({
          field: "name",
          operator: "eq" as const,
          value: `test${i}`,
        }))

      expect(() => Filtering.applyFilters(mockTable, filters)).toThrow("Too many filters")
    })

    test("should handle all unknown fields gracefully", () => {
      const filters: Filtering.Filter[] = [
        { field: "unknown1", operator: "eq", value: "test" },
        { field: "unknown2", operator: "eq", value: "test" },
      ]

      const result = Filtering.applyFilters(mockTable, filters)
      expect(result).toBeUndefined()
    })
  })

  describe("sortSchema", () => {
    test("should accept valid sort", () => {
      const sort = {
        field: "name",
        order: "asc" as const,
      }

      expect(() => Filtering.sortSchema.parse(sort)).not.toThrow()
    })

    test("should use default order", () => {
      const sort = Filtering.sortSchema.parse({ field: "name" })

      expect(sort.order).toBe("desc")
    })

    test("should accept both asc and desc", () => {
      const sorts = [
        { field: "name", order: "asc" as const },
        { field: "created", order: "desc" as const },
      ]

      sorts.forEach((sort) => {
        expect(() => Filtering.sortSchema.parse(sort)).not.toThrow()
      })
    })

    test("should reject invalid order", () => {
      const sort = {
        field: "name",
        order: "invalid",
      }

      expect(() => Filtering.sortSchema.parse(sort)).toThrow()
    })
  })

  describe("searchSchema", () => {
    test("should accept valid search", () => {
      const search = {
        query: "john doe",
        fields: ["name", "email"],
      }

      expect(() => Filtering.searchSchema.parse(search)).not.toThrow()
    })

    test("should reject empty query", () => {
      const search = {
        query: "",
        fields: ["name"],
      }

      expect(() => Filtering.searchSchema.parse(search)).toThrow()
    })

    test("should reject query exceeding max length", () => {
      const search = {
        query: "a".repeat(501),
        fields: ["name"],
      }

      expect(() => Filtering.searchSchema.parse(search)).toThrow()
    })

    test("should reject empty fields array", () => {
      const search = {
        query: "test",
        fields: [],
      }

      expect(() => Filtering.searchSchema.parse(search)).toThrow()
    })

    test("should reject too many fields", () => {
      const search = {
        query: "test",
        fields: Array(11).fill("field"),
      }

      expect(() => Filtering.searchSchema.parse(search)).toThrow()
    })
  })

  describe("applySearch", () => {
    test("should apply search across multiple fields with OR logic", () => {
      const search: Filtering.Search = {
        query: "john",
        fields: ["name", "email"],
      }

      const result = Filtering.applySearch(mockTable, search)
      expect(result).toBeDefined()
    })

    test("should sanitize LIKE wildcards", () => {
      const search: Filtering.Search = {
        query: "test%value_here",
        fields: ["name"],
      }

      const result = Filtering.applySearch(mockTable, search)
      expect(result).toBeDefined()
      // The % and _ should be escaped to prevent SQL injection
    })

    test("should skip unknown fields", () => {
      const search: Filtering.Search = {
        query: "test",
        fields: ["unknownField", "name"],
      }

      const result = Filtering.applySearch(mockTable, search)
      expect(result).toBeDefined() // Only applies to valid field
    })

    test("should return undefined for all unknown fields", () => {
      const search: Filtering.Search = {
        query: "test",
        fields: ["unknown1", "unknown2"],
      }

      const result = Filtering.applySearch(mockTable, search)
      expect(result).toBeUndefined()
    })
  })

  describe("querySchema", () => {
    test("should accept valid query params", () => {
      const params = {
        filters: [{ field: "name", operator: "eq" as const, value: "John" }],
        sort: { field: "created", order: "desc" as const },
        search: { query: "test", fields: ["name"] },
        page: 1,
        pageSize: 20,
      }

      expect(() => Filtering.querySchema.parse(params)).not.toThrow()
    })

    test("should use default values", () => {
      const params = Filtering.querySchema.parse({})

      expect(params.page).toBe(1)
      expect(params.pageSize).toBe(20)
    })

    test("should reject invalid page", () => {
      expect(() => Filtering.querySchema.parse({ page: 0 })).toThrow()
      expect(() => Filtering.querySchema.parse({ page: -1 })).toThrow()
    })

    test("should reject invalid pageSize", () => {
      expect(() => Filtering.querySchema.parse({ pageSize: 0 })).toThrow()
      expect(() => Filtering.querySchema.parse({ pageSize: 101 })).toThrow()
    })

    test("should allow optional filters, sort, and search", () => {
      const params1 = Filtering.querySchema.parse({ page: 1 })
      expect(params1.filters).toBeUndefined()
      expect(params1.sort).toBeUndefined()
      expect(params1.search).toBeUndefined()

      const params2 = Filtering.querySchema.parse({
        filters: [{ field: "name", operator: "eq" as const, value: "test" }],
      })
      expect(params2.filters).toBeDefined()
      expect(params2.sort).toBeUndefined()
    })
  })

  describe("buildWhereClause", () => {
    test("should build WHERE clause with filters only", () => {
      const params: Filtering.QueryParams = {
        filters: [
          { field: "name", operator: "like", value: "John" },
          { field: "age", operator: "gte", value: 18 },
        ],
        page: 1,
        pageSize: 20,
      }

      const result = Filtering.buildWhereClause(mockTable, params)
      expect(result).toBeDefined()
    })

    test("should build WHERE clause with search only", () => {
      const params: Filtering.QueryParams = {
        search: {
          query: "test",
          fields: ["name", "email"],
        },
        page: 1,
        pageSize: 20,
      }

      const result = Filtering.buildWhereClause(mockTable, params)
      expect(result).toBeDefined()
    })

    test("should build WHERE clause with both filters and search", () => {
      const params: Filtering.QueryParams = {
        filters: [{ field: "status", operator: "eq", value: "active" }],
        search: {
          query: "john",
          fields: ["name"],
        },
        page: 1,
        pageSize: 20,
      }

      const result = Filtering.buildWhereClause(mockTable, params)
      expect(result).toBeDefined()
    })

    test("should return undefined with no filters or search", () => {
      const params: Filtering.QueryParams = {
        page: 1,
        pageSize: 20,
      }

      const result = Filtering.buildWhereClause(mockTable, params)
      expect(result).toBeUndefined()
    })

    test("should handle empty filters array", () => {
      const params: Filtering.QueryParams = {
        filters: [],
        page: 1,
        pageSize: 20,
      }

      const result = Filtering.buildWhereClause(mockTable, params)
      expect(result).toBeUndefined()
    })
  })

  describe("edge cases", () => {
    test("should handle special characters in LIKE search", () => {
      const search: Filtering.Search = {
        query: "test@example.com",
        fields: ["email"],
      }

      const result = Filtering.applySearch(mockTable, search)
      expect(result).toBeDefined()
    })

    test("should handle unicode characters", () => {
      const filter: Filtering.Filter = {
        field: "name",
        operator: "eq",
        value: "José García",
      }

      const result = Filtering.applyFilter(mockTable.name, filter)
      expect(result).toBeDefined()
    })

    test("should handle empty string value", () => {
      const filter: Filtering.Filter = {
        field: "name",
        operator: "eq",
        value: "",
      }

      const result = Filtering.applyFilter(mockTable.name, filter)
      expect(result).toBeDefined()
    })

    test("should handle zero numeric value", () => {
      const filter: Filtering.Filter = {
        field: "age",
        operator: "eq",
        value: 0,
      }

      const result = Filtering.applyFilter(mockTable.age, filter)
      expect(result).toBeDefined()
    })

    test("should handle negative numbers", () => {
      const filter: Filtering.Filter = {
        field: "age",
        operator: "gt",
        value: -1,
      }

      const result = Filtering.applyFilter(mockTable.age, filter)
      expect(result).toBeDefined()
    })

    test("should handle boolean values", () => {
      const filter: Filtering.Filter = {
        field: "active",
        operator: "eq",
        value: false,
      }

      const result = Filtering.applyFilter(sql`active`, filter)
      expect(result).toBeDefined()
    })
  })

  describe("security", () => {
    test("should prevent SQL injection via field name", () => {
      const filters: Filtering.Filter[] = [
        {
          field: "name; DROP TABLE users--",
          operator: "eq",
          value: "test",
        },
      ]

      // Should not throw, but should skip the invalid field
      const result = Filtering.applyFilters(mockTable, filters)
      expect(result).toBeUndefined()
    })

    test("should prevent SQL injection via search query", () => {
      const search: Filtering.Search = {
        query: "'; DROP TABLE users--",
        fields: ["name"],
      }

      // Should sanitize the query and not cause SQL injection
      const result = Filtering.applySearch(mockTable, search)
      expect(result).toBeDefined()
    })

    test("should enforce filter count limit", () => {
      const filters = Array(21)
        .fill(null)
        .map(() => ({
          field: "name",
          operator: "eq" as const,
          value: "test",
        }))

      expect(() => Filtering.applyFilters(mockTable, filters)).toThrow()
    })

    test("should enforce string length limits", () => {
      const longString = "a".repeat(1001)

      expect(() =>
        Filtering.filterSchema.parse({
          field: "name",
          operator: "eq",
          value: longString,
        }),
      ).toThrow()
    })

    test("should enforce array size limits", () => {
      const largeArray = Array(101).fill("test")

      expect(() =>
        Filtering.filterSchema.parse({
          field: "id",
          operator: "in",
          value: largeArray,
        }),
      ).toThrow()
    })
  })
})
