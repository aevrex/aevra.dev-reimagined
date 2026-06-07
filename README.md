# Aevra

A small Go and HTMX website built in the same style as Fisher.

## Structure

- `main.go` contains the data, handlers, and routes.
- `templates/layout.html` contains the page shell.
- `templates/partials` contains the pages HTMX swaps into `#content`.
- Tailwind and HTMX load from their browser CDNs.

There is no database, build step, custom CSS, or custom JavaScript.

```powershell
go run .
```

Open `http://127.0.0.1:8091`.
