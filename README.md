## PADDOCK

**Paddock** is a fast, secure and cross-platform cli tool for automating F1 TV subscription token retrieval.

It uses [rod](https://github.com/go-rod/rod) to connect to a browser instance managed by the CLI. Paddock automates the interaction required to log in or out securely. Currently, Google Chrome and Microsoft Edge are supported.
<br>
<br>
When executing the program for the first time, it may take longer than usual (especially on Windows due to security scans). The program caches your token in your system's standard configuration directory. Subsequent runs retrieve the cached token instantly (~50ms) without launching the browser.

---

## 

## Installation

### MacOS & Linux

```bash
brew install --cask eepzii/tap/paddock
```

### Windows

```bash
scoop bucket add paddock https://github.com/eepzii/scoop-bucket
scoop install paddock
```

### Using curl

```bash
curl -sL https://raw.githubusercontent.com/eepzii/paddock/main/install.sh | sh
```

### Manual

1. head over to [releases page](https://github.com/eepzii/paddock/releases/) 
2. download the archive for your system
3. extract the binary

---

## Usage

### bash
```bash
export PASSWORD="your_password"
paddock --email="your_email" 

# or
PASSWORD="your_password" paddock --email="your_email"
``` 

### powershell
```bash
$env:PASSWORD="your_password"
paddock --email="your_email"
```

#### flags

| Flag | Type | Description |
| :--- | :--- | :--- |
| `--email` | `string` | **(required)** expects your F1 TV account email |
| `--headless` | `bool` | disables the visual ui of the chromium browser |
| `--logout` | `bool` | will log you out from your current session |
| `--force` | `bool` | forces a login with a new set `--email` |
| `--path` | `string` | **(not safe)** sets the path to your favorite chromium based browser |
| `--freshness` | `duration` | will check whether the token is in x future still valid |

<br>

#### experimental: using a proxy

With paddock version >=0.3.0 you can now use a proxy. It is recommended to only use residential proxies that are located relatively close from where your program is executed. This feature is in an early stage. There might follow some updates regarding reliability but for a fast login you should stick to a regular login without a proxy or use a proxy that provides low latency. <br><br>
Before running the standard cli prompt you need to add these envs in the terminal session: <br>
* `PROXY_HOST` (scheme must be `http://ip:port`)
* `PROXY_USER` 
* `PROXY_PASS` 

### bash
```bash
# PROXY_HOST scheme must be: http://ip:port 
export PROXY_HOST="http://ip:port"
export PROXY_USER="user"
export PROXY_PASS="password"
```

### powershell
```bash
# PROXY_HOST scheme must be: http://ip:port 
$env:PROXY_HOST="http://ip:port"
$env:PROXY_USER="user"
$env:PROXY_PASS="password"
```

**NOTE**: using a VPN or a publicly known IP from a data center might fail. <br>

---

## How it works

When entering the command `PASSWORD paddock --email=` you essentially made a login command. <br>
Now you have to wait until you are logged in. (the first run on windows will take a good amount longer) <br>
<br>
The flag `--headless` is on every run recommended.
<br>
When `--headless` is set it will take about 6s-7s to be logged in.
<br>
To log out you simply add the flag `--logout`.

### Output

Paddock will output a JSON object to `stdout`.

#### Success:

```bash
    {
      "success": true,
      "token": "j.w.t",
      # when logging out "token": ""
      "duration": "52.1ms",
      "error": ""
    }
```

#### Failure (Exit Code 1):

```bash
    {
      "success": false,
      "token": "",
      "duration": "8.2s",
      "error": "err"
    }
```

---

## Troubleshooting

* If the program times out a lot then the websites content probably changed. The program needs to be adjusted for the new content. Please contact [me](mailto:julianmaxromeis@gmail.com) or open an issue in the repo.  
* If you forgot your old `--email` you have set once and want to login with a new one, then you have to use the `--force` flag. This will log the old account out and log the new account in.

---

## Disclaimer

This project is unofficial and is not associated in any way with the Formula 1 companies. F1, F1-TV, FORMULA ONE, FORMULA 1, FIA FORMULA ONE WORLD CHAMPIONSHIP, GRAND PRIX and related marks are trade marks of Formula One Licensing B.V