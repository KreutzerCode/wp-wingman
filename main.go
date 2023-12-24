package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"wp-wingman/fileManager"
	"wp-wingman/overdrive"
	"wp-wingman/pluginSlugLoader"
	"wp-wingman/pluginVersion"
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
var maxStringLength int
var currentPluginInCheckIndex int = 0
var pluginNameListLength int = 0
var targetPluginTag string = "security"

func init() {
	flag.StringVar(&wpURL, "u", "", "wordpress url")
	flag.StringVar(&rValue, "r", "", "rate limit on target (default 0-1s)")
	flag.StringVar(&tFlagArgument, "t", "", "wordpress plugin tag (default securtiy but read the docs)")
	flag.IntVar(&wFlagArgument, "w", 10, "number of workers to execute playbook (only available in overdrive mode) (default 10)")
	flag.BoolVar(&overdriveActive, "overdrive", false, "executes playbook with the boys (very aggressiv)")
	flag.BoolVar(&savePlaybook, "save-playbook", false, "save collected plugins in file")
	flag.BoolVar(&saveResult, "save-result", false, "save plugins found on target in file")

	flag.Parse()
}

func main() {
	printLogo()

	if wpURL == "" {
		helpMenu()
		fmt.Println("\n\033[1;31mError: Missing -u argument\033[0m")
		os.Exit(0)
	}

	if rValue != "" && !overdriveActive {
		// set global variable named rate limit to the value provided
		rateLimit, _ = strconv.Atoi(rValue)
		fmt.Printf("\033[1;32mSet rate limit to: %s\033[0m\n", rValue)
	}

	if tFlagArgument != "" {
		// set global variable named target plugin tag to the value provided
		targetPluginTag = tFlagArgument
		fmt.Printf("\033[1;32mSet plugin tag to: %s\033[0m\n", targetPluginTag)
	}

	if wFlagArgument != 10 && overdriveActive {
		fmt.Printf("\033[1;32mSet number of workers to: %d\033[0m\n", wFlagArgument)
	}

	StartWingmanJob()
}

func helpMenu() {
	fmt.Println("\033[1;33mArguments:\n\t\033[1;31mrequired:\033[1;33m -u\t\t\twordpress url\033[1;33m\n\t\033[1;34moptional:\033[1;33m -t\t\t\twordpress plugin tag (default securtiy but read the docs)\t\t\t\n\t\033[1;34moptional:\033[1;33m -r\t\t\trate limit on target (default 0-1s)\n\t\033[1;34moptional:\033[1;33m -w\t\t\tnumber of workers to execute playbook (only available in overdrive mode)\n\t\033[1;34moptional:\033[1;33m --overdrive\t\texecutes playbook with the boys (very aggressiv)\n\t\033[1;34moptional:\033[1;33m --save-playbook\tsave collected plugins in file\n\t\033[1;34moptional:\033[1;33m --save-result\t\tsave plugins found on target in file\n\t\033[1;33m")
	fmt.Println("Send over Wingman:\n./scan.sh -u www.example.com -r 5 -t newsletter \033[1;32m")
}

func printLogo() {
	fmt.Println("\033[1;31m" +
		"__        ______   __        _____ _   _  ____ __  __    _    _   _ \n" +
		"\\ \\      / /  _ \\  \\ \\      / /_ _| \\ | |/ ___|  \\/  |  / \\  | \\ | |\n" +
		" \\ \\ /\\ / /| |_) |  \\ \\ /\\ / / | ||  \\| | |  _| |\\/| | / _ \\ |  \\| |\n" +
		"  \\ V  V / |  __/    \\ V  V /  | || |\\  | |_| | |  | |/ ___ \\| |\\  |\n" +
		"   \\_/\\_/  |_|        \\_/\\_/  |___|_| \\_|\\____|_|  |_/_/   \\_\\_| \\_|\n\n" +
		"\033[1;34m \t\t\t @kreutzercode \n" +
		"\033[0m")
}

func printResult(pluginsFoundOnTarget []types.PluginData) {
	fmt.Println("\n\n\n\033[1;32mDone.\n\033[0m")
	fmt.Println("\033[1;32mSummary:\n\033[0m")

	if len(pluginsFoundOnTarget) != 0 {
		for _, pluginData := range pluginsFoundOnTarget {
			fmt.Printf("\033[1;31m%-*s\033[0m \033[1;31m[found][%s]\033[0m\n", maxStringLength, pluginData.Name, pluginData.Version)
		}

		fmt.Println("\n\033[1;32mThese are my findings. Good luck sir!\033[0m")
	} else {
		fmt.Println("\033[1;32mNothing found. Good luck.\033[0m")
	}
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
	result := wordpressFinder.IsWordpressSite(wpURL)

	if result == false {
		fmt.Println("\033[1;31mThe URL is not a WordPress site.\033[0m")
		fmt.Println("\033[1;31m" + wpURL + "\033[0m")
		os.Exit(0)
	}

	fmt.Println("\033[1;32mWordPress site detected: " + wpURL + "\033[0m")

	pluginNameList := getPluginSlugList()
	pluginNameListLength = len(pluginNameList)
	maxStringLength = determineMaxStringLength(pluginNameList)

	if savePlaybook == true {
		fileName := fmt.Sprintf("wp-wingman-%s.txt", targetPluginTag)
		fileManager.SavePlaybookToFile(pluginNameList, fileName)
	}

	fmt.Println("\033[1;33mDo you want me to start? (y/n)\033[0m")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	if answer != "y\n" {
		fmt.Println("\033[1;32mPuuh, okey bye.\033[0m")
		os.Exit(0)
	}

	pluginsFoundOnTarget := checkPluginsAvailability(wpURL, pluginNameList)

	printResult(pluginsFoundOnTarget)

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
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		if answer == "y\n" {
			pluginSlugList = pluginSlugLoader.FetchPluginSlugsFromFile(targetPluginTag)

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

	pluginsFoundOnTarget := []types.PluginData{}

	if overdriveActive {
		pluginsFoundOnTarget = overdrive.CheckPluginsInOverdriveMode(url, maxStringLength, pluginNameList, wFlagArgument)
	} else {
		pluginsFoundOnTarget = checkPluginsInNormalMode(url, pluginNameList)
	}

	return pluginsFoundOnTarget
}

func checkPluginsInNormalMode(url string, pluginNameList []string) []types.PluginData {
	pluginsPrefix := "wp-content/plugins"
	pluginSuffix := "readme.txt"
	pluginsFoundOnTarget := []types.PluginData{}
	for _, pluginName := range pluginNameList {
		result, err := utils.FetchReadme(fmt.Sprintf("%s/%s/%s/%s", url, pluginsPrefix, pluginName, pluginSuffix))

		if err != nil {
			fmt.Println("\033[1;31mError checking Plugin: "+pluginName+"\033[0m\n", err)
			continue
		}

		if content, ok := result.(string); ok {
			versionData, found := pluginVersion.GetPluginVersion(content)
			if found {
				fmt.Printf("\033[1;31m%-*s\033[0m \033[1;31m[found][%s]\033\n", maxStringLength, pluginName, versionData.Number)
			} else {
				fmt.Printf("\033[1;31m%-*s\033[0m \033[1;31m[not found]\033\n", maxStringLength, pluginName)
			}
			pluginsFoundOnTarget = append(pluginsFoundOnTarget, types.PluginData{Name: pluginName, Version: versionData.Number, Found: true})
		} else {
			fmt.Printf("\033[K\033[1;34m%-*s\033[0m \033[1;34m[ok][%d/%d]\033[0m\r", maxStringLength, pluginName, currentPluginInCheckIndex+1, pluginNameListLength)
		}

		currentPluginInCheckIndex++
		// Introduce a rate limit between 0 and X seconds
		time.Sleep(time.Duration(rand.Intn(rateLimit)) * time.Second)
	}

	return pluginsFoundOnTarget
}
