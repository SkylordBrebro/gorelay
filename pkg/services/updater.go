package services

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// Updater handles downloading and updating game resources
type Updater struct {
	baseURL     string
	assetsPath  string
	versionsURL string
}

// NewUpdater creates a new updater instance
func NewUpdater(assetsPath string) *Updater {
	return &Updater{
		baseURL:     "https://static.drips.pw/rotmg/production",
		assetsPath:  assetsPath,
		versionsURL: "https://static.drips.pw/rotmg/production/current/version.txt",
	}
}

// UpdateInfo contains information about required updates
type UpdateInfo struct {
	NeedAssetUpdate  bool
	NeedClientUpdate bool
	AssetVersion     string
	ClientVersion    string
}

// CheckForUpdates checks if any updates are needed
func (u *Updater) CheckForUpdates() (*UpdateInfo, error) {
	// Get current version info
	resp, err := http.Get(u.versionsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get version info: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read version info: %v", err)
	}

	var versions struct {
		Assets  string `json:"assets"`
		Client  string `json:"client"`
		Current string `json:"current"`
	}

	if err := json.Unmarshal(data, &versions); err != nil {
		return nil, fmt.Errorf("failed to parse version info: %v", err)
	}

	// Compare with local versions
	localVersions, err := u.loadLocalVersions()
	if err != nil {
		return nil, err
	}

	return &UpdateInfo{
		NeedAssetUpdate:  versions.Assets != localVersions.AssetVersion,
		NeedClientUpdate: versions.Client != localVersions.ClientVersion,
		AssetVersion:     versions.Assets,
		ClientVersion:    versions.Client,
	}, nil
}

// PerformUpdate updates the necessary resources
func (u *Updater) PerformUpdate(info *UpdateInfo) error {
	if info.NeedAssetUpdate {
		if err := u.updateAssets(); err != nil {
			return fmt.Errorf("failed to update assets: %v", err)
		}
	}

	if info.NeedClientUpdate {
		if err := u.updateClient(); err != nil {
			return fmt.Errorf("failed to update client: %v", err)
		}
	}

	// Update local version info
	return u.saveVersions(info)
}

// Helper methods

func (u *Updater) loadLocalVersions() (*UpdateInfo, error) {
	path := filepath.Join(u.assetsPath, "versions.json")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &UpdateInfo{}, nil
		}
		return nil, fmt.Errorf("failed to read versions file: %v", err)
	}

	var info UpdateInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("failed to parse versions file: %v", err)
	}

	return &info, nil
}

func (u *Updater) saveVersions(info *UpdateInfo) error {
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal version info: %v", err)
	}

	path := filepath.Join(u.assetsPath, "versions.json")
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write versions file: %v", err)
	}

	return nil
}

func (u *Updater) updateAssets() error {
	// Download Objects.json
	if err := u.downloadFile("Objects.json", "Objects.json"); err != nil {
		return err
	}

	// Download Tiles.json
	if err := u.downloadFile("Tiles.json", "Tiles.json"); err != nil {
		return err
	}

	return nil
}

func (u *Updater) updateClient() error {
	// Download client.swf or other client files
	return u.downloadFile("client.swf", "client.swf")
}

func (u *Updater) downloadFile(remotePath, localPath string) error {
	url := fmt.Sprintf("%s/%s", u.baseURL, remotePath)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download %s: %v", remotePath, err)
	}
	defer resp.Body.Close()

	path := filepath.Join(u.assetsPath, localPath)
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", localPath, err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %v", localPath, err)
	}

	return nil
}
