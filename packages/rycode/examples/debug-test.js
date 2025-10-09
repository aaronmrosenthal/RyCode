/**
 * Example Node.js program for testing the RyCode debugger
 *
 * This program intentionally has bugs to demonstrate debugging features:
 * 1. Variable undefined issue
 * 2. Logic error in calculation
 * 3. Async/await flow
 */

function calculateTotal(items, discount) {
  console.log("Calculating total...");

  if (items.length === 0) {
    return 0;
  }

  // BUG: Discount is never applied!
  const total = items.reduce((sum, item) => {
    return sum + item.price;  // Should be: sum + (item.price * discount)
  }, 0);

  console.log(`Total: $${total}`);
  return total;
}

async function fetchUserData(userId) {
  console.log(`Fetching user ${userId}...`);

  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 100));

  // BUG: email is missing from response
  return {
    id: userId,
    name: "John Doe",
    // email: "john@example.com"  // <-- Missing!
  };
}

async function processOrder(userId) {
  console.log("Processing order...");

  const user = await fetchUserData(userId);

  // BUG: This will fail because email is undefined
  console.log(`Sending receipt to: ${user.email}`);

  const items = [
    { name: "Widget", price: 10 },
    { name: "Gadget", price: 20 },
    { name: "Doohickey", price: 30 }
  ];

  const discount = 0.9;  // 10% discount
  const total = calculateTotal(items, discount);

  console.log(`Order total for ${user.name}: $${total}`);
  console.log(`Expected: $54, Got: $${total}`);

  return {
    user,
    total,
    items
  };
}

// Main execution
async function main() {
  try {
    const result = await processOrder(12345);
    console.log("Order completed:", result);
  } catch (error) {
    console.error("Order failed:", error);
  }
}

main();
