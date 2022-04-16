package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/utkuufuk/entrello/pkg/trello"
	"github.com/utkuufuk/github-service/internal/config"
	"github.com/utkuufuk/github-service/internal/github"
)

func main() {
	client := github.GetClient()
	http.HandleFunc("/entrello", handleGetRequest(client.FetchAssignedIssues))
	http.HandleFunc("/entrello/prlo", handleGetRequest(client.FetchOtherPullRequests))
	http.HandleFunc("/entrello/prlme", handleGetRequest(client.FetchOtherPullRequestsAssignedToMe))
	http.HandleFunc("/entrello/prlmy", handleGetRequest(client.FetchMyPullRequests))
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
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
