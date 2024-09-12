# Shawty

**Shawty** is a URL shortener built using [Golang](https://go.dev) and
[Turso](https://turso.tech). It provides a simple interface for shortening
URLs, tracking their usage, and offering statistics about shortened URLs.

## Features
- **URL Click Statistics**:
Keep track of how many times each shortened URL is accessed.

- **Input Validation**:
Checks if the URL is a valid URL schema. It only allows `http://` and
`https://` URLs. And also checking if the URL contains a
[TLD](https://en.wikipedia.org/wiki/Top-level_domain).

## Getting Started

### Prerequisites

- Go (version 1.20 or higher) - [Download](https://go.dev/doc/install)
- Turso Account - [Website](https://turso.tech)

### Installation and Setup

1. **Clone the repository**:
   ```
   git clone https://github.com/wavly/shawty.git
   cd shawty
   ```
2. **Set ENV Variables**:
   Get the database URL and Token: [Turso Docs](https://docs.turso.tech/sdk/go/quickstart)

   ```
   cp .env .env.local
   ```
3. **Install the dependencies**:
   ```
   go mod tidy
   ```
4. **Run the server**:
   ```
   go run .
   ```

## Contributing

We welcome contributions to this project! For guidelines on how to contribute, please refer to the [CONTRIBUTING.md](.github/CONTRIBUTING.md) file.

## LICENSE

- Shawty is [License](LICENSE) under MPL-2.0.
