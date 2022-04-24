package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/utkuufuk/entrello/pkg/trello"
	"github.com/utkuufuk/github-service/internal/config"
	"github.com/utkuufuk/github-service/internal/github"
	"github.com/utkuufuk/github-service/internal/logger"
)

var cfg config.Config

func init() {
	var err error
	cfg, err = config.ParseServerConfig()
	if err != nil {
		logger.Error("Failed to parse config: %v", err)
		os.Exit(1)
	}
}

func main() {
	client := github.GetClient(cfg.GitHub)
	http.HandleFunc("/entrello", handleGetRequest(client.FetchAssignedIssues))
	http.HandleFunc("/entrello/prlo", handleGetRequest(client.FetchOtherPullRequests))
	http.HandleFunc("/entrello/prlme", handleGetRequest(client.FetchOtherPullRequestsAssignedToMe))
	http.HandleFunc("/entrello/prlmy", handleGetRequest(client.FetchMyPullRequests))
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
}

func handleGetRequest(
	fetchCards func(ctx context.Context) ([]trello.Card, error),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if cfg.Secret != "" && r.Header.Get("X-Api-Key") != cfg.Secret {
			w.WriteHeader(http.StatusUnauthorized)
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
