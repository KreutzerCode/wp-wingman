package pluginSlugLoader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"wp-wingman/fileManager"
	"wp-wingman/types"
)

func FetchPluginSlugsFromAPI(targetPluginTag string, overdriveActive bool) []string {
	fmt.Printf("\033[K\033[1;33m%s\033[0m \033[1;33m\033[0m\r", "Updating PlayBook...")

	targetAPIEndpoint := "https://api.wordpress.org/plugins/info/1.2/?action=query_plugins&request[tag]=" + targetPluginTag

	if targetPluginTag == "all" {
		targetAPIEndpoint = "https://api.wordpress.org/plugins/info/1.2/?action=query_plugins&request[browse]"
	}

	pluginNameList := ReturnPluginSlugsFromAPI(targetAPIEndpoint)

	if overdriveActive == true {
		fmt.Println("\n\033[1;32mDone.\033[0m", "\033[1;31m", len(pluginNameList), "found!\033[0m")
	} else {
		fmt.Println("\n\033[1;32mDone.\033[0m", "\033[1;32m", len(pluginNameList), "found.\033[0m")
	}

	// Display the plugin names
	//fmt.Println("Plugin Names:")
	//for _, name := range pluginNameList {
	//	fmt.Println(name)
	//}

	return pluginNameList
}

func ReturnPluginSlugsFromAPI(targetAPIEndpoint string) []string {
	response := fetchWordpressApi(targetAPIEndpoint)

	page := 2
	totalPages := response.Info.Pages
	pluginNameList := make([]string, 0)

	for _, plugin := range response.Plugins {
		pluginNameList = append(pluginNameList, plugin.Slug)
	}

	for page <= totalPages {
		response = fetchWordpressApi(targetAPIEndpoint + "&request[page]=" + strconv.Itoa(page))
		for _, plugin := range response.Plugins {
			pluginNameList = append(pluginNameList, plugin.Slug)
		}
		fmt.Printf("\033[K\033[1;33m%s\033[0m \033[1;33m[%d/%d]\033[0m\r", "Updating PlayBook...", page, totalPages)

		page++
	}

	return pluginNameList
}

func fetchWordpressApi(url string) types.PluginInfo {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var pluginInfo types.PluginInfo
	err = json.Unmarshal(body, &pluginInfo)
	if err != nil {
		panic(err)
	}

	return pluginInfo
}

func FetchPluginSlugsFromFile(targetPluginTag string, overdriveActive bool) []string {
	fmt.Println("\033[1;33mLoading Playbook from save file...\033[0m")
	fileName := fmt.Sprintf("wp-wingman-%s.txt", targetPluginTag)

	pluginNameList := fileManager.LoadPluginSlugsFromFile(fileName)
	pluginNameListLength := len(pluginNameList)

	if overdriveActive {
		fmt.Printf("\033[1;31mDone. %d found!!!\n\033[0m", pluginNameListLength)
	} else {
		fmt.Printf("\033[1;32mDone. %d found.\n\033[0m", pluginNameListLength)
	}

	return pluginNameList
}
