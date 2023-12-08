<h1 align="center">
  WP Wingman
  <br>
</h1>

<h4 align="center">The helper we all need some times.</h4>

<hr>

### WP Wingman is a small helper that lets you now if any popular security plugins are installed on YOUR wordpress site.

<br>

### Functionality

The Wingman uses a static list of known security plugins and checks if the plugin directory of that plugin contains a readme.txt and reports its findings back to you.

<br>

### Documentation

### install

```yaml
git clone https://github.com/KreutzerCode/wp-wingman.git
cd wp-wingman
chmod -R 777 scan.sh
./scan.sh
```

#### Usage

```yaml
┌──(you㉿kali)-[~/Desktop/wp-wingman]
└─$ ./scan.sh
__        ______   __        _____ _   _  ____ __  __    _    _   _
\ \      / /  _ \  \ \      / /_ _| \ | |/ ___|  \/  |  / \  | \ | |
 \ \ /\ / /| |_) |  \ \ /\ / / | ||  \| | |  _| |\/| | / _ \ |  \| |
  \ V  V / |  __/    \ V  V /  | || |\  | |_| | |  | |/ ___ \| |\  |
   \_/\_/  |_|        \_/\_/  |___|_| \_|\____|_|  |_/_/   \_\_| \_|

                            \e[1;34m  @kreutzercode
Arguments:
        -u              wordpress url

Send over Wingman:
./scan.sh -u www.example.com

```

<br>
