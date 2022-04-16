package entrello

import (
	"fmt"

	"github.com/google/go-github/v39/github"
	"github.com/utkuufuk/entrello/pkg/trello"
)

func CreateCardsFromIssues(issues []*github.Issue) ([]trello.Card, error) {
	cards := make([]trello.Card, 0, len(issues))
	for _, i := range issues {
		if *i.Repository.Name == "" || *i.User.Login == "" || *i.Title == "" || *i.HTMLURL == "" {
			return nil, fmt.Errorf(
				"could not create card from issue: title, repo name, url and author are mandatory",
			)
		}

		assignee := ""
		if i.Assignee != nil && *i.User.Login != *i.Assignee.Login {
			assignee = " -> " + *i.Assignee.Login
		}

		name := fmt.Sprintf("[%s] (%s%s) %s", *i.Repository.Name, *i.User.Login, assignee, *i.Title)
		c, err := trello.NewCard(name, *i.HTMLURL, nil)
		if err != nil {
			return nil, fmt.Errorf("could not create new card: %w", err)
		}

		cards = append(cards, c)
	}
	return cards, nil
}

func CreateCardsFromPullRequests(prs []*github.PullRequest) ([]trello.Card, error) {
	cards := make([]trello.Card, 0, len(prs))
	for _, i := range prs {
		if *i.Head.Repo.Name == "" || *i.User.Login == "" || *i.Title == "" || *i.HTMLURL == "" {
			return nil, fmt.Errorf(
				"could not create card from PR: title, repo name, url and author are mandatory",
			)
		}

		assignee := ""
		if i.Assignee != nil && *i.User.Login != *i.Assignee.Login {
			assignee = " -> " + *i.Assignee.Login
		}

		name := fmt.Sprintf("[%s] (%s%s) %s", *i.Head.Repo.Name, *i.User.Login, assignee, *i.Title)
		c, err := trello.NewCard(name, *i.HTMLURL, nil)
		if err != nil {
			return nil, fmt.Errorf("could not create new card: %w", err)
		}

		cards = append(cards, c)
	}
	return cards, nil
}
