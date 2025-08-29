
# 🚀 **red404**: Una red social differente

Welcome to the **red404** project! This is a collaborative effort by the members of our Discord programming community to build a toy social network. Our goal is to create a fun, functional application while learning and growing together in a supportive environment.

This project is designed to be **beginner-friendly**, focusing on practical application of modern web development technologies.

## ✨ Features (Current & Planned)

*   **User Management:** Register, Login, User Profiles.
*   **Post Creation:** Create and view text-based posts.
*   **Comments:** Add comments to posts.
*   **Follow System:** (Planned) Follow other users.
*   **Real-time Updates:** (Future consideration) Using WebSockets for live notifications/feeds.
*   **User Interface:** Clean, responsive UI with Tailwind CSS.

## 💻 Tech Stack

We're using a modern and robust stack to build our application:

| Category  | Technology       | Why we chose it                                                                 |
| :-------- | :--------------- | :------------------------------------------------------------------------------ |
| **Backend** | Go (net/http, chi) | Simple, performant, strong concurrency. `chi` for clean routing.                |
|           | Ent ORM          | Type-safe, code-generated ORM. Reduces boilerplate, simplifies DB interaction. |
|           | PostgreSQL       | Robust, reliable, widely used relational database.                              |
| **Frontend**| React            | Popular, component-based UI library.                                            |
|           | Vite             | Extremely fast dev server, quick builds, excellent developer experience.        |
|           | TypeScript       | Adds type safety to JavaScript, reducing bugs and improving code clarity.       |
|           | Tailwind CSS     | Utility-first CSS framework for rapid and consistent styling.                   |
|           | Zustand          | Lightweight, simple, and performant state management for React.                 |
| **Tools** | Docker Compose   | Simplifies environment setup, ensures consistency across dev machines.          |
|           | `air` (Go)       | Hot-reloading for Go backend for fast feedback loop during development.         |

## 📂 Project Structure

```
my-social-app/
├── backend/          # Go API server
│   ├── cmd/          # Main application entry point
│   ├── internal/     # Core application logic (handlers, services, models, config, etc.)
│   ├── ent/          # Ent ORM schemas and generated code
│   └── Dockerfile    # Docker build instructions for the Go backend
├── frontend/         # React application
│   ├── public/       # Static assets
│   ├── src/          # React components, pages, hooks, services, styles
│   └── vite.config.ts# Vite configuration (with API proxy)
├── docker-compose.yml# Orchestrates Docker services (DB, Backend)
├── Makefile          # Helper commands for common tasks (start, build, test, lint)
├── .env.example      # Example environment variables
├── .gitignore        # Files/directories to ignore in Git
├── README.md         # This file!
```

## 🚀 Getting Started

Follow these steps to get the project up and running on your local machine.

### Prerequisites

Make sure you have the following installed:

*   **Git:** For version control.
*   **Docker Desktop:** (Recommended) For running the database and backend services easily.
    *   [Download Docker Desktop](https://www.docker.com/products/docker-desktop/)
*   **Go:** Latest stable version (1.21+ recommended).
    *   [Download Go](https://go.dev/dl/)
*   **Node.js & npm/yarn/pnpm:** Latest LTS version.
    *   [Download Node.js](https://nodejs.org/en/download/)

### Setup Steps

1.  **Clone the Repository:**
    ```bash
    git clone https://github.com/your-org/my-social-app.git
    cd my-social-app
    ```

2.  **Configure Environment Variables:**
    ```bash
    cp .env.example .env
    ```
    Open the newly created `.env` file and review the variables. The defaults should work for local development.

3.  **Initialize Go Modules & Ent ORM:**
    *   First, enter the backend directory:
        ```bash
        cd backend
        ```
    *   Download Go modules:
        ```bash
        go mod tidy
        ```
    *   Generate Ent ORM code (this assumes you have Ent schemas defined in `backend/ent/schema`):
        ```bash
        go generate ./ent
        ```
    *   Go back to the root directory:
        ```bash
        cd ..
        ```

4.  **Start Services with Docker Compose:**
    We use `Makefile` commands to simplify this. The `make dev` command will start Docker services (PostgreSQL database and Go backend), and then launch the frontend and backend hot-reloaders.

    ```bash
    make dev
    ```
    *   This command will:
        *   Build and start the `db` (PostgreSQL) and `backend` services using Docker Compose.
        *   Start the React frontend development server (typically on `http://localhost:5173`).
        *   Start the Go backend hot-reloader (`air`), which will restart your Go app whenever you save changes.

5.  **Access the Application:**
    *   **Frontend:** Open your web browser and navigate to `http://localhost:5173` (or the port Vite indicates in the terminal).
    *   **Backend API:** The Go backend will be running on `http://localhost:8080`. The frontend will proxy requests to `/api` to this backend.

    You are now ready to start coding!

### Useful `make` Commands

The `Makefile` in the root directory contains convenient commands for development:

*   `make dev`: Starts all services and dev servers with hot-reload (recommended for daily dev).
*   `make up`: Starts Docker services (DB, Backend) in detached mode (`-d`).
*   `make down`: Stops and removes Docker services.
*   `make build-backend`: Builds the production Docker image for the backend.
*   `make build-frontend`: Builds the production JavaScript assets for the frontend.
*   `make test`: Runs all tests (Go backend, and eventually React frontend).
*   `make lint`: Runs linters and formatters for both backend and frontend.
*   `make clean`: Removes temporary files and Docker volumes.
*   `make help`: Displays a list of all available commands.

## 🤝 Contributing

We welcome contributions from everyone! This project is a fantastic opportunity to learn and apply new skills.

## 📄 License

This project is licensed under the MIT License - see the `LICENSE` file for details.

## 💬 Community & Support

Join our Discord server to chat with other contributors, ask questions, and share your progress!
[Your Discord Server Invite Link Here] (e.g., `https://discord.gg/S4wdMTPm`)

Happy coding!
