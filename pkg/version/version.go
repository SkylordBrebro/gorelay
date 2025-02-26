package version

import (
	"encoding/xml"
	"fmt"
	"gorelay/pkg/logger"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Add new function to fetch Unity build version
func FetchUnityBuildHash(logger *logger.Logger) (string, error) {
	baseURL := "https://www.realmofthemadgod.com/app/init"

	// Create HTTP client with appropriate headers
	client := &http.Client{}

	// Prepare form data
	data := make(url.Values)
	data.Set("platform", "standalonewindows64")
	data.Set("key", "9KnJFxtTvLu2frXv")
	data.Set("game_net", "Unity")
	data.Set("play_platform", "Unity")
	data.Set("game_net_user_id", "")
 
	// Create request
	req, err := http.NewRequest("POST", baseURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set Unity-specific headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Unity-Version", "2021.3.16f1")
	req.Header.Set("User-Agent", "UnityPlayer/2021.3.16f1 (UnityWebRequest/1.0, libcurl/7.84.0-DEV)")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Log response for debugging
	if logger != nil {
		logger.Debug("Client", "Init response: %s", string(body))
	}
	// Parse XML response
	doc := &struct {
		XMLName      xml.Name `xml:"AppSettings"`
		BuildHash    string   `xml:"BuildHash"`
	}{}

	if err := xml.Unmarshal(body, doc); err != nil {
		return "", fmt.Errorf("failed to parse XML response: %v", err)
	}

	if doc.BuildHash != "" {
		return doc.BuildHash, nil
	}

	return "", fmt.Errorf("neither BuildVersion nor BuildHash found in response")
}
