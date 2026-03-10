package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/olruss/copilot-usage/internal/analysis"
	"github.com/olruss/copilot-usage/internal/api"
	"github.com/olruss/copilot-usage/internal/auth"
	"github.com/olruss/copilot-usage/internal/display"
)

var trendsCmd = &cobra.Command{
	Use:   "trends",
	Short: "Show week-over-week trends",
	Long:  "Compare this week's Copilot usage with last week, showing deltas.",
	RunE:  runTrends,
}

func init() {
	rootCmd.AddCommand(trendsCmd)
}

func runTrends(cmd *cobra.Command, args []string) error {
	tok, err := auth.GetToken(token)
	if err != nil {
		return err
	}

	// Fetch last 21 days to ensure we have data for both weeks
	now := time.Now()
	s := now.AddDate(0, 0, -21).Format("2006-01-02")
	u := now.Format("2006-01-02")

	client := api.NewClient(tok)
	days, err := client.FetchMetrics(org, s, u)
	if err != nil {
		return err
	}

	if len(days) == 0 {
		fmt.Println("No data available.")
		return nil
	}

	thisWeek, lastWeek := analysis.WeekOverWeek(days)

	if jsonOut {
		return display.JSON(map[string]interface{}{
			"this_week": thisWeek,
			"last_week": lastWeek,
		})
	}

	t := &display.Table{
		Title:   "Week-over-Week Trends",
		Headers: []string{"Metric", "Last Week", "This Week", "Delta"},
	}

	t.Rows = append(t.Rows,
		trendRow("Days with data", lastWeek.Days, thisWeek.Days),
		trendRow("Avg Active Users", int(lastWeek.AvgActiveUsers), int(thisWeek.AvgActiveUsers)),
		trendRow("Suggestions", lastWeek.TotalSuggestions, thisWeek.TotalSuggestions),
		trendRow("Acceptances", lastWeek.TotalAcceptances, thisWeek.TotalAcceptances),
		trendRowPct("Acceptance Rate", lastWeek.AcceptanceRate, thisWeek.AcceptanceRate),
		trendRow("Chats", lastWeek.TotalChats, thisWeek.TotalChats),
	)

	t.Render()

	fmt.Printf("  Last week: %s to %s\n", lastWeek.StartDate, lastWeek.EndDate)
	fmt.Printf("  This week: %s to %s\n\n", thisWeek.StartDate, thisWeek.EndDate)

	return nil
}

func trendRow(label string, last, this int) []string {
	delta := this - last
	deltaStr := formatDelta(delta)
	return []string{label, formatNum(last), formatNum(this), deltaStr}
}

func trendRowPct(label string, last, this float64) []string {
	delta := this - last
	sign := "+"
	if delta < 0 {
		sign = ""
	}
	deltaStr := fmt.Sprintf("%s%.1f%%", sign, delta)
	return []string{label, fmt.Sprintf("%.1f%%", last), fmt.Sprintf("%.1f%%", this), deltaStr}
}

func formatDelta(d int) string {
	up := color.New(color.FgGreen)
	down := color.New(color.FgRed)

	if d > 0 {
		return up.Sprintf("+%s", formatNum(d))
	}
	if d < 0 {
		return down.Sprintf("%s", formatNum(d))
	}
	return "0"
}

