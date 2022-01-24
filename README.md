# rj

`rj` is a command line tool for printing HTTP Response as JSON.

![Screenshot](https://user-images.githubusercontent.com/10682/143975194-2a808418-a4fe-4570-8d16-6495a6d54b7a.png)

## Installation

```plain
go install github.com/yusukebe/rj/cmd/rj@latest
```

## Usage

```plain
rj http://example.com/
```

```plain
rj is a command line tool show the HTTP Response as JSON

Usage:
  rj [flags] [url]

Flags:
  -A, --agent string         User-Agent name (default "rj/v0.0.1")
  -u, --basic string         Basic Auth username:password
  -H, --header stringArray   HTTP Request Header
  -h, --help                 help for rj
      --http1.1              Use HTTP/1.1
      --http3                Use HTTP/3
  -X, --method string        HTTP Request method (default "GET")
```

support HTTP/3 with `--http3` option:

![Screenshot](https://user-images.githubusercontent.com/10682/143975571-3925c02d-113d-414f-b2cc-a445c54bbd18.png)

## Author

Yusuke Wada <https://github.com/yusukebe>

## License

MIT
