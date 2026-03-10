from __future__ import annotations

import json
import sys

from rich.console import Console
from rich.table import Table
from rich.text import Text

console = Console()

MAX_BAR_WIDTH = 40


def render_table(title: str, headers: list[str], rows: list[list[str]]):
    table = Table(title=title, show_lines=False)
    for i, h in enumerate(headers):
        justify = "left" if i == 0 else "right"
        table.add_column(h, justify=justify, style="cyan bold" if i == 0 else None)
    for row in rows:
        table.add_row(*row)
    console.print()
    console.print(table)


def render_kv(title: str, items: list[tuple[str, str]]):
    table = Table(title=title, show_header=False, show_lines=False)
    table.add_column("Key", style="bold white", min_width=20)
    table.add_column("Value", style="green")
    for key, value in items:
        table.add_row(key, value)
    console.print()
    console.print(table)


def render_bar_chart(title: str, items: list[tuple[str, int, str]]):
    """Render a horizontal bar chart.

    items: list of (label, value, extra_annotation)
    """
    if not items:
        return

    max_val = max(v for _, v, _ in items)
    max_label = min(max(len(label) for label, _, _ in items), 20)

    console.print()
    console.print(f"[bold cyan]{title}[/]")
    console.print("─" * (len(title) + 10))

    for label, value, extra in items:
        display_label = label[:17] + "..." if len(label) > 20 else label

        bar_len = 0
        if max_val > 0:
            bar_len = int(value / max_val * MAX_BAR_WIDTH)
        if bar_len == 0 and value > 0:
            bar_len = 1

        bar = "█" * bar_len
        line = Text()
        line.append(f"  {display_label:<{max_label}} ", style="bold white")
        line.append(bar, style="green")
        line.append(f" {value}", style="yellow")
        if extra:
            line.append(f" {extra}", style="dim")
        console.print(line)

    console.print()


def render_json(data):
    json.dump(data, sys.stdout, indent=2, default=str)
    sys.stdout.write("\n")
