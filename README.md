# Shawty

**Shawty** is a URL shortener built using [Golang](https://go.dev),
[Turso](https://turso.tech) and [HTMX](https://htmx.org). It provides a simple
web interface for shortening URLs, tracking their usage, and offering
statistics about shortened URLs.

## Features
- **Web UI with Golang's** `html/template`:
A minimalistic web interface built using Golang's standard `html/template`
package and [Tailwind](https://tailwindcss.com).

- **URL Click Statistics**:
Keep track of how many times each shortened URL is accessed.

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
5. **Access the web interface on port: `1234`:
   ```bash
   curl -i http://localhost:1234/
   ```

### Development

Use the `Makefile` to run/build the web server.

#### Requirements

- [Watchexec](https://github.com/watchexec/watchexec) - A file watcher for restarting and running the web server when the source files is updated
- [Bun](https://bun.sh) - Bun package manager (or `npm`,`pmpm`) for installing and watching static content for tailwind classes

#### Commands

- Run `make server` to start the server in watch mode
- Run `make tailwind` to watch for tailwind classes

## Contributing

We welcome contributions to this project! For guidelines on how to contribute, please refer to the [CONTRIBUTING.md](.github/CONTRIBUTING.md) file.

## LICENSE

- Shawty is [License](LICENSE) under MPL-2.0.
