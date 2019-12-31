# Overview

![Commits](https://img.shields.io/github/commit-activity/m/tdh8316/Investigo.svg?label=Commits&style=popout)
![Go](https://img.shields.io/badge/Go-%3E%3D1.8-blue.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![GitHub stars](https://img.shields.io/github/stars/tdh8316/Investigo.svg?style=social)](https://github.com/tdh8316/Investigo/stargazers)
<a href="https://travis-ci.com/tdh8316/Investigo"><img src="https://img.shields.io/travis/tdh8316/Investigo.svg" /></a>
[![Twitter](https://img.shields.io/twitter/url/https/github.com/tdh8316/Investigo.svg?style=social)](https://twitter.com/intent/tweet?text=WoW:&url=https%3A%2F%2Fgithub.com%2Ftdh8316%2FInvestigo)

+ [English](./README.MD) 
  
  ~~[한국어](./README-ko_kr.MD)~~

*Investigo* is a command-line tool written in `Go` for searching usernames on [social media](./sites.md).

![screenshot](./docs/GIF.gif)

# Installation
## Download Windows Binary
You can download the binary for Windows.
 - [Download latest release](https://github.com/tdh8316/Investigo/releases)

`$ investigo.exe -h`

## Build from source
To compile investigo, Go 1.8 or higher is required.

```cmd
$ git clone https://github.com/tdh8316/Investigo.git
$ cd investigo
$ go get github.com/fatih/color
$ go run investigo.go -h
```

## Build with docker


# Usage
```cmd
$ investigo [-h] [--no-color] [-v|--verbose] [-t|--tor] [--update] USERNAME [USERNAMES...]
positional arguments:
	USERNAMES             one or more usernames to investigate

optional arguments:
	-h, --help			  show this help message and exit
	-v, --verbose         output sites which is username was not found
	--no-color            disable colored stdout output
	--update              update datebase from github.com/tdh8316/investigo/
```

# Why Investigo?
Investigo is trying to build our own database while keeping compatibility with Sherlock.
This is because even though Sherlock project contains a vast database, adult sites are forbidden.

Furthermore, Investigo is more than twice as fast as the original implementation!

# TODO
## Manifesto
Investigo is another implementation of Sherlock. Accordingly, I would like not to make a competition with Sherlock.

 - [ ] Avoid deep packet inspection
 - [x] Support Tor
 - [ ] Implement regex check
 - [x] Command-line arguments
 - [x] Add tests

## Become a stargazer
★

# Addding/Fixing site
Investigo uses Sherlock's database.

## To add sites...
If you want to add sites to Investigo, make a Pull Request to [Sherlock project](https://github.com/sherlock-project/sherlock). Or if you want to add prohibited sites for Sherlock, make an issue or a Pull Request to this repository!

## To fix incompatible sites...
Make an issue [here](https://github.com/sherlock-project/sherlock/issues) to let Sherlock developers know.

# License
Licensed under the MIT License

Copyright 2019 Donghyeok Tak
