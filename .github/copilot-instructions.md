# MyMovies - HTMX Go Web Application

MyMovies is a Go web application that displays movie information using The Movie Database (TMDB) API. It features an HTMX frontend for dynamic content loading, responsive CSS styling with dark/light theme support, and embedded static assets for production deployment.

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### Dependencies and Prerequisites
- Go 1.25 is required (specified in go.mod)
- No additional dependencies needed - uses only Go standard library
- Internet access required for TMDB API calls (api.themoviedb.org)
- Application uses embedded TMDB API token (hardcoded in tmdb/tmdb.go)

### Building the Application
- Bootstrap and build the repository:
  - `go mod tidy` - Download dependencies (takes <1 second)
  - `go build .` - Build production version (takes <1 second). NEVER CANCEL: Set timeout to 30+ seconds for safety.
  - `go build -tags dev .` - Build development version (takes <1 second). NEVER CANCEL: Set timeout to 30+ seconds for safety.
- The build produces an executable named `htmxBackend`
- Production build embeds static assets using `//go:embed public`
- Development build serves files from disk for live reload during development

### Running the Application
- Production mode: `./htmxBackend` (uses embedded assets)
- Development mode: Build with `-tags dev` then run `./htmxBackend` (serves from public/ directory)
- Application starts on port 8080
- Logs will show either "using prod static handler" or "using templates" indicating the build mode

### Testing and Validation
- `go test ./...` - Run tests (takes <1 second, currently no test files exist). NEVER CANCEL: Set timeout to 30+ seconds for safety.
- `go fmt ./...` - Format code
- `go vet ./...` - Static analysis
- Test API endpoints (production build only):
  - `curl http://localhost:8080/api` should return "Hello World!" (fails in dev mode due to routing bug)
  - `curl http://localhost:8080/api/movies/list` works in both modes but requires TMDB API access
  - `curl http://localhost:8080/style.css` should return CSS content
  - `curl http://localhost:8080/` requires TMDB API access and may fail in restricted environments

### Expected Timing and Cancellation Warnings
- Build time: <1 second. NEVER CANCEL: Set timeout to 30+ seconds minimum.
- Test time: <1 second (no tests exist). NEVER CANCEL: Set timeout to 30+ seconds minimum.
- Application startup: immediate
- **CRITICAL**: NEVER CANCEL any build or test commands even if they appear fast - always use generous timeouts.

## Validation Scenarios

### Always test these scenarios after making changes:
1. **Build Validation**: Both production and development builds must succeed
   - `go build .` must complete without errors
   - `go build -tags dev .` must complete without errors
2. **Static File Serving**: Verify static assets are served correctly
   - Start the application: `./htmxBackend`
   - Test: `curl http://localhost:8080/style.css` should return CSS content
   - Test: `curl http://localhost:8080/index.js` should return JavaScript content
3. **API Endpoints**: Test basic API functionality
   - **Production build only**: `curl http://localhost:8080/api` should return "Hello World!"
   - Both modes: `curl http://localhost:8080/api/movies/list` (may fail without TMDB access)
   - Note: Development mode has a routing bug where `/api` endpoint doesn't work (only longer paths like `/api/movies/list`)
4. **Main Application**: Verify the root page functionality
   - Test: `curl -I http://localhost:8080/` (may return 500 without TMDB access - this is expected)
   - In environments with TMDB access, the page should load with movie data

### Manual Validation Requirements
- ALWAYS start the application and verify it listens on port 8080
- ALWAYS test static file serving (CSS, JS) to ensure assets are properly embedded/served
- Test both development and production builds if making changes to static files
- Verify build mode in logs: look for "using prod static handler" vs "using templates"

## Common Issues and Limitations

### TMDB API Access
- The application requires internet access to api.themoviedb.org
- In restricted environments, the main page will return HTTP 500 with error: "dial tcp: lookup api.themoviedb.org"
- Static files and API endpoints will still work
- TMDB token is embedded in tmdb/tmdb.go and may expire

### Build Modes
- Production build (`go build .`): Uses `//go:embed public` to embed static files
- Development build (`go build -tags dev .`): Serves files from public/ directory for live reload
- Always verify which build mode you're using by checking the log output when starting the server
- **Development mode limitation**: The `/api` endpoint has a routing bug and returns 404. Use longer paths like `/api/movies/list` for testing.
- Always verify which build mode you're using by checking the log output when starting the server

## Key Projects and File Structure

### Core Application Files
- `main.go` - Main application entry point with HTTP server setup
- `mainHandler.go` - Production mode file handler (uses embedded assets)
- `mainHandler_dev.go` - Development mode file handler (serves from disk)
- `go.mod` - Go module definition (htmxBackend module, Go 1.25)

### Frontend Assets (public/ directory)
- `public/index.gohtml` - Main HTML template with HTMX integration
- `public/style.css` - Application styling with dark/light theme support
- `public/reset.css` - CSS reset
- `public/index.js` - Dark mode toggle functionality

### Backend Modules
- `tmdb/tmdb.go` - The Movie Database API client with authentication
- `tmdb/types.go` - Data structures for movie API responses
- `templates/types.go` - Template data structures

### CI/CD
- `.github/workflows/release.yml` - GitHub Actions workflow for building and releasing binaries
- Builds for Linux and Windows platforms
- Generates checksums for releases

### API Endpoints
- `/` - Main movie listing page (requires TMDB access)
- `/api` - Returns "Hello World!"
- `/api/movies/list` - JSON API for movie listings
- `/hx/{action}` - HTMX endpoints for dynamic content
- Static files served from `/` path

## Additional Notes
- No external JavaScript frameworks - uses vanilla JS and HTMX
- CSS uses modern features like CSS custom properties and light-dark() function
- Application designed for both desktop and mobile with responsive design
- Error handling for TMDB API failures is built-in but may cause 500 errors on main page