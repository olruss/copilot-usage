package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/olruss/copilot-usage/internal/analysis"
	"github.com/olruss/copilot-usage/internal/display"
)

var languagesCmd = &cobra.Command{
	Use:     "languages",
	Aliases: []string{"langs"},
	Short:   "Show language breakdown",
	Long:    "Display Copilot code completion metrics broken down by programming language.",
	RunE:    runLanguages,
}

func init() {
	rootCmd.AddCommand(languagesCmd)
}

func runLanguages(cmd *cobra.Command, args []string) error {
	days, err := fetchMetrics()
	if err != nil {
		return err
	}

	if len(days) == 0 {
		fmt.Println("No data available for the selected period.")
		return nil
	}

	langs := analysis.ByLanguage(days)

	if jsonOut {
		return display.JSON(langs)
	}

	// Table view
	t := &display.Table{
		Title:   fmt.Sprintf("Language Breakdown (%s to %s)", days[0].Date, days[len(days)-1].Date),
		Headers: []string{"Language", "Suggestions", "Acceptances", "Rate", "Lines Accepted"},
	}

	for _, l := range langs {
		t.Rows = append(t.Rows, []string{
			l.Name,
			formatNum(l.Suggestions),
			formatNum(l.Acceptances),
			fmt.Sprintf("%.1f%%", l.AcceptanceRate),
			formatNum(l.LinesAccepted),
		})
	}
	t.Render()

	// Bar chart of acceptances
	items := make([]display.BarChartItem, 0, len(langs))
	for _, l := range langs {
		items = append(items, display.BarChartItem{
			Label: l.Name,
			Value: l.Acceptances,
			Extra: fmt.Sprintf("(%.1f%% rate)", l.AcceptanceRate),
		})
	}
	display.BarChart("Acceptances by Language", items)

	return nil
}
