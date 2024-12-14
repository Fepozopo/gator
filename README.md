# RSS Aggregator

A command-line RSS feed aggregator built in Go. This application allows users to add RSS feeds, follow feeds, fetch and display posts, and manage user accounts. It's designed to run continuously to scrape RSS feeds at regular intervals.

---

## Features

- User account management (register, login, reset)
- Add and manage RSS feeds
- Follow and unfollow feeds
- Continuous feed aggregation with scraping
- Browse posts from followed feeds

---

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)
- [Database Schema](#database-schema)

---

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/Fepozopo/gator
   cd gator
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Install `sqlc` (for SQL code generation):

   ```bash
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```

4. Install `Goose` (for database migrations):

   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

5. Set up the database:

   ```bash
   goose up
   ```

6. Generate SQL code:

   ```bash
   sqlc generate
   ```

7. Run the application:

   ```bash
   go run .
   ```

---

## Usage

The application is run through command-line commands. Users can register, add feeds, follow feeds, and browse posts.

### Example Workflow:

1. Register a user:
   ```bash
   go run . register your-username
   ```

2. Add a feed:
   ```bash
   go run . addfeed "Feed Name" https://example.com/rss
   ```

3. Follow a feed:
   ```bash
   go run . follow https://example.com/rss
   ```

4. Start the aggregator:
   ```bash
   go run . agg 1m
   ```

5. Browse posts:
   ```bash
   go run . browse 5
   ```

---

## Commands

| Command        | Description                                                                                       |
|----------------|---------------------------------------------------------------------------------------------------|
| `register`     | Register a new user.                                                                              |
| `login`        | Log in as an existing user.                                                                       |
| `reset`        | Reset the database (deletes all users, feeds, and posts).                                         |
| `addfeed`      | Add a new RSS feed to the database.                                                               |
| `feeds`        | List all RSS feeds along with their owners.                                                       |
| `follow`       | Follow an RSS feed (by URL).                                                                      |
| `unfollow`     | Unfollow an RSS feed (by URL).                                                                    |
| `agg`          | Start the aggregator service. Continuously fetch posts from all feeds.                            |
| `browse`       | Display posts from followed feeds, optionally limiting the number displayed (default: 2).         |

---

## Database Schema

### Tables

#### `users`
- `id` (UUID, primary key)
- `created_at` (timestamp)
- `updated_at` (timestamp)
- `name` (unique, string)

#### `feeds`
- `id` (UUID, primary key)
- `created_at` (timestamp)
- `updated_at` (timestamp)
- `name` (string)
- `url` (unique, string)
- `user_id` (foreign key, references `users`, `ON DELETE CASCADE`)
- `last_fetched_at` (nullable, timestamp)

#### `feed_follows`
- `id` (UUID, primary key)
- `created_at` (timestamp)
- `updated_at` (timestamp)
- `user_id` (foreign key, references `users`, `ON DELETE CASCADE`)
- `feed_id` (foreign key, references `feeds`, `ON DELETE CASCADE`)
- Unique constraint on (`user_id`, `feed_id`)

#### `posts`
- `id` (UUID, primary key)
- `created_at` (timestamp)
- `updated_at` (timestamp)
- `title` (string)
- `url` (unique, string)
- `description` (nullable, string)
- `published_at` (nullable, timestamp)
- `feed_id` (foreign key, references `feeds`, `ON DELETE CASCADE`)

