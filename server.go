package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func copyHeaders(from, to http.Header) {
	for header, items := range from {
		for _, item := range items {
			to.Add(header, item)
		}
	}
}

// One of the path elements will start and end with @. If we strip that element
// out, we will get a path + tag/whatever.
// ex: github.com/msiebuhr/@master/foo.git
//
// TODO: We should probably have a generic parser that results in (VCS, host,
// commitish)
func splitPathAndCommitish(path string) (string, string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	c := ""
	p_arr := make([]string, 0, len(parts))

	for _, part := range parts {
		if strings.HasPrefix(part, "@") {
			c = part
		} else {
			p_arr = append(p_arr, part)
		}
	}

	return strings.Join(p_arr, "/"), strings.Trim(c, "@")
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		// Is it a go-get request? And why should I care?

		// TODO: Check upstream exists!
		w.WriteHeader(200)
		w.Write([]byte("<html><head>\n"))
		fmt.Fprintf(
			w,
			"<meta name=\"go-import\" content=\"127.0.0.1:8080%s git http://127.0.0.1:8080/_git%s\"></meta>\n",
			r.URL.Path,
			r.URL.Path,
		)
		w.Write([]byte("</head><body>foobar</body></html>"))
	})

	// Magic GIT imports
	http.HandleFunc("/_git/github.com/", func(w http.ResponseWriter, r *http.Request) {
		path, commitish := splitPathAndCommitish(strings.TrimPrefix(r.URL.Path, "/_git"))
		baseUrl := fmt.Sprintf("https://%s", path)
		fullUrl := baseUrl
		if r.URL.RawQuery != "" {
			fullUrl = fmt.Sprintf("%v?%v", baseUrl, r.URL.RawQuery)
		}

		fmt.Println(fullUrl, commitish)

		// Create a new request and send it off
		client := &http.Client{}
		req, _ := http.NewRequest(r.Method, fullUrl, r.Body)
		copyHeaders(r.Header, req.Header)
		res, err := client.Do(req)
		if err != nil {
			//fmt.Println(err)
			http.Error(w, fmt.Sprintf("Proxy error: %v", err), 500)
			return
		}

		// If if it is an info/refs thing, then we want to modify the body before it goes back
		if strings.HasSuffix(path, "info/refs") {
			body, _ := parseGitUploadPack(res.Body)

			err := body.SetMaster(commitish)

			if (err != nil) {
				fmt.Println("ERROR:", err);
				w.WriteHeader(404);
			}

			// Send back response headers
			copyHeaders(res.Header, w.Header())
			w.WriteHeader(res.StatusCode)

			w.Write([]byte(body.String()))
			//w.Close()
		} else {
			// Copy over response
			copyHeaders(res.Header, w.Header())
			w.WriteHeader(res.StatusCode)
			io.Copy(w, res.Body)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
