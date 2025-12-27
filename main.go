package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func findDisk(label string) (string, error) {
	labelPath := filepath.Join("/dev/disk/by-label", label)

	if _, err := os.Stat(labelPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("disk with label %q not connected", label)
		}
		return "", fmt.Errorf("failed to stat label path: %w", err)
	}

	device, err := filepath.EvalSymlinks(labelPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve symlink for %s: %w", labelPath, err)
	}

	return device, nil
}

func findMountPoint(device string) (string, bool) {
	f, err := os.Open("/proc/self/mounts")
	if err != nil {
		log.Fatalf("failed to open mounts file: %v", err)
		return "", false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 && fields[0] == device {
			return fields[1], true
		}
	}

	return "", false
}

func mountDisk(device, mountPoint string) error {
	log.Printf("Mounting device %s at %s", device, mountPoint)

	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return fmt.Errorf("failed to create mount point %s: %w", mountPoint, err)
	}

	cmd := exec.Command("mount", device, mountPoint)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mount command failed: %w", err)
	}

	log.Println("Disk mounted successfully")
	return nil
}

func unmountDisk(mountPoint string) error {
	log.Printf("Unmounting %s", mountPoint)

	cmd := exec.Command("umount", mountPoint)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("warning: failed to unmount %s: %v", mountPoint, err)
	}

	log.Println("Disk unmounted successfully")
	return nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	label := "XPG_NVME"
	mountPoint := "/mnt/backup_disk"
	mountedByScript := false

	device, err := findDisk(label)
	if err != nil {
		log.Fatalf("Backup disk unavailable: %v", err)
		os.Exit(1)
	}

	log.Printf("Backup disk found: %s", device)

	if mp, ok := findMountPoint(device); ok {
		log.Printf("Backup disk already mounted at %s", mp)
	} else {
		log.Println("Backup disk not mounted")
		if err := mountDisk(device, mountPoint); err != nil {
			log.Fatalf("Failed to mount disk: %v", err)
			mountedByScript = false
			os.Exit(1)
		}
		mountedByScript = true
	}

	if mountedByScript {
		defer unmountDisk(mountPoint)
	}

	log.Println("Backup preparation complete")
}
