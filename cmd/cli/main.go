package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/utkuufuk/entrello/pkg/trello"
	"github.com/utkuufuk/github-service/internal/config"
	"github.com/utkuufuk/github-service/internal/github"
	"github.com/utkuufuk/github-service/internal/logger"
)

var gitHubConfig config.GitHubConfig

func init() {
	var err error
	gitHubConfig, err = config.ParseGitHubConfig()
	if err != nil {
		logger.Error("Failed to parse config: %v", err)
		os.Exit(1)
	}
}

func main() {
	client := github.GetClient(gitHubConfig)

	if len(os.Args) == 1 {
		displayCards(client.FetchAssignedIssues)
		return
	}

	switch os.Args[1] {
	case "prlo":
		displayCards(client.FetchOtherPullRequests)
	case "prlmy":
		displayCards(client.FetchMyPullRequests)
	case "prlme":
		displayCards(client.FetchOtherPullRequestsAssignedToMe)
	default:
		logger.Error("Uncrecognized command: %s", os.Args[2])
	}
}

func displayCards(query func(ctx context.Context) ([]trello.Card, error)) {
	cards, err := query(context.Background())
	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, c := range cards {
		fmt.Printf("%s - %s\n", c.Name, c.Desc)
	}
}
