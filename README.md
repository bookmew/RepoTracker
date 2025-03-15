# RepoTracker
GitHub Repository Statistics Tracker

## Tech Stack

- Go Language
- Gin Web Framework
- PostgreSQL Database
- RESTful API Design

## Project Structure 
```
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── .env                  # Environment variables
├── src/                  # Source code directory
│   ├── main.go           # Application entry point
│   ├── controller/       # Route controllers
│   ├── handler/          # Request handlers
│   ├── model/            # Data models
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic
│   └── util/             # Utility functions
│       └── github.go     # GitHub API client
└── common/               # Common resources
    └── sql/              # SQL scripts
        └── table/        # Table definitions
```

## Features

- Fetch and store GitHub repository statistics (stars, forks, contributors)
- Periodic updates of repository statistics (hourly)
- RESTful API for accessing repository data
- Support for multiple repositories

## Installation and Running

### Prerequisites

- Go 1.16+
- PostgreSQL 12+
- GitHub API token (for higher rate limits)

### Installation Steps

1. Clone the repository
   ```bash
   git clone https://github.com/bookmew/RepoTracker.git
   cd github-repo-tracker
   ```

2. Install dependencies
   ```bash
   go mod tidy
   ```

3. Configure the database
   - Create a PostgreSQL database
   - Run the SQL scripts in `common/sql/table/` to create the necessary tables

4. Configure environment variables
   - Copy `.env.example` to `.env`
   - Update the values in `.env` with your configuration
   - Make sure to set your GitHub API token for higher rate limits

5. Run the application
   ```bash
   go run src/main.go
   ```

6. Access the API
   The server will run at `http://localhost:8080`

## Automatic Updates

The application automatically updates repository statistics every hour.
## License

[MIT](LICENSE)
