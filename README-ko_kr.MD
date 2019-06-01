# 개요

![Commits](https://img.shields.io/github/commit-activity/m/tdh8316/Investigo.svg?label=Commits&style=popout)
![Go](https://img.shields.io/badge/Go-%3E%3D1.8-blue.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![GitHub stars](https://img.shields.io/github/stars/tdh8316/Investigo.svg?style=social)](https://github.com/tdh8316/Investigo/stargazers)

  [English](./README.MD) 
  
+ [한국어](./README-ko_kr.MD)

*Investigo*는 *[Sherlock](https://sherlock-project.github.io/)* 의 `Go`언어 구현입니다.

[여기](./sites.md)에서 지원되는 사이트 목록을 확인할 수 있습니다.

![screenshot](./docs/sample.png)

# 다운로드
investigo의 최신 버전을 다운로드하세요.
 - [Releases](https://github.com/tdh8316/Investigo/releases)

# 사용법
`$ investigo.exe USERNAMES... [--verbose] [--disable-color] [--site NAME]`
## 직접 컴파일하기
investigo를 컴파일하기 위해서는 Go 1.8 이상이 필요합니다.
```
$ git clone https://github.com/tdh8316/Investigo.git
$ cd investigo
$ go run investigo.go USERNAMES... [--verbose] [--disable-color] [--site NAME]
```

## 선택적 인자
```
Investigo - Investigate User Across Social Networks.

Usage: go run investigo.go [-h] [--no-color] [--verbose] [--update] [--site SITE_NAME] USERNAMES

 positional arguments:
        USERNAMES          찾을 ID를 지정합니다.

 optional arguments:
        --site             찾을 사이트를 지정합니다. (e.g --site github)
        --disable-color
        --verbose          해당 ID가 검색되지 않은 사이트 목록을 출력합니다.
        --update           사이트 목록을 업데이트합니다.
```

# License
Licensed under the MIT License

Copyright 2019 Donghyeok Tak
