# logsnap

A lightweight CLI toolshots across deployments.

---

## Installation

```bash
go install github.com/yourname/logsnap@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/logsnap.git && cd logsnap && go build -o logsnap .
```

---

## Usage

Capture a snapshot of your current logs:

```bash
logsnap capture --source ./logs/app.log --out snapshot-v1.snap
```

Diff two snapshots across deployments:

```bash
logsnap diff snapshot-v1.snap snapshot-v2.snap
```

Example output:

```
[+] NEW     error: "database connection timeout" (x3)
[-] REMOVED warn:  "cache miss rate high"
[~] CHANGED info:  "startup time" 1.2s → 0.9s
```

Filter by log level:

```bash
logsnap capture --source ./logs/app.log --level error --out snapshot-v3.snap
```

---

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--source` | Path to log file or directory | `stdin` |
| `--out` | Output snapshot file | `snapshot.snap` |
| `--level` | Filter by log level | `all` |
| `--format` | Log format (`json`, `logfmt`) | `json` |

---

## License

MIT © yourname