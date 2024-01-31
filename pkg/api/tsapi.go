package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ModDetailsResponse struct {
	Downloads     int    `json:"downloads"`
	RatingScore   int    `json:"rating_score"`
	LatestVersion string `json:"latest_version"`
}

func FetchModDetails(modAuthor, modName string) (*ModDetailsResponse, error) {
	url := fmt.Sprintf("https://thunderstore.io/api/v1/package-metrics/%s/%s", modAuthor, modName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var modInfo ModDetailsResponse
	if err := json.Unmarshal(body, &modInfo); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return &modInfo, nil
}

// DownloadModPackage downloads the specified mod package as a zip file to a temporary location.
func DownloadModPackage(modAuthor, modName, modVersion string) (string, error) {
	downloadURL := fmt.Sprintf("https://gcdn.thunderstore.io/live/repository/packages/%s-%s-%s.zip", modAuthor, modName, modVersion)

	// Attempt to download the mod with retries for handling rate limits or temporary network issues.
	maxRetries := 5
	for attempt := 1; attempt <= maxRetries; attempt++ {
		tmpFileName, err := tryDownloadMod(downloadURL)
		if err == nil {
			return tmpFileName, nil
		}

		if attempt < maxRetries && err == errTooManyRequests {
			time.Sleep(2 * time.Second) // Sleep before retrying.
			continue
		}

		return "", err // Return the error if max retries reached or other error occurred.
	}

	return "", fmt.Errorf("failed to download mod after %d attempts", maxRetries)
}

// tryDownloadMod performs a single attempt to download the mod from the given URL.
func tryDownloadMod(downloadURL string) (string, error) {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("error making download request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return saveModToFile(resp.Body)
	case http.StatusTooManyRequests:
		return "", errTooManyRequests // Custom error to indicate retry.
	default:
		return "", fmt.Errorf("received non-OK response status: %s", resp.Status)
	}
}

// errTooManyRequests is a custom error to indicate that the API has returned a rate limit response.
var errTooManyRequests = fmt.Errorf("received too many requests response")

// saveModToFile saves the mod from the response body to a temporary file.
func saveModToFile(body io.Reader) (string, error) {
	tmpFile, err := os.CreateTemp("", "mod-*.zip")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, body); err != nil {
		return "", fmt.Errorf("error writing to temp file: %w", err)
	}

	return tmpFile.Name(), nil
}
