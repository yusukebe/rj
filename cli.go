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
	method string
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
		param := param{
			method: method,
		}
		request(url, param)
	},
}

func init() {
	rootCmd.Flags().StringP("method", "X", "GET", "HTTP Request method")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func request(url string, param param) {
	req, _ := http.NewRequest("GET", url, nil)
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer res.Body.Close()

	kv := make(map[string]string)

	for key, value := range res.Header {
		headkey := strings.ToLower(key)
		headValue := strings.Join(value, ", ")
		kv[headkey] = headValue
	}

	bytes, err := json.Marshal(kv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(bytes))
}
