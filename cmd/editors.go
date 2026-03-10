package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/olruss/copilot-usage/internal/analysis"
	"github.com/olruss/copilot-usage/internal/display"
)

var editorsCmd = &cobra.Command{
	Use:   "editors",
	Short: "Show editor breakdown",
	Long:  "Display Copilot usage metrics broken down by IDE/editor.",
	RunE:  runEditors,
}

func init() {
	rootCmd.AddCommand(editorsCmd)
}

func runEditors(cmd *cobra.Command, args []string) error {
	days, err := fetchMetrics()
	if err != nil {
		return err
	}

	if len(days) == 0 {
		fmt.Println("No data available for the selected period.")
		return nil
	}

	editors := analysis.ByEditor(days)

	if jsonOut {
		return display.JSON(editors)
	}

	t := &display.Table{
		Title:   fmt.Sprintf("Editor Breakdown (%s to %s)", days[0].Date, days[len(days)-1].Date),
		Headers: []string{"Editor", "Users", "Suggestions", "Acceptances", "Rate", "Chats"},
	}

	for _, e := range editors {
		t.Rows = append(t.Rows, []string{
			e.Name,
			fmt.Sprintf("%d", e.EngagedUsers),
			formatNum(e.Suggestions),
			formatNum(e.Acceptances),
			fmt.Sprintf("%.1f%%", e.AcceptanceRate),
			fmt.Sprintf("%d", e.Chats),
		})
	}
	t.Render()

	// Bar chart of suggestions
	items := make([]display.BarChartItem, 0, len(editors))
	for _, e := range editors {
		items = append(items, display.BarChartItem{
			Label: e.Name,
			Value: e.Suggestions,
			Extra: fmt.Sprintf("(%d users)", e.EngagedUsers),
		})
	}
	display.BarChart("Suggestions by Editor", items)

	return nil
}
