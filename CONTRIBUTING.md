# Contributing to Go Incident CLI

Thank you for your interest in contributing to `goincidentcli`! This document outlines the process and conventions we follow to maintain a high-quality codebase.

## Development Setup

To contribute to `goincidentcli`, you need to set up your local environment.

### Prerequisites

- Go >= 1.21 installed on your machine.
- Git.
- `make` utility.

### Getting Started

1. Fork the repository on GitHub.
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/goincidentcli.git
   cd goincidentcli
   ```
3. Copy the environment variables template and configure your Slack token if testing integration:
   ```bash
   cp .env.example .env
   ```

## Code Style and Conventions

We follow standard Go idioms and formatting:
- **Formatting:** Run `gofmt` or `goimports` before committing. Your code must be formatted correctly.
- **Naming:** Follow recommendations in [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).
- **Documentation:** All exported types and functions must have standard `godoc` comments.
- **Language:** Ensure all text, variables, methods, and comments use English exclusively.

## Running Tests and Linting

Before pushing your changes, verify that your code passes all checks:

```bash
# Run all tests
make test

# Run linter and spell check
make lint
```

## Pull Request Process

1. **Branch Naming:** Create a branch for your feature or bug fix. Use a descriptive name like `feat/slack-retry` or `fix/report-formatting`.
2. **Commit Convention:** Use clear and descriptive commit messages. Provide context on what was changed and why.
3. **Review Process:** Once your PR is created, tests will run automatically via GitHub Actions. Maintainers will review your code.
4. **Approval:** Address any feedback. After approval, the maintainer will merge your PR.

## Where to Contribute

Contributions are welcome in various areas:
- Enhancing integration test cases.
- Adding generic export formats (JSON, PDF) via the `export` command.
- Contributing to the ongoing Phase 5 features (Terminal UI dashboards, metrics).

Please check the issue tracker for items labeled `good first issue` or `help wanted`.
