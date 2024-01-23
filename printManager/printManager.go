package printManager

import (
	"fmt"
	"wp-wingman/types"
)

func HelpMenu() {
    fmt.Println("\033[1;33mArguments:")
    fmt.Println("\t\033[1;31mrequired:\033[1;33m -u\t\t\twordpress url")
    fmt.Println("\t\033[1;34moptional:\033[1;33m -t\t\t\twordpress plugin tag (default security but read the docs)")
    fmt.Println("\t\033[1;34moptional:\033[1;33m -r\t\t\trate limit on target (default 0s)")
    fmt.Println("\t\033[1;34moptional:\033[1;33m -w\t\t\tnumber of workers to execute playbook (default 10)")
    fmt.Println("\t\033[1;34moptional:\033[1;33m --save-playbook\tsave collected plugins in file")
    fmt.Println("\t\033[1;34moptional:\033[1;33m --save-result\t\tsave plugins found on target in file")
    fmt.Println("\t\033[1;34moptional:\033[1;33m --user-agent\t\tuse random user agent for every request")
    fmt.Println("\nSend over Wingman:")
    fmt.Println("./wp-wingman -u www.example.com -r 5 -t newsletter \033[1;32m")
}

func PrintLogo() {
	fmt.Println("\033[1;31m" +
		"__        ______   __        _____ _   _  ____ __  __    _    _   _ \n" +
		"\\ \\      / /  _ \\  \\ \\      / /_ _| \\ | |/ ___|  \\/  |  / \\  | \\ | |\n" +
		" \\ \\ /\\ / /| |_) |  \\ \\ /\\ / / | ||  \\| | |  _| |\\/| | / _ \\ |  \\| |\n" +
		"  \\ V  V / |  __/    \\ V  V /  | || |\\  | |_| | |  | |/ ___ \\| |\\  |\n" +
		"   \\_/\\_/  |_|        \\_/\\_/  |___|_| \\_|\\____|_|  |_/_/   \\_\\_| \\_|\n\n" +
		"\033[1;34m \t\t\t @kreutzercode \n" +
		"\033[0m")
}

func PrintResult(pluginsFoundOnTarget []types.PluginData, maxNameLength int) {
	fmt.Println("\n\n\n\033[1;32mDone.\n\033[0m")
	fmt.Println("\033[1;32mSummary:\n\033[0m")

	if len(pluginsFoundOnTarget) != 0 {
		for _, pluginData := range pluginsFoundOnTarget {
			if len(pluginData.Version) == 0 {
				fmt.Printf("\033[1;31m%-*s\033[0m \033[1;31m[%s]\033[0m\n", maxNameLength, pluginData.Name, pluginData.DetectionMethod)
			} else {
				fmt.Printf("\033[1;31m%-*s\033[0m \033[1;31m[%s][%s]\033[0m\n", maxNameLength, pluginData.Name, pluginData.DetectionMethod, pluginData.Version)
			}
		}

		fmt.Println("\n\033[1;32mThese are my findings. Good luck sir!\033[0m")
	} else {
		fmt.Println("\033[1;32mNothing found. Good luck.\033[0m")
	}
}