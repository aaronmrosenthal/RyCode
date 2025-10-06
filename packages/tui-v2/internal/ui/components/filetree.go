package components

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/theme"
	"github.com/charmbracelet/lipgloss"
)

// GitStatus represents the git status of a file
type GitStatus int

const (
	GitUntracked GitStatus = iota
	GitModified
	GitAdded
	GitDeleted
	GitRenamed
	GitClean
	GitIgnored
)

// String returns the string representation of GitStatus
func (g GitStatus) String() string {
	switch g {
	case GitUntracked:
		return "?"
	case GitModified:
		return "M"
	case GitAdded:
		return "A"
	case GitDeleted:
		return "D"
	case GitRenamed:
		return "R"
	case GitClean:
		return "✓"
	case GitIgnored:
		return "•"
	default:
		return " "
	}
}

// Icon returns the colored icon for GitStatus
func (g GitStatus) Icon() string {
	style := lipgloss.NewStyle()
	switch g {
	case GitUntracked:
		return style.Foreground(theme.NeonYellow).Render("?")
	case GitModified:
		return style.Foreground(theme.NeonOrange).Render("M")
	case GitAdded:
		return style.Foreground(theme.MatrixGreen).Render("A")
	case GitDeleted:
		return style.Foreground(theme.NeonPink).Render("D")
	case GitRenamed:
		return style.Foreground(theme.NeonCyan).Render("R")
	case GitClean:
		return style.Foreground(theme.MatrixGreenDim).Render("✓")
	case GitIgnored:
		return style.Foreground(theme.MatrixGreenDark).Render("•")
	default:
		return " "
	}
}

// TreeNode represents a file or directory in the tree
type TreeNode struct {
	Path      string
	Name      string
	IsDir     bool
	Expanded  bool
	Selected  bool
	Level     int
	Children  []*TreeNode
	GitStatus GitStatus
}

// NewTreeNode creates a new tree node
func NewTreeNode(path string, level int) (*TreeNode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	name := filepath.Base(path)
	if name == "." {
		name = filepath.Base(filepath.Dir(path))
	}

	return &TreeNode{
		Path:      path,
		Name:      name,
		IsDir:     info.IsDir(),
		Expanded:  false,
		Selected:  false,
		Level:     level,
		Children:  []*TreeNode{},
		GitStatus: GitClean,
	}, nil
}

// Icon returns the file type icon
func (n *TreeNode) Icon() string {
	if n.IsDir {
		if n.Expanded {
			return "📂"
		}
		return "📁"
	}

	// File type icons based on extension
	ext := strings.ToLower(filepath.Ext(n.Name))
	switch ext {
	case ".go":
		return "🔷"
	case ".js", ".jsx", ".ts", ".tsx":
		return "📜"
	case ".py":
		return "🐍"
	case ".rs":
		return "🦀"
	case ".json":
		return "📋"
	case ".yaml", ".yml":
		return "⚙️"
	case ".md":
		return "📝"
	case ".git":
		return "🔀"
	case ".env":
		return "🔐"
	case ".docker", ".dockerfile":
		return "🐳"
	default:
		return "📄"
	}
}

// FileTree represents a file tree component
type FileTree struct {
	Root          *TreeNode
	FlatList      []*TreeNode // Flattened list for rendering
	SelectedIndex int
	ScrollOffset  int
	Width         int
	Height        int
	RootPath      string
	ShowHidden    bool
	GitStatusMap  map[string]GitStatus
}

// NewFileTree creates a new file tree
func NewFileTree(rootPath string, width, height int) *FileTree {
	ft := &FileTree{
		Root:          nil,
		FlatList:      []*TreeNode{},
		SelectedIndex: 0,
		ScrollOffset:  0,
		Width:         width,
		Height:        height,
		RootPath:      rootPath,
		ShowHidden:    false,
		GitStatusMap:  make(map[string]GitStatus),
	}

	// Build initial tree
	ft.Refresh()
	return ft
}

// Refresh rebuilds the tree from the root path
func (ft *FileTree) Refresh() error {
	root, err := NewTreeNode(ft.RootPath, 0)
	if err != nil {
		return err
	}

	root.Expanded = true
	ft.Root = root

	// Load children for root
	if err := ft.loadChildren(root); err != nil {
		return err
	}

	// Rebuild flat list
	ft.rebuildFlatList()

	// Update git status
	ft.updateGitStatus()

	return nil
}

// loadChildren loads children for a directory node
func (ft *FileTree) loadChildren(node *TreeNode) error {
	if !node.IsDir {
		return nil
	}

	entries, err := os.ReadDir(node.Path)
	if err != nil {
		return err
	}

	node.Children = []*TreeNode{}

	for _, entry := range entries {
		// Skip hidden files unless ShowHidden is true
		if !ft.ShowHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		childPath := filepath.Join(node.Path, entry.Name())
		child, err := NewTreeNode(childPath, node.Level+1)
		if err != nil {
			continue
		}

		node.Children = append(node.Children, child)
	}

	// Sort: directories first, then alphabetically
	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].IsDir != node.Children[j].IsDir {
			return node.Children[i].IsDir
		}
		return strings.ToLower(node.Children[i].Name) < strings.ToLower(node.Children[j].Name)
	})

	return nil
}

// rebuildFlatList creates a flattened list for rendering
func (ft *FileTree) rebuildFlatList() {
	ft.FlatList = []*TreeNode{}
	ft.flattenNode(ft.Root)
}

// flattenNode recursively flattens the tree
func (ft *FileTree) flattenNode(node *TreeNode) {
	if node == nil {
		return
	}

	ft.FlatList = append(ft.FlatList, node)

	if node.Expanded && node.IsDir {
		for _, child := range node.Children {
			ft.flattenNode(child)
		}
	}
}

// updateGitStatus updates git status for all files
func (ft *FileTree) updateGitStatus() {
	// TODO: Implement actual git status parsing
	// For now, just set some example statuses
	for _, node := range ft.FlatList {
		if status, ok := ft.GitStatusMap[node.Path]; ok {
			node.GitStatus = status
		} else {
			node.GitStatus = GitClean
		}
	}
}

// ToggleExpanded toggles the expanded state of the selected node
func (ft *FileTree) ToggleExpanded() error {
	if ft.SelectedIndex >= len(ft.FlatList) {
		return nil
	}

	node := ft.FlatList[ft.SelectedIndex]
	if !node.IsDir {
		return nil
	}

	node.Expanded = !node.Expanded

	// Load children if expanding
	if node.Expanded && len(node.Children) == 0 {
		if err := ft.loadChildren(node); err != nil {
			return err
		}
	}

	ft.rebuildFlatList()
	return nil
}

// SelectNext moves selection down
func (ft *FileTree) SelectNext() {
	if ft.SelectedIndex < len(ft.FlatList)-1 {
		ft.SelectedIndex++
		ft.ensureVisible()
	}
}

// SelectPrev moves selection up
func (ft *FileTree) SelectPrev() {
	if ft.SelectedIndex > 0 {
		ft.SelectedIndex--
		ft.ensureVisible()
	}
}

// SelectFirst moves selection to the first item
func (ft *FileTree) SelectFirst() {
	ft.SelectedIndex = 0
	ft.ScrollOffset = 0
}

// SelectLast moves selection to the last item
func (ft *FileTree) SelectLast() {
	ft.SelectedIndex = len(ft.FlatList) - 1
	ft.ensureVisible()
}

// GoToParent collapses current directory or moves to parent
func (ft *FileTree) GoToParent() {
	if ft.SelectedIndex >= len(ft.FlatList) {
		return
	}

	node := ft.FlatList[ft.SelectedIndex]

	// If expanded directory, collapse it
	if node.IsDir && node.Expanded {
		node.Expanded = false
		ft.rebuildFlatList()
		return
	}

	// Otherwise, find and select parent
	parentLevel := node.Level - 1
	for i := ft.SelectedIndex - 1; i >= 0; i-- {
		if ft.FlatList[i].Level == parentLevel {
			ft.SelectedIndex = i
			ft.ensureVisible()
			return
		}
	}
}

// GetSelected returns the currently selected node
func (ft *FileTree) GetSelected() *TreeNode {
	if ft.SelectedIndex >= 0 && ft.SelectedIndex < len(ft.FlatList) {
		return ft.FlatList[ft.SelectedIndex]
	}
	return nil
}

// ensureVisible ensures the selected item is visible
func (ft *FileTree) ensureVisible() {
	visibleHeight := ft.Height - 2 // Account for borders

	if ft.SelectedIndex < ft.ScrollOffset {
		ft.ScrollOffset = ft.SelectedIndex
	}

	if ft.SelectedIndex >= ft.ScrollOffset+visibleHeight {
		ft.ScrollOffset = ft.SelectedIndex - visibleHeight + 1
	}
}

// Render renders the file tree
func (ft *FileTree) Render() string {
	if ft.Root == nil {
		return theme.MatrixTheme.Error.Render("No directory loaded")
	}

	visibleHeight := ft.Height - 2 // Account for borders
	lines := []string{}

	// Calculate visible range
	start := ft.ScrollOffset
	end := ft.ScrollOffset + visibleHeight
	if end > len(ft.FlatList) {
		end = len(ft.FlatList)
	}

	// Render visible nodes
	for i := start; i < end; i++ {
		node := ft.FlatList[i]
		lines = append(lines, ft.renderNode(node, i == ft.SelectedIndex))
	}

	// Fill remaining space
	for len(lines) < visibleHeight {
		lines = append(lines, "")
	}

	content := strings.Join(lines, "\n")

	// Wrap in border
	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.MatrixGreen).
		Width(ft.Width - 2).
		Height(ft.Height - 2)

	return borderStyle.Render(content)
}

// renderNode renders a single tree node
func (ft *FileTree) renderNode(node *TreeNode, selected bool) string {
	// Indentation
	indent := strings.Repeat("  ", node.Level)

	// Expand/collapse indicator
	expandIndicator := " "
	if node.IsDir {
		if node.Expanded {
			expandIndicator = "▼"
		} else {
			expandIndicator = "▶"
		}
	}

	// File/folder icon
	icon := node.Icon()

	// Git status
	gitIcon := node.GitStatus.Icon()

	// Name
	name := node.Name
	nameStyle := lipgloss.NewStyle()

	if selected {
		nameStyle = nameStyle.
			Foreground(theme.MatrixGreen).
			Bold(true).
			Background(theme.MatrixGreenDark)
	} else if node.IsDir {
		nameStyle = nameStyle.Foreground(theme.NeonCyan)
	} else {
		nameStyle = nameStyle.Foreground(theme.MatrixGreenDim)
	}

	// Compose the line
	parts := []string{
		indent,
		expandIndicator,
		" ",
		icon,
		" ",
		nameStyle.Render(name),
		" ",
		gitIcon,
	}

	line := strings.Join(parts, "")

	// Truncate if too long
	maxWidth := ft.Width - 4 // Account for borders and padding
	if lipgloss.Width(line) > maxWidth {
		// Simple truncation for now
		if len(line) > maxWidth-3 {
			line = line[:maxWidth-3] + "..."
		}
	}

	return line
}

// SetWidth updates the width
func (ft *FileTree) SetWidth(width int) {
	ft.Width = width
}

// SetHeight updates the height
func (ft *FileTree) SetHeight(height int) {
	ft.Height = height
}

// ToggleHidden toggles showing hidden files
func (ft *FileTree) ToggleHidden() {
	ft.ShowHidden = !ft.ShowHidden
	ft.Refresh()
}

// SetGitStatus sets the git status for a path
func (ft *FileTree) SetGitStatus(path string, status GitStatus) {
	ft.GitStatusMap[path] = status
	ft.updateGitStatus()
}
