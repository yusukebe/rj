package rj

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type response struct {
	Code    int                `json:"code"`
	Status  string             `json:"status"`
	Headers map[string]string  `json:"header"`
	Timing  map[string]float64 `json:"timing"`
}

func TestRequest(t *testing.T) {

	buffer := &bytes.Buffer{}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("X-Custom", "Foo")
		fmt.Fprintln(w, "Hello, rj")
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	param := param{
		method:       "GET",
		outputWriter: buffer,
	}

	request(ts.URL, param)

	output := buffer.String()

	var resp response
	err := json.Unmarshal([]byte(output), &resp)
	if err != nil {
		t.Fatalf("Can't decode JSON")
	}

	if resp.Code != 200 {
		t.Errorf("Expected: %d, Actual: %d", 200, resp.Code)
	}

	if resp.Headers["content-type"] != "text/plain" {
		t.Errorf("Expected: %s, Actual: %s", "text/plain", resp.Headers["content-type"])
	}

	if resp.Headers["x-custom"] != "Foo" {
		t.Errorf("Expected: %s, Actual: %s", "Foo", resp.Headers["x-custom"])
	}

	if resp.Timing["total"] == 0 {
		t.Errorf("timing.total will be greater than 0")
	}

}

func TestNewCmd(t *testing.T) {
	buffer := new(bytes.Buffer)
	cmd := newCmd()
	cmd.SetOut(buffer)

	cmd.SetArgs([]string{"--version"})
	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	actual := strings.TrimSpace(buffer.String())
	expect := "rj version " + Version
	if actual != expect {
		t.Errorf("Expected: %s, Actual: %s", expect, actual)
	}
}
