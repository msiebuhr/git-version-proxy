package main

import (
	"testing"
)

func TestSplitPathAndCommitish(t *testing.T) {
	var tests = []struct {
		in string
		op string
		oc string
	}{
		{"/github.com/foo/@master/bar.git", "github.com/foo/bar.git", "master"},
		{"/github.com/foo/bar.git/@master", "github.com/foo/bar.git", "master"},
		{"/github.com/foo/bar@master", "github.com/foo/bar", "master"},
		{"/github.com/coreos/etcd@v0.1.0", "github.com/coreos/etcd", "v0.1.0"},
	}

	for _, tt := range tests {
		p, c := splitPathAndCommitish(tt.in)

		if p != tt.op || c != tt.oc {
			t.Errorf(
				"Expected %v to have path %v and commitish %v, got %v and %v",
				tt.in, tt.op, tt.oc, p, c,
			)
		}
	}
}
