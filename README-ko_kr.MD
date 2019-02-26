# 개요
[English](./README.MD) [한국어](./README-ko_kr.MD)

*Investigo*는 *[Sherlock](https://github.com/TheYahya/sherlock)* 의 `Go`언어 구현입니다.

[여기](./sites.md)에서 지원되는 사이트 목록을 확인할 수 있습니다.

![screenshot](./docs/sample.gif)

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
--site             검색할 사이트를 지정합니다. (e.g --site github)
--disable-color    출력 색상을 비활성화합니다.
--verbose          USERNAME 이 발견되지 않은 사이트도 출력합니다.
```

# License
Licensed under the MIT License

Copyright 2019 Donghyeok Tak
