package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"wp-wingman/fileManager"
	"wp-wingman/passivDetector"
	"wp-wingman/pluginFinder/normalMode"
	"wp-wingman/pluginFinder/overdriveMode"
	"wp-wingman/pluginSlugLoader"
	"wp-wingman/printManager"
	"wp-wingman/store"
	"wp-wingman/types"
	"wp-wingman/utils"
	"wp-wingman/wordpressFinder"
)

var (
	wpURL           string
	rValue          string
	tFlagArgument   string
	wFlagArgument   int
	overdriveActive bool
	savePlaybook    bool
	saveResult      bool
)
var rateLimit int = 1
var workerCount int = 10
var targetPluginTag string = "security"
var useRandomUserAgent bool
var usingPlaybookFromFile bool = false

func init() {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
    flagSet.SetOutput(io.Discard) // Suppress default error messages

    flagSet.StringVar(&wpURL, "u", "", "wordpress url")
    flagSet.StringVar(&tFlagArgument, "t", "", "wordpress plugin tag (default securtiy but read the docs)")
    flagSet.StringVar(&rValue, "r", "", "rate limit on target (default 0-1s)")
    flagSet.IntVar(&wFlagArgument, "w", 10, "number of workers to execute playbook (only available in overdrive mode) (default 10)")
    flagSet.BoolVar(&overdriveActive, "overdrive", false, "executes playbook with the boys (very aggressiv)")
    flagSet.BoolVar(&savePlaybook, "save-playbook", false, "save collected plugins in file")
    flagSet.BoolVar(&saveResult, "save-result", false, "save plugins found on target in file")
    flagSet.BoolVar(&useRandomUserAgent, "user-agent", false, "use random user agent for every request")

    flagSet.Usage = func() {
		printManager.PrintLogo()
		printManager.HelpMenu()
		fmt.Printf("\n\033[1;31mError: check input for invalid arguments\033[0m\n")
        os.Exit(0)
    }

    flagSet.Parse(os.Args[1:])
}

func main() {
	printManager.PrintLogo()

	if wpURL == "" {
		printManager.HelpMenu()
		fmt.Println("\n\033[1;31mError: Missing -u argument\033[0m")
		os.Exit(0)
	}

	if rValue != "" && !overdriveActive {
		rateLimit, _ = strconv.Atoi(rValue)
		fmt.Printf("\033[1;32mSet rate limit to: %s\033[0m\n", rValue)
	}

	if tFlagArgument != "" {
		targetPluginTag = tFlagArgument
		fmt.Printf("\033[1;32mSet plugin tag to: %s\033[0m\n", targetPluginTag)
	}

	if overdriveActive {
		workerCount = wFlagArgument
		if workerCount < 2 {
			workerCount = 2
		}
		
		fmt.Printf("\033[1;32mSet number of workers to: %d\033[0m\n", workerCount)
	}

	if useRandomUserAgent {
		fmt.Printf("\033[1;32mUse random user agent: %s\033[0m\n", strconv.FormatBool(useRandomUserAgent))
	}

	StartWingmanJob()
}

func determineMaxStringLength(list []string) int {
	maxStringLength := 0
	for _, pluginName := range list {
		if len(pluginName) > maxStringLength {
			maxStringLength = len(pluginName)
		}
	}

	return maxStringLength
}

func StartWingmanJob() {
	result := wordpressFinder.IsWordpressSite(wpURL, useRandomUserAgent)

	if !result {
		fmt.Println("\033[1;31mThe URL is not a WordPress site.\033[0m")
		fmt.Println("\033[1;31m" + wpURL + "\033[0m")
		os.Exit(0)
	}

	fmt.Println("\033[1;32mWordPress site detected: " + wpURL + "\033[0m")

	pluginNameList := getPluginSlugList()
	store.MaxStringLength = determineMaxStringLength(pluginNameList)

	if savePlaybook && !usingPlaybookFromFile {
		fileName := fmt.Sprintf("wp-wingman-%s.txt", targetPluginTag)
		fileManager.SavePlaybookToFile(pluginNameList, fileName)
	}

	fmt.Println("\033[1;33mDo you want me to start? (y/n)\033[0m")
	if !utils.GetUserInputYesNo() {
		fmt.Println("\033[1;32mPuuh, okey bye.\033[0m")
		os.Exit(0)
	}

	var pluginsFoundOnTarget = checkPluginsAvailability(wpURL, pluginNameList)

	fmt.Println("\n\n\033[1;33mCkeck additional plugins via content? (passiv-detection*) (y/n)\033[0m")
	if utils.GetUserInputYesNo() {
		missingPlugins := passivDetector.FindPluginsInContent(wpURL, pluginsFoundOnTarget, useRandomUserAgent, rateLimit)
		// Append the missing plugins to the pluginsFoundOnTarget slice
		pluginsFoundOnTarget = append(pluginsFoundOnTarget, missingPlugins...)
	}

	printManager.PrintResult(pluginsFoundOnTarget)

	if saveResult {
		fileName := strings.Split(strings.Split(wpURL, "//")[1], "/")[0]
		fileManager.SaveResultToFile(pluginsFoundOnTarget, fileName)
	}

	os.Exit(0)
}

func getPluginSlugList() []string {
	fileName := "wp-wingman-" + targetPluginTag + ".txt"
	pluginSlugList := []string{}
	if fileManager.CheckIfSaveFileExists(fileName) {
		fmt.Println("\033[1;33mSave file found - should i use it? (y/n)\033[0m")
		if utils.GetUserInputYesNo() {
			pluginSlugList = pluginSlugLoader.FetchPluginSlugsFromFile(targetPluginTag)
			usingPlaybookFromFile = true
		} else {
			pluginSlugList = pluginSlugLoader.FetchPluginSlugsFromAPI(targetPluginTag)
		}
	} else {
		pluginSlugList = pluginSlugLoader.FetchPluginSlugsFromAPI(targetPluginTag)
	}

	return pluginSlugList
}

func checkPluginsAvailability(url string, pluginNameList []string) []types.PluginData {
	fmt.Println("\n\033[1;33m[+] Let me check this for you:\n\033")

	var pluginsFoundOnTarget = []types.PluginData{}

	if overdriveActive {
		pluginsFoundOnTarget = overdriveMode.CheckPluginsInOverdriveMode(url, pluginNameList, workerCount, useRandomUserAgent)
	} else {
		pluginsFoundOnTarget = normalMode.CheckPluginsInNormalMode(url, pluginNameList, useRandomUserAgent, rateLimit, "aggressive")
	}

	return pluginsFoundOnTarget
}
