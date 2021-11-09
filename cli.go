package rj

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type param struct {
	method    string
	userAgent string
	headers   []string
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
		param := param{
			method:    method,
			userAgent: userAgent,
			headers:   headers,
		}
		request(url, param)
	},
}

func init() {
	rootCmd.Flags().StringP("method", "X", "GET", "HTTP Request method")
	rootCmd.Flags().StringP("agent", "A", "rj/v0.0.1", "User-Agent name")
	rootCmd.Flags().StringArrayP("header", "H", nil, "HTTP Request Header")
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
			req.Header.Set(kv[0], kv[1])
		}
	}

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer res.Body.Close()

	r := make(map[string]interface{})

	r["status"] = res.Status
	r["code"] = res.StatusCode

	headers := make(map[string]interface{})

	for key, value := range res.Header {
		headkey := strings.ToLower(key)
		headValue := strings.Join(value, ", ") // XXX
		headers[headkey] = headValue
	}

	r["header"] = headers

	bytes, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(bytes))
}
