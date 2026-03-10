package display

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const maxBarWidth = 40

// BarChartItem represents a single bar in the chart.
type BarChartItem struct {
	Label string
	Value int
	Extra string // optional annotation (e.g., percentage)
}

// BarChart renders a horizontal ASCII bar chart.
func BarChart(title string, items []BarChartItem) {
	if len(items) == 0 {
		return
	}

	fmt.Println()
	headerColor.Println(title)
	fmt.Println(strings.Repeat("─", len(title)+10))

	// Find max label width and max value
	maxLabel := 0
	maxVal := 0
	for _, item := range items {
		if len(item.Label) > maxLabel {
			maxLabel = len(item.Label)
		}
		if item.Value > maxVal {
			maxVal = item.Value
		}
	}

	if maxLabel > 20 {
		maxLabel = 20
	}

	barColor := color.New(color.FgGreen)
	numColor := color.New(color.FgYellow)

	for _, item := range items {
		label := item.Label
		if len(label) > 20 {
			label = label[:17] + "..."
		}

		barLen := 0
		if maxVal > 0 {
			barLen = int(float64(item.Value) / float64(maxVal) * maxBarWidth)
		}
		if barLen == 0 && item.Value > 0 {
			barLen = 1
		}

		labelColor.Printf("  %-*s ", maxLabel, label)
		barColor.Print(strings.Repeat("█", barLen))
		fmt.Print(" ")
		numColor.Printf("%d", item.Value)
		if item.Extra != "" {
			dimColor.Printf(" %s", item.Extra)
		}
		fmt.Println()
	}
	fmt.Println()
}
