package responsive

import (
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

// Rect represents a rectangular area
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Contains checks if point is within rectangle
func (r Rect) Contains(x, y int) bool {
	return x >= r.X &&
		x < r.X+r.Width &&
		y >= r.Y &&
		y < r.Y+r.Height
}

// Intersects checks if two rectangles overlap
func (r Rect) Intersects(other Rect) bool {
	return r.X < other.X+other.Width &&
		r.X+r.Width > other.X &&
		r.Y < other.Y+other.Height &&
		r.Y+r.Height > other.Y
}

// CoordinateMapper tracks element positions in rendered output
type CoordinateMapper struct {
	elements map[string]Rect
	zIndex   map[string]int
}

// NewCoordinateMapper creates a coordinate mapper
func NewCoordinateMapper() *CoordinateMapper {
	return &CoordinateMapper{
		elements: make(map[string]Rect),
		zIndex:   make(map[string]int),
	}
}

// Register registers an element's position
func (cm *CoordinateMapper) Register(id string, rect Rect, z int) {
	cm.elements[id] = rect
	cm.zIndex[id] = z
}

// Unregister removes an element
func (cm *CoordinateMapper) Unregister(id string) {
	delete(cm.elements, id)
	delete(cm.zIndex, id)
}

// Clear removes all elements
func (cm *CoordinateMapper) Clear() {
	cm.elements = make(map[string]Rect)
	cm.zIndex = make(map[string]int)
}

// HitTest finds element at coordinates (respects z-index)
func (cm *CoordinateMapper) HitTest(x, y int) string {
	var found string
	maxZ := -999999

	for id, rect := range cm.elements {
		if rect.Contains(x, y) {
			z := cm.zIndex[id]
			if z > maxZ {
				maxZ = z
				found = id
			}
		}
	}

	return found
}

// GetBounds returns element bounds
func (cm *CoordinateMapper) GetBounds(id string) (Rect, bool) {
	rect, exists := cm.elements[id]
	return rect, exists
}

// GetAllElements returns all registered elements
func (cm *CoordinateMapper) GetAllElements() map[string]Rect {
	return cm.elements
}

// LayoutMeasurer measures rendered text and builds coordinate map
type LayoutMeasurer struct {
	mapper *CoordinateMapper
}

// NewLayoutMeasurer creates a layout measurer
func NewLayoutMeasurer() *LayoutMeasurer {
	return &LayoutMeasurer{
		mapper: NewCoordinateMapper(),
	}
}

// MeasureText analyzes rendered text and extracts positions
func (lm *LayoutMeasurer) MeasureText(rendered string) *CoordinateMapper {
	lm.mapper.Clear()

	lines := strings.Split(rendered, "\n")
	y := 0

	for _, line := range lines {
		// Strip ANSI codes for accurate width measurement
		cleanLine := stripANSI(line)
		width := visualWidth(cleanLine)

		// Look for element markers in format: <!--id:element-name-->
		// These would be embedded in the render output
		lm.scanLineForMarkers(line, y)

		y++
	}

	return lm.mapper
}

// scanLineForMarkers scans a line for position markers
func (lm *LayoutMeasurer) scanLineForMarkers(line string, y int) {
	// Look for HTML-style comments that mark element positions
	// Format: <!--ELEMENT:id:x:width:height:z-->

	if !strings.Contains(line, "<!--ELEMENT:") {
		return
	}

	parts := strings.Split(line, "<!--ELEMENT:")
	for _, part := range parts[1:] {
		endIdx := strings.Index(part, "-->")
		if endIdx == -1 {
			continue
		}

		marker := part[:endIdx]
		lm.parseMarker(marker, y)
	}
}

// parseMarker parses position marker
func (lm *LayoutMeasurer) parseMarker(marker string, currentY int) {
	// Format: id:x:width:height:z
	parts := strings.Split(marker, ":")
	if len(parts) < 5 {
		return
	}

	id := parts[0]
	x := atoi(parts[1])
	width := atoi(parts[2])
	height := atoi(parts[3])
	z := atoi(parts[4])

	rect := Rect{
		X:      x,
		Y:      currentY,
		Width:  width,
		Height: height,
	}

	lm.mapper.Register(id, rect, z)
}

// PositionedElement represents an element with known position
type PositionedElement struct {
	ID      string
	Content string
	Rect    Rect
	ZIndex  int
}

// RenderWithPosition renders element and marks its position
func RenderWithPosition(id string, content string, x, y, width, height, z int) string {
	// Embed invisible marker
	marker := formatMarker(id, x, width, height, z)

	// The marker is in the content but won't be visible
	// (terminals ignore unknown escape sequences)
	return marker + content
}

// formatMarker creates position marker
func formatMarker(id string, x, width, height, z int) string {
	// Use escape sequence that terminals ignore
	// OSC (Operating System Command) for custom data
	return "\x1b]1337;PositionMarker=" + id + ";" +
		itoa(x) + ";" + itoa(width) + ";" + itoa(height) + ";" + itoa(z) + "\x07"
}

// MeasureString returns visual width of string (handles multi-byte chars)
func MeasureString(s string) int {
	clean := stripANSI(s)
	return visualWidth(clean)
}

// visualWidth calculates visual width (handling wide chars)
func visualWidth(s string) int {
	width := 0
	for _, r := range s {
		// Simplified: count most chars as 1, CJK as 2
		if r >= 0x1100 {
			width += 2 // Wide character
		} else {
			width += 1
		}
	}
	return width
}

// stripANSI removes ANSI escape codes
func stripANSI(s string) string {
	result := strings.Builder{}
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inEscape = true
			continue
		}

		if inEscape {
			if (s[i] >= 'A' && s[i] <= 'Z') ||
				(s[i] >= 'a' && s[i] <= 'z') ||
				s[i] == 'm' {
				inEscape = false
			}
			continue
		}

		result.WriteByte(s[i])
	}

	return result.String()
}

// LayoutBuilder helps build layouts with automatic position tracking
type LayoutBuilder struct {
	currentX int
	currentY int
	mapper   *CoordinateMapper
	maxWidth int
	elements []PositionedElement
}

// NewLayoutBuilder creates a layout builder
func NewLayoutBuilder(maxWidth int) *LayoutBuilder {
	return &LayoutBuilder{
		mapper:   NewCoordinateMapper(),
		maxWidth: maxWidth,
		elements: []PositionedElement{},
	}
}

// AddElement adds an element and tracks its position
func (lb *LayoutBuilder) AddElement(id, content string, width, height, z int) {
	rect := Rect{
		X:      lb.currentX,
		Y:      lb.currentY,
		Width:  width,
		Height: height,
	}

	lb.mapper.Register(id, rect, z)

	lb.elements = append(lb.elements, PositionedElement{
		ID:      id,
		Content: content,
		Rect:    rect,
		ZIndex:  z,
	})

	// Update position for next element
	lb.currentX += width
}

// NewLine moves to next line
func (lb *LayoutBuilder) NewLine() {
	lb.currentX = 0
	lb.currentY++
}

// NewLines moves down multiple lines
func (lb *LayoutBuilder) NewLines(n int) {
	lb.currentX = 0
	lb.currentY += n
}

// SetPosition sets current position
func (lb *LayoutBuilder) SetPosition(x, y int) {
	lb.currentX = x
	lb.currentY = y
}

// GetMapper returns the coordinate mapper
func (lb *LayoutBuilder) GetMapper() *CoordinateMapper {
	return lb.mapper
}

// GetElements returns all positioned elements
func (lb *LayoutBuilder) GetElements() []PositionedElement {
	return lb.elements
}

// Render renders all elements with proper positioning
func (lb *LayoutBuilder) Render() string {
	if len(lb.elements) == 0 {
		return ""
	}

	// Calculate total height
	maxY := 0
	for _, elem := range lb.elements {
		if elem.Rect.Y+elem.Rect.Height > maxY {
			maxY = elem.Rect.Y + elem.Rect.Height
		}
	}

	// Create canvas
	lines := make([][]rune, maxY)
	for i := range lines {
		lines[i] = make([]rune, lb.maxWidth)
		for j := range lines[i] {
			lines[i][j] = ' '
		}
	}

	// Sort elements by z-index (lower first)
	sortedElements := make([]PositionedElement, len(lb.elements))
	copy(sortedElements, lb.elements)
	// Simple bubble sort by z-index
	for i := 0; i < len(sortedElements); i++ {
		for j := i + 1; j < len(sortedElements); j++ {
			if sortedElements[j].ZIndex < sortedElements[i].ZIndex {
				sortedElements[i], sortedElements[j] = sortedElements[j], sortedElements[i]
			}
		}
	}

	// Render each element
	for _, elem := range sortedElements {
		lb.renderElement(lines, elem)
	}

	// Convert to string
	result := strings.Builder{}
	for _, line := range lines {
		result.WriteString(string(line))
		result.WriteString("\n")
	}

	return result.String()
}

// renderElement renders an element onto the canvas
func (lb *LayoutBuilder) renderElement(canvas [][]rune, elem PositionedElement) {
	contentLines := strings.Split(stripANSI(elem.Content), "\n")

	for lineIdx, line := range contentLines {
		y := elem.Rect.Y + lineIdx
		if y >= len(canvas) {
			break
		}

		x := elem.Rect.X
		for _, r := range line {
			if x >= len(canvas[y]) {
				break
			}
			canvas[y][x] = r
			x++
		}
	}
}

// Helper functions

func atoi(s string) int {
	result := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	return result
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}

	negative := i < 0
	if negative {
		i = -i
	}

	digits := []rune{}
	for i > 0 {
		digits = append([]rune{rune('0' + i%10)}, digits...)
		i /= 10
	}

	if negative {
		digits = append([]rune{'-'}, digits...)
	}

	return string(digits)
}

// calculateTextBounds calculates bounds of text block
func calculateTextBounds(text string, x, y int) Rect {
	lines := strings.Split(text, "\n")

	maxWidth := 0
	for _, line := range lines {
		width := MeasureString(line)
		if width > maxWidth {
			maxWidth = width
		}
	}

	return Rect{
		X:      x,
		Y:      y,
		Width:  maxWidth,
		Height: len(lines),
	}
}

// WrapFocusableWithPosition wraps a focusable element with position tracking
type PositionedFocusable struct {
	inner    FocusableElement
	rect     Rect
	zIndex   int
	mapper   *CoordinateMapper
}

// NewPositionedFocusable creates a focusable with position tracking
func NewPositionedFocusable(inner FocusableElement, rect Rect, z int, mapper *CoordinateMapper) *PositionedFocusable {
	pf := &PositionedFocusable{
		inner:  inner,
		rect:   rect,
		zIndex: z,
		mapper: mapper,
	}

	// Register position
	pf.mapper.Register(inner.ID(), rect, z)

	return pf
}

// Implement FocusableElement interface
func (pf *PositionedFocusable) ID() string           { return pf.inner.ID() }
func (pf *PositionedFocusable) IsFocused() bool      { return pf.inner.IsFocused() }
func (pf *PositionedFocusable) Focus()               { pf.inner.Focus() }
func (pf *PositionedFocusable) Blur()                { pf.inner.Blur() }
func (pf *PositionedFocusable) HandleKey(key string) tea.Cmd {
	return pf.inner.HandleKey(key)
}

func (pf *PositionedFocusable) Render(theme *theme.Theme) string {
	content := pf.inner.Render(theme)

	// Update position in mapper
	pf.mapper.Register(pf.ID(), pf.rect, pf.zIndex)

	return content
}

// GetRect returns element bounds
func (pf *PositionedFocusable) GetRect() Rect {
	return pf.rect
}

// SetRect updates element bounds
func (pf *PositionedFocusable) SetRect(rect Rect) {
	pf.rect = rect
	pf.mapper.Register(pf.ID(), rect, pf.zIndex)
}
