package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"wp-wingman/pluginVersion"
	"wp-wingman/wordpressFinder"
)

type PluginInfo struct {
	Info struct {
		Pages int `json:"pages"`
	} `json:"info"`
	Plugins []struct {
		Slug string `json:"slug"`
	} `json:"plugins"`
}

type PluginData struct {
	Name    string
	Version string
}

var (
	wpURL           string
	rValue          string
	tFlagArgument   string
	overdriveActive bool
	savePlaybook    bool
	saveResult      bool
)
var rateLimit int = 1
var maxStringLength int
var currentPluginInCheckIndex int = 0
var pluginNameListLength int = 0
var targetPluginTag string = "security"

func helpMenu() {
	fmt.Println("\033[1;33mArguments:\n\t\033[1;31mrequired:\033[1;33m -u\t\t\twordpress url\033[1;33m\n\t\033[1;34moptional:\033[1;33m -t\t\t\twordpress plugin tag (default securtiy)\t\t\t\n\t\033[1;34moptional:\033[1;33m -r\t\t\trate limit on target (default 0-1s)\n\t\033[1;34moptional:\033[1;33m --overdrive\t\tcheck all public plugins on target (very aggressiv)\n\t\033[1;34moptional:\033[1;33m --save-playbook\tsave collected plugins in file\n\t\033[1;34moptional:\033[1;33m --save-result\t\tsave plugins found on target in file\n\t\033[1;33m")
	fmt.Println("Send over Wingman:\n./scan.sh -u www.example.com -r 5 -t newsletter \033[1;32m")
}

func fetchReadme(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return string(body), nil
}

func init() {
	flag.StringVar(&wpURL, "u", "", "wordpress url")
	flag.StringVar(&rValue, "r", "", "rate limit on target (default 0-1s)")
	flag.StringVar(&tFlagArgument, "t", "", "wordpress plugin tag (default securtiy but read the docs)")
	flag.BoolVar(&overdriveActive, "overdrive", false, "check all public plugins on target (very aggressiv)")
	flag.BoolVar(&savePlaybook, "save-playbook", false, "save collected plugins in file")
	flag.BoolVar(&saveResult, "save-result", false, "save plugins found on target in file")

	flag.Parse()
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

func main() {
	printLogo()

	flag.Parse()

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

	if wpURL != "" {
		StartWingmanJob()
	}
}

func CheckIfSaveFileExists() bool {
	fileName := "wp-wingman-" + targetPluginTag + ".txt"
	if targetPluginTag == "all" {
		fileName = "wp-wingman-all.txt"
	}

	dir, _ := os.Getwd()
	filePath := filepath.Join(dir, fileName)

	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func StartWingmanJob() {
	result := wordpressFinder.IsWordpressSite(wpURL)

	if result == false {
		fmt.Println("\033[1;31mThe URL is not a WordPress site.\033[0m")
		fmt.Println("\033[1;31m" + wpURL + "\033[0m")
		os.Exit(0)
	}

	pluginNameList := []string{}
	fmt.Println("\033[1;32mWordPress site detected: " + wpURL + "\033[0m")

	if CheckIfSaveFileExists() {
		fmt.Println("\033[1;33mSave file found - should i use it? (y/n)\033[0m")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		if answer == "y\n" {
			pluginNameList = LoadPluginNamesFromSaveFile()
		} else {
			pluginNameList = FetchPluginsByTag()
		}
	} else {
		pluginNameList = FetchPluginsByTag()
	}

	fmt.Println("\033[1;33mDo you want me to start? (y/n)\033[0m")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	if answer != "y\n" {
		fmt.Println("\033[1;32mPuuh, okey bye.\033[0m")
		os.Exit(0)
	}

	CheckPluginsAvailability(wpURL, pluginNameList)
}

func CheckPluginsAvailability(url string, pluginNameList []string) {
	fmt.Println("\n\033[1;33m[+] Let me check this for you:\n\033")
	pluginsPrefix := "wp-content/plugins"
	pluginSuffix := "readme.txt"
	pluginsFoundOnTarget := []PluginData{}
	for _, pluginName := range pluginNameList {
		result, err := fetchReadme(fmt.Sprintf("%s/%s/%s/%s", url, pluginsPrefix, pluginName, pluginSuffix))

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
			pluginsFoundOnTarget = append(pluginsFoundOnTarget, PluginData{pluginName, versionData.Number})
		} else {
			fmt.Printf("\033[K\033[1;34m%-*s\033[0m \033[1;34m[ok][%d/%d]\033[0m\r", maxStringLength, pluginName, currentPluginInCheckIndex+1, pluginNameListLength)
		}

		currentPluginInCheckIndex++
		// Introduce a rate limit between 0 and X seconds
		// Only when not in OVERDRIVE!!!
		if !overdriveActive {
			time.Sleep(time.Duration(rand.Intn(rateLimit)) * time.Second)
		}
	}

	PrintResult(pluginsFoundOnTarget)

	if saveResult {
		SaveResultToFile(pluginsFoundOnTarget)
	}

	os.Exit(0)
}

func PrintResult(pluginsFoundOnTarget []PluginData) {
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

func SaveResultToFile(pluginsFoundOnTarget []PluginData) {
	fmt.Println("\033[1;33mSaving Result...\033[0m")
	timestamp := time.Now().Format("20060102150405")
	fileName := strings.Split(strings.Split(wpURL, "//")[1], "/")[0]

	file, err := os.OpenFile(fmt.Sprintf("wp-wingman-%s-%s.txt", fileName, timestamp), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, pluginData := range pluginsFoundOnTarget {
		if _, err := file.WriteString(pluginData.Name + " " + pluginData.Version + "\n"); err != nil {
			panic(err)
		}
	}

	fmt.Println("\033[1;32mDone. Have a great day!\033[0m")
}

func LoadPluginNamesFromSaveFile() []string {
	fmt.Println("\033[1;33mLoading Playbook from save file...\033[0m")
	fileName := fmt.Sprintf("wp-wingman-%s.txt", targetPluginTag)
	if targetPluginTag == "all" {
		fileName = "wp-wingman-all.txt"
	}

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	pluginNameList := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pluginNameList = append(pluginNameList, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	pluginNameListLength = len(pluginNameList)
	for _, pluginName := range pluginNameList {
		if len(pluginName) > maxStringLength {
			maxStringLength = len(pluginName)
		}
	}

	if overdriveActive {
		fmt.Printf("\033[1;31mDone. %d found!!!\n\033[0m", pluginNameListLength)
	} else {
		fmt.Printf("\033[1;32mDone. %d found.\n\033[0m", pluginNameListLength)
	}

	return pluginNameList
}

func SavePlaybookToFile(pluginNameList []string) {
	fmt.Println("\033[1;33mSaving Playbook...\033[0m")
	fileName := fmt.Sprintf("wp-wingman-%s.txt", targetPluginTag)
	if targetPluginTag == "all" {
		fileName = "wp-wingman-all.txt"
	}

	// Remove the existing file if it exists
	if _, err := os.Stat(fileName); err == nil {
		os.Remove(fileName)
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, str := range pluginNameList {
		if _, err := file.WriteString(str + "\n"); err != nil {
			panic(err)
		}
	}
}

func FetchPluginsByTag() []string {
	fmt.Printf("\033[K\033[1;33m%-*s\033[0m \033[1;33m\033[0m\r", maxStringLength, "Updating PlayBook...")

	targetAPIEndpoint := "https://api.wordpress.org/plugins/info/1.2/?action=query_plugins&request[tag]=" + targetPluginTag

	if targetPluginTag == "all" {
		targetAPIEndpoint = "https://api.wordpress.org/plugins/info/1.2/?action=query_plugins&request[browse]"
	}

	response := makeHTTPRequest(targetAPIEndpoint)

	page := 2
	totalPages := response.Info.Pages
	pluginNameList := make([]string, 0)

	for _, plugin := range response.Plugins {
		pluginNameList = append(pluginNameList, plugin.Slug)
	}

	for page <= totalPages {
		response = makeHTTPRequest(targetAPIEndpoint + "&request[page]=" + strconv.Itoa(page))
		for _, plugin := range response.Plugins {
			pluginNameList = append(pluginNameList, plugin.Slug)
		}
		fmt.Printf("\033[K\033[1;33m%-*s\033[0m \033[1;33m[%d/%d]\033[0m\r", maxStringLength, "Updating PlayBook...", page, totalPages)

		page++
	}

	pluginNameListLength = len(pluginNameList)
	for _, pluginName := range pluginNameList {
		if len(pluginName) > maxStringLength {
			maxStringLength = len(pluginName)
		}
	}

	if savePlaybook == true {
		SavePlaybookToFile(pluginNameList)
	}

	if overdriveActive == true {
		fmt.Println("\n\033[1;32mDone.\033[0m", "\033[1;31m", pluginNameListLength, "found!\033[0m")
	} else {
		fmt.Println("\n\033[1;32mDone.\033[0m", "\033[1;32m", pluginNameListLength, "found.\033[0m")
	}

	// Display the plugin names
	//fmt.Println("Plugin Names:")
	//for _, name := range pluginNameList {
	//	fmt.Println(name)
	//}

	return pluginNameList
}

func makeHTTPRequest(url string) PluginInfo {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var pluginInfo PluginInfo
	err = json.Unmarshal(body, &pluginInfo)
	if err != nil {
		panic(err)
	}

	return pluginInfo
}
