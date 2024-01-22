package overdriveMode

import (
	"fmt"
	"sync"
	"wp-wingman/pluginVersion"
	"wp-wingman/store"
	"wp-wingman/types"
	"wp-wingman/utils"
)

var wpURL string
var useRandomUserAgent bool
var numWorkers = 10
var detectionMode string

func checkURL(pluginName string, resultsChannel chan<- types.PluginData, pluginsFoundOnTarget *[]types.PluginData) {
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
        *pluginsFoundOnTarget = append(*pluginsFoundOnTarget, types.PluginData{Name: pluginName, Version: versionData.Number, DetectionMethod: detectionMode, Found: true})
    }

    resultsChannel <- pluginData
}

func worker(urlsToCheck <-chan string, resultsChannel chan<- types.PluginData, pluginsFoundOnTarget *[]types.PluginData) {
    for url := range urlsToCheck {
        checkURL(url, resultsChannel, pluginsFoundOnTarget)
    }
}

func CheckPluginsInOverdriveMode(url string, pluginNameList []string, numberOfWorkers int, randomUserAgent bool, targetDetectionMode string) []types.PluginData {
    pluginsFoundOnTarget := []types.PluginData{} // Initialize it here
    urlsToCheck := pluginNameList
    numWorkers = numberOfWorkers
    listLength := len(urlsToCheck)
    var waitGroup sync.WaitGroup
    resultsChannel := make(chan types.PluginData, listLength)
    urlsToCheckChannel := make(chan string, listLength)
    maxStringLength := store.MaxStringLength
    wpURL = url
    useRandomUserAgent = randomUserAgent
    detectionMode = targetDetectionMode

    for i := 0; i < numWorkers; i++ {
        go worker(urlsToCheckChannel, resultsChannel, &pluginsFoundOnTarget)
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
                fmt.Printf("\033[1;31m%-*s\033[0m \033[1;31m[%d/%d][found][%s][%s]\033\n", maxStringLength, pluginData.Name, index, listLength, detectionMode, pluginData.Version)
            } else {
                fmt.Printf("\033[1;31m%-*s\033[0m \033[1;31m[%d/%d][found][%s]\033\n", maxStringLength, pluginData.Name, index, listLength, detectionMode)
            }
        }

        waitGroup.Done()
    }

    return pluginsFoundOnTarget
}