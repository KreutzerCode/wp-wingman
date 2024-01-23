package pluginFinder

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"wp-wingman/pluginVersion"
	"wp-wingman/types"
	"wp-wingman/utils"
)

var wpURL string
var useRandomUserAgent bool
var numWorkers = 10
var detectionMode string

func checkURL(pluginName string, resultsChannel chan<- types.PluginData, pluginsFoundOnTarget *[]types.PluginData) {
    pluginsPrefix := "wp-content/plugins"
    pluginIndexPhpFileExists, err := utils.DoesRemoteFileExist(fmt.Sprintf("%s/%s/%s/%s.php", wpURL, pluginsPrefix, pluginName, pluginName), useRandomUserAgent)

    if err != nil {
        fmt.Println("\033[1;31mError checking Plugin: "+pluginName+"\033[0m\n", err)
        return
    }

    pluginData := types.PluginData{}
    pluginData.Name = pluginName

    if !pluginIndexPhpFileExists {
        resultsChannel <- pluginData
        return
    }

    pluginData.Found = true

    versionData, found := returnPluginVersion(wpURL, pluginsPrefix, pluginName, useRandomUserAgent)
    if found {
        pluginData.Version = versionData.Number
    }

    *pluginsFoundOnTarget = append(*pluginsFoundOnTarget, types.PluginData{Name: pluginName, Version: versionData.Number, DetectionMethod: detectionMode, Found: true})
    
    resultsChannel <- pluginData
}

func worker(urlsToCheck <-chan string, resultsChannel chan<- types.PluginData, pluginsFoundOnTarget *[]types.PluginData) {
    for url := range urlsToCheck {
        checkURL(url, resultsChannel, pluginsFoundOnTarget)
    }
}

func CheckPluginSlugsOnTarget(url string, pluginNameList []string, numberOfWorkers int, randomUserAgent bool, targetDetectionMode string, rateLimit int) []types.PluginData {
    pluginsFoundOnTarget := []types.PluginData{}
    urlsToCheck := pluginNameList
    numWorkers = numberOfWorkers
    listLength := len(urlsToCheck)
    var waitGroup sync.WaitGroup
    resultsChannel := make(chan types.PluginData, listLength)
    urlsToCheckChannel := make(chan string, listLength)
    maxStringLength := utils.DetermineMaxStringLength(urlsToCheck)
    wpURL = url
    useRandomUserAgent = randomUserAgent
    detectionMode = targetDetectionMode

    var ticker *time.Ticker
    if rateLimit > 0 {
        ticker = time.NewTicker(time.Duration(rand.Intn(rateLimit)) * time.Millisecond)
        defer ticker.Stop()
    }

    for i := 0; i < numWorkers; i++ {
        go worker(urlsToCheckChannel, resultsChannel, &pluginsFoundOnTarget)
    }

    for _, pluginName := range urlsToCheck {
        waitGroup.Add(1)
        if rateLimit > 0 {
            <-ticker.C
        }
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

func returnPluginVersion(url string, pluginsPrefix string, pluginName string, randomUserAgent bool) (types.VersionNumber, bool){
	version := types.VersionNumber{}
	versionFound := false
	result, err := utils.FetchReadme(fmt.Sprintf("%s/%s/%s/%s", url, pluginsPrefix, pluginName, "readme.txt"), randomUserAgent)

	if err != nil {
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