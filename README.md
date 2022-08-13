# Ohana

Cloud-native file storage and sharing, built with security, convenience, and
resilience in mind.

## Getting Started

To run a multi-node deployment, clone this repository and run `make prod-up`:

```bash
# Clone the repository
git clone git@github.com:OhanaFS/ohana.git

# Enter the cloned directory
cd ohana

# Start
make prod-up
```

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

### Cert Generation

If you get a
`panic: open certificates/main_GLOBAL_CERTIFICATE.pem: no such file or directory`
error, you need to generate the certificates to allow the servers to securely
communicate with each other.

1. Make a copy of certshost.example.yaml, name it certshost.yaml, and fill it
   out with all the hostnames or IP addresses of the servers in the cluster. You
   can use wildcards, but they cannot be too generic (e.g. \*.hosts is fine,
   but \* is not). See wildcards
   [here](https://en.wikipedia.org/wiki/Wildcard_certificate) for more.

2. From the ohana directory, run
   `./bin/ohana --gen-ca --gen-certs -hosts certhosts.example.yaml`

   1. `--gen-ca` will generate a new CA package that includes the following:
      1. `main_csr.json` - this is the CSR for the CA. This is required for
         making new client and node certificates
      2. `main_PRIVATE_KEY.pem` - this is the private key for the CA. This is
         required for making new client and node certificates.
      3. `main_GLOBAL_CERTIFICATE.pem` - this is the certificate for the CA.
         This is required for making new client and node certificates, **and for
         validating the certificates of the other servers.**
   2. `--gen-certs` will generate a new Cert Package that includes the
      following:
      1. `output_cert.pem` - this is the certificate for the server.
      2. `output_key.pem` - this is the private key for the server.

**Cert Generation Parameters**

- You can use `-num-of-certs` to specify the number of certificates to generate.
- After generating the CA once, you can generate more certificates by specifying
  the paths of the CA files using the following parameters and run
  `./bin/ohana --gen-certs ` with the following parameters:
  - `--csr-path` for `main_csr.json`
  - `--cert-path` for `main_GLOBAL_CERTIFICATE.pem`
  - `--pk-path` for `main_PRIVATE_KEY.pem`
- You can find more information about the parameters by running
  `./bin/ohana --help`

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
