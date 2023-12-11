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
pluginNameListLength=0
max_string_length=0
currentPluginInCheckIndex=0
rateLimit=1
targetPluginTag="security"
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

function fetch_plugins_by_tag(){
    echo -e "\e[1;33mUpdating PlayBook...\e[0m"
    local api_url="https://api.wordpress.org/plugins/info/1.2/"
    local response=$(curl -g -s -A "${user_agents[RANDOM % ${#user_agents[@]}]}" "${api_url}?action=query_plugins&request[tag]=$targetPluginTag")
    
    local page=2
    local total_pages=$(echo "$response" | jq -r '.info.pages')
    local plugin_names=($(echo "$response" | jq -r '.plugins[].slug'))

    while [ "$page" -le "$total_pages" ]; do
        response=$(curl -g -s -A "${user_agents[RANDOM % ${#user_agents[@]}]}" "${api_url}?action=query_plugins&request[tag]=$targetPluginTag&request[page]=${page}")
        names_on_page=($(echo "$response" | jq -r '.plugins[].slug'))
        plugin_names+=("${names_on_page[@]}")

        ((page++))
    done

    pluginNameList=("${plugin_names[@]}")
    pluginNameListLength=${#pluginNameList[@]}
    max_string_length=$(printf "%s\n" "${pluginNameList[@]}" | awk '{ if (length > x) x = length } END { print x }')

    echo -e "\e[1;32mDone. $pluginNameListLength found.\e[0m"

    # Display the plugin names
    #echo "Plugin Names:"
    #for name in "${plugin_names[@]}"; do
    #    echo "$name"
    #done
} 

helpMenu(){
    echo -e "\e[1;33mArguments:\n\t\e[1;31mrequired:\e[1;33m -u\t\twordpress url\e[1;33m\n\t\e[1;34moptional:\e[1;33m -t\t\twordpress plugin tag (default securtiy)\t\n\t\e[1;34moptional:\e[1;33m -r\t\trate limit on target (default 0-1s)\n\t\e[1;33m"
    echo -e "Send over Wingman:\n./scan.sh -u www.example.com -r 5 -t newsletter \e[1;32m"
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

printFindings(){
    local isFound=$1
    local pluginName=$2

    if [ "$isFound" == "true" ]; then
       echo -e "\e[1;31m$(printf "%-${max_string_length}s" "$pluginName")\e[0m \e[1;31m[found]\e[0m"
       allClear=false
    else
        printf "\e[1;34m%-${max_string_length}s\e[0m \e[1;34m[ok][${currentPluginInCheckIndex} / ${pluginNameListLength}]\e[0m\r" "$pluginName"
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
        
        printFindings "$result" "$pluginName"
        ((currentPluginInCheckIndex++))


        # Introduce a ratelimit between 0 and X seconds
        # Add one to get desired value in sek 4 = 3
        sleep $(($RANDOM % $rateLimit + 1))
    done

    if [ "$allClear" == "true" ]; then
        echo -e "\n\e[1;32mNothing found, good luck.\e[0m\n"
    else
        echo -e "\n\e[1;31mFound something mate, good luck.\e[0m\n"
    fi

    exit
}

startWingmanJob(){
    local WP_URL=$1
    result=$(testUrl "$WP_URL/wp-login.php")
    if [ "$result" == "true" ]; then
        echo -e "\e[1;32mWordPress site detected: $WP_URL\e[0m"
        fetch_plugins_by_tag

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
} 

args=()

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -u)
            shift
            WP_URL="$1"
            args+=("-u" "$WP_URL")
            ;;
        -r)
            shift
            R_VALUE="$1"
            args+=("-r" "$R_VALUE")
            ;;
        -t)
            shift
            T_VALUE="$1"
            args+=("-t" "$T_VALUE")
            ;;
        -*)
            echo -e "\n\e[1;31mInvalid argument: $1\e[0m\n"
            helpMenu
            exit 1
            ;;
    esac
    shift
done

if [ ${#args[@]} -eq 0 ]; then
    helpMenu
    exit 1
fi

# Check if -u argument is missing
if [ -z "$WP_URL" ]; then
    echo -e "\n\e[1;31mError: Missing -u argument\e[0m\n"
    helpMenu
    exit 1
fi

# move the -u argument to the end so the configuration takes place first
u_index=$(printf '%s\n' "${args[@]}" | grep -n '^\-u' | cut -f1 -d:)
[ -n "$u_index" ] && args=("${args[@]:0:$u_index-1}" "${args[@]:$u_index+1}" "-u" "$WP_URL")

for ((i=0; i<${#args[@]}; i+=2)); do
    option="${args[i]}"
    value="${args[i+1]}"
    if [ "$option" == "-r" ]; then
        rateLimit="$value"
        echo -e "\e[1;32mSet rate limit to: $value\e[0m"
    fi

    if [ "$option" == "-t" ]; then
        targetPluginTag="$value"
        echo -e "\e[1;32mSet plugin tag to: $value\e[0m"
    fi

    if [ "$option" == "-u" ]; then
        startWingmanJob "$value"
    fi
done