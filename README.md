# rj

`rj` is a command line tool show the HTTP Response as JSON

![Screenshot](https://user-images.githubusercontent.com/10682/143195610-9c596786-f347-477b-9839-30c9c22f02e1.png)

## Installation

```
$ go install github.com/yusukebe/rj/cmd/rj@latest
```

## Usage

```
$ rj http://example.com/
```

```
rj is a command line tool show the HTTP Response as JSON

Usage:
  rj [url] [flags]

Flags:
  -A, --agent string         User-Agent name (default "rj/v0.0.1")
  -H, --header stringArray   HTTP Request Header
  -h, --help                 help for rj
  -b, --include-body         Include Response body
  -X, --method string        HTTP Request method (default "GET")
```

with `-b` Option:

![Screenshot](https://user-images.githubusercontent.com/10682/143204005-7fc8c8df-0a55-4905-9ad3-ef9cb76e268d.png)

## Author

Yusuke Wada <https://github.com/yusukebe>

## License

MIT