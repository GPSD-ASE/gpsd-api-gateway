@author: Kaaviya Ramkumar
@date_created: 04.04.2024

## Functions of API Gateway

### 1. **Request routing or proxy for backend services**

API gateway routes incoming requests to appropriate backend services. In the `routes.go` file, different endpoints are that forward requests to:

- gpsd-user-mgmt
- gpsd-map-mgmt
- gpsd-incident-mgmt

### 2. **Authentication & authorization**

API gateway handles authentication with JSON tokens through:

- User registration routes (`/register`, `/register-admin`)
- Sign-in and sign-out (`/signin`, `/signout`)
- Token verification (`/verify`)

### 3. Supports **CORS**

In `main.go`, API Gateway includes CORS middleware to handle cross-origin requests.

### 4. **Request logging**

`middleware.go` logs information about each request including which is useful for debugging and monitoring:

- HTTP method
- URL path
- Protocol
- Client's IP address
- User-Agent
- Response duration

## Robust DevOps pipeline

### **Continuous Integration/Continuous Deployment (CI/CD)**

Implemented a GitHub Actions workflow (`release.yml`) that provides:

1. **Automated testing** - runs Go tests with coverage reporting
2. **Code quality** - uses golangci-lint for Go code and Helm linting for charts
3. **Version management** - automatically detects when version bumps are needed based on commit types
4. **Docker image building** - builds and pushes Docker images to Docker Hub
5. **Security scanning** - Uses Trivy to scan Docker images for vulnerabilities
6. **Helm chart publishing** - packages and publishes Helm charts to GitHub Pages

### **Versioning System**

1. **Semantic versioning** - `bump-version.sh` implements semantic versioning like major.minor.patch
2. **Commit-based versioning** - determines version bump type from commit messages:
    - "BREAKING CHANGE" commits trigger major version bumps
    - "feat" commits trigger minor version bumps
    - All other commits trigger patch version bumps like “fix”, “test”, etc
3. **Automated CHANGELOG** - `update-changelog.sh` automatically generates changelog entries categorized by:
    - Added (features)
    - Fixed (bug fixes)
    - Breaking Changes

## Kubernetes deployment

1. **Health monitoring** - Implemented liveness and readiness probes
2. **Traffic management** - Set up service and ingress for routing traffic with Traefik acting as reverse proxy. Accepts only HTTPS requests from outside and forwards HTTP requests to API Gateway
3. **TLS security** - Configured the system to use HTTPS with Let's Encrypt certificates with automatic renewal
4. (In each repo) created helm file templates with deployments and services to provide easy customization and deployment of each service on the server

## Development workflow

Makefile (in each repo also) provides convenient commands for:

1. Local development commands for building, running, and testing locally
2. Building, push, and run Docker images
3. Install, upgrade, and uninstall Helm charts
4. Publish Helm charts to GitHub Pages repository
5. Automates the release process

## Vault integration

Integrated with HashiCorp Vault for:

1. **Secret management** - stores and retrieves sensitive credentials for each service
2. **Health monitoring** - checks Vault's health as part of API Gateway's health checks to see if required infrastructure is ready
3. **Authentication** - Supports multiple auth methods (token, Kubernetes)
