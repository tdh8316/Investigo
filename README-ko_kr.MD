# 개요
[[English]](./README.MD) [[한국어]](./README-ko_kr.MD)

*Investigo*는 *[Sherlock](https://sherlock-project.github.io/)* 의 `Go`언어 구현입니다.

[여기](./sites.md)에서 지원되는 사이트 목록을 확인할 수 있습니다.

![screenshot](./docs/sample.png)

# 사용법
## With Go compiler
Go 1.8 이상이 필요합니다.
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
