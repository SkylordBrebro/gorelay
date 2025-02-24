package models

import "strings"

// API endpoints
const (
	// ServerEndpoint is used to retrieve server list and character information
	ServerEndpoint = "https://realmofthemadgodhrd.appspot.com/char/list"

	// AssetEndpoint is used to retrieve the latest resources
	AssetEndpoint = "https://static.drips.pw/rotmg/production"

	// GithubContentEndpoint is used to retrieve file contents from the main repository
	GithubContentEndpoint = "https://api.github.com/repos/thomas-crane/nrelay/contents"

	// ClientVersionEndpoint is used to check the version of the latest client
	ClientVersionEndpoint = "https://www.realmofthemadgod.com/version.txt"

	// ClientDownloadEndpoint is used to retrieve the latest client
	// Replace {{version}} with the current version before use
	ClientDownloadEndpoint = "https://www.realmofthemadgod.com/AssembleeGameClient{{version}}.swf"
)

// GetClientDownloadURL returns the download URL for a specific client version
func GetClientDownloadURL(version string) string {
	return strings.Replace(ClientDownloadEndpoint, "{{version}}", version, 1)
}
