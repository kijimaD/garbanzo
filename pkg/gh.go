//go:generate mockgen -source=gh.go -destination=gh_mock.go -package=garbanzo

package garbanzo

import (
	"context"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

type clientI interface {
	getNotifications() error
}
type GitHub struct {
	Client *github.Client
}

func newGitHub() *GitHub {
	token := os.Getenv("GH_TOKEN")
	ctx := context.Background()
	sts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
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

// notificationsの情報を補足してeventに変換する
// 処理し終わったら配列から削除する
func (gh *GitHub) processNotification(s *store) error {
	for _, n := range s.notifications {
		event, err := gh.getEvent(n)
		if err != nil {
			return err
		}

		id := notificationID(*n.ID)
		s.events[id] = event
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (gh *GitHub) getEvent(n *github.Notification) (*Event, error) {
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
	event := newEvent(
		*n.ID,
		*comment.User.Login,
		*comment.User.AvatarURL,
		*n.Subject.Title,
		*comment.Body,
		*n.Subject.LatestCommentURL,
		*n.UpdatedAt,
	)

	return event, nil
}
