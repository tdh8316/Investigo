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
> Note that these binary releases are not the latest version.

- [Download latest release](https://github.com/tdh8316/Investigo/releases/latest)

`$ investigo.exe -h`

## Build from source

To compile investigo, Go 1.12 or higher is required.

```text
git clone https://github.com/tdh8316/Investigo.git
cd investigo
go build
```

# Usage

```text
$ investigo -h
Investigo - Investigate User Across Social Networks.

usage: investigo USERNAME [USERNAMES...] flags options
perform test: investigo --test

positional arguments:
        USERNAMES             one or more usernames to investigate

flags:
        -h, --help            show this help message and exit
        --no-color            disable colored stdout output
        --update              update database before run from Sherlock repository
        -t, --tor             use tor proxy
        -v, --verbose         verbose output
        -d, --download        download the contents of site if available

options:
        --database DATABASE   use custom database
        --sites SITES         specific sites to investigate (Separated by comma)
        --timeout SECONDS     specific http requeset timeout
```

# Addding/Fixing site

Investigo is compatible with [Sherlock database](https://github.com/sherlock-project/sherlock).

## To add sites

To add sites to Investigo, make a Pull Request to [Sherlock project](https://github.com/sherlock-project/sherlock). Whenever the Sherlock repository updates, the Investigo database will update automatically!

If you want to add prohibited sites for Sherlock, make an issue or a Pull Request to this repository!

## To fix incompatible sites

Make an issue [here](https://github.com/sherlock-project/sherlock/issues) to let Sherlock developers know.

# Contributors

Thanks to all contributors!

- [x0rzkov](https://github.com/x0rzkov)
- [lucmski](https://github.com/lucmski)

Contributions are more than welcome!

# License

Licensed under the MIT License

Copyright 2022 `Donghyeok Tak`
