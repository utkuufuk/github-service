package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v39/github"
	"github.com/utkuufuk/entrello/pkg/trello"
	"github.com/utkuufuk/github-service/internal/config"
	"github.com/utkuufuk/github-service/internal/entrello"
	"golang.org/x/oauth2"
)

type Client struct {
	client          *github.Client
	subscribedRepos []string
	orgName         string
	userName        string
}

func GetClient(cfg config.GitHubConfig) Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.PersonalAccessToken})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)
	return Client{
		client,
		cfg.SubscribedRepos,
		cfg.OrgName,
		cfg.UserName,
	}
}

func (c Client) FetchAssignedIssues(ctx context.Context) ([]trello.Card, error) {
	issues, _, err := c.client.Issues.List(ctx, false, nil)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve self-assigned issues: %w", err)
	}

	nonPullRequestIssues := make([]*github.Issue, 0)
	for _, i := range issues {
		if !i.IsPullRequest() {
			nonPullRequestIssues = append(nonPullRequestIssues, i)
		}
	}
	return entrello.CreateCardsFromIssues(nonPullRequestIssues)
}

func (c Client) FetchOtherPullRequests(ctx context.Context) ([]trello.Card, error) {
	pullRequests := make([]*github.PullRequest, 0)
	for _, repo := range c.subscribedRepos {
		prs, _, err := c.client.PullRequests.List(ctx, c.orgName, repo, nil)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve pull requests from %s/%s: %w", c.orgName, repo, err)
		}
		pullRequests = append(pullRequests, prs...)
	}

	otherPullRequests := make([]*github.PullRequest, 0)
	for _, i := range pullRequests {
		if !*i.Draft && *i.User.Login != c.userName && (i.Assignee == nil || *i.Assignee.Login != c.userName) {
			otherPullRequests = append(otherPullRequests, i)
		}
	}
	return entrello.CreateCardsFromPullRequests(otherPullRequests)
}

func (c Client) FetchOtherPullRequestsAssignedToMe(ctx context.Context) ([]trello.Card, error) {
	assignedIssues, _, err := c.client.Issues.ListByOrg(ctx, c.orgName, &github.IssueListOptions{
		Filter: "assigned",
	})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve assigned %s issues: %w", c.orgName, err)
	}

	assignedPullRequests := make([]*github.Issue, 0)
	for _, i := range assignedIssues {
		if i.IsPullRequest() && *i.User.Login != c.userName {
			assignedPullRequests = append(assignedPullRequests, i)
		}
	}
	return entrello.CreateCardsFromIssues(assignedPullRequests)
}

func (c Client) FetchMyPullRequests(ctx context.Context) ([]trello.Card, error) {
	createdIssues, _, err := c.client.Issues.ListByOrg(ctx, c.orgName, &github.IssueListOptions{
		Filter: "created",
	})
	if err != nil {
		return nil, fmt.Errorf("could not retrieve created Portchain issues: %w", err)
	}

	createdPullRequests := make([]*github.Issue, 0)
	for _, i := range createdIssues {
		if i.IsPullRequest() {
			createdPullRequests = append(createdPullRequests, i)
		}
	}
	return entrello.CreateCardsFromIssues(createdPullRequests)
}
