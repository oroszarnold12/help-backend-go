# HelP (Hybrid E-Learning Platform) backend in Go

A platform designed to support both teachers and students in educational activities.

Teachers can create courses, invite students to participate, and upload relevant resources on the subject matter.

Additionally, they can assign and grade assignments, and engage in discussions about the grading criteria. The platform also features quizzes that are automatically graded for enhanced learning assessment.

Both students and teachers receive real-time push notifications, ensuring they stay informed about important updates, such as comments on assignments or grading results. This helps keep everyone engaged and up-to-date with key events.

## Tech

- **Go** — primary language
- **MySQL** — relational database (`go-sql-driver/mysql`)
- **Gorilla Mux** — HTTP router
- **golang-migrate** — database migrations
- **JWT** — authentication (`golang-jwt/jwt`)
- **bcrypt** — password hashing (`golang.org/x/crypto`)
- **go-playground/validator** — request validation
- **godotenv** — environment variable loading
- **rs/cors** — CORS middleware

## Makefile commands

| Command | Description |
|---|---|
| `make build` | Compile the application and output the binary to `bin/help` |
| `make run` | Build and run the application |
| `make migration <name>` | Create a new SQL migration file with the given name |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Roll back the last migration |
| `make test` | Run all tests. Pass `TEST=<name>` to run a specific test by name |
