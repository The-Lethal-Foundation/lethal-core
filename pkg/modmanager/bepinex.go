package modmanager

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/The-Lethal-Foundation/lethal-core/pkg/config"
	"github.com/The-Lethal-Foundation/lethal-core/pkg/filesystem"
)

// BepInExRelease represents the structure of a GitHub release.
type BepInExRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

const bepInExRepoURL = "https://api.github.com/repos/BepInEx/BepInEx/releases/latest"
const bepInExCacheDir = filesystem.DefaultCacheDir

// FetchLatestBepInExVersion fetches the latest stable release version of BepInEx from GitHub.
func FetchLatestBepInExVersion() (string, error) {
	resp, err := http.Get(bepInExRepoURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release BepInExRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

// DownloadAndCacheBepInEx downloads and caches the latest BepInEx release.
func DownloadAndCacheBepInEx(basePath string, version string) (string, error) {
	cachePath := filepath.Join(basePath, bepInExCacheDir)
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return "", err
	}

	// Fetch release information.
	resp, err := http.Get(bepInExRepoURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release BepInExRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	// Download the release zip file.
	zipResp, err := http.Get(release.Assets[0].BrowserDownloadURL)
	if err != nil {
		return "", err
	}
	defer zipResp.Body.Close()

	zipPath := filepath.Join(cachePath, "BepInEx.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	_, err = io.Copy(zipFile, zipResp.Body)
	if err != nil {
		return "", err
	}

	return zipPath, nil
}

// IsBepInExUpToDate checks if the cached version of BepInEx matches the latest stable version.
func IsBepInExUpToDate(basePath string, currentConfig *config.Config) (bool, error) {
	latestVersion, err := FetchLatestBepInExVersion()
	if err != nil {
		return false, err
	}

	return currentConfig.CachedBepInExVersion == latestVersion, nil
}

func IsBepInExCached(basePath string) bool {
	cachePath := filepath.Join(basePath, bepInExCacheDir, "BepInEx.zip")
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// UnpackBepInEx unpacks the BepInEx zip into the specified directory.
func UnpackBepInEx(targetDir string) error {
	zipPath := filepath.Join(filesystem.GetDefaultPath(), bepInExCacheDir, "BepInEx.zip")

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(targetDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
