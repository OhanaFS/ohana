# Ohana

Cloud-native file storage and sharing, built with security, convenience, and
resilience in mind.

## Getting Started

TODO

## Development

To get started with development, you'll need:

- [Go 1.18](https://go.dev/)
- [Docker and Docker Compose](https://www.docker.com/)
- GNU Make
- [upx](https://github.com/upx/upx) (optional)
- Node

Once you have the tools installed, simply run `make`:

```bash
# Clone the repository
git clone git@github.com:OhanaFS/ohana.git

# Enter the cloned directory
cd ohana

# Build
make

# Run
CONFIG_FILE=config.example.yaml make run
```

Then try out one of the routes! Go to https://127.0.0.1:8000/v1/_health in your
browser and you should see something.

To edit the front-end, go to web folder:
```bash
cd admin-front

# Install dependencies
yarn

# Start dev server
yarn dev
```

The front-end will launch at http://localhost:3000/

### Conventions

- Code formatting

  To ensure style consistency, make sure you install
  [gofmt](https://pkg.go.dev/cmd/gofmt) to automatically format your code. Most
  editors should include this by default in their language extensions.
  [Enable formatting in VSCode](https://code.visualstudio.com/docs/languages/go#_formatting).

- Try to wrap comments to at most 80 characters wide.
- We also try to follow
  [Uber's Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).
