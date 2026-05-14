# envlens

> A CLI utility for diffing and auditing environment variable files across staging and production configs.

---

## Installation

```bash
go install github.com/yourusername/envlens@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envlens.git && cd envlens && go build -o envlens .
```

---

## Usage

Compare two environment files and surface missing, added, or mismatched keys:

```bash
# Diff a staging env file against production
envlens diff .env.staging .env.production

# Audit a single file for empty or suspicious values
envlens audit .env.production

# Output results as JSON
envlens diff .env.staging .env.production --format json
```

**Example output:**

```
[MISSING]  DATABASE_URL        found in staging, not in production
[ADDED]    NEW_RELIC_KEY       found in production, not in staging
[MISMATCH] LOG_LEVEL           staging=debug  production=info
```

---

## Flags

| Flag | Description |
|------|-------------|
| `--format` | Output format: `text` (default) or `json` |
| `--ignore` | Comma-separated list of keys to ignore |
| `--strict` | Exit with non-zero code if any differences are found |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)