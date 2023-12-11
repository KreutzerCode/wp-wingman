<div align="center">

![wp-wingman](./img/logo.png)

</div>

<p align="center">
The WordPress Plugin Scanner designed for identifying any plugins on WordPress sites.
</p>

## Features

- Fetches plugin names based on a specified tag from the WordPress Plugins API.
- Supports rate limiting to avoid excessive requests to the target site.
- Checks for the existence of each plugin on the target WordPress site.
- Provides user-friendly prompts and outputs for easy interaction.

## Functionality

The script utilizes the WordPress Plugins API to fetch plugin names based on a specified tag. The user can provide a WordPress site URL, an optional rate limit for requests, and an optional plugin tag. The script then checks for the existence of each plugin by attempting to access its readme.txt file on the target system.

## Requirements

Ensure that you have the following necessary dependencies installed:

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
        required: -u              wordpress url
        optional: -t              wordpress plugin tag (default securtiy)
        optional: -r              rate limit on target (default 0-1s)
        optional: --overdrive     checks all public plugins on target (very aggressiv)

Send over Wingman:
./scan.sh -u www.example.com -r 5 -t newsletter

Happy scanning!
```

### Overdrive

In overdrive mode, the script gathers and evaluates all plugins accessible through the WordPress plugin API on the specified target. The collection process may take some time. During this mode, any default or custom rate limits are deactivated.

# TODO

- local storage file for tags with use or update
