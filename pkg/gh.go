package garbanzo

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"

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

func (*GitHub) getNotifications() error {
	gh := newGitHub()
	ctx := context.Background()
	notifications, _, err := gh.Client.Activity.ListRepositoryNotifications(ctx, "golang", "go", nil)
	if err != nil {
		return err
	}

	for _, n := range notifications {
		u, err := url.Parse(*n.Subject.LatestCommentURL)
		if err != nil {
			return err
		}
		commentID := path.Base(u.Path)

		IDint64, err := strconv.ParseInt(commentID, 10, 64)
		if err != nil {
			return err
		}

		comment, _, err := gh.Client.Issues.GetComment(ctx, "golang", "go", IDint64)
		if err != nil {
			return err
		}
		fmt.Println(*comment.Body)
	}
	return nil
}
