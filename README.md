# MagicSort

MagicSort is a high-performance, extension-agnostic file organizer and background daemon written in Go.

Unlike traditional file organizers that rely on easily spoofed or missing file extensions, such as checking whether a file ends with `.pdf`, MagicSort reads the first 512 bytes of a file to determine its true identity using **Magic Bytes**, also known as hexadecimal file signatures.

This provides accurate file type detection even when a file is downloaded as `.tmp`, `.crdownload`, has an incorrect extension, or has no extension at all.

---

## Key Features

* **True File Type Detection:** Uses magic byte analysis to determine file types securely. A disguised executable cannot bypass the system simply by changing its extension to `.jpg`.

* **Zero-CPU Background Monitoring:** Uses `fsnotify` to listen directly to operating system kernel events such as `Create`, `Rename`, and `CloseWrite`. When no file events occur, the daemon sleeps and consumes near-zero CPU.


* **Zero Dependencies:** Compiles into a single standalone binary. No external runtime environment is required on the host machine.

* **Dual Architecture: CLI & Daemon:** Powered by Cobra, MagicSort works both as a 24/7 background service and as an on-demand command-line tool.

---

## Installation

Make sure Go is installed on your system.

Clone the repository and build the binary:

```bash
git clone https://github.com/barisaydogdu/magic-byte-organizer.git
cd magic-byte-organizer
go build -o magicsort ./cmd/main.go
```

Optional: move the binary to `/usr/local/bin` so it can be used globally:

```bash
sudo mv magicsort /usr/local/bin/
```

Verify the installation:

```bash
magicsort --help
```

---

## Usage

MagicSort has two primary commands:

* `start`: continuously monitors a directory in real time
* `scan`: organizes an existing directory once and exits

---

## 1. Continuous Watcher: `start`

The `start` command runs MagicSort as a watcher/daemon and monitors the target directory in real time.

```bash
magicsort start --dir /home/user/Downloads
```

### Flags

| Flag                   | Description                                                                   |
| ---------------------- | ----------------------------------------------------------------------------- |
| `-d`, `--dir string`   | Target directory to watch                                                     |
| `-t`, `--delay string` | Wait time before moving the file. Examples: `10m`, `30s`, `1h`. Default: `0s` |
| `--dry-run`            | Simulation mode. Logs what would be moved without changing files on disk      |

### Example

```bash
magicsort start -d /home/user/Downloads -t 10m --dry-run
```

---

## 2. One-Time Scanner: `scan`

The `scan` command scans and organizes all existing files in the given directory once. It does not start a continuous watcher.

```bash
magicsort scan --dir /home/user/Downloads
```

### Flags

| Flag                 | Description                                                   |
| -------------------- | ------------------------------------------------------------- |
| `-d`, `--dir string` | Target directory to scan                                      |
| `--dry-run`          | Simulation mode. Logs the execution plan without moving files |

### Example

```bash
magicsort scan -d /home/user/Downloads --dry-run
```

---

## Configuration

MagicSort uses a `config.json` file to map detected MIME types to target folders.

The application includes smart path resolution, so it can locate the configuration file whether it is executed with `go run` or as a compiled background service.

### Example `config.json`

```json
{
  "image/png": "~/Pictures",
  "image/jpeg": "~/Pictures",
  "application/pdf": "~/Documents",
  "text/plain": "~/Desktop/Texts"
}
```

### Note About Plain Text Files

Text-based files such as `.txt`, `.csv`, `.md`, and `.go` usually do not contain magic bytes. They may be detected as:

```text
text/plain
```

Map `text/plain` carefully to avoid unintentionally moving source code, notes, or project files.

---



## Dry Run Recommendation

Before running MagicSort on important directories, test the behavior with `--dry-run`.

Example:

```bash
magicsort scan -d /home/user/Downloads --dry-run
```

This shows which files would be moved without actually changing the file system.
