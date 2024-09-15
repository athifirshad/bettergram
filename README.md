# Bettergram

Bettergram is a photo-sharing application API built with Go, PostgreSQL, and Docker. It allows users to register, authenticate, upload photos, like photos, and comment on them. The project follows best practices for API development, including proper error handling, context management, and secure authentication tokens.

## Features

- **User Registration & Authentication**
  - Register new users with unique usernames and emails.
  - Secure authentication using Bearer tokens.
  
- **Photo Management**
  - Upload photos with captions.
  - Retrieve all photos or photos by specific users.
  - Search photos based on keywords.

- **Interactions**
  - Like and unlike photos.
  - Add and retrieve comments on photos.
  - Search photos based on username or caption.

- **API Documentation**
  - Interactive API documentation available via Swagger UI.

## Tech Stack

- **Backend:** Go, Chi, pgxpool
- **Database:** PostgreSQL
- **Containerization:** Docker, Docker Compose
- **API Documentation:** Swagger UI
- **Logging:** `slog` with `tint` for colorful output

## Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Getting Started

### Clone the Repository

```bash

git clone https://github.com/yourusername/bettergram.git
```

### Run the Application

Navigate to the project directory:

```bash
cd bettergram
```
```bash
docker-compose up --build
```

This command will:

1. Build the Go application.
2. Start the PostgreSQL database.
3. Run database migrations.
4. Launch the API server.
5. Serve Swagger UI for API documentation.

### Accessing the Application

- **API Server:** [http://localhost:4000](http://localhost:4000)
- **Swagger UI:** [http://localhost:8080](http://localhost:8080)
- **PostgreSQL:** Accessible on port `5432`.

## API Documentation

The API is documented using OpenAPI 3.0 and Swagger UI. Once the application is running, navigate to [http://localhost:8080](http://localhost:8080) to explore the API endpoints interactively.

## Project Structure
```
bettergram/
├── cmd/
│ └── api/
│ ├── main.go
│ ├── server.go
│ ├── routes.go
│ ├── handlers.go
│ ├── users.go
│ ├── photos.go
│ ├── interaction.go
│ ├── tokens.go
│ ├── middleware.go
│ ├── context.go
│ └── errors.go
├── internal/
│ └── data/
│ ├── models.go
│ ├── user.go
│ ├── photo.go
│ ├── like.go
│ ├── comment.go
│ └── tokens.go
├── migrations/
│ ├── 000001_create_users_table.up.sql
│ ├── 000001_create_users_table.down.sql
│ ├── 000002_create_photos_table.up.sql
│ ├── 000002_create_photos_table.down.sql
│ ├── 000003_create_likes_table.up.sql
│ ├── 000003_create_likes_table.down.sql
│ ├── 000004_create_comments_table.up.sql
│ └── 000004_create_comments_table.down.sql
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

The project is organized as follows:

- **cmd/api/**: Contains the main application code, controllers, and routes.
- **internal/**: Contains the models, handlers, database driver config and other internal components.
- **migrations/**: Contains the SQL migration files.
- **Dockerfile**: Defines the Docker container for the application.
- **docker-compose.yml**: Defines the Docker Compose setup for the application.

## Model

**Models** handle the data and business logic of the application. They interact with the PostgreSQL database to perform CRUD operations.

- **Location:** `internal/data/`
- **Key Models:**
  - `UserModel` (`internal/data/user.go`): Manages user data and authentication.
  - `PhotoModel` (`internal/data/photo.go`): Handles photo uploads and retrieval.
  - `LikeModel` (`internal/data/like.go`): Manages likes on photos.
  - `CommentModel` (`internal/data/comment.go`): Handles comments on photos.
  - `TokenModel` (`internal/data/tokens.go`): Manages authentication tokens.

## Controller

**Controllers** handle HTTP requests, process input, interact with models, and determine the appropriate responses.

- **Location:** `cmd/api/`
- **Key Controllers:**
  - `users.go`: Handles user registration and authentication.
  - `photos.go`: Manages photo uploads and retrieval.
  - `interaction.go`: Manages likes and comments.
  - `tokens.go`: Handles token creation and validation.
  - `routes.go`: Defines API endpoints and associates them with controllers.

**Cleaning Up**

To clean up the project, run the following command:

```bash
docker-compose down
```

This command will stop and remove the containers, networks created by the application.

If you want to remove the images and volumes, run the following command:

```bash
docker-compose down --rmi all --volumes
```

This command will remove the images from your local Docker registry.

