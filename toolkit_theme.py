"""
Neural Black Theme for Toolkit-CLI Rich Output
===============================================

This theme provides the distinctive "AI in control" aesthetic for all
toolkit-cli Rich console output with high-contrast colors optimized
for terminal readability.

Color Palette:
- Background: #0B0C10 (Deep matte black)
- Primary: #00FF88 (Electric mint - success, alignment, intelligence)
- Secondary: #00AEEF (Cool cyan - system, context, thought)
- Warning: #FFD166 (Gold - drift, caution)
- Error: #FF5555 (Neon red - failure or interruption)
- Text: #E0E0E0 (Soft white)
- Muted: #7B7F8B (Muted gray - secondary info)
- Purple: #B26FFF (Educational or persona commentary)

Usage:
    from rich.console import Console
    from toolkit_theme import toolkit_theme

    console = Console(theme=toolkit_theme)
    console.print("✅ Phase 3 Complete", style="success")
"""

from rich.theme import Theme

toolkit_theme = Theme({
    # Status indicators
    "success": "bold #00FF88",
    "warning": "bold #FFD166",
    "error": "bold #FF5555",
    "info": "bold #00AEEF",

    # Contextual elements
    "context": "#00AEEF",
    "comment": "dim #7B7F8B",
    "meta": "italic #7B7F8B",
    "educational": "#B26FFF",

    # Diff output
    "diff.add": "black on #00FF88",
    "diff.remove": "black on #FF4C6A",
    "diff.context": "#E0E0E0",

    # UI elements
    "prompt": "bold #00AEEF",
    "title": "bold #00AEEF",
    "subtitle": "#7B7F8B",
    "file": "#00AEEF",
    "path": "#00FF88",

    # Code syntax
    "code.keyword": "#00AEEF",
    "code.function": "#00FF88",
    "code.string": "#FFD166",
    "code.number": "#B26FFF",
    "code.comment": "#7B7F8B",
    "code.type": "#00AEEF",
    "code.operator": "#00FF88",

    # Progress and status
    "progress.description": "#E0E0E0",
    "progress.percentage": "#00FF88",
    "progress.download": "#00AEEF",
    "progress.remaining": "#7B7F8B",

    # Tables
    "table.header": "bold #00AEEF",
    "table.border": "#3A3A3A",
    "table.cell": "#E0E0E0",

    # Special elements
    "link": "underline #00FF88",
    "highlight": "reverse #00FF88",
    "emphasis": "#FFD166",
    "strong": "bold #00FF88",
})

# Status emoji mappings for consistent output
STATUS_ICONS = {
    "success": "✅",
    "warning": "⚠️",
    "error": "❌",
    "info": "💡",
    "context": "🧩",
    "learning": "💡",
    "drift": "⚠️",
    "ship": "🚀",
    "thought": "🧠",
    "lambda": "λ",  # Prompt symbol
}

# Custom progress bar style
PROGRESS_BAR_STYLE = {
    "bar.complete": "#00FF88",
    "bar.finished": "#00FF88",
    "bar.pulse": "#00AEEF",
    "progress.elapsed": "#7B7F8B",
}

def get_console():
    """
    Get a pre-configured Rich console with Neural Black theme.

    Returns:
        Console: Configured Rich console instance
    """
    from rich.console import Console
    return Console(theme=toolkit_theme, highlight=False)

def print_status(message: str, status: str = "info"):
    """
    Print a status message with appropriate icon and styling.

    Args:
        message: The message to display
        status: Status type (success, warning, error, info)
    """
    console = get_console()
    icon = STATUS_ICONS.get(status, "")
    style = status
    console.print(f"{icon} {message}", style=style)

def print_section(title: str, body: str = ""):
    """
    Print a section header with optional body.

    Args:
        title: Section title
        body: Optional body text
    """
    console = get_console()
    console.print(f"\n{title}", style="title")
    if body:
        console.print(body, style="context")

if __name__ == "__main__":
    # Theme demonstration
    console = get_console()

    console.print("\n🧠 Neural Black Theme Demo\n", style="title")

    console.print("✅ Phase 3 Complete", style="success")
    console.print("⚠️ Spec alignment drift detected", style="warning")
    console.print("❌ Test failed with 3 errors", style="error")
    console.print("💡 Context refresh recommended", style="info")

    console.print("\nDiff Example:", style="subtitle")
    console.print("- Old value removed", style="diff.remove")
    console.print("+ New value added", style="diff.add")

    console.print("\nCode Example:", style="subtitle")
    console.print("def process_data(items: List[str]) -> Dict:", style="code.function")
    console.print("    # Process the items", style="code.comment")
    console.print('    return {"status": "success"}', style="code.keyword")

    console.print("\n🚀 Theme loaded successfully\n", style="success")
