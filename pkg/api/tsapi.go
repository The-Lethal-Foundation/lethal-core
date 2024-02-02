package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

type OrderingType string

const (
	LastUpdated    OrderingType = "last-updated"
	Newest         OrderingType = "newest"
	MostDownloaded OrderingType = "most-downloaded"
	TopRated       OrderingType = "top-rated"
)

type SectionType string

const (
	Mods              SectionType = "mods"
	AssetReplacements SectionType = "asset-replacements"
	Libraries         SectionType = "libraries"
	Modpacks          SectionType = "modpacks"
)

type GlobalModView struct {
	ModAuthor  string `json:"mod_author"`
	ModName    string `json:"mod_name"`
	ModPicture string `json:"mod_picture"`
}

// GlobalListMods fetches and parses the mod list from Thunderstore.
func GlobalListMods(ordering OrderingType, sectionType SectionType, query string, page int) ([]GlobalModView, error) {
	document, err := fetchModsDocument(ordering, sectionType, query, page)
	if err != nil {
		return nil, fmt.Errorf("error fetching mods document: %w", err)
	}

	return parseModsDocument(document), nil
}

// fetchModsDocument retrieves the HTML document from Thunderstore.
func fetchModsDocument(ordering OrderingType, sectionType SectionType, query string, page int) (*goquery.Document, error) {
	encodedQuery := url.QueryEscape(query)
	reqUrl := fmt.Sprintf("https://thunderstore.io/c/lethal-company/?q=%s&ordering=%s&section=%s&page=%d", encodedQuery, ordering, sectionType, page)

	response, err := http.Get(reqUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching data: %s", response.Status)
	}

	return goquery.NewDocumentFromReader(response.Body)
}

// parseModsDocument parses the goquery document to extract mod information.
func parseModsDocument(document *goquery.Document) []GlobalModView {
	var mods []GlobalModView
	document.Find("div.col-md-4").Each(func(index int, element *goquery.Selection) {
		modName := element.Find("div > h5").Text()
		modPicture, _ := element.Find("div > a > img").Attr("src")
		modAuthor := strings.Trim(element.Find("div:nth-child(2) > div:nth-child(3) > a").Text(), " \n")

		if modName != "" {
			mods = append(mods, GlobalModView{
				ModAuthor:  modAuthor,
				ModName:    modName,
				ModPicture: modPicture,
			})
		}
	})

	return mods
}
