// pulls lists pull requests on a github repo.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"text/tabwriter"

	"github.com/google/go-github/github"
)

var usage = `usage: gh COMMAND

GitHub from the command line

commands:

  pulls		list open pull requests on the repo
  issues	list open issues for the repo

This command expects the following environment variables to be set:

  GITHUB_OWNER	owner of the repo
  GITHUB_REPO	the repo on GitHub
`

func main() {
	if len(os.Args) != 1+1 || len(os.Args) >= 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Fprintf(os.Stderr, usage)
		os.Exit(1)
	}
	cmd := os.Args[1]

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	client := github.NewClient(nil)
	owner := mustGetenv("GITHUB_OWNER")
	repo := mustGetenv("GITHUB_REPO")
	ctx := context.Background()

	var err error
	switch cmd {
	case "pulls":
		err = pulls(ctx, client, w, owner, repo)
	case "issues":
		err = issues(ctx, client, w, owner, repo)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		os.Exit(1)
	}
	if err != nil {
		log.Fatal(err)
	}
	w.Flush()
}

func pulls(ctx context.Context, client *github.Client, w io.Writer, owner, repo string) error {
	prs, _, err := client.PullRequests.List(ctx, owner, repo, nil)
	if err != nil {
		return err
	}
	var pr *github.PullRequest
	for _, pr = range prs {
		if *pr.State == "open" {
			fmt.Fprintf(w, "%s\t%s\t%s\n", *pr.User.Login, *pr.HTMLURL, *pr.Title)
		}
	}
	return nil
}

func issues(ctx context.Context, client *github.Client, w io.Writer, owner, repo string) error {
	issues, _, err := client.Issues.ListByRepo(ctx, owner, repo, nil)
	if err != nil {
		return err
	}
	var i *github.Issue
	for _, i = range issues {
		if *i.State == "open" {
			whom := ""
			if i.Assignee != nil {
				whom = *i.Assignee.Login
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", whom, *i.HTMLURL, *i.Title)
		}
	}
	return nil
}

func mustGetenv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		fmt.Fprintf(os.Stderr, "%s environment variable must be set\n", name)
		os.Exit(1)
	}
	return val
}
