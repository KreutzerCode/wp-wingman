package pluginVersion

import (
	"regexp"
	"sort"
	"strings"
	"wp-wingman/types"
)

func GetPluginVersion(readmeContent string) (types.VersionNumber, bool) {
	versionData, found := versionNumbers(readmeContent)
	return versionData, found
}

// VersionNumbers extracts version numbers from the body and returns the first one found along with its source and confidence.
// It also returns a boolean indicating whether a version number was found.
func versionNumbers(body string) (types.VersionNumber, bool) {
	number := fromStableTag(body)
	if number != "" {
		return types.VersionNumber{Number: number, FoundBy: "Stable Tag"}, true
	}

	number = fromChangelogSection(body)
	if number != "" {
		return types.VersionNumber{Number: number, FoundBy: "ChangeLog Section"}, true
	}

	return types.VersionNumber{}, false
}

// FromStableTag extracts the version number from the stable tag in the body.
// It returns the version number if found, or an empty string if not found.
func fromStableTag(body string) string {
	re := regexp.MustCompile(`(?i)\b(?:stable tag|version):\s*([0-9a-z.-]+)`)
	matches := re.FindStringSubmatch(body)
	if len(matches) < 2 || matches[1] == "trunk" {
		return ""
	}

	number := matches[1]
	if strings.ContainsAny(number, "0123456789") {
		return number
	}

	return ""
}

func fromChangelogSection(body string) string {
	re := regexp.MustCompile(`(?i)^=+\s+(?:v(?:ersion)?\s*)?([0-9.-]+)[^=]*=+$`)
	matches := re.FindAllStringSubmatch(body, -1)

	var extractedVersions []string
	for _, match := range matches {
		if len(match) > 1 && strings.ContainsAny(match[1], "0123456789") {
			extractedVersions = append(extractedVersions, match[1])
		}
	}

	if len(extractedVersions) == 0 {
		return ""
	}

	sort.Strings(extractedVersions)

	return extractedVersions[len(extractedVersions)-1]
}
