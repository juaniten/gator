# gator: A CLI RSS Aggregator

## Overview

`gator` is a command-line interface (CLI) tool designed to aggregate and manage RSS feeds for different users. Users can log in to follow feeds and browse the feeds they are following. Built using Go and backed by a PostgreSQL database, `gator` provides a streamlined experience for collecting and accessing RSS feed data.

## Prerequisites

### PostgreSQL Installation

`gator` requires PostgreSQL as its database system. Ensure PostgreSQL is installed on your machine before running the program.

- **Installation Guide**: You can find instructions for installing PostgreSQL [here](https://www.postgresql.org/download/).

To check if PostgreSQL is installed successfully, run:

```bash
psql --version
```

### Go Installation

`gator` is developed in Go, so you need the Go toolchain installed on your system.

- **Installation Guide**: Download and install Go from the official [Go website](https://golang.org/dl/).
- **Optional Installation via Webi**: You can also install Go using Webi:
  ```bash
  curl -sS https://webi.sh/golang | sh
  ```

To verify your Go installation, run:

```bash
go version
```

## Installing gator CLI

Install `gator` using the `go install` command:

```bash
go install github.com/juaniten/gator@latest
```

This command installs the `gator` binary into your Go bin directory.

### Adding Go Bin to PATH

Ensure that your Go bin directory is included in your PATH. You can check and update your PATH with:

```bash
echo $PATH
```

If necessary, add the Go bin directory:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Verifying Installation

After installation, confirm `gator` is installed correctly by running:

```bash
gator
```

This should display an error that a command name is needed and provide usage help.

## Setting Up the Config File

`gator` requires a configuration file named `.gatorconfig.json` to run. This file should be placed in the HOME directory of the user.

### Config File Format

- **JSON Example**:
  ```json
  {
    "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
    "current_user_name": "juan"
  }
  ```

**Purpose**: The configuration file provides essential settings, like database connection details and user information, for `gator` to function.

## Running the CLI

To start the RSS aggregator, use:

```bash
gator run
```

### Available Commands

#### Database Management

- **reset**:
  ```bash
  gator reset
  ```
  Resets user information.

#### User Management

- **login**:

  ```bash
  gator login
  ```

  Logs in a user.

- **register**:

  ```bash
  gator register
  ```

  Registers a new user.

- **users**:
  ```bash
  gator users
  ```
  Lists all users.

#### Feeds

- **agg**:

  ```bash
  gator agg
  ```

  Aggregates RSS feeds.

- **addfeed**:

  ```bash
  gator addfeed
  ```

  Adds a new RSS feed (requires login).

- **feeds**:

  ```bash
  gator feeds
  ```

  Lists all available RSS feeds.

- **follow**:

  ```bash
  gator follow
  ```

  Follows a new feed (requires login).

- **unfollow**:

  ```bash
  gator unfollow
  ```

  Unfollows a feed (requires login).

- **following**:

  ```bash
  gator following
  ```

  Lists all followed feeds (requires login).

- **browse**:
  ```bash
  gator browse
  ```
  Browses RSS feeds (requires login).
