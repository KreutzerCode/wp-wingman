package fileManager

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"wp-wingman/types"
)

func SaveResultToFile(pluginsFoundOnTarget []types.PluginData, fileName string) {
	fmt.Println("\033[1;33mSaving Result...\033[0m")
	timestamp := time.Now().Format("20060102150405")

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

func LoadPluginSlugsFromFile(fileName string) []string {
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

	return pluginNameList
}

func SavePlaybookToFile(pluginNameList []string, fileName string) {
	fmt.Println("\033[1;33mSaving Playbook...\033[0m")

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

func CheckIfSaveFileExists(fileName string) bool {
	dir, _ := os.Getwd()
	filePath := filepath.Join(dir, fileName)

	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
