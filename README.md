# Chirpy

### Motivation

This is a project focused on building an HTTP web server capable of handling RESTful requests across multiple endpoints.

#### Goal

Chirpy is built for learning purposes. The goal of Chirpy is to simulate a web server, that can support messages (chirps) between users. It is a RESTful API, uses PostgreSQL for its databases, has authentication functionality via JWT (access token) and refresh tokens.

### Getting Started

#### Prerequisites
- Go 1.22+
- PostgreSQL

#### Setup
1. Clone the repository:
   ```bash
   git clone https://github.com/CrstKng/Chirpy.git
   cd Chirpy
   ```
2. Configure environment variables (create a `.env` file):
   ```env
   DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
   JWT_SECRET=your_jwt_secret
   POLKA_KEY=your_polka_key
   ```
3. Start the server:
   ```bash
   go run main.go
   ```

Note: you need to use your own username and password for postgress to run the migrations yourself

### Key Features & Endpoints
- `POST /api/users` - Create a user
- `POST /api/login` - Authenticate and obtain JWT
- `POST /api/chirps` - Create a chirp (authenticated)
- `GET /api/chirps` - List all chirps (optional parameters: author_id: get only chirps registered by a particular user, sort: sort shown chirps by their creation time)
- `GET /api/chirps/{chirpID}` - Retrieve a chirp specified by ID
- `POST /api/refresh` - Expects a refresh token in the headers in client request and retrieves an access token
- `POST /api/revoke` - Revokes an active refresh token
- `PUT /api/users` - Updates user's email and password
- `DELETE /api/chirps/{chirpID}` - Deletes chirp specified by id
- `POST /api/polka/webhooks` - Handles event based webhooks from 3rd party servers and has API Key authentication (POLKA_KEY in `.env`)

### Additional

If you want to build new queries I recommend using `sqlc` to safely convert the queries to Go code. You can install it by using `go install`:
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Also, if you want to run additional migration, I recommend using `goose`. You can install it by using `go install`:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

To run those migrations run from `sql/schema` directory:
For Linux or WSL on Windows:
```bash
goose postgres "postgres://postgres_username:postgres_password@localhost:5432/chirpy" up
```
For macOS:
```bash
goose postgres "postgres://your_username:@localhost:5432/chirpy" up
```

