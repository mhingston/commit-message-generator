# Commit Message Generator

This small Go application helps you generate concise and meaningful git commit messages using Azure OpenAI hosted models.

## Features

- Generates commit messages based on the staged changes (git diff).
- Prepends the formatted branch name (if it starts with "feature/") to the commit message.
- Allows you to use an existing commit message if it's not empty.

## Prerequisites

- Go installed and configured.
- Git installed and configured.
- `azoai` Go package installed: `go get github.com/mhingston/azoai`
- Git hook for `prepare-commit-msg` ([.git/hooks/prepare-commit-msg](https://git-scm.com/docs/githooks#_prepare_commit_msg))
- Environment variables set:
  - `AZURE_OPENAI_API_KEY`: Your Azure OpenAI API key.
  - `AZURE_OPENAI_ENDPOINT`: Your Azure OpenAI endpoint.
  - `AZURE_OPENAI_API_VERSION`: The API version for the Azure OpenAI service.
  - `AZURE_OPENAI_API_DEPLOYMENT`: The deployment name for the model.

## Installation

1. Clone the repository: `git clone https://github.com/mhingston/commit-message-generator.git`
2. Build the application: `go build main.go`

## Usage

1. Stage your changes using `git add`.
2. The git hook will trigger the application e.g. `commit-message-generator .git/COMMIT_EDITMSG`
3. The application will generate a commit message based on the staged changes.