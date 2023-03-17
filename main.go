package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const metadataBase = "http://169.254.169.254/latest/meta-data/"
const defaultPortSpec = ":8080"

// AWS metadata service is very local, we can hard-timeout much sooner
const awsHTTPTimeout = 2 * time.Second

var options struct {
	portspec string
}

func init() {
	flag.StringVar(&options.portspec, "port", defaultPortSpec, "port to listen on for HTTP requests")
}

func addSection(ctx context.Context, w io.Writer, path string) error {
	u := metadataBase + path
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
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

func showError(w io.Writer, path string, err error) {
	fmt.Fprintf(w, "\n<h3 class=\"error\">%s</h3>\n<div class=\"error errmsg\">%s</div>\n",
		template.HTMLEscapeString(path), template.HTMLEscapeString(err.Error()))
}

func AddSection(ctx context.Context, w io.Writer, path string) {
	if err := addSection(ctx, w, path); err != nil {
		showError(w, path, err)
	}
}

func rootHandle(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "<html><head><title>AWS Info Dumper</title></head><body><h1>AWS Info Dumper</h1>\n")
	childCtx, cancel := context.WithTimeout(req.Context(), awsHTTPTimeout)
	defer cancel()
	if p := os.Getenv("ECS_CONTAINER_METADATA_FILE"); p != "" {
		io.WriteString(w, "<h2>ECS metadata from file</h2>\n")
		contents, err := os.ReadFile(p)
		if err != nil {
			showError(w, p, err)
		} else {
			fmt.Fprintf(w, "\n<h3>%s</h3>\n", template.HTMLEscapeString(p))
			template.HTMLEscape(w, contents)
		}
	} else {
		io.WriteString(w, "<h2>AWS metadata service (HTTP requests)</h2>\n")
		for _, section := range []string{
			"hostname",
			"placement/availability-zone",
			"iam/info",
		} {
			AddSection(childCtx, w, section)
			if childCtx.Err() != nil {
				// any context expiration has _almost_ certainly been shown in the output of the
				// AddSection error handling; there's a few nanoseconds race, so rather than
				// risk aborting early without saying so, just explicitly say "hey we're done".
				showError(w, "timeout", fmt.Errorf("terminated early"))
				break
			}
		}
	}
}

func parseFlagsSanely() {
	envPort := os.Getenv("PORT")
	if envPort != "" {
		options.portspec = envPort
	}
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
