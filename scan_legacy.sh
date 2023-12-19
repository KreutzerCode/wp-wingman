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
plugin_name_list=()
plugin_name_list_length=0
max_string_length=0
current_plugin_in_check_index=0
rate_limit=1
target_plugin_tag="security"
wp_url=""
plugins_found_on_target=()
overdrive_active=false
save_playbook=false
save_result=false

function FetchPluginsByTag() {
    echo -e "\e[1;33mUpdating PlayBook...\e[0m"
    local target_api_endpoint="https://api.wordpress.org/plugins/info/1.2/?action=query_plugins&request[tag]=${target_plugin_tag}"

    if [ "$overdrive_active" == true ]; then
        echo -e "\e[1;31mThat takes a while because of...OVERDRIVE!!!\e[0m"
        target_api_endpoint="https://api.wordpress.org/plugins/info/1.2/?action=query_plugins&request[browse]"
    fi

    local response=$(curl -g -s -A "${user_agents[RANDOM % ${#user_agents[@]}]}" "${target_api_endpoint}")

    local page=2
    local total_pages=$(jq -r '.info.pages' <<<"$response")
    mapfile -t plugin_name_list < <(jq -r '.plugins[].slug' <<<"$response")

    while [ "$page" -le "$total_pages" ]; do
        response=$(curl -g -s -A "${user_agents[RANDOM % ${#user_agents[@]}]}" "${target_api_endpoint}&request[page]=${page}")
        plugin_name_list+=($(echo "$response" | jq -r '.plugins[].slug'))
        ((page++))
    done

    plugin_name_list_length=${#plugin_name_list[@]}
    max_string_length=$(printf "%s\n" "${plugin_name_list[@]}" | awk '{ if (length > x) x = length } END { print x }')

    if [ "$save_playbook" == true ]; then
        SavePlaybookToFile
    fi

    if [ "$overdrive_active" == true ]; then
        echo -e "\e[1;31mDone. $plugin_name_list_length found!!!\e[0m"
    else
        echo -e "\e[1;32mDone. $plugin_name_list_length found.\e[0m"
    fi

    # Display the plugin names
    #echo "Plugin Names:"
    #for name in "${plugin_name_list[@]}"; do
    #    echo "$name"
    #done
}

function HelpMenu() {
    echo -e "\e[1;33mArguments:\n\t\e[1;31mrequired:\e[1;33m -u\t\t\twordpress url\e[1;33m\n\t\e[1;34moptional:\e[1;33m -t\t\t\twordpress plugin tag (default securtiy)\t\t\t\n\t\e[1;34moptional:\e[1;33m -r\t\t\trate limit on target (default 0-1s)\n\t\e[1;33m\e[1;34moptional:\e[1;33m --overdrive\t\tcheck all public plugins on target (very aggressiv)\n\t\e[1;33m\e[1;34moptional:\e[1;33m\e[1;33m --save-playbook\tsave collected plugins in file\n\t\e[1;33m\e[1;34moptional:\e[1;33m\e[1;33m --save-result\t\tsave plugins found on target in file\n\t\e[1;33m"
    echo -e "Send over Wingman:\n./scan.sh -u www.example.com -r 5 -t newsletter \e[1;32m"
}

function TestUrlForAvailability() {
    local url=$1
    local check_url=$(curl -s -A "${user_agents[RANDOM % ${#user_agents[@]}]}" -o /dev/null --head --write-out '%{http_code}' "$url")
    if [ "$check_url" -eq 200 ]; then
        echo "true"
    else
        echo "false"
    fi
}

function PrintFindings() {
    local is_found=$1
    local plugin_name=$2

    if [ "$is_found" == "true" ]; then
        echo -e "\e[1;31m$(printf "%-${max_string_length}s" "$plugin_name")\e[0m \e[1;31m[found]\e[0m\033[K"

        plugins_found_on_target+=($plugin_name)
    else
        printf "\e[1;34m%-${max_string_length}s\e[0m \e[1;34m[ok][%d/%d]\e[0m\r" "$plugin_name" "$((current_plugin_in_check_index + 1))" "$plugin_name_list_length"
    fi
}

function CheckPluginsAvailablity() {
    local url=$1
    echo -e "\n\e[1;33m[+] Let me check this for you:\e[0m\n"
    local plugins_prefix="wp-content/plugins"
    local plugin_suffix="readme.txt"

    for plugin_name in "${plugin_name_list[@]}"; do
        result=$(TestUrlForAvailability "$url/$plugins_prefix/$plugin_name/$plugin_suffix")

        PrintFindings "$result" "$plugin_name"
        ((current_plugin_in_check_index++))

        # Introduce a rate_limit between 0 and X seconds
        # Add one to get desired value in sek 4 = 3
        # Only when not in OVERDRIVE!!!
        if [ "$overdrive_active" == false ]; then
            sleep $(($RANDOM % $rate_limit + 1))
        fi
    done

    PrintResult

    if [ "$save_result" == true ]; then
        SaveResultToFile
    fi

    exit
}

function SaveResultToFile() {
    echo -e "\e[1;33mSaving Result...\e[0m"
    local timestamp=$(date +"%Y%m%d%H%M%S")
    local file_name=$(echo "$url" | sed -e 's|^.*://||' -e 's|/.*$||' | awk -F[/:.] '{print $(NF-1)}')

    for string in "${plugins_found_on_target[@]}"; do
        echo "$string" >>"wp-wingman-$file_name-$timestamp.txt"
    done

    echo -e "\e[1;32mDone. Have a great day!\e[0m\n"
}

function PrintResult() {
    echo -e "\n\n\n\e[1;32mDone.\e[0m\n"
    echo -e "\e[1;32mSummary:\e[0m\n"

    if [ "${#plugins_found_on_target[@]}" -ne 0 ]; then

        for plugin_name in "${plugins_found_on_target[@]}"; do
            echo -e "\e[1;31m$(printf "%-${max_string_length}s" "$plugin_name")\e[0m \e[1;31m[found]\e[0m"
        done

        echo -e "\n\e[1;32mThese are my findings. Good luck sir!\e[0m\n"
    else
        echo -e "\e[1;32mNothing found. Good luck.\e[0m\n"
    fi
}

function StartWingmanJob() {
    result=$(TestUrlForAvailability "$wp_url/wp-login.php")
    if [ "$result" == "true" ]; then
        echo -e "\e[1;32mWordPress site detected: $wp_url\e[0m"

        if CheckIfSaveFileExists; then
            echo -e "\e[1;33mSave file found - should i use it? (y/n)\e[0m"
            read answer
            if [ "$answer" == "y" ]; then
                LoadPluginNamesFromSaveFile
            else
                FetchPluginsByTag
            fi
        else
            FetchPluginsByTag
        fi

        echo -e "\e[1;33mDo you want me to start? (y/n)\e[0m"
        read answer
        if [ "$answer" != "y" ]; then
            echo -e "\e[1;32mPuuh, okey bye.\e[0m\n"
            exit
        fi

        CheckPluginsAvailablity $wp_url
    else
        echo -e "\e[1;31mThe URL is not a WordPress site.\e[0m"
        echo -e "\e[1;31m$wp_url\e[0m"
        exit 1
    fi
}

function LoadPluginNamesFromSaveFile() {
    echo -e "\e[1;33mLoading Playbook from save file...\e[0m"
    local file_name="wp-wingman-${target_plugin_tag}.txt"
    if [ "$overdrive_active" == true ]; then
        file_name="wp-wingman-overdrive.txt"
    fi

    while IFS= read -r line || [ -n "$line" ]; do
        plugin_name_list+=("$line")
    done <$file_name

    plugin_name_list_length=${#plugin_name_list[@]}
    max_string_length=$(printf "%s\n" "${plugin_name_list[@]}" | awk '{ if (length > x) x = length } END { print x }')

    if [ "$overdrive_active" == true ]; then
        echo -e "\e[1;31mDone. $plugin_name_list_length found!!!\e[0m"
    else
        echo -e "\e[1;32mDone. $plugin_name_list_length found.\e[0m"
    fi
}

function SavePlaybookToFile() {
    echo -e "\e[1;33mSaving Playbook...\e[0m"
    local file_name="wp-wingman-${target_plugin_tag}.txt"
    if [ "$overdrive_active" == true ]; then
        file_name="wp-wingman-overdrive.txt"
    fi

    # Remove the existing file if it exists
    if [ -e "$file_name" ]; then
        rm "$file_name"
    fi

    for string in "${plugin_name_list[@]}"; do
        echo "$string" >>"$file_name"
    done
}

function CheckIfSaveFileExists() {
    local file_name="wp-wingman-${target_plugin_tag}.txt"
    if [ "$overdrive_active" == true ]; then
        file_name="wp-wingman-overdrive.txt"
    fi
    file_path="$(dirname "$(readlink -f "$0")")/${file_name}"

    [ -e "$file_path" ] && true || false
}

args=()

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
    -u)
        shift
        wp_url="$1"
        args+=("-u" "$wp_url")
        ;;
    -r)
        shift
        r_value="$1"
        args+=("-r" "$r_value")
        ;;
    -t)
        shift
        t_value="$1"
        args+=("-t" "$t_value")
        ;;
    --overdrive)
        overdrive_active=true
        ;;
    --save-playbook)
        save_playbook=true
        ;;
    --save-result)
        save_result=true
        ;;
    -*)
        echo -e "\n\e[1;31mInvalid argument: $1\e[0m\n"
        HelpMenu
        exit 1
        ;;
    esac
    shift
done

if [ ${#args[@]} -eq 0 ]; then
    HelpMenu
    exit 1
fi

# Check if -u argument is missing
if [ -z "$wp_url" ]; then
    echo -e "\n\e[1;31mError: Missing -u argument\e[0m\n"
    HelpMenu
    exit 1
fi

# move the -u argument to the end so the configuration takes place first
u_index=$(printf '%s\n' "${args[@]}" | grep -n '^\-u' | cut -f1 -d:)
[ -n "$u_index" ] && args=("${args[@]:0:$u_index-1}" "${args[@]:$u_index+1}" "-u" "$wp_url")

for ((i = 0; i < ${#args[@]}; i += 2)); do
    option="${args[i]}"
    value="${args[i + 1]}"
    if [ "$option" == "-r" ] && [ "$overdrive_active" == false ]; then
        rate_limit="$value"
        echo -e "\e[1;32mSet rate limit to: $value\e[0m"
    fi

    if [ "$option" == "-t" ] && [ "$overdrive_active" == false ]; then
        target_plugin_tag="$value"
        echo -e "\e[1;32mSet plugin tag to: $value\e[0m"
    fi

    if [ "$option" == "-u" ]; then
        StartWingmanJob
    fi
done
