package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url" // Added import
	"testing"
	"time"

	ghlib "github.com/google/go-github/v58/github"
)

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux to control the server response.
// Caller must close the server.
func setup() (ghClient *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests cannot accidentally make network calls
	// so we create a bogus client and then set its transport to be the
	// test server.
	server := httptest.NewServer(mux)

	ghc := ghlib.NewClient(nil)
	parsedURL, _ := url.Parse(server.URL + "/")
	ghc.BaseURL = parsedURL
	ghc.UploadURL = parsedURL

	ghClient = &Client{
		ghClient: ghc,
		ctx:      context.Background(),
	}

	return ghClient, mux, server.URL, func() {
		server.Close()
	}
}

func TestNewClient(t *testing.T) {
	token := "test-token"
	client := NewClient(token)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.ghClient == nil {
		t.Error("NewClient.ghClient is nil")
	}
	// We can't directly inspect the token used by ghClient easily without reflection,
	// so we rely on the client creation process to be correct.
}

func TestGetPullRequest(t *testing.T) {
	ghClient, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/owner/repo/pulls/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"number": 123, 
			"title": "Test PR", 
			"base": {"ref": "main"}, 
			"head": {"ref": "feature"}
		}`)
	})

	pr, err := ghClient.GetPullRequest("owner", "repo", 123)
	if err != nil {
		t.Fatalf("GetPullRequest returned error: %v", err)
	}

	if pr.GetNumber() != 123 {
		t.Errorf("GetPullRequest returned wrong PR number: got %d, want %d", pr.GetNumber(), 123)
	}
	if pr.GetTitle() != "Test PR" {
		t.Errorf("GetPullRequest returned wrong PR title: got %s, want %s", pr.GetTitle(), "Test PR")
	}
}

func TestGetPullRequestFiles(t *testing.T) {
	ghClient, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/owner/repo/pulls/123/files", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{"filename": "file1.txt", "status": "added"},
			{"filename": "file2.go", "status": "modified"}
		]`)
	})

	files, err := ghClient.GetPullRequestFiles("owner", "repo", 123)
	if err != nil {
		t.Fatalf("GetPullRequestFiles returned error: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("GetPullRequestFiles returned %d files, want 2", len(files))
	}
	if *files[0].Filename != "file1.txt" {
		t.Errorf("GetPullRequestFiles returned wrong filename: got %s, want %s", *files[0].Filename, "file1.txt")
	}
}

func TestGetRecentMergedPRs(t *testing.T) {
	ghClient, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/owner/repo/pulls", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		// Simulate a recently merged PR and an older one
		now := time.Now()
		mergedRecent := now.Add(-1 * time.Hour)
		mergedOld := now.Add(-48 * time.Hour)

		fmt.Fprintf(w, `[
			{"number": 1, "title": "Old PR", "merged_at": "%s"},
			{"number": 2, "title": "Recent PR", "merged_at": "%s"}
		]`, mergedOld.Format(time.RFC3339), mergedRecent.Format(time.RFC3339))
	})

	prs, err := ghClient.GetRecentMergedPRs("owner", "repo", "main")
	if err != nil {
		t.Fatalf("GetRecentMergedPRs returned error: %v", err)
	}

	if len(prs) != 1 {
		t.Errorf("GetRecentMergedPRs returned %d PRs, want 1", len(prs))
	}
	if prs[0].GetNumber() != 2 {
		t.Errorf("GetRecentMergedPRs returned wrong PR: got %d, want %d", prs[0].GetNumber(), 2)
	}
}

func TestGetOpenPRs(t *testing.T) {
	ghClient, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/owner/repo/pulls", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[
			{"number": 10, "title": "Open PR 1", "state": "open"},
			{"number": 11, "title": "Open PR 2", "state": "open"}
		]`)
	})

	prs, err := ghClient.GetOpenPRs("owner", "repo", "main")
	if err != nil {
		t.Fatalf("GetOpenPRs returned error: %v", err)
	}

	if len(prs) != 2 {
		t.Errorf("GetOpenPRs returned %d PRs, want 2", len(prs))
	}
	if prs[0].GetNumber() != 10 {
		t.Errorf("GetOpenPRs returned wrong PR: got %d, want %d", prs[0].GetNumber(), 10)
	}
}

// testMethod is a helper function to assert the HTTP method.
func testMethod(t *testing.T, r *http.Request, expected string) {
	t.Helper()
	if r.Method != expected {
		t.Errorf("Request method got: %v, want %v", r.Method, expected)
	}
}
