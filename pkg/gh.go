//go:generate mockgen -source=gh.go -destination=gh_mock.go -package=garbanzo

package garbanzo

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/google/go-github/v48/github"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
)

const timezone = "Asia/Tokyo"
const timeformat = "2006-01-02 15:04"

var PROXY_BASE string

type Env struct {
	ProxyHost   string `envconfig:"PROXY_BASE" default:"http://localhost"`
	ProxyPort   uint16 `envconfig:"PROXY_PORT" default:"8081"`
	GitHubToken string `envconfig:"GH_TOKEN" required:"true"`
}

var env Env

func init() {
	err := envconfig.Process("", &env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't parse environment variables: %s\n", err.Error())
		os.Exit(1)
	}
	PROXY_BASE = env.ProxyHost + ":" + strconv.FormatUint(uint64(env.ProxyPort), 10)
}

type clientI interface {
	getNotifications() error
}
type GitHub struct {
	Client        *github.Client
	notifications map[string]*github.Notification
}

func newGitHub() *GitHub {
	ctx := context.Background()
	sts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: env.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, sts)
	client := github.NewClient(tc)
	return &GitHub{
		Client:        client,
		notifications: make(map[string]*github.Notification),
	}
}

// issueが開かれたときに対応してない。その場合は、LatestCommentURLにコメントIDではなく、issue IDが入る。
func (gh *GitHub) getNotifications() error {
	ctx := context.Background()
	ns, _, err := gh.Client.Activity.ListNotifications(ctx, nil)
	if err != nil {
		return err
	}

	for _, n := range ns {
		gh.notifications[*n.ID] = n
	}

	return nil
}

const ISSUES_EVENT_TYPE = "issues"
const COMMENTS_EVENT_TYPE = "comments"
const PULLREQUESTS_EVENT_TYPE = "pulls"
const RELEASES_EVENT_TYPE = "releases"

// notificationsの情報を補足してeventに変換する
// 処理し終わったら配列から削除する
func (gh *GitHub) processNotification(r *room) error {
	// notificationsを日付順にソートしてからループを実行する
	keys := make([]string, 0, len(r.events))
	for key := range gh.notifications {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		// 過去 → 未来の順
		return gh.notifications[keys[i]].UpdatedAt.Before(*gh.notifications[keys[j]].UpdatedAt)
	})
	sorted := make([]*github.Notification, 0, len(gh.notifications))
	for _, k := range keys {
		sorted = append(sorted, gh.notifications[k])
	}

	for _, n := range sorted {
		if _, exists := r.events[*n.ID]; exists {
			continue
		}
		if *n.Subject.Type == "Discussion" {
			// discussionはなぜかURLが空になっていて、辿ることができない
			// https://github.com/orgs/community/discussions/15252
			continue
		}
		if *n.Subject.Type == "CheckSuite" {
			// コミットへの通知? URLを持っていない
			continue
		}

		var originURL string
		if *n.Subject.Type == "PullRequest" && n.Subject.LatestCommentURL == nil {
			// PRオープンやクローズ、レビュー送信の場合はLatestCommentURLがない
			originURL = *n.Subject.URL
		} else {
			originURL = *n.Subject.LatestCommentURL
		}

		u, err := url.Parse(originURL)
		if err != nil {
			return err
		}
		elements := strings.Split(u.Path, "/")
		secondLastElement := elements[len(elements)-2]
		thirdLastElement := elements[len(elements)-3]
		// pull:              /pulls/xxxxx
		// issue:             /issue/xxxxx
		// comment(Issue+PR): /issues/comments/xxxxx
		// commit comment:    /comments/xxxxx
		// release            /releases/xxxxx

		if secondLastElement == PULLREQUESTS_EVENT_TYPE {
			event, err := gh.getPullRequestEvent(n)
			if err != nil {
				return err
			}
			r.fetch <- event
		} else if secondLastElement == ISSUES_EVENT_TYPE {
			// issue open
			event, err := gh.getIssueEvent(n)
			if err != nil {
				return err
			}
			r.fetch <- event
		} else if thirdLastElement == ISSUES_EVENT_TYPE &&
			secondLastElement == COMMENTS_EVENT_TYPE {
			// comment
			event, err := gh.getIssueCommentEvent(n)
			if err != nil {
				return err
			}
			r.fetch <- event
		} else if secondLastElement == COMMENTS_EVENT_TYPE {
			// commit comment
		} else if secondLastElement == RELEASES_EVENT_TYPE {
			event, err := gh.getReleaseEvent(n)
			if err != nil {
				return err
			}
			r.fetch <- event
		} else {
			log.Println("URLパースを通過した", *n.Subject.LatestCommentURL, *n.Subject.Title)
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

func (gh *GitHub) getPullRequestEvent(n *github.Notification) (*Event, error) {
	ctx := context.Background()

	u, err := url.Parse(*n.Subject.URL)
	if err != nil {
		return nil, err
	}
	pullID := path.Base(u.Path)
	pullIDint, _ := strconv.Atoi(pullID)
	pull, _, err := gh.Client.PullRequests.Get(ctx, *n.Repository.Owner.Login, *n.Repository.Name, pullIDint)
	if err != nil {
		return nil, err
	}

	// ホストをプロキシサーバにする
	proxyURL, err := genProxyURL(*pull.HTMLURL)
	if err != nil {
		return nil, err
	}

	htmlBody := mdToHTML([]byte(*pull.Body))
	htmlTitle := mdToHTML([]byte(*pull.Title))

	event := newEvent(
		*n.ID,
		*pull.User.Login,
		*pull.User.AvatarURL,
		*pull.Title,
		string(htmlTitle),
		string(htmlBody),
		*pull.HTMLURL,
		proxyURL,
		*n.Repository.FullName,
		genTimeWithTZ(n.UpdatedAt),
		"PR",
		*n.UpdatedAt,
	)

	return event, nil
}

func (gh *GitHub) getIssueEvent(n *github.Notification) (*Event, error) {
	ctx := context.Background()

	u, err := url.Parse(*n.Subject.URL)
	if err != nil {
		return nil, err
	}
	issueID := path.Base(u.Path)
	issueIDint, _ := strconv.Atoi(issueID)
	issue, _, err := gh.Client.Issues.Get(ctx, *n.Repository.Owner.Login, *n.Repository.Name, issueIDint)
	if err != nil {
		return nil, err
	}

	// ホストをプロキシサーバにする
	proxyURL, err := genProxyURL(*issue.HTMLURL)
	if err != nil {
		return nil, err
	}

	htmlBody := mdToHTML([]byte(*issue.Body))
	htmlTitle := mdToHTML([]byte(*issue.Title))

	event := newEvent(
		*n.ID,
		*issue.User.Login,
		*issue.User.AvatarURL,
		*issue.Title,
		string(htmlTitle),
		string(htmlBody),
		*issue.HTMLURL,
		proxyURL,
		*n.Repository.FullName,
		genTimeWithTZ(n.UpdatedAt),
		"Issue",
		*n.UpdatedAt,
	)

	return event, nil
}

func (gh *GitHub) getIssueCommentEvent(n *github.Notification) (*Event, error) {
	ctx := context.Background()

	u, err := url.Parse(*n.Subject.LatestCommentURL)
	if err != nil {
		return nil, err
	}
	commentID := path.Base(u.Path)

	IDint64, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		return nil, err
	}
	comment, _, err := gh.Client.Issues.GetComment(ctx, *n.Repository.Owner.Login, *n.Repository.Name, IDint64)
	if err != nil {
		return nil, err
	}

	// ホストをプロキシサーバにする
	proxyURL, err := genProxyURL(*comment.HTMLURL)
	if err != nil {
		return nil, err
	}

	htmlBody := mdToHTML([]byte(*comment.Body))
	htmlTitle := mdToHTML([]byte(*n.Subject.Title))

	event := newEvent(
		*n.ID,
		*comment.User.Login,
		*comment.User.AvatarURL,
		*n.Subject.Title,
		string(htmlTitle),
		string(htmlBody),
		*comment.HTMLURL,
		proxyURL,
		*n.Repository.FullName,
		genTimeWithTZ(n.UpdatedAt),
		"Comment",
		*n.UpdatedAt,
	)

	return event, nil
}

func (gh *GitHub) getReleaseEvent(n *github.Notification) (*Event, error) {
	ctx := context.Background()

	u, err := url.Parse(*n.Subject.LatestCommentURL)
	if err != nil {
		return nil, err
	}
	releaseID := path.Base(u.Path)

	IDint64, err := strconv.ParseInt(releaseID, 10, 64)
	if err != nil {
		return nil, err
	}
	release, _, err := gh.Client.Repositories.GetRelease(ctx, *n.Repository.Owner.Login, *n.Repository.Name, IDint64)
	if err != nil {
		return nil, err
	}

	// ホストをプロキシサーバにする
	proxyURL, err := genProxyURL(*release.HTMLURL)
	if err != nil {
		return nil, err
	}

	htmlBody := mdToHTML([]byte(*release.Body))
	htmlTitle := mdToHTML([]byte(*n.Subject.Title))

	event := newEvent(
		*n.ID,
		*n.Repository.Owner.Login,     // リリースにはユーザがないのでとりあえずownerを設定する
		*n.Repository.Owner.AvatarURL, // リリースにはユーザがないのでとりあえずownerを設定する
		*n.Subject.Title,
		string(htmlTitle),
		string(htmlBody),
		*release.HTMLURL,
		proxyURL,
		*n.Repository.FullName,
		genTimeWithTZ(n.UpdatedAt),
		"Release",
		*n.UpdatedAt,
	)

	return event, nil
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

// URLのホストをプロキシサーバにする
func genProxyURL(u string) (string, error) {
	h, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	proxyURL := PROXY_BASE + h.Path + "#" + h.Fragment
	return proxyURL, nil
}

func genTimeWithTZ(t *time.Time) string {
	jst := time.FixedZone(timezone, 9*60*60)
	nowJST := t.In(jst)
	updatedAt := nowJST.Format(timeformat)
	return updatedAt
}
