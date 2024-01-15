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

func CheckPluginsInNormalMode(url string, pluginNameList []string, randomUserAgent bool, rateLimit int) []types.PluginData {
	pluginsPrefix := "wp-content/plugins"
	pluginSuffix := "readme.txt"
	pluginsFoundOnTarget := []types.PluginData{}
	maxStringLength := store.MaxStringLength
	pluginNameListLength := len(pluginNameList)
	currentPluginInCheckIndex := 0

	for _, pluginName := range pluginNameList {
		result, err := utils.FetchReadme(fmt.Sprintf("%s/%s/%s/%s", url, pluginsPrefix, pluginName, pluginSuffix), randomUserAgent)

		if err != nil {
			fmt.Println("\033[1;31mError checking Plugin: "+pluginName+"\033[0m\n", err)
			continue
		}

		if content, ok := result.(string); ok {
			versionData, found := pluginVersion.GetPluginVersion(content)
			if found {
				fmt.Printf("\033[1;31m%-*s\033[0m \033[1;31m[found][%s]\033\n", maxStringLength, pluginName, versionData.Number)
			} else {
				fmt.Printf("\033[1;31m%-*s\033[0m \033[1;31m[found]\033\n", maxStringLength, pluginName)
			}
			pluginsFoundOnTarget = append(pluginsFoundOnTarget, types.PluginData{Name: pluginName, Version: versionData.Number, Found: true})
		} else {
			fmt.Printf("\033[K\033[1;34m%-*s\033[0m \033[1;34m[ok][%d/%d]\033[0m\r", maxStringLength, pluginName, currentPluginInCheckIndex+1, pluginNameListLength)
		}

		currentPluginInCheckIndex++

		if rateLimit > 0 {
			// Introduce a rate limit between 0 and X seconds
			time.Sleep(time.Duration(rand.Intn(rateLimit)) * time.Second)
		}
	}

	return pluginsFoundOnTarget
}