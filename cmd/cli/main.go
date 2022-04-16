package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/utkuufuk/entrello/pkg/trello"
	"github.com/utkuufuk/github-service/internal/github"
)

func main() {
	client := github.GetClient()

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
		log.Fatalf("Uncrecognized command: %s", os.Args[2])
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
