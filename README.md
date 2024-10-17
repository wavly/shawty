# Shawty

**Shawty** is a URL shortener built using [Golang](https://go.dev),
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
   git clone https://github.com/wavly/shawty.git
   cd shawty
   ```
2. **Set ENV Variables**:
   Get the database URL and Token: [Turso Docs](https://docs.turso.tech/sdk/go/quickstart)
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
5. **Access the web interface on port**: `1234`:
   ```bash
   curl -i http://localhost:1234/
   ```

### Development

Use the `Makefile` to run/build the web server.

#### Requirements

- [Watchexec](https://github.com/watchexec/watchexec) - A file watcher for restarting and running the web server when the source files are updated.
- [Bun](https://bun.sh) - Bun package manager (or `npm`,`pmpm`) for installing and watching static content for tailwind classes.
- [Sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) - Generating type-safe code from SQL.

#### Commands

- `make server` to start the server in watch mode
- `make tailwind` to watch for tailwind classes
- `make tailmini` to minify the generated tailwind CSS file

## Contributing

We welcome any contributions to this project! For major changes, please open an issue first to discuss what you would like to change.

## LICENSE

- Shawty is [Licensed](LICENSE) under MIT
