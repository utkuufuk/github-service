# github-service
A simple service to query GitHub issues & pull requests.
* See [Server Mode](#server-mode) for using it as an [entrello](https://github.com/utkuufuk/entrello) service
* See [CLI Mode](#cli-mode) for using it as a CLI tool

## Configuration
Put your environment variables in a file called `.env`. See `.env.example` for reference.

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
