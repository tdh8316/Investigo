# Overview

![Commits](https://img.shields.io/github/commit-activity/m/tdh8316/Investigo.svg?label=Commits&style=popout)
![Go](https://img.shields.io/badge/Go-%3E%3D1.8-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

+ [English](./README.MD) 
  
  [한국어](./README-ko_kr.MD)

*Investigo* is the `Go` implementation of *[Sherlock](https://sherlock-project.github.io/)*.

Investigo supports [several Social Media](./sites.md).

![screenshot](./docs/sample.png)

# Usage
## With Go compiler
Go 1.8 or higher is required.
```
$ git clone https://github.com/tdh8316/Investigo.git
$ cd investigo
$ go run investigo.go USERNAMES... [--verbose] [--disable-color] [--site NAME]
```

## Optional Arguments
```
Investigo - Investigate User Across Social Networks.

Usage: go run investigo.go [-h] [--no-color] [--verbose] [--update] [--site SITE_NAME] USERNAMES

 positional arguments:
        USERNAMES          Usernames to investigate

 optional arguments:
        --site             Specify sites to search. (e.g --site github)
        --disable-color
        --verbose          Print sites USERNAME is not found.
        --update           Update list automatically
```

# License
Licensed under the MIT License

Copyright 2019 Donghyeok Tak
