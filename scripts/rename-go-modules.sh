#!/bin/bash

set -e

echo "🔍 Finding and updating Go modules..."
echo ""

updated_count=0
total_files=0

# Update go.mod files
echo "📦 Updating go.mod files..."
while IFS= read -r file; do
  total_files=$((total_files + 1))
  echo "Processing: $file"

  # Check if file contains opencode references
  if grep -q "github.com/sst/opencode" "$file"; then
    # Update module path
    sed -i.bak 's|github.com/sst/opencode-sdk-go|github.com/aaronmrosenthal/rycode-sdk-go|g' "$file"
    sed -i.bak 's|github.com/sst/opencode|github.com/aaronmrosenthal/rycode|g' "$file"

    # Clean up backup
    rm "${file}.bak"

    echo "  ✓ Updated module references"
    updated_count=$((updated_count + 1))
  else
    echo "  ⏭️  No changes needed"
  fi
  echo ""
done < <(find . -name "go.mod" -not -path "*/node_modules/*")

# Update Go import statements in .go files
echo "📝 Updating Go import statements..."
go_files_updated=0
go_files_total=0

while IFS= read -r file; do
  go_files_total=$((go_files_total + 1))

  # Check if file contains opencode imports
  if grep -q "github.com/sst/opencode" "$file"; then
    sed -i.bak 's|github.com/sst/opencode-sdk-go|github.com/aaronmrosenthal/rycode-sdk-go|g' "$file"
    sed -i.bak 's|github.com/sst/opencode|github.com/aaronmrosenthal/rycode|g' "$file"

    rm "${file}.bak"
    go_files_updated=$((go_files_updated + 1))
  fi
done < <(find . -name "*.go" -not -path "*/node_modules/*")

echo ""
echo "📊 Summary:"
echo "   go.mod files:"
echo "     Total: $total_files"
echo "     Updated: $updated_count"
echo "   .go files:"
echo "     Total: $go_files_total"
echo "     Updated: $go_files_updated"
echo ""

# Run go mod tidy in directories with go.mod
echo "🧹 Running go mod tidy..."
while IFS= read -r file; do
  dir=$(dirname "$file")
  echo "  Tidying: $dir"
  (cd "$dir" && go mod tidy)
done < <(find . -name "go.mod" -not -path "*/node_modules/*")

echo ""
echo "✅ Go modules updated successfully"
