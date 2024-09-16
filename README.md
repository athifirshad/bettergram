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

## API Routes

### User Management
- `POST /users`: Register a new user
  - Request body: `{ "username": string, "email": string, "password": string }`
  - Response: User object

- `POST /users/login`: Login user
  - Request body: `{ "email": string, "password": string }`
  - Response: Authentication token

- `GET /users/profile`: Get user profile (requires authentication)
  - Response: User object

### Photo Management
- `POST /photos`: Upload a new photo (requires authentication)
  - Request body: Multipart form data with "photo" file and "caption" text
  - Response: Photo object

- `GET /photos`: Get all photos
  - Response: Array of Photo objects

- `GET /photos/{id}`: Get a specific photo
  - Response: Photo object

- `GET /users/photos`: Get photos of the authenticated user (requires authentication)
  - Response: Array of Photo objects

- `GET /photos/search`: Search photos
  - Query parameter: `q` (search query)
  - Response: Array of Photo objects

### Interactions
- `POST /photos/{id}/like`: Like a photo (requires authentication)
  - Response: Like object

- `DELETE /photos/{id}/like`: Unlike a photo (requires authentication)
  - Response: No content

- `POST /photos/{id}/comments`: Add a comment to a photo (requires authentication)
  - Request body: `{ "content": string }`
  - Response: Comment object

- `GET /photos/{id}/comments`: Get comments for a photo
  - Response: Array of Comment objects

### Authentication
- `POST /tokens`: Create authentication token
  - Request body: `{ "email": string, "password": string }`
  - Response: Authentication token

### Miscellaneous
- `GET /status`: Get API status
  - Response: Status object


## Project Structure

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

