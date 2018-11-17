// pulls lists pull requests on a github repo.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/google/go-github/github"
)

func main() {
	client := github.NewClient(nil)
	owner := mustGetenv("GITHUB_OWNER")
	repo := mustGetenv("GITHUB_REPO")
	prs, _, err := client.PullRequests.List(context.Background(), owner, repo, nil)
	if err != nil {
		log.Fatal(err)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	var pr *github.PullRequest
	for _, pr = range prs {
		if *pr.State == "open" {
			fmt.Fprintf(w, "%s\t%s\t%s\n", *pr.User.Login, *pr.HTMLURL, *pr.Title)
		}
	}
	w.Flush()
}

func mustGetenv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		fmt.Fprintf(os.Stderr, "%s environment variable must be set\n", name)
		os.Exit(1)
	}
	return val
}
