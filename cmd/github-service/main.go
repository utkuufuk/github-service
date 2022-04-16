package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/utkuufuk/entrello/pkg/trello"
	"github.com/utkuufuk/github-service/internal/config"
	"github.com/utkuufuk/github-service/internal/github"
)

func main() {
	terminalMode := flag.Bool("t", false, "Run in terminal mode")
	flag.Parse()

	client := github.GetClient()
	funRepo := map[string]func(ctx context.Context) ([]trello.Card, error){
		"":      client.FetchAssignedIssues,
		"prlo":  client.FetchOtherPullRequests(),
		"prlmy": client.FetchMyPullRequests,
		"prlme": client.FetchOtherPullRequestsAssignedToMe,
	}

	if !*terminalMode {
		http.HandleFunc("/entrello", handleGetRequest(funRepo[""]))
		http.HandleFunc("/entrello/prlo", handleGetRequest(funRepo["prlo"]))
		http.HandleFunc("/entrello/prlme", handleGetRequest(funRepo["prlme"]))
		http.HandleFunc("/entrello/prlmy", handleGetRequest(funRepo["prlmy"]))
		http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
		return
	}

	if len(os.Args) == 2 {
		displayCards(funRepo[""])
		return
	}

	switch os.Args[2] {
	case "prlo":
		displayCards(funRepo["prlo"])
	case "prlmy":
		displayCards(funRepo["prlmy"])
	case "prlme":
		displayCards(funRepo["prlme"])
	default:
		log.Fatalf("Uncrecognized command: %s", os.Args[2])
	}
}

func handleGetRequest(
	fetchCards func(ctx context.Context) ([]trello.Card, error),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		cards, err := fetchCards(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "could not fetch github cards for entrello: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cards)
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
