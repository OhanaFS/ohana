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
- Node.js and Yarn
- [upx](https://github.com/upx/upx)

Once you have the tools installed, simply run `make`:

```bash
# Clone the repository
git clone git@github.com:OhanaFS/ohana.git

# Enter the cloned directory
cd ohana

# Build
make

# Run backend in development mode
CONFIG_FILE=config.example.yaml make dev

# Set up and run frontend development server
cd web
yarn
yarn dev
```

Then try it out! Go to https://127.0.0.1:8000/ in your browser and you should
see the React app. The login URL is `/auth/login` and the default credentials is
`admin:password`, as defined in the `.dev/docker-compose.yaml` file (the
`USERS_CONFIGURATION_INLINE` part).

When you're done, don't forget to tear down redis and the postgres database:

```
make dev-down
```

### Conventions

- Code formatting

  To ensure style consistency, make sure you install
  [gofmt](https://pkg.go.dev/cmd/gofmt) to automatically format your code. Most
  editors should include this by default in their language extensions.
  [Enable formatting in VSCode](https://code.visualstudio.com/docs/languages/go#_formatting).

- Try to wrap comments to at most 80 characters wide.
- We also try to follow
  [Uber's Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).
