# rj

`rj` is a command line tool for printing HTTP Response as JSON.

![Carbon](https://user-images.githubusercontent.com/10682/150918329-5f516b83-2fcc-42b7-8126-7af5dc1a1f9f.png)

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

![Screenshot](https://user-images.githubusercontent.com/10682/150918416-6c5b73ce-7875-4582-a9d4-08327a5342e8.png)

### HTTP/3

Now, support HTTP/3 with `--http3` option:

![Screenshot](https://user-images.githubusercontent.com/10682/143975571-3925c02d-113d-414f-b2cc-a445c54bbd18.png)

## Related projects

- [reorx/httpstat](https://github.com/reorx/httpstat)
- [jaygooby/ttfb.sh](https://github.com/jaygooby/ttfb.sh)
- [gfx/hj](https://github.com/gfx/hj)

## Author

Yusuke Wada <https://github.com/yusukebe>

## License

MIT
