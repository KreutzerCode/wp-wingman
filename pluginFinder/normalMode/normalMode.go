package normalMode

import (
	"fmt"
	"math/rand"
	"time"
	"wp-wingman/pluginVersion"
	"wp-wingman/store"
	"wp-wingman/types"
	"wp-wingman/utils"
)

func ReturnPluginVersion(url string, pluginsPrefix string, pluginName string, randomUserAgent bool, maxStringLength int) (types.VersionNumber, bool){
	version := types.VersionNumber{}
	versionFound := false
	result, err := utils.FetchReadme(fmt.Sprintf("%s/%s/%s/%s", url, pluginsPrefix, pluginName, "readme.txt"), randomUserAgent)

	if err != nil {
		fmt.Printf("\033[1;31m%-*s\033[0m \033[1;31m[found]\033\n", maxStringLength, pluginName)
		return version, versionFound
	}

	if content, ok := result.(string); ok {
		versionData, found := pluginVersion.GetPluginVersion(content)
		if found {
			version = versionData
			versionFound = true
		}
	}

	return version, versionFound
}

func CheckPluginsInNormalMode(url string, pluginNameList []string, randomUserAgent bool, rateLimit int) []types.PluginData {
	pluginsPrefix := "wp-content/plugins"
	pluginsFoundOnTarget := []types.PluginData{}
	maxStringLength := store.MaxStringLength
	pluginNameListLength := len(pluginNameList)

	for index, pluginName := range pluginNameList {

		if rateLimit > 0 {
			// Introduce a rate limit between 0 and X seconds
			time.Sleep(time.Duration(rand.Intn(rateLimit)) * time.Second)
		}

		pluginIndexPhpFileExists, err := utils.DoesRemoteFileExist(fmt.Sprintf("%s/%s/%s/%s.php", url, pluginsPrefix, pluginName, pluginName), randomUserAgent)

		if err != nil {
			fmt.Println("\033[1;31mError checking Plugin: "+pluginName+"\033[0m\n", err)
			continue
		}

		if !pluginIndexPhpFileExists {
			fmt.Printf("\033[K\033[1;34m%-*s\033[0m \033[1;34m[ok][%d/%d]\033[0m\r", maxStringLength, pluginName, index+1, pluginNameListLength)
			continue
		}

		versionData, found := ReturnPluginVersion(url, pluginsPrefix, pluginName, randomUserAgent, maxStringLength)
		if found {
			fmt.Printf("\033[1;31m%-*s\033[0m \033[1;31m[found][%s]\033\n", maxStringLength, pluginName, versionData.Number)
		} else {
			fmt.Printf("\033[1;31m%-*s\033[0m \033[1;31m[found]\033\n", maxStringLength, pluginName)
		}

		pluginsFoundOnTarget = append(pluginsFoundOnTarget, types.PluginData{Name: pluginName, Version: versionData.Number, Found: true})
	}

	return pluginsFoundOnTarget
}