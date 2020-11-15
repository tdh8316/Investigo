![Works on my machine](https://img.shields.io/badge/works-on%20my%20machine-green)

# Overview

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Stargazers](https://img.shields.io/github/stars/tdh8316/Investigo.svg)](https://github.com/tdh8316/Investigo/stargazers)
[![Go](https://img.shields.io/badge/Go-%3E%3D1.12-blue.svg)](https://github.com/tdh8316/Investigo)
[![Tweet](https://img.shields.io/twitter/url/https/github.com/tdh8316/Investigo.svg?style=social)](https://twitter.com/intent/tweet?text=WoW:&url=https%3A%2F%2Fgithub.com%2Ftdh8316%2FInvestigo)

*Investigo* is a command-line tool written in `Go` to investigate usernames and [download](https://github.com/tdh8316/Investigo/tree/master/downloader) their data on [social media](./sites.md).

![screenshot](./docs/GIF.gif)

# Installation
## Download Windows Binary
You can download the binary for Windows.
> Note that the frequency of update of release is too slow.
 - [Download latest release](https://github.com/tdh8316/Investigo/releases/latest)

`$ investigo.exe -h`

## Build from source
To compile investigo, Go 1.12 or higher is required.

```
$ git clone https://github.com/tdh8316/Investigo.git
$ cd investigo
$ go build
```

## Build with docker
If you don't have Docker/Docker-Compose check **Setup Docker** section

<details>
<summary><b>Setup Docker</b></summary>
<p>

## Docker
macOS: <a href="https://docs.docker.com/docker-for-mac/install/"> https://docs.docker.com/docker-for-mac/install/ </a>

linux: <a href="https://docs.docker.com/install/linux/docker-ce/ubuntu/"> https://docs.docker.com/install/linux/docker-ce/ubuntu/ </a>

## Docker Compose

linux: <a href="https://docs.docker.com/compose/install/"> https://docs.docker.com/compose/install/ </a>
</p>
</details>

### Investigo (standalone mode, no screenshots)

<details>
<summary><b>Example with docker</b></summary>
<p>
	<code> 
		$ docker run x0rzkov/investigo:latest investigo USERNAME
	</code>
</p>
</details>

<details>
<summary><b>Example with docker-compose</b></summary>
<p>
	<code> 
		$ docker-compose run investigo investigo USERNAME
	</code>
</p>
</details>

### Investigo (headless chromium mode, with screenshots)

<details>
<summary><b>Example run with docker</b></summary>
<p>
	<code> 
		$ docker run -v $(pwd)/data:/screenshots x0rzkov/investigo-chromium:latest -s USERNAME
	</code>
</p>
</details>

<details>
<summary><b>Example run with docker-compose</b></summary>
<p>
	<code> 
		$ docker-compose run investigo-chromium -s USERNAME
	</code>
</p>
</details>

# Usage
```
$ investigo -h
Investigo - Investigate User Across Social Networks.

usage: investigo USERNAME [USERNAMES...] flags options
perform test: investigo [--test]

positional arguments:
        USERNAMES             one or more usernames to investigate

flags:
        -h, --help            show this help message and exit
        --no-color            disable colored stdout output
        --update              update database before run from Sherlock repository
        -t, --tor             use tor proxy
        -s, --screenshot      take a screenshot of each matched urls
        -v, --verbose         verbose output
        -d, --download        download the contents of site if available

options:
        --database DATABASE   use custom database
        --site SITE           specific site to investigate
```

# Addding/Fixing site
Investigo uses [Sherlock project](https://github.com/sherlock-project/sherlock)'s database.

## To add sites...
To add sites to Investigo, make a Pull Request to [Sherlock project](https://github.com/sherlock-project/sherlock). Whenever the Sherlock repository updates, the Investigo database will update automatically!

If you want to add prohibited sites for Sherlock, make an issue or a Pull Request to this repository!

## To fix incompatible sites...
Make an issue [here](https://github.com/sherlock-project/sherlock/issues) to let Sherlock developers know.

# Contributors
Thanks to all contributors!
 - [x0rzkov](https://github.com/x0rzkov)
 - [lucmski](https://github.com/lucmski)

Contributions are more than welcome!

# License
Licensed under the MIT License

Copyright 2020 `Donghyeok Tak`
