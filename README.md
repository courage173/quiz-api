# quiz-api

# Getting Started with the Quiz

This project is a CLI-based quiz application with a REST API backend in Golang, designed to allow users to answer quiz questions and view their scores and comparisons with other participants.

## Prerequisites

Ensure you have the following installed:

- [Go](https://golang.org/doc/install)
- [Air](https://github.com/cosmtrek/air) (optional, for live reloading during development)

## Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/courage173/quiz-api.git
   cd fasttrack-code-test-quiz
   ```

2. **Install Dependencies**:

   Run this command to install the required Go modules.

   ```bash
   go mod tidy
   ```

3. **Optional**: Install Air for live reloading.

   ```bash
   go install github.com/cosmtrek/air@latest
   ```

   If you install Air, make sure Go’s `bin` directory is in your `$PATH` to access the `air` command.

## Starting the Server

There are two ways to run the server:

1. **With Air (Live Reloading)**:
   Run this command to start the server with live reloading. The server will restart automatically when changes are made.

   ```bash
   air
   ```

2. **With Go (Manual Start)**:
   Alternatively, you can start the server manually with:

   ```bash
   go run ./cmd/server
   ```

The server should now be running on `http://localhost:4000`.

## Using the CLI

The project includes a CLI, built with Cobra, to interact with the API. The CLI has commands to fetch quiz questions, submit answers, and retrieve scores.

### Available CLI Commands

1. **Get Questions**:
   Retrieve all quiz questions from the server.

   ```bash
   go run main.go get
   ```

2. **Submit Answers**:
   Submit answers for a specified user. Replace `<user>` with one of the valid users (`user1`, `user2`, or `user3`).

   ```bash
   go run main.go submit --user <user>
   ```

   Example:

   ```bash
   go run main.go submit --user user1
   ```

3. **Get Score and Comparison**:
   View the user’s score and how it compares with other quiz participants.

   ```bash
   go run main.go score --user user1
   ```

### Example Workflow

1. **Start the server** using either Air or Go as described above.
2. **Run the CLI commands**:
   - Use `get` to retrieve the list of quiz questions.
   - Use `answer` with the `--user` flag to submit answers for a specific user.
   - Use `score` with the `--user` to see the score and comparison results.

This setup should get you started with running and interacting with the project. For any questions, feel free to reach out!
