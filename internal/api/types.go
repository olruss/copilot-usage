package api

// DayMetrics represents a single day's Copilot usage metrics.
type DayMetrics struct {
	Date               string              `json:"date"`
	TotalActiveUsers   int                 `json:"total_active_users"`
	TotalEngagedUsers  int                 `json:"total_engaged_users"`
	IDECodeCompletions *IDECodeCompletions `json:"copilot_ide_code_completions"`
	IDEChat            *IDEChat            `json:"copilot_ide_chat"`
	DotcomChat         *DotcomChat         `json:"copilot_dotcom_chat"`
	DotcomPullRequests *DotcomPullRequests `json:"copilot_dotcom_pull_requests"`
}

// IDECodeCompletions contains code completion metrics broken down by editor.
type IDECodeCompletions struct {
	TotalEngagedUsers int                        `json:"total_engaged_users"`
	Editors           []IDECodeCompletionsEditor  `json:"editors"`
}

// IDECodeCompletionsEditor contains code completion metrics for a single editor.
type IDECodeCompletionsEditor struct {
	Name              string                      `json:"name"`
	TotalEngagedUsers int                         `json:"total_engaged_users"`
	Models            []IDECodeCompletionsModel   `json:"models"`
}

// IDECodeCompletionsModel contains code completion metrics for a single model.
type IDECodeCompletionsModel struct {
	Name                   string                       `json:"name"`
	IsCustomModel          bool                         `json:"is_custom_model"`
	TotalEngagedUsers      int                          `json:"total_engaged_users"`
	Languages              []IDECodeCompletionsLanguage `json:"languages"`
}

// IDECodeCompletionsLanguage contains code completion metrics for a single language.
type IDECodeCompletionsLanguage struct {
	Name              string `json:"name"`
	TotalEngagedUsers int    `json:"total_engaged_users"`
	TotalCodeSuggestions  int `json:"total_code_suggestions"`
	TotalCodeAcceptances  int `json:"total_code_acceptances"`
	TotalCodeLinesSuggested int `json:"total_code_lines_suggested"`
	TotalCodeLinesAccepted  int `json:"total_code_lines_accepted"`
}

// IDEChat contains IDE chat metrics broken down by editor.
type IDEChat struct {
	TotalEngagedUsers int              `json:"total_engaged_users"`
	Editors           []IDEChatEditor  `json:"editors"`
}

// IDEChatEditor contains chat metrics for a single editor.
type IDEChatEditor struct {
	Name              string          `json:"name"`
	TotalEngagedUsers int             `json:"total_engaged_users"`
	Models            []IDEChatModel  `json:"models"`
}

// IDEChatModel contains chat metrics for a single model.
type IDEChatModel struct {
	Name              string `json:"name"`
	IsCustomModel     bool   `json:"is_custom_model"`
	TotalEngagedUsers int    `json:"total_engaged_users"`
	TotalChats        int    `json:"total_chats"`
	TotalChatInsertionEvents int `json:"total_chat_insertion_events"`
	TotalChatCopyEvents      int `json:"total_chat_copy_events"`
}

// DotcomChat contains Copilot Chat on github.com metrics.
type DotcomChat struct {
	TotalEngagedUsers int               `json:"total_engaged_users"`
	Models            []DotcomChatModel `json:"models"`
}

// DotcomChatModel contains dotcom chat metrics for a single model.
type DotcomChatModel struct {
	Name              string `json:"name"`
	IsCustomModel     bool   `json:"is_custom_model"`
	TotalEngagedUsers int    `json:"total_engaged_users"`
	TotalChats        int    `json:"total_chats"`
}

// DotcomPullRequests contains PR summary metrics.
type DotcomPullRequests struct {
	TotalEngagedUsers int                         `json:"total_engaged_users"`
	Repositories      []DotcomPullRequestsRepo    `json:"repositories"`
}

// DotcomPullRequestsRepo contains PR metrics for a single repository.
type DotcomPullRequestsRepo struct {
	Name              string                      `json:"name"`
	TotalEngagedUsers int                         `json:"total_engaged_users"`
	Models            []DotcomPullRequestsModel   `json:"models"`
}

// DotcomPullRequestsModel contains PR metrics for a single model.
type DotcomPullRequestsModel struct {
	Name                   string `json:"name"`
	IsCustomModel          bool   `json:"is_custom_model"`
	TotalEngagedUsers      int    `json:"total_engaged_users"`
	TotalPRSummaries       int    `json:"total_pr_summaries"`
}
