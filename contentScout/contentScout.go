package contentScout

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"wp-wingman/pluginFinder"
	"wp-wingman/types"
)

func FindPluginsInContent(url string, pluginsFoundWithAggressiveMode []types.PluginData, useRandomUserAgent bool, rateLimit int, workerCount int)[]types.PluginData {
	uniqueMatches := findPluginsInContent(url)
	missingPlugins := returnPluginsThatAreNotFoundAlready(uniqueMatches, pluginsFoundWithAggressiveMode)
	return pluginFinder.CheckPluginSlugsOnTarget(url, missingPlugins, workerCount, useRandomUserAgent, "content", rateLimit)
}

func returnPluginsThatAreNotFoundAlready(uniqueMatches map[string]bool, pluginsFoundWithAggressiveMode []types.PluginData)[]string {
	missingPlugins := []string{}

    pluginDataMap := make(map[string]bool)
    for _, pluginData := range pluginsFoundWithAggressiveMode {
        pluginDataMap[pluginData.Name] = true
    }

    for match := range uniqueMatches {
        if !pluginDataMap[match] {
			missingPlugins = append(missingPlugins, match)
        }
    }

    return missingPlugins
}

func findPluginsInContent(url string)map[string]bool {
	matches := getPluginSlugsFromContent(url)

    // Use a map to filter out duplicate entries
    uniqueMatches := make(map[string]bool)
	for match := range matches {
		// Remove the '/wp-content/plugins/' part from each entry
		strippedMatch := strings.Replace(match, "/wp-content/plugins/", "", 1)
		uniqueMatches[strippedMatch] = true
	}

	return uniqueMatches
}


func getPluginSlugsFromContent(url string) (map[string]bool) {
	matches, _ := findPluginSlugsInContent(url)
    // Use a map to filter out duplicate entries
    uniqueMatches := make(map[string]bool)
    for _, match := range matches {
        // Remove the '/wp-content/plugins/' part from each entry
        strippedMatch := strings.Replace(match, "/wp-content/plugins/", "", 1)
        uniqueMatches[strippedMatch] = true
    }

	return uniqueMatches
}

func findPluginSlugsInContent(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []string{}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	re := regexp.MustCompile(`(/wp-content/plugins/[^/]+)`)

	// Find all substrings that match the pattern
	matches := re.FindAllString(string(body), -1)

	return matches, nil
}
