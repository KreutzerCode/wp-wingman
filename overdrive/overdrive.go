package overdrive

import (
	"fmt"
	"sync"
	"wp-wingman/pluginVersion"
	"wp-wingman/types"
	"wp-wingman/utils"
)

var pluginsFoundOnTarget []types.PluginData
var wpURL string
var useRandomUserAgent bool
var maxStringLength int = 0
var numWorkers = 10

func checkURL(pluginName string, resultsChannel chan<- types.PluginData) {
	pluginsPrefix := "wp-content/plugins"
	pluginSuffix := "readme.txt"
	result, err := utils.FetchReadme(fmt.Sprintf("%s/%s/%s/%s", wpURL, pluginsPrefix, pluginName, pluginSuffix), useRandomUserAgent)

	if err != nil {
		fmt.Println("\033[1;31mError checking Plugin: "+pluginName+"\033[0m\n", err)
		return
	}

	pluginData := types.PluginData{}
	pluginData.Name = pluginName

	if content, ok := result.(string); ok {
		pluginData.Found = true
		versionData, found := pluginVersion.GetPluginVersion(content)
		if found {
			pluginData.Version = versionData.Number
		}
		pluginsFoundOnTarget = append(pluginsFoundOnTarget, types.PluginData{Name: pluginName, Version: versionData.Number, Found: true})
	}

	resultsChannel <- pluginData
}

func worker(urlsToCheck <-chan string, resultsChannel chan<- types.PluginData) {
	for url := range urlsToCheck {
		checkURL(url, resultsChannel)
	}
}

func CheckPluginsInOverdriveMode(url string, maxPluginNameLength int, pluginNameList []string, numberOfWorkers int, randomUserAgent bool) []types.PluginData {
	urlsToCheck := pluginNameList
	numWorkers = numberOfWorkers
	listLength := len(urlsToCheck)
	var waitGroup sync.WaitGroup
	resultsChannel := make(chan types.PluginData, listLength)
	urlsToCheckChannel := make(chan string, listLength)
	maxStringLength = maxPluginNameLength
	wpURL = url
	useRandomUserAgent = randomUserAgent

	for i := 0; i < numWorkers; i++ {
		go worker(urlsToCheckChannel, resultsChannel)
	}

	for _, pluginName := range urlsToCheck {
		waitGroup.Add(1)
		urlsToCheckChannel <- pluginName
	}

	close(urlsToCheckChannel)

	go func() {
		waitGroup.Wait()
		close(resultsChannel)
	}()

	index := 0
	for pluginData := range resultsChannel {
		index++
		if !pluginData.Found {
			fmt.Printf("\033[K\033[1;34m%-*s\033[0m \033[1;34m[%d/%d][ok]\033[0m\r", maxStringLength, pluginData.Name, index, listLength)
		} else {
			if pluginData.Version != "" {
				fmt.Printf("\033[1;31m%-*s\033[0m \033[1;31m[%d/%d][found][%s]\033\n", maxStringLength, pluginData.Name, index, listLength, pluginData.Version)
			} else {
				fmt.Printf("\033[1;31m%-*s\033[0m \033[1;31m[%d/%d][found]\033\n", maxStringLength, pluginData.Name, index, listLength)
			}
		}

		waitGroup.Done()
	}

	return pluginsFoundOnTarget
}
