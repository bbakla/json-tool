# JSON Formatter & Key Finder

A small Go + Gin web app that formats JSON, searches keys/values, minifies JSON, converts JSON to YAML, and extracts key/value pairs as JSON. Uses Go templates with a base layout and minimal inline styling.

## Features
- Format JSON with pretty-print.
- Minify JSON (compact).
- Search for values by key and keys by value.
- Extract a key and its value(s) as JSON.
- Convert JSON to YAML.
- Copy formatted output to clipboard.

## Endpoints
- `GET /` – main page.
- `POST /format` – format JSON (default form action).
- `POST /minify` – compact JSON.
- `POST /toyaml` – JSON to YAML.
- `POST /find/key` – find values for a key.
- `POST /find/value` – find keys for a value.
- `POST /extract/key` – extract the key/value(s) as JSON.
- `GET /healthz` – health check.

## Running
From the project root:
```bash
cd /Volumes/External/src/workspace_projects/json_formatter
# standard run
go run .
```
Then open `http://localhost:8888`.

## Hot reload (Air)
If you installed Air (configured via `.air.toml`):
```bash
cd /Volumes/External/src/workspace_projects/json_formatter
air -c .air.toml
```

## Tests
Run all tests:
```bash
cd /Volumes/External/src/workspace_projects/json_formatter
go test ./...
```

## Layout
- `main.go`: Gin setup and handlers.
- `handlers/`: core JSON/YAML logic (format, minify, search, extract).
- `templates/layout.html`: base template wrapper.
- `templates/index.html`: page content template.

## Notes
- Invalid JSON or missing inputs return HTTP 400 while still rendering the page with an error banner.
- Uses CDN Bootstrap/Tailwind for styling; no local static assets are required.

