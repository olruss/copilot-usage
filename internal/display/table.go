package display

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	headerColor = color.New(color.FgCyan, color.Bold)
	labelColor  = color.New(color.FgWhite, color.Bold)
	valueColor  = color.New(color.FgGreen)
	dimColor    = color.New(color.FgHiBlack)
)

// Table renders a simple aligned table to stdout.
type Table struct {
	Headers []string
	Rows    [][]string
	Title   string
}

// Render prints the table to stdout.
func (t *Table) Render() {
	if t.Title != "" {
		fmt.Println()
		headerColor.Println(t.Title)
		fmt.Println(strings.Repeat("─", len(t.Title)+10))
	}

	if len(t.Headers) == 0 || len(t.Rows) == 0 {
		dimColor.Println("  No data available")
		return
	}

	// Calculate column widths
	widths := make([]int, len(t.Headers))
	for i, h := range t.Headers {
		widths[i] = len(h)
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print header
	fmt.Println()
	for i, h := range t.Headers {
		if i > 0 {
			fmt.Print("  ")
		}
		headerColor.Printf("%-*s", widths[i], h)
	}
	fmt.Println()

	// Separator
	for i, w := range widths {
		if i > 0 {
			fmt.Print("  ")
		}
		dimColor.Print(strings.Repeat("─", w))
	}
	fmt.Println()

	// Print rows
	for _, row := range t.Rows {
		for i, cell := range row {
			if i >= len(widths) {
				break
			}
			if i > 0 {
				fmt.Print("  ")
			}
			if i == 0 {
				labelColor.Printf("%-*s", widths[i], cell)
			} else {
				fmt.Printf("%-*s", widths[i], cell)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// KV renders a key-value summary.
type KV struct {
	Title string
	Items []KVItem
}

// KVItem is a single key-value pair.
type KVItem struct {
	Key   string
	Value string
}

// Render prints the key-value summary.
func (kv *KV) Render() {
	if kv.Title != "" {
		fmt.Println()
		headerColor.Println(kv.Title)
		fmt.Println(strings.Repeat("─", len(kv.Title)+10))
	}

	maxKeyLen := 0
	for _, item := range kv.Items {
		if len(item.Key) > maxKeyLen {
			maxKeyLen = len(item.Key)
		}
	}

	for _, item := range kv.Items {
		labelColor.Printf("  %-*s  ", maxKeyLen, item.Key)
		valueColor.Println(item.Value)
	}
	fmt.Println()
}
