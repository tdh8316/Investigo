# Overview

![Commits](https://img.shields.io/github/commit-activity/m/tdh8316/Investigo.svg?label=Commits&style=popout)
![Go](https://img.shields.io/badge/Go-%3E%3D1.8-blue.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![GitHub stars](https://img.shields.io/github/stars/tdh8316/Investigo.svg?style=social)](https://github.com/tdh8316/Investigo/stargazers)

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

# Usage
```cmd
$ investigo USERNAMES [-v|--verbose] [--no-color]
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

# Site support
Investigo uses Sherlock's database.

Accordingly, If you want to add sites to Investigo, make a PR to [Sherlock project](https://github.com/sherlock-project/sherlock). Or if you want to add prohibited sites for Sherlock, make an issue or a PR to this repository!

# License
Licensed under the MIT License

Copyright 2019 Donghyeok Tak
