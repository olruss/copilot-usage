package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/olruss/copilot-usage/internal/analysis"
	"github.com/olruss/copilot-usage/internal/display"
)

var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Show day-by-day breakdown",
	Long:  "Display Copilot usage metrics for each day in the selected period.",
	RunE:  runDaily,
}

func init() {
	rootCmd.AddCommand(dailyCmd)
}

func runDaily(cmd *cobra.Command, args []string) error {
	days, err := fetchMetrics()
	if err != nil {
		return err
	}

	if len(days) == 0 {
		fmt.Println("No data available for the selected period.")
		return nil
	}

	breakdown := analysis.DailyBreakdown(days)

	if jsonOut {
		return display.JSON(breakdown)
	}

	t := &display.Table{
		Title:   "Daily Copilot Usage",
		Headers: []string{"Date", "Active", "Engaged", "Suggestions", "Acceptances", "Rate", "Lines", "Chats"},
	}

	for _, d := range breakdown {
		t.Rows = append(t.Rows, []string{
			d.Date,
			fmt.Sprintf("%d", d.ActiveUsers),
			fmt.Sprintf("%d", d.EngagedUsers),
			formatNum(d.Suggestions),
			formatNum(d.Acceptances),
			fmt.Sprintf("%.1f%%", d.AcceptanceRate),
			formatNum(d.LinesAccepted),
			fmt.Sprintf("%d", d.Chats),
		})
	}

	t.Render()
	return nil
}
