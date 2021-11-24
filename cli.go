package rj

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptrace"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dyatlov/go-htmlinfo/htmlinfo"
	"github.com/spf13/cobra"
)

type param struct {
	method      string
	userAgent   string
	headers     []string
	includeBody bool
}

var rootCmd = &cobra.Command{
	Use:   "rj [url]",
	Short: "rj is a command line tool show the HTTP Response as JSON",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		url := args[0]
		method, _ := cmd.Flags().GetString("method")
		userAgent, _ := cmd.Flags().GetString("agent")
		headers, _ := cmd.Flags().GetStringArray("header")
		includeBody, _ := cmd.Flags().GetBool("include-body")
		param := param{
			method:      method,
			userAgent:   userAgent,
			headers:     headers,
			includeBody: includeBody,
		}
		request(url, param)
	},
}

func init() {
	rootCmd.Flags().StringP("method", "X", "GET", "HTTP Request method")
	rootCmd.Flags().StringP("agent", "A", "rj/v0.0.1", "User-Agent name")
	rootCmd.Flags().StringArrayP("header", "H", nil, "HTTP Request Header")
	rootCmd.Flags().BoolP("include-body", "b", false, "Include Response body")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func request(url string, param param) {
	req, _ := http.NewRequest("GET", url, nil)

	if param.userAgent != "" {
		req.Header.Set("User-Agent", param.userAgent)
	}

	if param.headers != nil {
		for _, h := range param.headers {
			kv := strings.Split(h, ":")
			req.Header.Set(strings.TrimSpace(kv[0]), kv[1])
		}
	}

	var start, connect, dns, tlsHandshake, wait time.Time
	var dnsMs, sslMs, connectionMs, ttfbMs, totalMs float64

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			dnsMs = timeToMs(dns)
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			connectionMs = timeToMs(connect)
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			sslMs = timeToMs(tlsHandshake)
			wait = time.Now()
		},

		GotFirstResponseByte: func() {
			ttfbMs = timeToMs(wait)
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	start = time.Now()

	client := new(http.Client)
	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer res.Body.Close()
	totalMs = timeToMs(start)

	r := make(map[string]interface{})

	r["status"] = res.Status
	r["code"] = res.StatusCode

	timing := make(map[string]interface{})
	timing["dns_lookup"] = dnsMs
	timing["tcp_connection"] = connectionMs
	timing["tls_handshake"] = sslMs
	timing["ttfb"] = ttfbMs
	timing["total"] = totalMs
	r["timing"] = timing

	headers := make(map[string]string)

	for key, value := range res.Header {
		headkey := strings.ToLower(key)
		headValue := strings.Join(value, ", ") // XXX
		headers[headkey] = headValue
	}

	r["header"] = headers

	if param.includeBody {
		if contentType, ok := headers["content-type"]; ok {
			if matchRegexp(contentType, `text/html`) {
				info := htmlinfo.NewHTMLInfo()
				err = info.Parse(res.Body, &url, &contentType)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				r["body"] = info
			} else if matchRegexp(contentType, `application/json`) {
				var data map[string]interface{}
				err := json.NewDecoder(res.Body).Decode(&data)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				r["body"] = data
			}
		}
	}

	bytes, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(bytes))
}

func timeToMs(t time.Time) float64 {
	return floor(time.Since(t).Seconds())
}

func floor(f float64) float64 {
	return math.Floor(f*100000) / 100000
}

func matchRegexp(str string, regString string) bool {
	reg := regexp.MustCompile(regString)
	return reg.MatchString(str)
}
