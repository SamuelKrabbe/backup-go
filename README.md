# Backup Tool (Go)

A simple, reliable backup utility written in Go that automatically detects an external disk by **filesystem label**, mounts it if necessary, and prepares it for backup operations.

Designed for personal backups (e.g. family data) with an emphasis on **safety**, **explicit behavior**, and **extensibility**.

---

## Features

* Detects external disks by **filesystem label**
* Resolves real device paths (`/dev/sdX`)
* Checks whether the disk is already mounted
* Automatically mounts the disk if needed
* Cleanly unmounts disks mounted by the tool
* Clear, timestamped logging
* Minimal dependencies (standard library only)

---

## Requirements

* Linux
* Go 1.20+
* Root privileges (required for mounting)
* A disk with a filesystem label (e.g. `SSD_NVME`)

---

## Disk Preparation Example

Example using `exfat`:

```bash
sudo mkfs.exfat -n SSD_NVME /dev/sda1
```

Verify the label:

```bash
lsblk -f
```

Expected output:

```
sda1 exfat SSD_NVME
```

---

## How It Works

1. Looks for the disk at:

   ```
   /dev/disk/by-label/<LABEL>
   ```
2. Resolves the symlink to the real device
3. Checks `/proc/self/mounts` to see if it is already mounted
4. If not mounted:

   * Creates the mount point
   * Mounts the disk
   * Remembers it was mounted by the program
5. Automatically unmounts on exit (if mounted by the tool)

---

## Usage

Run with root privileges:

```bash
sudo go run main.go
```

Typical output:

```
2025/12/27 10:42:10 Disk found: /dev/sda1
2025/12/27 10:42:10 Disk not mounted
2025/12/27 10:42:10 Mounting device /dev/sda1 at /mnt/backup_disk
2025/12/27 10:42:11 Disk mounted successfully
```

---

## Configuration

Currently hardcoded (will be configurable later):

```go
label      := "SSD_NVME"
mountPoint := "/mnt/backup_disk"
```

---

## Safety Notes

* The program **never formats disks**
* It only mounts disks explicitly identified by label
* It does not unmount disks it did not mount
* Errors are logged clearly and terminate safely when needed

---

## Roadmap

Planned improvements:

* Backup profiles (per user / directory)
* Incremental backups
* Hash verification
* Exclude rules
* Dry-run mode
* Configuration file support
* Structured (JSON) logging
* Systemd / cron integration

---

## Philosophy

This tool is intentionally:

* Explicit over clever
* Predictable over magical
* Safe over fast

Backups should be boring — and trustworthy.

---

## License

MIT
