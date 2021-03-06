package hooks

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/autograde/quickfeed/ci"
	"github.com/autograde/quickfeed/database"
	logq "github.com/autograde/quickfeed/log"
	"github.com/autograde/quickfeed/scm"
	"github.com/google/go-cmp/cmp"
)

const (
	gitHubTestOrg = "autograder-test"
	secret        = "the-secret-autograder-test"
)

// To enable this test, please see instructions in the developer guide (dev.md).
// You will also need access to the autograder-test organization; you may request
// access by sending your GitHub username to hein.meling at uis.no.

// On macOS, get ngrok using `brew cask install ngrok`.
// See steps to follow [here](https://groob.io/tutorial/go-github-webhook/).

// To run this test, use the following (replace the forwarding URL with your own):
//
// NGROK_FWD=https://53c51fa9.ngrok.io go test -v -run TestGitHubWebHook
//
// This will block waiting for a push event from GitHub; meaning that you
// will manually have to create a push event to the 'tests' repository.

func TestGitHubWebHook(t *testing.T) {
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if len(accessToken) < 1 {
		t.Skip("This test requires a 'GITHUB_ACCESS_TOKEN' and access to the 'autograder-test' GitHub organization")
	}
	serverURL := os.Getenv("NGROK_FWD")
	if len(serverURL) < 1 {
		t.Skip("This test requires a 'NGROK_FWD' and access to the 'autograder-test' GitHub organization")
	}

	logger := logq.Zap(true).Sugar()
	defer func() { _ = logger.Sync() }()

	var s scm.SCM
	s, err := scm.NewSCMClient(logger, "github", accessToken)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	opt := &scm.CreateHookOptions{
		URL:        serverURL,
		Secret:     secret,
		Repository: &scm.Repository{Owner: gitHubTestOrg, Path: "tests"},
	}
	err = s.CreateHook(ctx, opt)
	if err != nil {
		t.Fatal(err)
	}

	hooks, err := s.ListHooks(ctx, opt.Repository, "")
	if err != nil {
		t.Fatal(err)
	}
	for _, hook := range hooks {
		t.Logf("hook: %v", hook)
	}

	var db database.Database
	var runner ci.Runner
	webhook := NewGitHubWebHook(logger, db, runner, secret)

	log.Println("starting webhook server")
	http.HandleFunc("/webhook", webhook.Handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func TestExtractChanges(t *testing.T) {
	modifiedFiles := []string{
		"go.mod",
		"go.sum",
		"exercise.go",
		"README.md",
		"lab2/fib.go",
		"lab3/detector/fd.go",
		"paxos/proposer.go",
		"/hallo",
		"",
	}
	want := map[string]bool{
		"lab2":  true,
		"lab3":  true,
		"paxos": true,
	}
	got := make(map[string]bool)
	extractChanges(modifiedFiles, got)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("content mismatch (-want +got):\n%s", diff)
	}
}
