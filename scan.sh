#!/bin/bash
# kreutzercode

# Print Banner.
echo -e "\e[1;31m
__        ______   __        _____ _   _  ____ __  __    _    _   _ 
\ \      / /  _ \  \ \      / /_ _| \ | |/ ___|  \/  |  / \  | \ | |
 \ \ /\ / /| |_) |  \ \ /\ / / | ||  \| | |  _| |\/| | / _ \ |  \| |
  \ V  V / |  __/    \ V  V /  | || |\  | |_| | |  | |/ ___ \| |\  |
   \_/\_/  |_|        \_/\_/  |___|_| \_|\____|_|  |_/_/   \_\_| \_|
                                                        
                            \e[1;34m  @kreutzercode 
"

pluginNameList=()

user_agents=("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
            "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
            "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
            "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:57.0) Gecko/20100101 Firefox/57.0"
            "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:57.0) Gecko/20100101 Firefox/57.0"
            "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:57.0) Gecko/20100101 Firefox/57.0"
            "Mozilla/5.0 (X11; Linux x86_64; rv:57.0) Gecko/20100101 Firefox/57.0"
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64; Trident/7.0; AS; rv:11.0) like Gecko"
            "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko"
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/80.0.361.109"
            "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/80.0.361.109"
            "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Edge/80.0.361.109"
            "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/80.0.361.109")

function fetch_security_plugins(){
    echo -e "\e[1;33mUpdating PlayBook...\e[0m"
    local api_url="https://api.wordpress.org/plugins/info/1.2/"
    local response=$(curl -g -s -A "${user_agents[RANDOM % ${#user_agents[@]}]}" "${api_url}?action=query_plugins&request[tag]=security")
    
    local page=2
    local total_pages=$(echo "$response" | jq -r '.info.pages')
    local plugin_names=($(echo "$response" | jq -r '.plugins[].slug'))

    while [ "$page" -le "$total_pages" ]; do
        response=$(curl -g -s -A "${user_agents[RANDOM % ${#user_agents[@]}]}" "${api_url}?action=query_plugins&request[tag]=security&request[page]=${page}")
        names_on_page=($(echo "$response" | jq -r '.plugins[].slug'))
        plugin_names+=("${names_on_page[@]}")

        ((page++))
    done

    pluginNameList=("${plugin_names[@]}")
    array_length=${#pluginNameList[@]}
    echo -e "\e[1;32mDone. $array_length found.\e[0m"

    # Display the security-related plugin names
    #echo "Security-Related Plugin Names:"
    #for name in "${plugin_names[@]}"; do
    #    echo "$name"
    #done
} 

# Help
helpMenu(){
    echo -e "\e[1;33mArguments:\n\t-u\t\twordpress url\n\t\n"
    echo -e "Send over Wingman:\n./scan.sh -u www.example.com\n \e[1;32m"
}

testUrl() {
    local url=$1 
    CHECK_URL=$(curl -s -A "${user_agents[RANDOM % ${#user_agents[@]}]}" -o /dev/null --head --write-out '%{http_code}' "$url")
    if [ "$CHECK_URL" -eq 200 ]; then
        echo "true"
    else
        echo "false"
    fi
}

guardEnum(){
    local url=$1
    local allClear=true
    echo -e "\n\e[1;33m[+] Let me check this for you:\e[0m\n"
    pluginsPrefix="wp-content/plugins"
    pluginSuffix="readme.txt"

    for pluginName in "${pluginNameList[@]}"; do
        result=$(testUrl "$url/$pluginsPrefix/$pluginName/$pluginSuffix")
        if [ "$result" == "true" ]; then
            echo -e "\e[1;31m$pluginName\e[0m"
            allClear=false
        else
            #echo -ne "\e[1;34m$pluginName\e[0m \033[0K\r"
            echo -e "\e[1;34m$pluginName\e[0m"
        fi

        # Introduce a random delay between 0 and 3 seconds TODO add ass script argument <3
        sleep $(($RANDOM % 4))
    done

    if [ "$allClear" == "true" ]; then
        echo -e "\n\e[1;32mLooking good, have fun.\e[0m\n"
    else
        echo -e "\n\e[1;31mNot locking to good mate, take care.\e[0m\n"
    fi

    exit
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -u)
            shift
            WP_URL="$1"
            result=$(testUrl "$WP_URL/wp-login.php")
            if [ "$result" == "true" ]; then
                echo -e "\e[1;32mWordPress site detected: $WP_URL\e[0m"
                fetch_security_plugins

                echo -e "\e[1;33mDo you want me to start? (y/n)\e[0m"
                read answer
                if [ "$answer" != "y" ]; then
                    echo -e "\e[1;32mPuuh, okey bye.\e[0m\n"
                    exit
                fi

                guardEnum $WP_URL
            else
                echo -e "\e[1;31mThe URL is not a WordPress site.\e[0m"
                echo -e "\e[1;31m$WP_URL\e[0m"
                exit 1
            fi
            ;;
        -*)
            echo -e "\n\e[1;31mInvalid argument: $1\e[0m\n"
            ;;
    esac
    shift
done

# Check if URL is provided
if [ -z "$WP_URL" ]; then
    helpMenu
fi