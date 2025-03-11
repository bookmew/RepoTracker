# RepoTracker
AI Agent Meme GitHub Repository Tracker

## Tech Stack

- Go Language
- Gin Web Framework
- MySQL Database
- RESTful API Design

## Project Structure 
├── go.mod # Go module definition
├── src/ # Source code directory
├── main.go # Application entry point
│ ├── controller/ # Route controllers
│ │ └── routes.go # API route definitions
│ ├── handler/ # Request handlers
│ ├── repository/ # Data access layer
│ └── util/ # Utility functions
└── log/ # Log files

## Installation and Running

### Prerequisites

- Go 1.16+
- MySQL 5.7+

### Installation Steps

1. Clone the repository
   ```bash
   git clone https://github.com/bookmew/GithubResource-Demo.git
   cd web3-ai-meme
   ```

2. Install dependencies
   ```bash
   go mod tidy
   ```

3. Configure the database
   - Create a MySQL database

4. Run the application
   ```bash
   go run main.go
   ```

5. Access the API
   The server will run at `http://localhost:8080`

## API Documentation

## Database Structure

## License

[MIT](LICENSE)