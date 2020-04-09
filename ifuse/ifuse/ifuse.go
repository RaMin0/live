package ifuse

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	ifuseBinPath = "/usr/local/bin/ifuse"
)

type App struct {
	ID      string
	Version string
	Name    string
}

func (a App) String() string {
	return fmt.Sprintf("%s (id=%s, version=%s)", a.Name, a.ID, a.Version)
}

func ListApps(serialNumber string) ([]App, error) {
	appsStr, err := Run("-u", serialNumber, "--list-apps")
	if err != nil {
		return nil, err
	}

	var apps []App
	appsArr := strings.Split(strings.TrimSpace(appsStr), "\n")
	for _, appStr := range appsArr {
		appArr := strings.Split(appStr, ", ")
		apps = append(apps, App{
			ID:      appArr[0],
			Version: strings.Trim(appArr[1], "\""),
			Name:    strings.Trim(appArr[2], "\""),
		})
	}
	return apps, nil
}

func MountAppDocuments(serialNumber string, a App) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	mountPath := path.Join(homeDir, ".fuse", "mnt")
	if err := os.MkdirAll(mountPath, 0644); err != nil {
		return "", err
	}
	if _, err := Run(
		"-u", serialNumber,
		"--documents", a.ID,
		mountPath,
	); err != nil {
		return "", err
	}
	return mountPath, nil
}

func Run(args ...string) (string, error) {
	cmd := exec.Command(ifuseBinPath, args...)
	var b bytes.Buffer
	cmd.Stdout, cmd.Stderr = &b, &b
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("run: %v", err)
	}
	return b.String(), nil
}
