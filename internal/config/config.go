package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type GitHubConfig struct {
	OrgName             string
	PersonalAccessToken string
	SubscribedRepos     []string
	UserName            string
}

type Config struct {
	// GitHub config
	GitHub GitHubConfig

	// additional server mode config
	Port   int
	Secret string
}

func ParseGitHubConfig() (cfg GitHubConfig, err error) {
	godotenv.Load()

	cfg.PersonalAccessToken = os.Getenv("PERSONAL_ACCESS_TOKEN")
	if cfg.PersonalAccessToken == "" {
		return cfg, fmt.Errorf("PERSONAL_ACCESS_TOKEN not set")
	}

	cfg.OrgName = os.Getenv("ORG_NAME")
	if cfg.OrgName == "" {
		return cfg, fmt.Errorf("ORG_NAME not set")
	}

	cfg.UserName = os.Getenv("USER_NAME")
	if cfg.UserName == "" {
		return cfg, fmt.Errorf("USER_NAME not set")
	}

	cfg.SubscribedRepos = strings.Split(os.Getenv("SUBSCRIBED_REPOS"), ",")

	return cfg, nil
}

func ParseServerConfig() (cfg Config, err error) {
	ghCfg, err := ParseGitHubConfig()
    cfg.GitHub = ghCfg

	port := os.Getenv("PORT")
	cfg.Port, err = strconv.Atoi(port)
	if err != nil {
		return cfg, fmt.Errorf("PORT not set")
	}

	cfg.Secret = os.Getenv("SECRET")
	return cfg, nil
}
