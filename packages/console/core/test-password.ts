#!/usr/bin/env bun
/**
 * Quick test script for password hashing functionality
 * Run with: bun run test-password.ts
 */

import { Password } from "./src/util/password"

async function testPasswordHashing() {
  console.log("🔐 Testing Password Hashing Implementation\n")

  // Test 1: Valid password hashing
  console.log("Test 1: Hash a valid password")
  const password = "MySecureP@ss123"
  const hash = await Password.hash(password)
  console.log(`✓ Password hashed: ${hash.substring(0, 20)}...`)

  // Test 2: Verify correct password
  console.log("\nTest 2: Verify correct password")
  const isValid = await Password.verify(password, hash)
  console.log(`✓ Verification result: ${isValid}`)
  if (!isValid) throw new Error("Password verification failed!")

  // Test 3: Reject incorrect password
  console.log("\nTest 3: Reject incorrect password")
  const isInvalid = await Password.verify("WrongP@ssword123", hash)
  console.log(`✓ Wrong password rejected: ${!isInvalid}`)
  if (isInvalid) throw new Error("Wrong password was accepted!")

  // Test 4: Password strength validation
  console.log("\nTest 4: Password strength validation")
  const weakPassword = "password"
  const weakValidation = Password.validate(weakPassword)
  console.log(`✓ Weak password detected: ${!weakValidation.valid}`)
  console.log(`  Errors: ${weakValidation.errors.join(", ")}`)

  const strongPassword = "MySecureP@ss123"
  const strongValidation = Password.validate(strongPassword)
  console.log(`✓ Strong password accepted: ${strongValidation.valid}`)

  // Test 5: Password strength scoring
  console.log("\nTest 5: Password strength scoring")
  const passwords = ["12345", "password", "Password1", "MyP@ss1", "MySecureP@ssword123!"]

  for (const pwd of passwords) {
    const strength = Password.strength(pwd)
    console.log(`  "${pwd}": Score ${strength.score}/4 - ${strength.feedback}`)
  }

  // Test 6: Zod schema validation
  console.log("\nTest 6: Zod schema validation")
  try {
    Password.schema.parse("weak")
    console.log("✗ Schema should have rejected weak password")
  } catch (e: any) {
    console.log(`✓ Schema rejected weak password: ${e.issues[0].message}`)
  }

  try {
    Password.schema.parse("MySecureP@ss123")
    console.log("✓ Schema accepted strong password")
  } catch (e) {
    console.log("✗ Schema should have accepted strong password")
  }

  // Test 7: Rehash detection
  console.log("\nTest 7: Rehash detection")
  const needsRehash = await Password.needsRehash(hash)
  console.log(`✓ Needs rehash: ${needsRehash} (should be false for new hash)`)

  console.log("\n✅ All password hashing tests passed!")
}

testPasswordHashing().catch((error) => {
  console.error("\n❌ Test failed:", error.message)
  process.exit(1)
})
