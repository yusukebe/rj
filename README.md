# rj

`rj` is a command line tool show the HTTP Response as JSON

![Screenshot](https://user-images.githubusercontent.com/10682/143975194-2a808418-a4fe-4570-8d16-6495a6d54b7a.png)

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
      --http1.1              Use HTTP/1.1
      --http3                Use HTTP/3
  -b, --include-body         Include Response body
  -X, --method string        HTTP Request method (default "GET")
```

with `-b` Option:

![Screenshot](https://user-images.githubusercontent.com/10682/143975402-6cb0d463-acd6-4ccc-ba0b-439998414ae4.png)

support HTTP/3 with `--http3` option:

![Screenshot](
https://user-images.githubusercontent.com/10682/143975571-3925c02d-113d-414f-b2cc-a445c54bbd18.png)


## Author

Yusuke Wada <https://github.com/yusukebe>

## License

MIT