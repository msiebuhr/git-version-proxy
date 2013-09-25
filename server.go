package main

import (
	"fmt"
	"io"
	//"io/ioutil"
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

func main() {
	http.HandleFunc("/github/", func(w http.ResponseWriter, r *http.Request) {
		// Parse the path: /github/<username>/<commitish>/<project/subdir/...>
		parts := strings.SplitN(r.URL.Path, "/", 5)
		baseUrl := fmt.Sprintf("https://github.com/%s/%s", parts[2], parts[4])
		fullUrl := baseUrl
		if r.URL.RawQuery != "" {
			fullUrl = fmt.Sprintf("%v?%v", baseUrl, r.URL.RawQuery)
		}
		commitish := parts[3]

		// Is it for info/refs?
		// TODO: Should parse the "Smart" reply and only send the relevant REFs
		if strings.HasSuffix(parts[4], "info/refs") {
			// TODO: If we get a commit SHA, use that as REF / master

			// Fetch and parse the advanced info/refs from github and send a simple one down...
			res, _ := http.Get(fullUrl)
			//body, _ := ioutil.ReadAll(res.Body)
			body, _ := readPktLine(res.Body)

			// Go through the lines and filter out irelevant refs
			newBody := []string{}
			for _, line := range body {
				if strings.HasPrefix(line, "#") || strings.HasSuffix(line, commitish) || len(line) == 0 || len(strings.Fields(line)) > 2 {
					fmt.Println(commitish, line)
					newBody = append(newBody, line)
				}
			}

			// TODO: Look only for suffixes on the form heads/<commitish> and tags/<commitish>.

			// TODO: We need to treat the line with <sha> HEAD<gitstuff>
			// specially - namely, to seems we have to use whatever <sha> we
			// need to represent as the HEAD.
			//
			// I suspect it does this instead of making a separate request for
			// /HEAD on the repository

			copyHeaders(res.Header, w.Header())
			writePktLine(w, newBody)
			return
		}

		/*
			// Is it for HEAD?
			if strings.HasSuffix(parts[4], "HEAD") {
				fmt.Fprint(w, "ref: refs/heads/master\n");
				return;
			}
		*/

		fmt.Println(fullUrl, commitish)

		// Create a new request and send it off
		client := &http.Client{}
		req, _ := http.NewRequest(r.Method, fullUrl, r.Body)
		copyHeaders(r.Header, req.Header)
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			http.Error(w, fmt.Sprintf("Proxy error: %v", err), 500)
			return
		}

		fmt.Println(res.Status, res.ContentLength)

		// Send back response headers
		copyHeaders(res.Header, w.Header())
		w.WriteHeader(res.StatusCode)

		// Copy over response
		io.Copy(w, res.Body)

		//fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
