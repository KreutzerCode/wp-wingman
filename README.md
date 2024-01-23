<div align="center">

![wp-wingman](./img/logo.png)

</div>

<p align="center">
The WordPress Plugin Scanner designed for identifying any plugins on WordPress sites.
</p>

## Features

- Fetches up-to-date plugin slugs from the WordPress plugins API.
- Supports rate limiting to avoid excessive requests to the target site.
- Checks for the existence of each plugin on the target WordPress site.
- Possibility to save the plugin slugs collected via the Wordpress API in a file.
- Provides a summary and the option to save the results in a file.
- Supports the parallel execution of plugin slug checks.
- Detecting unknown plugins in the web content of the target site.

## Functionality

The Purpose of this devensive Penetation testing tool to check every installed Plugin on a target system. It utilizes the WordPress Plugins API to fetch plugin slugs based on a specified tag. And runs it against the target system. The user can specify optional arguments to change the plugin search or plugin check procedure. In addition, during the check, the user has the option of having the content of the target site checked for additional plugins, which are most likely to be premium or custom plugins.

## Intentions

This script is intended for security testing. Any use should be approved by the owner of the target website. Use the rate limiting options if you are concerned about the server load caused.

## Download

Visit [the lastest release](https://github.com/KreutzerCode/wp-wingman/releases/latest) and download the build for your environment (currently only available for Linux).

## Install (Linux)

```yaml
$ cd /path/to/binary
$ chmod -R 777  wp-wingman-linux-amd64
$ ./wp-wingman-linux-amd64
```

## Usage

### Linux

```yaml
┌──(you㉿linux)-[~/Desktop]
└─$ ./wp-wingman-linux-amd64
__        ______   __        _____ _   _  ____ __  __    _    _   _
\ \      / /  _ \  \ \      / /_ _| \ | |/ ___|  \/  |  / \  | \ | |
 \ \ /\ / /| |_) |  \ \ /\ / / | ||  \| | |  _| |\/| | / _ \ |  \| |
  \ V  V / |  __/    \ V  V /  | || |\  | |_| | |  | |/ ___ \| |\  |
   \_/\_/  |_|        \_/\_/  |___|_| \_|\____|_|  |_/_/   \_\_| \_|

                            @kreutzercode
Arguments:
        required: -u                    wordpress url
        optional: -t                    wordpress plugin tag (default securtiy)
        optional: -r                    rate limit on target (default 0s)
        optional: -w                    number of workers to execute playbook (default 10)
        optional: --save-playbook       save collected plugins in file
        optional: --save-result         save plugins found on target in file

Send over Wingman:
./wp-wingman -u www.example.com -r 5 -t newsletter

Happy scanning!
```

### Url

#### Argument: `-u`

With the -u argument you specify the url of the target system. Note that this is the only **required** argument.

**Examples:**

```
$ ./wp-wingman -u www.example.com

$ ./wp-wingman -u https://www.example.com
```

### Tag

#### Argument: `-t`

With the -t argument you can specify the target plugin group by searching for a specific tag like _security_, _newsletter_ etc.  
If you want to fetch all public plugins, you can add the all argument `-t all`. Be in mind that this takes a while.

**Tip**: use in combination with `--save-playbook` to skip waiting time in the next run

### Rate limit

#### Argument: `-r`

**Important**: This argument will reset the used *workers* to **1**.

With the -r argument, you can specify the number of seconds to wait before the next plugin slug is checked on the target system. The timeout is set randomly between *0* and *your entry*.

### Workers

#### Argument: `-w`

With the -w argument, you can specify the number of workers that process the playbook in parallel.

- Min Value: 1

### Save Playbook

#### Argument: `--save-playbook`

With the `--save-playbook` argument, the plugin slugs collected from the wordpress api are saved to a `wp-wingman-x.txt` file.

The script automatically determines whether a save file exists for the current mode or tag and asks whether it should be used.

### Save Result

#### Argument: `--save-result`

With the argument `--save-result` the found plugins are saved in a file `wp-wingman-x-x.txt` after a successful operation.

## Example output

![example output](./img/example_usage.png)