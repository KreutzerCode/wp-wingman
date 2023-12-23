<div align="center">

![wp-wingman](./img/logo.png)

</div>

<p align="center">
The WordPress Plugin Scanner designed for identifying any plugins on WordPress sites.
</p>

## Features

- Fetches up-to-date plugin slugs from the WordPress Plugins API.
- Supports rate limiting to avoid excessive requests to the target site.
- Checks for the existence of each plugin on the target WordPress site.
- Possibility to save the plugin slugs collected via the Wordpress API in a file.
- Provides a summary and the option to save the results in a file

## Functionality

The Purpose of this devensive Penetation testing tool to check every installed Plugin on a target system. It utilizes the WordPress Plugins API to fetch plugin slugs based on a specified tag. And runs it against the target system. The user can provide an optional rate limit for requests, and an optional plugin tag and other options to save the collected data or run the script in a specific mode.

## Intentions

This script is intended for security testing. Any use should be approved by the owner of the target website. Use the rate limiting options if you are concerned about the server load caused.

## Requirements

Ensure that you have the following dependencies installed:

- curl
- jq

## Install

```yaml
git clone https://github.com/KreutzerCode/wp-wingman.git
cd wp-wingman
chmod -R 777 scan.sh
./scan.sh
```

## Usage

```yaml
┌──(you㉿kali)-[~/Desktop/wp-wingman]
└─$ ./scan.sh
__        ______   __        _____ _   _  ____ __  __    _    _   _
\ \      / /  _ \  \ \      / /_ _| \ | |/ ___|  \/  |  / \  | \ | |
 \ \ /\ / /| |_) |  \ \ /\ / / | ||  \| | |  _| |\/| | / _ \ |  \| |
  \ V  V / |  __/    \ V  V /  | || |\  | |_| | |  | |/ ___ \| |\  |
   \_/\_/  |_|        \_/\_/  |___|_| \_|\____|_|  |_/_/   \_\_| \_|

                            @kreutzercode
Arguments:
        required: -u                    wordpress url
        optional: -t                    wordpress plugin tag (default securtiy)
        optional: -r                    rate limit on target (default 0-1s)
        optional: --overdrive           check all public plugins on target (very aggressiv)
        optional: --save-playbook       save collected plugins in file
        optional: --save-result         save plugins found on target in file

Send over Wingman:
./scan.sh -u www.example.com -r 5 -t newsletter

Happy scanning!
```

### Overdrive

#### Argument: `--overdrive`

In overdrive mode, the script gathers and evaluates all plugins accessible through the WordPress plugin API on the specified target. The collection process may take some time. During this mode, any default or custom rate limits are deactivated.

### Save Playbook

#### Argument: `--save-playbook`

With the `--save-playbook` argument, the plugin slugs collected from the wordpress api are saved to a `wp-wingman-x.txt` file.

The script automatically determines whether a save file exists for the current mode or tag and asks whether it should be used.

### Save Result

#### Argument: `--save-result`

With the argument `--save-result` the found plugins are saved in a file `wp-wingman-x-x.txt` after a successful operation.

## Example output

![example output](./img/example_usage.png)

# TODO

- extend -t flag to fetch all plugins wia "all" argument
- rework overdrive mode to run multi threadet
- add custom error messages for invalud / missing arguments
