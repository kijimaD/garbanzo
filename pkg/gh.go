//go:generate mockgen -source=gh.go -destination=gh_mock.go -package=garbanzo

package garbanzo

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
)

var PROXY_URL string

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
	PROXY_URL = env.ProxyHost + ":" + strconv.FormatUint(uint64(env.ProxyPort), 10)
}

type clientI interface {
	getNotifications() error
}
type GitHub struct {
	Client *github.Client
}

func newGitHub() *GitHub {
	ctx := context.Background()
	sts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: env.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, sts)
	client := github.NewClient(tc)
	return &GitHub{Client: client}
}

// issueが開かれたときに対応してない。その場合は、LatestCommentURLにコメントIDではなく、issue IDが入る。
func (gh *GitHub) getNotifications(s *store) error {
	ctx := context.Background()
	ns, _, err := gh.Client.Activity.ListRepositoryNotifications(ctx, "golang", "go", nil)
	if err != nil {
		return err
	}

	for _, n := range ns {
		id := notificationID(*n.ID)
		s.notifications[id] = n
	}

	return nil
}

const ISSUES_EVENT_TYPE = "issues"
const COMMENTS_EVENT_TYPE = "comments"

// notificationsの情報を補足してeventに変換する
// 処理し終わったら配列から削除する
func (gh *GitHub) processNotification(s *store) error {
	for _, n := range s.notifications {
		u, err := url.Parse(*n.Subject.LatestCommentURL)
		if err != nil {
			return err
		}
		elements := strings.Split(u.Path, "/")
		// 最後から2番目の要素を取得する
		secondLastElement := elements[len(elements)-2]
		// comment /issues/comments/xxxxxxxx
		// issue open /issue/xxxxx

		if secondLastElement == ISSUES_EVENT_TYPE {
			event, err := gh.getIssueEvent(n)
			if err != nil {
				return err
			}
			id := notificationID(*n.ID)
			s.events[id] = event
			time.Sleep(100 * time.Millisecond)
		} else if secondLastElement == COMMENTS_EVENT_TYPE {
			event, err := gh.getCommentEvent(n)
			if err != nil {
				return err
			}
			id := notificationID(*n.ID)
			s.events[id] = event
			time.Sleep(100 * time.Millisecond)
		}
	}

	return nil
}

func (gh *GitHub) getIssueEvent(n *github.Notification) (*Event, error) {
	ctx := context.Background()

	u, err := url.Parse(*n.Subject.URL)
	if err != nil {
		return nil, err
	}
	issueID := path.Base(u.Path)
	issueIDint, _ := strconv.Atoi(issueID)
	issue, _, err := gh.Client.Issues.Get(ctx, "golang", "go", issueIDint)
	if err != nil {
		return nil, err
	}

	// ホストをプロキシサーバにする
	h, err := url.Parse(*issue.HTMLURL)
	if err != nil {
		return nil, err
	}
	htmlURL := PROXY_URL + h.Path + "#" + h.Fragment

	event := newEvent(
		*n.ID,
		*issue.User.Login,
		*issue.User.AvatarURL,
		*issue.Title,
		*issue.Body,
		htmlURL,
		*n.Subject.LatestCommentURL,
		*n.Repository.FullName,
		*n.UpdatedAt,
	)

	return event, nil
}

func (gh *GitHub) getCommentEvent(n *github.Notification) (*Event, error) {
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

	comment, _, err := gh.Client.Issues.GetComment(ctx, "golang", "go", IDint64)
	if err != nil {
		return nil, err
	}

	// ホストをプロキシサーバにする
	h, err := url.Parse(*comment.HTMLURL)
	if err != nil {
		return nil, err
	}
	htmlURL := PROXY_URL + h.Path + "#" + h.Fragment

	event := newEvent(
		*n.ID,
		*comment.User.Login,
		*comment.User.AvatarURL,
		*n.Subject.Title,
		*comment.Body,
		htmlURL,
		*n.Subject.LatestCommentURL,
		*n.Repository.FullName,
		*n.UpdatedAt,
	)

	return event, nil
}
