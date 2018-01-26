package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const metadataBase = "http://169.254.169.254/latest/meta-data/"
const defaultPortSpec = ":8080"

var options struct {
	portspec string
}

func init() {
	flag.StringVar(&options.portspec, "port", defaultPortSpec, "port to listen on for HTTP requests")
}

func addSection(w io.Writer, path string) error {
	u := metadataBase + path
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "\n<h3>%s</h3>\n", template.HTMLEscapeString(path))
	if err == nil {
		template.HTMLEscape(w, body)
		// no error return, oops
	}
	return err
}

func AddSection(w io.Writer, path string) {
	if err := addSection(w, path); err != nil {
		fmt.Fprintf(w, "\n<h3 class=\"error\">%s</h3>\n<div class=\"error errmsg\">%s</div>\n",
			template.HTMLEscapeString(path), template.HTMLEscapeString(err.Error()))
	}
}

func rootHandle(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "<html><head><title>AWS Info Dumper</title></head><body><h1>AWS Info Dumper</h1>\n")
	for _, section := range []string{
		"hostname",
		"placement/availability-zone",
		"iam/info",
	} {
		AddSection(w, section)
	}
}

func parseFlagsSanely() {
	flag.Parse()
	if options.portspec == "" {
		options.portspec = defaultPortSpec
	}
	if !strings.Contains(options.portspec, ":") {
		options.portspec = ":" + options.portspec
	}
}

func main() {
	parseFlagsSanely()
	http.HandleFunc("/", rootHandle)
	log.Fatal(http.ListenAndServe(options.portspec, nil))
}
