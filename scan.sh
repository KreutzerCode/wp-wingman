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

function fetch_security_plugins() {
    echo -e "\e[1;33mUpdating PlayBook...\e[0m"
    local webpage_url="https://wordpress.org/plugins/tags/security/"
    
    local html_content=$(curl -s "$webpage_url")
    local lastPluginPageUrl=$(echo "$html_content" | grep -o '<a[^>]*class="page-numbers"[^>]*>[^<]*<\/a>' | awk 'NR==2 { match($0, /href="([^"]*)"/); url = substr($0, RSTART+6, RLENGTH-7); gsub(/\/$/, "", url); print url }')
    local lastPluginPageNumber="${lastPluginPageUrl##*/}"
    mapfile -t firstPagePluginNames <<< "$(echo "$html_content" | grep -o '<h3[^>]*class="entry-title"[^>]*>.*<\/h3>' | sed -n -e 's/.*<a[^>]*href="\([^"]*\)".*/\1/p' | sed 's:.*/\([^/]*\)/[^/]*$:\1:')"
    pluginNameList=("${firstPagePluginNames[@]}")

    for ((i = 2; i <= 3; i++)); do
        page_suffix="/page/$i"
        sub_page_url="${webpage_url%/}${page_suffix}/"
        local sub_page_html_content=$(curl -s "$sub_page_url")
        mapfile -t pagePluginNames <<< "$(echo "$sub_page_html_content" | grep -o '<h3[^>]*class="entry-title"[^>]*>.*<\/h3>' | sed -n -e 's/.*<a[^>]*href="\([^"]*\)".*/\1/p' | sed 's:.*/\([^/]*\)/[^/]*$:\1:')"
        pluginNameList=("${pluginNameList[@]}" "${pagePluginNames[@]}")
    done

    array_length=${#pluginNameList[@]}
    echo -e "\e[1;32mDone. $array_length found.\e[0m"
}

# Help
helpMenu(){
    echo -e "\e[1;33mArguments:\n\t-u\t\twordpress url\n\t\n"
    echo -e "Send over Wingman:\n./scan.sh -u www.example.com\n \e[1;32m"
}

testUrl() {
    local url=$1 
    CHECK_URL=$(curl -o /dev/null --silent --head --write-out '%{http_code}' "$url")
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
            echo -e "\e[1;34m$pluginName\e[0m"
        fi
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