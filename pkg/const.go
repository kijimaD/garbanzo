package garbanzo

type GitHubEventType string

const (
	IssuesEventType       GitHubEventType = "issues"
	CommentsEventType     GitHubEventType = "comments"
	PullrequestsEventType GitHubEventType = "pulls"
	ReleasesEventType     GitHubEventType = "releases"
)

type GitHubSubjectType string

const (
	DiscussionSubjectType  GitHubSubjectType = "Discussion"
	CheckSuitSubjectType   GitHubSubjectType = "CheckSuite"
	PullRequestSubjectType GitHubSubjectType = "PullRequest"
)
