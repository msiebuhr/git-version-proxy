package main

import (
	"io"
	"strconv"
	"fmt"
	"strings"
)

func readPktLine(r io.ReadCloser) ([]string, error) {
	defer r.Close()
	out := []string{}
	offsetBytes := make([]byte, 4)
	for {
		// Get and parse offset
		_, err := r.Read(offsetBytes)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		offset, _ := strconv.ParseUint(string(offsetBytes), 16, 16)
		if offset == 0 {
			out = append(out, "")
			continue
		}
		//if (offsetBytes == []byte{0, 0, 0, 4}) { continue }

		// Read remainer
		rest := make([]byte, offset-4)
		_, err = r.Read(rest)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		out = append(out, strings.TrimSpace(string(rest)))
	}
	return out, nil
}

func writePktLine(w io.Writer, d []string) error {
	for _, line := range d {
		if len(line) == 0 {
			fmt.Fprintf(w, "0000")
			continue
		}
		fmt.Printf("Encoding: %04x%s\n", len(line)+4, line)
		fmt.Fprintf(w, "%04x%s", len(line)+4, line)
	}
	return nil
}
