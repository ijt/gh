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

func main() {
	if len(os.Args) != 1+1 {
		fmt.Fprintf(os.Stderr, `usage: gh COMMAND

GitHub from the command line

commands:

  pulls		list open pull requests
`)
		os.Exit(1)
	}
	cmd := os.Args[1]

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	client := github.NewClient(nil)
	owner := mustGetenv("GITHUB_OWNER")
	repo := mustGetenv("GITHUB_REPO")

	switch cmd {
	case "pulls":
		pulls(client, w, owner, repo)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		os.Exit(1)
	}
	w.Flush()
}

func pulls(client *github.Client, w io.Writer, owner, repo string) {
	prs, _, err := client.PullRequests.List(context.Background(), owner, repo, nil)
	if err != nil {
		log.Fatal(err)
	}
	var pr *github.PullRequest
	for _, pr = range prs {
		if *pr.State == "open" {
			fmt.Fprintf(w, "%s\t%s\t%s\n", *pr.User.Login, *pr.HTMLURL, *pr.Title)
		}
	}
}

func mustGetenv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		fmt.Fprintf(os.Stderr, "%s environment variable must be set\n", name)
		os.Exit(1)
	}
	return val
}
