package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/olruss/copilot-usage/internal/api"
	"github.com/olruss/copilot-usage/internal/auth"
)

var (
	org       string
	token     string
	jsonOut   bool
	since     string
	until     string
	period    string
)

var rootCmd = &cobra.Command{
	Use:   "copilot-usage",
	Short: "GitHub Copilot usage analytics for your organization",
	Long:  "Analyze GitHub Copilot usage metrics including code completions, chat, and PR summaries.",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&org, "org", "o", "", "GitHub organization name (required)")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "GitHub token (default: gh auth token)")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "Output as JSON")
	rootCmd.PersistentFlags().StringVar(&since, "since", "", "Start date (YYYY-MM-DD)")
	rootCmd.PersistentFlags().StringVar(&until, "until", "", "End date (YYYY-MM-DD)")
	rootCmd.PersistentFlags().StringVarP(&period, "period", "p", "", "Named period: today, yesterday, week, month, last-month, year")

	rootCmd.MarkPersistentFlagRequired("org")
}

// fetchMetrics resolves auth and fetches metrics with the current flags.
func fetchMetrics() ([]api.DayMetrics, error) {
	tok, err := auth.GetToken(token)
	if err != nil {
		return nil, err
	}

	s, u := resolveDateRange()

	client := api.NewClient(tok)
	return client.FetchMetrics(org, s, u)
}

// resolveDateRange converts --period / --since / --until flags into date strings.
func resolveDateRange() (string, string) {
	if since != "" || until != "" {
		return since, until
	}

	now := time.Now()
	today := now.Format("2006-01-02")

	switch period {
	case "today":
		return today, today
	case "yesterday":
		y := now.AddDate(0, 0, -1).Format("2006-01-02")
		return y, y
	case "week":
		return now.AddDate(0, 0, -7).Format("2006-01-02"), today
	case "month":
		firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return firstOfMonth.Format("2006-01-02"), today
	case "last-month":
		firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		lastOfPrev := firstOfThisMonth.AddDate(0, 0, -1)
		firstOfPrev := time.Date(lastOfPrev.Year(), lastOfPrev.Month(), 1, 0, 0, 0, 0, now.Location())
		return firstOfPrev.Format("2006-01-02"), lastOfPrev.Format("2006-01-02")
	case "year":
		firstOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		// API limits to 100 days, clamp if needed
		daysSinceJan1 := int(now.Sub(firstOfYear).Hours() / 24)
		if daysSinceJan1 > 100 {
			return now.AddDate(0, 0, -100).Format("2006-01-02"), today
		}
		return firstOfYear.Format("2006-01-02"), today
	default:
		// Default: last 28 days
		return now.AddDate(0, 0, -28).Format("2006-01-02"), today
	}
}

func formatNum(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(n)/1000000)
}
