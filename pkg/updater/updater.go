package updater

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	buildServerURL  = "https://rotmg-build.decagames.com/build-release/"
	buildServerDir  = "/rotmg-exalt-win-64/"
)

// BuildFile represents a file in the build manifest
type BuildFile struct {
	File     string `json:"file"`
	Size     int    `json:"size"`
	Checksum string `json:"checksum"`
}

// BuildClient handles downloading game resources
type BuildClient struct {
	BuildHash string
	client    *http.Client
}

// NewBuildClient creates a new build client instance
func NewBuildClient(buildHash string) *BuildClient {
	return &BuildClient{
		BuildHash: buildHash,
		client:    &http.Client{},
	}
}

// GetBuildFilesList fetches and parses the checksum.json file
func (b *BuildClient) GetBuildFilesList() ([]BuildFile, error) {
	url := fmt.Sprintf("%s%s%schecksum.json", buildServerURL, b.BuildHash, buildServerDir)
	
	resp, err := b.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch checksum.json: %v", err)
	}
	defer resp.Body.Close()

	var data struct {
		Files []BuildFile `json:"files"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode checksum.json: %v", err)
	}

	/* Log file information
	for _, file := range data.Files {
		log.Printf("File: %s, Checksum: %s, Size: %d", file.File, file.Checksum, file.Size)
	}*/

	return data.Files, nil
}

// GetResource downloads a specific resource file
func (b *BuildClient) GetResource(path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s%s%s.gz", buildServerURL, b.BuildHash, buildServerDir, path)
	log.Printf("Downloading %s", url)

	resp, err := b.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download resource: %v", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// DownloadResourceAssets downloads and decompresses the resources.assets file
func DownloadGameFile(buildHash string, filename string) ([]byte, error) {
	client := NewBuildClient(buildHash)

	// Get the file list to verify size
	files, err := client.GetBuildFilesList()
	if err != nil {
		return nil, fmt.Errorf("failed to get build files list: %w", err)
	}

	// Find resources.assets in the file list
	var resourceFile *BuildFile
	for _, file := range files {
		if strings.HasSuffix(file.File, filename) {
			resourceFile = &file
			break
		}
	}

	if resourceFile == nil {
		return nil, fmt.Errorf("resources.assets not found in build files list")
	}

	// Download the compressed file
	compressedBytes, err := client.GetResource(resourceFile.File)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %w", filename, err)
	}

	// Verify compressed size
	if len(compressedBytes) != resourceFile.Size {
		return nil, fmt.Errorf("file size %d does not match checksum size %d", len(compressedBytes), resourceFile.Size)
	}

	// Decompress the file
	decompressed, err := decompress(compressedBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress %s: %w", filename, err)
	}

	log.Printf("Successfully downloaded and decompressed %s", filename)
	return decompressed, nil
}

func decompress(compressedData []byte) ([]byte, error) {
	reader := bytes.NewReader(compressedData)
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, gzipReader); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

//this is the main function that will download the required files and return parsed data
func DoUpdate(buildHash string) (string, map[string]string, error) {
	globalMetadata, err := DownloadGameFile(buildHash, "global-metadata.dat")
	if err != nil {
		log.Fatal(err)
		return "", nil, err
	}
	
	version, err := FindVersionFromMetadata(globalMetadata)
	if err != nil {
		log.Fatal(err)
		return "", nil, err
	}
	
	resourceData, err := DownloadGameFile(buildHash, "resources.assets")
	if err != nil {
		log.Fatal(err)
		return "", nil, err
	}
	
	xmls, err := ExtractXmlsFromResources(resourceData)
	if err != nil {
		log.Fatal(err)
		return "", nil, err
	}
	
	return version, xmls, nil
}

func FindVersionFromMetadata(metadata []byte) (string, error) {
	// Convert bytes to string for easier searching
	text := string(metadata)
	
	// Find the last occurrence of "127.0.0.1"
	localhost := "127.0.0.1"
	lastIndex := strings.LastIndex(text, localhost)
	
	if lastIndex == -1 {
		return "", fmt.Errorf("unable to find '127.0.0.1' in global-metadata.dat")
	}
	
	// Start searching after "127.0.0.1"
	versionStart := lastIndex + len(localhost) + 1
	
	// Build the version string
	var version strings.Builder
	readingVersion := false
	
	// Only search ahead 50 characters maximum
	maxPos := versionStart + 50
	if maxPos > len(text) {
		maxPos = len(text)
	}
	
	for i := versionStart; i < maxPos; i++ {
		c := text[i]
		if isVersionChar(c) {
			version.WriteByte(c)
			if !readingVersion {
				readingVersion = true
			}
		} else if readingVersion {
			break
		}
	}
	
	if version.Len() == 0 {
		return "", fmt.Errorf("no version found after '127.0.0.1' marker")
	}
	
	return version.String(), nil
}

// Helper function to check if a character is valid for version string
func isVersionChar(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '.'
}

func ExtractXmlsFromResources(resources []byte) (map[string]string, error) {
	text := string(resources)
	xmlHeader := "<?xml"
	headers := allIndexesOf(text, xmlHeader)
	xmlDict := make(map[string][]string)
	tags := []string{"Objects", "GroundTypes"/*, "Enchantments", "DungeonModifiers"*/}

	// Process each XML header found
	for index := 0; index < len(headers); index++ {
		headerPosition := headers[index]
		
		var cap int
		if index == len(headers)-1 {
			cap = len(text)
		} else {
			cap = headers[index+1]
		}

		foundHeaderEnd := -1
		foundTagStart := -1
		foundTagEnd := -1
		endTag := ""

		// Parse XML structure
		for i := headerPosition; i < cap; i++ {
			current := text[i]

			if foundTagEnd != -1 {
				// Search for closing tag (</Objects>, etc)
				if i+len(endTag) <= len(text) && text[i:i+len(endTag)] == endTag {
					tag := strings.NewReplacer("<", "", ">", "", "/", "", "\\", "").Replace(text[i:i+len(endTag)])
					
					if !contains(tags, tag) {
						break
					}

					xml := text[foundTagStart : i+len(endTag)]
					xml = strings.Replace(xml, "<"+tag+">", "", 1)
					xml = strings.Replace(xml, "</"+tag+">", "", 1)

					xmlDict[tag] = append(xmlDict[tag], xml)
					break
				}
				continue
			}

			if foundTagStart != -1 {
				// Find end of opening tag
				if current == '>' {
					tagName := text[foundTagStart:i+1]
					endTag = strings.Replace(tagName, "<", "</", 1)
					foundTagEnd = i
				}
				continue
			}

			if foundHeaderEnd != -1 {
				// Find start of tag name
				if current == '<' {
					foundTagStart = i
				}
				continue
			}

			if current == '>' {
				foundHeaderEnd = i
			}
		}
	}

	// Combine and validate XMLs
	result := make(map[string]string)
	for _, tag := range tags {
		xmlList, exists := xmlDict[tag]
		if !exists {
			return nil, fmt.Errorf("missing required XML tag: %s", tag)
		}

		combined := combineXMLs(tag, xmlList)
		result[tag] = combined
	}

	return result, nil
}

func combineXMLs(tag string, xmls []string) string {
	var sb strings.Builder
	sb.WriteString("<")
	sb.WriteString(tag)
	sb.WriteString(">")
	for _, xml := range xmls {
		sb.WriteString(xml)
		sb.WriteString("\n")
	}
	sb.WriteString("</")
	sb.WriteString(tag)
	sb.WriteString(">")
	return sb.String()
}

func allIndexesOf(str, substr string) []int {
	if substr == "" {
		return nil
	}

	var indexes []int
	for i := 0; ; {
		idx := strings.Index(str[i:], substr)
		if idx == -1 {
			break
		}
		indexes = append(indexes, i+idx)
		i += idx + len(substr)
	}
	return indexes
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}