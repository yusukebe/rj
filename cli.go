package rj

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"net/http/httptrace"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"

	"github.com/dyatlov/go-htmlinfo/htmlinfo"
	"github.com/spf13/cobra"
)

type param struct {
	method      string
	userAgent   string
	headers     []string
	includeBody bool
	http1_1     bool
	http3       bool
	basicAuth   string
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
		http1_1, _ := cmd.Flags().GetBool("http1.1")
		http3, _ := cmd.Flags().GetBool("http3")
		basicAuth, _ := cmd.Flags().GetString("basic")

		param := param{
			method:      method,
			userAgent:   userAgent,
			headers:     headers,
			includeBody: includeBody,
			http1_1:     http1_1,
			http3:       http3,
			basicAuth:   basicAuth,
		}
		request(url, param)
	},
}

func init() {
	rootCmd.Flags().StringP("method", "X", "GET", "HTTP Request method")
	rootCmd.Flags().StringP("agent", "A", "rj/v0.0.1", "User-Agent name")
	rootCmd.Flags().StringArrayP("header", "H", nil, "HTTP Request Header")
	rootCmd.Flags().BoolP("include-body", "b", false, "Include Response body")
	rootCmd.Flags().BoolP("http1.1", "", false, "Use HTTP/1.1")
	rootCmd.Flags().BoolP("http3", "", false, "Use HTTP/3")
	rootCmd.Flags().StringP("basic", "u", "", "Basic Auth username:password")
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

	if param.basicAuth != "" {
		kv := strings.Split(param.basicAuth, ":")
		req.SetBasicAuth(kv[0], kv[1])
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

	var client *http.Client
	if param.http3 {
		client = &http.Client{
			Transport: http3RoundTripper(),
		}
	} else if param.http1_1 {
		client = &http.Client{
			Transport: http1_1Transport(),
		}
	} else {
		client = &http.Client{}
	}

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
	r["protocol"] = res.Proto

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

	// XXX
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

func http1_1Transport() *http.Transport {
	transport := &http.Transport{
		DialContext:       (&net.Dialer{}).DialContext,
		ForceAttemptHTTP2: false,
	}
	return transport
}

func http3RoundTripper() *http3.RoundTripper {
	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}
	var qconf quic.Config
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
		QuicConfig: &qconf,
	}
	defer roundTripper.Close()
	return roundTripper
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
