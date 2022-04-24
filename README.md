# github-service
A simple service to query GitHub issues & pull requests.
* See [Server Mode](#server-mode) for using it as an [entrello](https://github.com/utkuufuk/entrello) service.
* See [CLI Mode](#cli-mode) for using it as a CLI tool.

## Server Mode
Start the server:
```sh
go run ./cmd/server
```

### List of Endpoints
#### `GET <SERVER_URL>/entrello`
Fetch self-assigned issues from personal repositories.

#### `GET <SERVER_URL>/entrello/prlo`
Fetch pull requests that meet all of the following conditions:
- belongs to the configured organization
- belongs to one of the configured subscribed repositories
- neither created by nor assigned to the configured user
- not draft

#### `GET <SERVER_URL>/entrello/prlmy`
Fetch pull requests that meet all of the following conditions:
- belongs to the configured organization
- created by the configured user

#### `GET <SERVER_URL>/entrello/prlme`
Fetch pull requests that meet all of the following conditions:
- belongs to the configured organization
- assigned to the configured user
- not created by configured user

## CLI Mode
The result set of each endpoint listed above can be alternatively retrieved using the corresponding CLI command.

### List of Commands
```sh
go run ./cmd/cli
go run ./cmd/cli prlo
go run ./cmd/cli prlmy
go run ./cmd/cli prlme
```

## Configuration
Put your environment variables in a file called `.env`, based on `.env.example`.

| Environment Variable | Description |
|-|-|
| `PERSONAL_ACCESS_TOKEN` | GitHub personal access token |
| `ORG_NAME` | GitHub organization name (required for `prlo`, `prlmy` and `prlme`) |
| `USER_NAME` | GitHub user name (required for `prlo` and `prlme`) |
| `SUBSCRIBED_REPOS` | Subscribed GitHub repositories (required for `prlo`) |
| `PORT` | HTTP port (server mode only) |
| `SECRET` | API secret (server mode only, optional) |


## Running With Docker
A new [Docker image](https://github.com/utkuufuk?tab=packages&repo_name=github-service) will be created upon each [release](https://github.com/utkuufuk/github-service/releases).

1. Authenticate with the GitHub container registry (only once):
    ```sh
    echo $GITHUB_ACCESS_TOKEN | docker login ghcr.io -u GITHUB_USERNAME --password-stdin
    ```

2. Pull the latest Docker image:
    ```sh
    docker pull ghcr.io/utkuufuk/github-service/image:latest
    ```

3. Start a container:
    ```sh
    # server
    docker run -d \
        -p <PORT>:<PORT> \
        --env-file </absolute/path/to/.env> \
        --restart unless-stopped \
        --name github-service \
        ghcr.io/utkuufuk/github-service/image:latest

    # CLI
    docker run --rm \
        --env-file </absolute/path/to/.env> \
        ghcr.io/utkuufuk/github-service/image:latest \
        ./cli
    ```

