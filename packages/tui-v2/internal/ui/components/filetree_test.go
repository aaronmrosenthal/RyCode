package components

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper function to create test directory structure
func createTestDir(t *testing.T) string {
	tmpDir := t.TempDir()

	// Create directory structure
	dirs := []string{
		"src",
		"src/components",
		"src/models",
		"test",
		".git",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
	}

	// Create files
	files := []string{
		"README.md",
		"go.mod",
		"src/main.go",
		"src/components/button.go",
		"src/components/input.go",
		"src/models/user.go",
		"test/main_test.go",
		".gitignore",
	}

	for _, file := range files {
		path := filepath.Join(tmpDir, file)
		err := os.WriteFile(path, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	return tmpDir
}

func TestNewTreeNode(t *testing.T) {
	tmpDir := createTestDir(t)

	node, err := NewTreeNode(tmpDir, 0)
	if err != nil {
		t.Fatalf("Failed to create tree node: %v", err)
	}

	if !node.IsDir {
		t.Error("Expected node to be a directory")
	}
	if node.Level != 0 {
		t.Errorf("Expected level 0, got %d", node.Level)
	}
	if node.Expanded {
		t.Error("Expected node to be collapsed initially")
	}
	if node.Selected {
		t.Error("Expected node to be unselected initially")
	}
}

func TestTreeNode_Icon(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{"Go file", "main.go", "🔷"},
		{"JavaScript file", "app.js", "📜"},
		{"TypeScript file", "app.tsx", "📜"},
		{"Python file", "script.py", "🐍"},
		{"Rust file", "main.rs", "🦀"},
		{"JSON file", "config.json", "📋"},
		{"YAML file", "config.yml", "⚙️"},
		{"Markdown file", "README.md", "📝"},
		{"Env file", ".env", "🔐"},
		{"Other file", "unknown.xyz", "📄"},
	}

	tmpDir := createTestDir(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.filename)
			os.WriteFile(path, []byte("test"), 0644)

			node, err := NewTreeNode(path, 0)
			if err != nil {
				t.Fatalf("Failed to create node: %v", err)
			}

			icon := node.Icon()
			if icon != tt.expected {
				t.Errorf("Expected icon %s, got %s", tt.expected, icon)
			}
		})
	}
}

func TestTreeNode_DirectoryIcon(t *testing.T) {
	tmpDir := createTestDir(t)

	node, err := NewTreeNode(tmpDir, 0)
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}

	// Collapsed directory
	icon := node.Icon()
	if icon != "📁" {
		t.Errorf("Expected collapsed icon 📁, got %s", icon)
	}

	// Expanded directory
	node.Expanded = true
	icon = node.Icon()
	if icon != "📂" {
		t.Errorf("Expected expanded icon 📂, got %s", icon)
	}
}

func TestGitStatus_String(t *testing.T) {
	tests := []struct {
		status   GitStatus
		expected string
	}{
		{GitUntracked, "?"},
		{GitModified, "M"},
		{GitAdded, "A"},
		{GitDeleted, "D"},
		{GitRenamed, "R"},
		{GitClean, "✓"},
		{GitIgnored, "•"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.status.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestNewFileTree(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)

	if ft.Root == nil {
		t.Fatal("Expected root to be initialized")
	}
	if ft.Width != 80 {
		t.Errorf("Expected width 80, got %d", ft.Width)
	}
	if ft.Height != 24 {
		t.Errorf("Expected height 24, got %d", ft.Height)
	}
	if ft.SelectedIndex != 0 {
		t.Error("Expected selected index to be 0")
	}
	if !ft.Root.Expanded {
		t.Error("Expected root to be expanded")
	}
}

func TestFileTree_LoadChildren(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)

	// Root should have children loaded
	if len(ft.Root.Children) == 0 {
		t.Error("Expected root to have children")
	}

	// Children should be sorted (directories first, then alphabetically)
	foundSrc := false
	for _, child := range ft.Root.Children {
		if child.Name == "src" && child.IsDir {
			foundSrc = true
			break
		}
	}

	if !foundSrc {
		t.Error("Expected to find 'src' directory in children")
	}
}

func TestFileTree_ShowHidden(t *testing.T) {
	tmpDir := createTestDir(t)

	// Without hidden files
	ft := NewFileTree(tmpDir, 80, 24)
	initialCount := len(ft.Root.Children)

	// Count hidden files
	hiddenCount := 0
	for _, child := range ft.Root.Children {
		if strings.HasPrefix(child.Name, ".") {
			hiddenCount++
		}
	}

	if hiddenCount > 0 {
		t.Error("Expected no hidden files when ShowHidden is false")
	}

	// With hidden files
	ft.ShowHidden = true
	ft.Refresh()

	newCount := len(ft.Root.Children)
	if newCount <= initialCount {
		t.Error("Expected more children when showing hidden files")
	}

	// Should now have .git and .gitignore
	foundGit := false
	foundGitignore := false
	for _, child := range ft.Root.Children {
		if child.Name == ".git" {
			foundGit = true
		}
		if child.Name == ".gitignore" {
			foundGitignore = true
		}
	}

	if !foundGit {
		t.Error("Expected to find .git when ShowHidden is true")
	}
	if !foundGitignore {
		t.Error("Expected to find .gitignore when ShowHidden is true")
	}
}

func TestFileTree_FlatList(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)

	// Flat list should be populated
	if len(ft.FlatList) == 0 {
		t.Error("Expected flat list to be populated")
	}

	// First item should be root
	if ft.FlatList[0] != ft.Root {
		t.Error("Expected first item to be root")
	}

	// All items should have increasing or equal levels
	for i := 1; i < len(ft.FlatList); i++ {
		levelDiff := ft.FlatList[i].Level - ft.FlatList[i-1].Level
		if levelDiff > 1 {
			t.Errorf("Level jump too large at index %d", i)
		}
	}
}

func TestFileTree_ToggleExpanded(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)

	// Find the 'src' directory
	srcIndex := -1
	for i, node := range ft.FlatList {
		if node.Name == "src" && node.IsDir {
			srcIndex = i
			break
		}
	}

	if srcIndex == -1 {
		t.Fatal("Could not find 'src' directory")
	}

	// Select and expand
	ft.SelectedIndex = srcIndex
	initialCount := len(ft.FlatList)

	err := ft.ToggleExpanded()
	if err != nil {
		t.Fatalf("Failed to toggle expanded: %v", err)
	}

	// Should have more items now
	if len(ft.FlatList) <= initialCount {
		t.Error("Expected more items after expanding")
	}

	// Should be expanded
	if !ft.FlatList[srcIndex].Expanded {
		t.Error("Expected node to be expanded")
	}

	// Toggle again to collapse
	ft.SelectedIndex = srcIndex
	err = ft.ToggleExpanded()
	if err != nil {
		t.Fatalf("Failed to toggle collapsed: %v", err)
	}

	// Should have original count
	if len(ft.FlatList) != initialCount {
		t.Errorf("Expected %d items after collapsing, got %d", initialCount, len(ft.FlatList))
	}

	// Should be collapsed
	if ft.FlatList[srcIndex].Expanded {
		t.Error("Expected node to be collapsed")
	}
}

func TestFileTree_SelectNext(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)

	initialIndex := ft.SelectedIndex
	ft.SelectNext()

	if ft.SelectedIndex != initialIndex+1 {
		t.Errorf("Expected index %d, got %d", initialIndex+1, ft.SelectedIndex)
	}
}

func TestFileTree_SelectPrev(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)

	ft.SelectedIndex = 2
	ft.SelectPrev()

	if ft.SelectedIndex != 1 {
		t.Errorf("Expected index 1, got %d", ft.SelectedIndex)
	}

	// Should not go below 0
	ft.SelectedIndex = 0
	ft.SelectPrev()

	if ft.SelectedIndex != 0 {
		t.Error("Expected index to stay at 0")
	}
}

func TestFileTree_SelectFirst(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)

	ft.SelectedIndex = 5
	ft.SelectFirst()

	if ft.SelectedIndex != 0 {
		t.Errorf("Expected index 0, got %d", ft.SelectedIndex)
	}
	if ft.ScrollOffset != 0 {
		t.Errorf("Expected scroll offset 0, got %d", ft.ScrollOffset)
	}
}

func TestFileTree_SelectLast(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)

	ft.SelectLast()

	expectedIndex := len(ft.FlatList) - 1
	if ft.SelectedIndex != expectedIndex {
		t.Errorf("Expected index %d, got %d", expectedIndex, ft.SelectedIndex)
	}
}

func TestFileTree_GoToParent(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)

	// Find and expand 'src' directory
	srcIndex := -1
	for i, node := range ft.FlatList {
		if node.Name == "src" && node.IsDir {
			srcIndex = i
			break
		}
	}

	if srcIndex == -1 {
		t.Fatal("Could not find 'src' directory")
	}

	ft.SelectedIndex = srcIndex
	ft.ToggleExpanded()

	// Now select a child
	ft.SelectedIndex = srcIndex + 1

	// Go to parent should select 'src'
	ft.GoToParent()

	if ft.SelectedIndex != srcIndex {
		t.Errorf("Expected to select parent at %d, got %d", srcIndex, ft.SelectedIndex)
	}
}

func TestFileTree_GetSelected(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)

	selected := ft.GetSelected()
	if selected == nil {
		t.Fatal("Expected to get selected node")
	}

	if selected != ft.FlatList[ft.SelectedIndex] {
		t.Error("Selected node does not match")
	}
}

func TestFileTree_SetGitStatus(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)

	testPath := filepath.Join(tmpDir, "README.md")
	ft.SetGitStatus(testPath, GitModified)

	if status, ok := ft.GitStatusMap[testPath]; !ok || status != GitModified {
		t.Error("Failed to set git status")
	}

	// Find the node and check status
	for _, node := range ft.FlatList {
		if node.Path == testPath {
			if node.GitStatus != GitModified {
				t.Errorf("Expected GitModified status, got %v", node.GitStatus)
			}
			return
		}
	}

	t.Error("Could not find node with test path")
}

func TestFileTree_Render(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)

	output := ft.Render()

	if output == "" {
		t.Error("Expected non-empty render output")
	}

	// Should contain some file/directory names
	if !strings.Contains(output, "src") {
		t.Error("Expected output to contain 'src'")
	}
}

func TestFileTree_SetWidth(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)
	ft.SetWidth(100)

	if ft.Width != 100 {
		t.Errorf("Expected width 100, got %d", ft.Width)
	}
}

func TestFileTree_SetHeight(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)
	ft.SetHeight(30)

	if ft.Height != 30 {
		t.Errorf("Expected height 30, got %d", ft.Height)
	}
}

func TestFileTree_EnsureVisible(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 10) // Small height

	// Select an item beyond visible range
	ft.SelectedIndex = 15
	ft.ensureVisible()

	visibleHeight := ft.Height - 2
	if ft.SelectedIndex < ft.ScrollOffset || ft.SelectedIndex >= ft.ScrollOffset+visibleHeight {
		t.Error("Selected item is not visible after ensureVisible")
	}
}

func TestFileTree_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir() // Empty directory

	ft := NewFileTree(tmpDir, 80, 24)

	if ft.Root == nil {
		t.Fatal("Expected root to be initialized even for empty directory")
	}

	if len(ft.Root.Children) != 0 {
		t.Error("Expected no children for empty directory")
	}
}

func TestFileTree_ToggleHidden(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)

	initialHidden := ft.ShowHidden
	ft.ToggleHidden()

	if ft.ShowHidden == initialHidden {
		t.Error("Expected ShowHidden to toggle")
	}
}

func TestFileTree_SelectNextAtEnd(t *testing.T) {
	tmpDir := createTestDir(t)

	ft := NewFileTree(tmpDir, 80, 24)

	// Go to last item
	ft.SelectLast()
	lastIndex := ft.SelectedIndex

	// Try to go next
	ft.SelectNext()

	// Should stay at last index
	if ft.SelectedIndex != lastIndex {
		t.Error("Expected to stay at last index")
	}
}
