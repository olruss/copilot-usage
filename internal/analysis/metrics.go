package analysis

import (
	"sort"
	"time"

	"github.com/olruss/copilot-usage/internal/api"
)

// Summary holds aggregated metrics for a period.
type Summary struct {
	DaysCount          int
	TotalActiveUsers   int
	TotalEngagedUsers  int
	AvgActiveUsers     float64
	AvgEngagedUsers    float64
	MaxActiveUsers     int
	TotalSuggestions   int
	TotalAcceptances   int
	TotalLinesSuggested int
	TotalLinesAccepted  int
	AcceptanceRate     float64
	TotalChats         int
	TotalChatInsertions int
	TotalChatCopies    int
	TotalDotcomChats   int
	TotalPRSummaries   int
}

// LanguageStat holds aggregated metrics for a single language.
type LanguageStat struct {
	Name            string
	Suggestions     int
	Acceptances     int
	LinesSuggested  int
	LinesAccepted   int
	AcceptanceRate  float64
}

// EditorStat holds aggregated metrics for a single editor.
type EditorStat struct {
	Name            string
	EngagedUsers    int
	Suggestions     int
	Acceptances     int
	AcceptanceRate  float64
	Chats           int
}

// DaySummary holds key metrics for a single day.
type DaySummary struct {
	Date            string
	ActiveUsers     int
	EngagedUsers    int
	Suggestions     int
	Acceptances     int
	AcceptanceRate  float64
	LinesAccepted   int
	Chats           int
}

// WeekSummary holds aggregated metrics for a week.
type WeekSummary struct {
	Label           string
	StartDate       string
	EndDate         string
	Days            int
	AvgActiveUsers  float64
	TotalSuggestions int
	TotalAcceptances int
	AcceptanceRate  float64
	TotalChats      int
}

// Summarize aggregates metrics across all days.
func Summarize(days []api.DayMetrics) Summary {
	s := Summary{DaysCount: len(days)}
	if len(days) == 0 {
		return s
	}

	for _, d := range days {
		s.TotalActiveUsers += d.TotalActiveUsers
		s.TotalEngagedUsers += d.TotalEngagedUsers
		if d.TotalActiveUsers > s.MaxActiveUsers {
			s.MaxActiveUsers = d.TotalActiveUsers
		}

		addCodeCompletions(&s, d.IDECodeCompletions)
		addIDEChat(&s, d.IDEChat)
		addDotcomChat(&s, d.DotcomChat)
		addPRSummaries(&s, d.DotcomPullRequests)
	}

	s.AvgActiveUsers = float64(s.TotalActiveUsers) / float64(len(days))
	s.AvgEngagedUsers = float64(s.TotalEngagedUsers) / float64(len(days))
	if s.TotalSuggestions > 0 {
		s.AcceptanceRate = float64(s.TotalAcceptances) / float64(s.TotalSuggestions) * 100
	}

	return s
}

func addCodeCompletions(s *Summary, cc *api.IDECodeCompletions) {
	if cc == nil {
		return
	}
	for _, editor := range cc.Editors {
		for _, model := range editor.Models {
			for _, lang := range model.Languages {
				s.TotalSuggestions += lang.TotalCodeSuggestions
				s.TotalAcceptances += lang.TotalCodeAcceptances
				s.TotalLinesSuggested += lang.TotalCodeLinesSuggested
				s.TotalLinesAccepted += lang.TotalCodeLinesAccepted
			}
		}
	}
}

func addIDEChat(s *Summary, chat *api.IDEChat) {
	if chat == nil {
		return
	}
	for _, editor := range chat.Editors {
		for _, model := range editor.Models {
			s.TotalChats += model.TotalChats
			s.TotalChatInsertions += model.TotalChatInsertionEvents
			s.TotalChatCopies += model.TotalChatCopyEvents
		}
	}
}

func addDotcomChat(s *Summary, chat *api.DotcomChat) {
	if chat == nil {
		return
	}
	for _, model := range chat.Models {
		s.TotalDotcomChats += model.TotalChats
	}
}

func addPRSummaries(s *Summary, pr *api.DotcomPullRequests) {
	if pr == nil {
		return
	}
	for _, repo := range pr.Repositories {
		for _, model := range repo.Models {
			s.TotalPRSummaries += model.TotalPRSummaries
		}
	}
}

// ByLanguage aggregates metrics by programming language, sorted by acceptances descending.
func ByLanguage(days []api.DayMetrics) []LanguageStat {
	langMap := make(map[string]*LanguageStat)

	for _, d := range days {
		if d.IDECodeCompletions == nil {
			continue
		}
		for _, editor := range d.IDECodeCompletions.Editors {
			for _, model := range editor.Models {
				for _, lang := range model.Languages {
					ls, ok := langMap[lang.Name]
					if !ok {
						ls = &LanguageStat{Name: lang.Name}
						langMap[lang.Name] = ls
					}
					ls.Suggestions += lang.TotalCodeSuggestions
					ls.Acceptances += lang.TotalCodeAcceptances
					ls.LinesSuggested += lang.TotalCodeLinesSuggested
					ls.LinesAccepted += lang.TotalCodeLinesAccepted
				}
			}
		}
	}

	result := make([]LanguageStat, 0, len(langMap))
	for _, ls := range langMap {
		if ls.Suggestions > 0 {
			ls.AcceptanceRate = float64(ls.Acceptances) / float64(ls.Suggestions) * 100
		}
		result = append(result, *ls)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Acceptances > result[j].Acceptances
	})

	return result
}

// ByEditor aggregates metrics by editor, sorted by suggestions descending.
func ByEditor(days []api.DayMetrics) []EditorStat {
	editorMap := make(map[string]*EditorStat)

	for _, d := range days {
		if d.IDECodeCompletions != nil {
			for _, editor := range d.IDECodeCompletions.Editors {
				es, ok := editorMap[editor.Name]
				if !ok {
					es = &EditorStat{Name: editor.Name}
					editorMap[editor.Name] = es
				}
				for _, model := range editor.Models {
					for _, lang := range model.Languages {
						es.Suggestions += lang.TotalCodeSuggestions
						es.Acceptances += lang.TotalCodeAcceptances
					}
				}
			}
		}
		if d.IDEChat != nil {
			for _, editor := range d.IDEChat.Editors {
				es, ok := editorMap[editor.Name]
				if !ok {
					es = &EditorStat{Name: editor.Name}
					editorMap[editor.Name] = es
				}
				for _, model := range editor.Models {
					es.Chats += model.TotalChats
				}
			}
		}
	}

	// Track max engaged users per editor per day (take max across days, not sum)
	for _, d := range days {
		if d.IDECodeCompletions != nil {
			for _, editor := range d.IDECodeCompletions.Editors {
				if es, ok := editorMap[editor.Name]; ok {
					if editor.TotalEngagedUsers > es.EngagedUsers {
						es.EngagedUsers = editor.TotalEngagedUsers
					}
				}
			}
		}
	}

	result := make([]EditorStat, 0, len(editorMap))
	for _, es := range editorMap {
		if es.Suggestions > 0 {
			es.AcceptanceRate = float64(es.Acceptances) / float64(es.Suggestions) * 100
		}
		result = append(result, *es)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Suggestions > result[j].Suggestions
	})

	return result
}

// DailyBreakdown returns per-day summaries.
func DailyBreakdown(days []api.DayMetrics) []DaySummary {
	result := make([]DaySummary, 0, len(days))

	for _, d := range days {
		ds := DaySummary{
			Date:        d.Date,
			ActiveUsers: d.TotalActiveUsers,
			EngagedUsers: d.TotalEngagedUsers,
		}

		if d.IDECodeCompletions != nil {
			for _, editor := range d.IDECodeCompletions.Editors {
				for _, model := range editor.Models {
					for _, lang := range model.Languages {
						ds.Suggestions += lang.TotalCodeSuggestions
						ds.Acceptances += lang.TotalCodeAcceptances
						ds.LinesAccepted += lang.TotalCodeLinesAccepted
					}
				}
			}
		}

		if d.IDEChat != nil {
			for _, editor := range d.IDEChat.Editors {
				for _, model := range editor.Models {
					ds.Chats += model.TotalChats
				}
			}
		}

		if ds.Suggestions > 0 {
			ds.AcceptanceRate = float64(ds.Acceptances) / float64(ds.Suggestions) * 100
		}

		result = append(result, ds)
	}

	return result
}

// WeekOverWeek computes weekly summaries and deltas.
// Returns [thisWeek, lastWeek] summaries.
func WeekOverWeek(days []api.DayMetrics) (thisWeek, lastWeek WeekSummary) {
	now := time.Now()

	// Find Monday of current week
	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	thisMonday := now.AddDate(0, 0, -int(weekday-time.Monday))
	lastMonday := thisMonday.AddDate(0, 0, -7)
	lastSunday := thisMonday.AddDate(0, 0, -1)

	thisMonStr := thisMonday.Format("2006-01-02")
	lastMonStr := lastMonday.Format("2006-01-02")
	lastSunStr := lastSunday.Format("2006-01-02")
	todayStr := now.Format("2006-01-02")

	var thisDays, lastDays []api.DayMetrics
	for _, d := range days {
		if d.Date >= thisMonStr && d.Date <= todayStr {
			thisDays = append(thisDays, d)
		} else if d.Date >= lastMonStr && d.Date <= lastSunStr {
			lastDays = append(lastDays, d)
		}
	}

	thisWeek = weekSummary("This week", thisMonStr, todayStr, thisDays)
	lastWeek = weekSummary("Last week", lastMonStr, lastSunStr, lastDays)

	return thisWeek, lastWeek
}

func weekSummary(label, start, end string, days []api.DayMetrics) WeekSummary {
	s := Summarize(days)
	return WeekSummary{
		Label:            label,
		StartDate:        start,
		EndDate:          end,
		Days:             len(days),
		AvgActiveUsers:   s.AvgActiveUsers,
		TotalSuggestions: s.TotalSuggestions,
		TotalAcceptances: s.TotalAcceptances,
		AcceptanceRate:   s.AcceptanceRate,
		TotalChats:       s.TotalChats + s.TotalDotcomChats,
	}
}
