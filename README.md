# Permacomputing Web Template

A production-ready, full-stack web application template designed for radical simplicity and zero maintenance for 50+ years.

## Philosophy

This template embodies **permacomputing principles**:
- Zero external dependencies at runtime
- Single binary deployment
- No build tools required
- No npm, no webpack, no transpilers
- Standard library first
- Simple, readable code

## Stack

- **Backend**: Go 1.22+ (using new `http.ServeMux` features)
- **Database**: SQLite with WAL mode
- **Frontend**: HTMX for dynamic interactions
- **Templating**: Go `html/template`
- **Styling**: Modern CSS with semantic HTML

## Features

- ✅ Single binary compilation via `//go:embed`
- ✅ Automatic database migrations
- ✅ HTMX-powered dynamic UI (no JavaScript frameworks)
- ✅ Logging middleware
- ✅ Panic recovery middleware
- ✅ Responsive design with dark mode support
- ✅ Zero configuration required

## Project Structure

```
.
├── main.go                      # Entry point
├── internal/
│   ├── db/
│   │   └── db.go               # SQLite setup and migrations
│   └── handlers/
│       └── routes.go           # HTTP handlers
├── ui/
│   ├── html/
│   │   ├── base.layout.tmpl   # Base template
│   │   └── home.page.tmpl     # Home page
│   └── static/
│       ├── css/
│       │   └── style.css      # Semantic CSS
│       └── js/
│           └── htmx.min.js    # HTMX library
└── app.db                      # SQLite database (auto-created)
```

## Getting Started

### Prerequisites

- Go 1.22 or higher
- That's it!

### Download HTMX

Before building, download the HTMX library:

```bash
curl -sL https://unpkg.com/htmx.org@1.9.10/dist/htmx.min.js -o ui/static/js/htmx.min.js
```

Or with wget:

```bash
wget https://unpkg.com/htmx.org@1.9.10/dist/htmx.min.js -O ui/static/js/htmx.min.js
```

### Install Dependencies

```bash
go mod download
```

### Run in Development

```bash
go run main.go
```

The application will start on `http://localhost:8080`

### Build for Production

Build a single binary with all assets embedded:

```bash
go build -o perma-app main.go
```

The resulting binary is completely self-contained. Copy it anywhere and run:

```bash
./perma-app
```

### Build Optimized Binary

For a smaller binary size:

```bash
go build -ldflags="-s -w" -o perma-app main.go
```

Optional: Compress with UPX:

```bash
upx --best --lzma perma-app
```

## Configuration

The application uses sensible defaults. You can modify these constants in `main.go`:

```go
const (
    dbPath  = "./app.db"     // Database file path
    port    = "8080"         // Server port
    timeout = 30 * time.Second
)
```

## Database

SQLite is configured with:
- **WAL mode** for better concurrency
- **Automatic migrations** on first run
- **Connection pooling** with reasonable limits

The database file `app.db` is created automatically on first run.

## HTMX Demo

The application includes a "click-to-edit" user profile example demonstrating:
- Inline editing without page reloads
- HTML over the wire (no JSON APIs needed)
- Graceful degradation

Try it:
1. Visit `http://localhost:8080`
2. Click "Edit" on any user card
3. Modify the user's name or email
4. Click "Save" or "Cancel"

## Middleware

Two middleware functions are implemented:

1. **Logging**: Logs all HTTP requests with duration
2. **Panic Recovery**: Catches panics and returns 500 errors gracefully

## CSS Features

The CSS uses a semantic/classless approach:
- CSS custom properties (variables) for easy theming
- Automatic dark mode support via `prefers-color-scheme`
- Mobile-responsive design
- No CSS frameworks, no build steps

## Production Deployment

1. Build the binary: `go build -ldflags="-s -w" -o perma-app main.go`
2. Copy the binary to your server
3. Run it: `./perma-app`
4. Optional: Set up systemd service for auto-restart

### Example Systemd Service

Create `/etc/systemd/system/perma-app.service`:

```ini
[Unit]
Description=Permacomputing Web Application
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/perma-app
ExecStart=/opt/perma-app/perma-app
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable perma-app
sudo systemctl start perma-app
```

## Extending the Application

### Adding a New Route

1. Define the handler in `internal/handlers/routes.go`
2. Register it in `main.go`:

```go
mux.HandleFunc("GET /your-route", app.YourHandler)
```

### Adding a New Template

1. Create `ui/html/yourpage.page.tmpl`
2. Use `{{template "base" .}}` to include the base layout
3. Render it from your handler:

```go
app.Templates.ExecuteTemplate(w, "yourpage.page.tmpl", data)
```

### Adding Database Tables

1. Add migration SQL to `internal/db/db.go` in the `migrations` slice
2. Delete `app.db` to re-run migrations (or write a proper migration system)

## Performance

This template is designed for low resource usage:
- Binary size: ~10-15MB (uncompressed)
- Memory usage: <50MB for typical workloads
- Startup time: <100ms

## Security Considerations

- Use HTTPS in production (reverse proxy with Caddy/Nginx recommended)
- Implement authentication/authorization as needed
- Validate all user input
- Use prepared statements (already done via `database/sql`)
- Keep Go updated for security patches

## Why This Stack?

### Go
- Fast, compiled language
- Excellent standard library
- Single binary deployment
- Strong backward compatibility

### SQLite
- Zero configuration
- File-based (no separate server)
- ACID compliant
- Battle-tested and reliable

### HTMX
- <15KB library
- No build step required
- Hypermedia-driven (REST-ful)
- Works without JavaScript (progressive enhancement)

### No Frameworks
- Less code to maintain
- Fewer dependencies to break
- Better performance
- Easier to understand and modify

## License

This template is released into the public domain. Use it however you want.

## Contributing

This is a template project. Fork it and make it your own!

## Support

For issues or questions, please open an issue on the repository.

## Donate 

To support the developer, you can <a href="https://buy.stripe.com/14AaEWfRXcm5aPrfSag360b">make a donation here<a>.

---

**Built with ❤️ for longevity and simplicity**
