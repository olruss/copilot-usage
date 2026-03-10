package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/olruss/copilot-usage/internal/analysis"
	"github.com/olruss/copilot-usage/internal/display"
)

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show usage summary for a period",
	Long:  "Display aggregated Copilot usage metrics for the selected time period.",
	RunE:  runSummary,
}

func init() {
	rootCmd.AddCommand(summaryCmd)
}

func runSummary(cmd *cobra.Command, args []string) error {
	days, err := fetchMetrics()
	if err != nil {
		return err
	}

	if len(days) == 0 {
		fmt.Println("No data available for the selected period.")
		return nil
	}

	s := analysis.Summarize(days)

	if jsonOut {
		return display.JSON(map[string]interface{}{
			"period": map[string]string{
				"from": days[0].Date,
				"to":   days[len(days)-1].Date,
			},
			"days":                  s.DaysCount,
			"avg_active_users":      s.AvgActiveUsers,
			"avg_engaged_users":     s.AvgEngagedUsers,
			"max_active_users":      s.MaxActiveUsers,
			"total_suggestions":     s.TotalSuggestions,
			"total_acceptances":     s.TotalAcceptances,
			"acceptance_rate":       s.AcceptanceRate,
			"total_lines_suggested": s.TotalLinesSuggested,
			"total_lines_accepted":  s.TotalLinesAccepted,
			"total_chats":           s.TotalChats,
			"total_dotcom_chats":    s.TotalDotcomChats,
			"total_chat_insertions": s.TotalChatInsertions,
			"total_chat_copies":     s.TotalChatCopies,
			"total_pr_summaries":    s.TotalPRSummaries,
		})
	}

	periodLabel := fmt.Sprintf("%s to %s (%d days)", days[0].Date, days[len(days)-1].Date, s.DaysCount)

	kv := &display.KV{
		Title: "Copilot Usage Summary",
		Items: []display.KVItem{
			{Key: "Period", Value: periodLabel},
			{Key: "Avg Active Users", Value: fmt.Sprintf("%.1f", s.AvgActiveUsers)},
			{Key: "Avg Engaged Users", Value: fmt.Sprintf("%.1f", s.AvgEngagedUsers)},
			{Key: "Max Active Users", Value: fmt.Sprintf("%d", s.MaxActiveUsers)},
			{Key: "Suggestions", Value: formatNum(s.TotalSuggestions)},
			{Key: "Acceptances", Value: formatNum(s.TotalAcceptances)},
			{Key: "Acceptance Rate", Value: fmt.Sprintf("%.1f%%", s.AcceptanceRate)},
			{Key: "Lines Suggested", Value: formatNum(s.TotalLinesSuggested)},
			{Key: "Lines Accepted", Value: formatNum(s.TotalLinesAccepted)},
			{Key: "IDE Chats", Value: formatNum(s.TotalChats)},
			{Key: "Dotcom Chats", Value: formatNum(s.TotalDotcomChats)},
			{Key: "Chat Insertions", Value: formatNum(s.TotalChatInsertions)},
			{Key: "Chat Copies", Value: formatNum(s.TotalChatCopies)},
			{Key: "PR Summaries", Value: formatNum(s.TotalPRSummaries)},
		},
	}
	kv.Render()

	return nil
}
