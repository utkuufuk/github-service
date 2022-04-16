package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/utkuufuk/github-service/internal/logger"
)

var (
	Port                int
	PersonalAccessToken string
	OrgName             string
	UserName            string
	SubscribedRepos     []string
)

func init() {
	var err error
	godotenv.Load()

	port := os.Getenv("PORT")
	Port, err = strconv.Atoi(port)
	if err != nil {
		logger.Error("PORT not set")
		os.Exit(1)
	}

	PersonalAccessToken = os.Getenv("PERSONAL_ACCESS_TOKEN")
	if PersonalAccessToken == "" {
		logger.Error("PERSONAL_ACCESS_TOKEN not set")
		os.Exit(1)
	}

	OrgName = os.Getenv("ORG_NAME")
	if OrgName == "" {
		logger.Error("ORG_NAME not set")
		os.Exit(1)
	}

	UserName = os.Getenv("USER_NAME")
	if UserName == "" {
		logger.Error("USER_NAME not set")
		os.Exit(1)
	}

	SubscribedRepos = strings.Split(os.Getenv("SUBSCRIBED_REPOS"), ",")
}
