# Surf

**Surf** is a URL shortener built using [Golang](https://go.dev),
[Turso](https://turso.tech) and [HTMX](https://htmx.org). It provides a simple
web interface for shortening URLs, tracking their usage, and offering
statistics about shortened URLs.

## Features
- **Web UI with Golang's** `html/template`:
A minimalistic web interface built using Go's standard `html/template`
package and [Tailwind](https://tailwindcss.com).

- **URL Click Statistics**:
Keep track of how many times each shortened URL is accessed.

- **Caching**: Using [go-cache](https://github.com/patrickmn/go-cache) for
faster redirections and fewer database calls, by caching the results of
redirecting requests.

- **Input Validation**:
Checks if the URL is a valid URL schema. It only allows `https://` URLs. And
also checks if the URL contains a valid
[TLD](https://en.wikipedia.org/wiki/Top-level_domain).

## Getting Started

### Prerequisites

- Go (version 1.20 or higher) - [Download](https://go.dev/doc/install)
- Turso Account - [Website](https://turso.tech)

### Installation and Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/wavly/surf.git
   cd surf
   ```
2. **Set ENV Variables**:

   Get the database URL and Token from: [Turso Docs](https://docs.turso.tech/sdk/go/quickstart).
   Only needed if you're going to run the server in `prod` mode else the server
   would create a temporary `sqlite3` database in the project directory
   ```bash
   cp .env .env.local
   ```
3. **Install the dependencies**:
   ```bash
   go mod tidy
   ```
4. **Run the server**:
   ```bash
   go run .
   ```
5. **Access the web interface on port**: `1920`:
   ```bash
   xdg-open http://localhost:1920
   ```

### Development

Use the `make` command to run/build the web server.

> [!NOTE]
> Make sure the `ENVIROMENT` variable in `.env.local` is set to `dev` in order run the server in development mode.

#### Requirements

Tools you'll be needing for development:

- [Watchexec](https://github.com/watchexec/watchexec) - A file watcher for restarting and running the web server when the source files are updated.
- [Sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) - Generating type-safe code from SQL.
- [Bun](https://bun.sh) - Bun package manager (or `npm`,`pmpm`) for JavaScript dependencies.
- [Node](https://nodejs.org/en) - For live-reloading web pages whenever the source files changes.

#### Make Commands

- `make server` to start the server in watch mode
- `make tailwind` to watch for tailwind classes
- `make tailmini` to minify the generated tailwind CSS file
- `make sqlc` to generate type safe **SQL** Go-code.

## Contributing

We welcome any contributions to this project! For major changes, please open an issue first to discuss what you would like to change.

## LICENSE

- Surf is [Licensed](LICENSE) under MIT
