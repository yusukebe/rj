# rj

`rj` is a command line tool for printing HTTP Response as JSON.

![Carbon](https://user-images.githubusercontent.com/10682/150896161-dd811e3b-fd43-4f69-bab3-451cca853ef5.png)

## Installation

### Homebrew

You can also install via hombrew on macOS:

```plain
$ brew install yusukebe/tap/rj
```

### Binary

Download the binary from [GitHub Releases](https://github.com/yusukebe/rj/releases) and install it somewhere in your $PATH. rj currently provides pre-built binaries for Linux, macOS and Windows.

### Source

To install from the source, use go install:

```plain
$ go install github.com/yusukebe/rj/cmd/rj@latest
```

## Usage

The usage:

```plain
$ rj [url] [flags]
```

Available options:

```plain
  -A, --agent string         User-Agent name (default "rj/{{ Version }}")
  -u, --basic string         Basic Auth username:password
  -H, --header stringArray   HTTP Request Header
  -h, --help                 help for rj
      --http1.1              Use HTTP/1.1
      --http3                Use HTTP/3
  -X, --method string        HTTP Request method (default "GET")
  -v, --version              version for rj
```

### Screenshots

With `jq`:

![Screenshot](https://user-images.githubusercontent.com/10682/150894939-337d2431-1c91-4b36-8e5d-a60d255ed576.png)

### HTTP/3

Now, support HTTP/3 with `--http3` option:

![Screenshot](https://user-images.githubusercontent.com/10682/143975571-3925c02d-113d-414f-b2cc-a445c54bbd18.png)

## Related projects

- [reorx/httpstat](https://github.com/reorx/httpstat)
- [jaygooby/ttfb.sh](https://github.com/jaygooby/ttfb.sh)

## Author

Yusuke Wada <https://github.com/yusukebe>

## License

MIT
