![Works on my machine](https://img.shields.io/badge/works-on%20my%20machine-green)

# Investigo

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Stargazers](https://img.shields.io/github/stars/tdh8316/Investigo.svg)](https://github.com/tdh8316/Investigo/stargazers)
[![Go](https://img.shields.io/badge/Go-%3E%3D1.12-blue.svg)](https://github.com/tdh8316/Investigo)
[![Tweet](https://img.shields.io/twitter/url/https/github.com/tdh8316/Investigo.svg?style=social)](https://twitter.com/intent/tweet?text=WoW:&url=https%3A%2F%2Fgithub.com%2Ftdh8316%2FInvestigo)

## Overview

*Investigo* is a command-line tool written in `Go` to investigate usernames and [download](https://github.com/tdh8316/Investigo/tree/master/downloaders) their data on [social media](./sites.md).

![screenshot](./docs/GIF.gif)

## Build from source

To compile Investigo, Go 1.12 or higher is required.

```text
git clone https://github.com/tdh8316/Investigo.git
cd Investigo
go build
```

## Usage

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
        --no-output           disable file output
        --update              update database before run from Sherlock repository
        -t, --tor             use tor proxy
        -v, --verbose         verbose output
        -d, --download        download the contents of site if available

options:
        --database DATABASE   use custom database
        --sites SITES         specific sites to investigate (Separated by comma)
        --timeout SECONDS     specific http request timeout
```

## Database
Investigo uses [Sherlock database](https://github.com/sherlock-project/sherlock).

Therefore, in case you want to add a new site to the database, you can make an issue or a pull request on the [Sherlock repository](https://github.com/sherlock-project/sherlock).

## License

Licensed under the MIT License

Copyright 2024 Donghyeok Tak
