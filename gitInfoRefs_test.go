package main

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestGitInfoRefs(t *testing.T) {
	example := strings.Join([]string{
		"001e# service=git-upload-pack",
		"000000b3c7d3d3371baa35587fb66d8a79c6d999a4dafd8e HEAD\000multi_ack thin-pack side-band side-band-64k ofs-delta shallow no-progress include-tag multi_ack_detailed no-done agent=git/1.8.4",
		"003cd58c4a91450924a963d2cc7407dfa3e38866cb06 refs/heads/0.2",
		"004ca8fde08d941efc61322ae0302bf8bafc13e2275c refs/heads/add-version-to-join",
		"003fc7d3d3371baa35587fb66d8a79c6d999a4dafd8e refs/heads/master",
		"004448da4910b78e24d8d3a831839cc751700ddc6e10 refs/heads/update-docs",
		"003ee2f04208620f3bc9d01cc3fb216b92fa4e4a5767 refs/pull/1/head",
		"003f9b36b682ebbd7bd224b621fb90864821726b11b3 refs/pull/1/merge",
		"003f1eb0be10fe9ebf6e99a6c16abd3e583a68533dbd refs/pull/10/head",
		"0000",
	}, "\n")

	in := strings.NewReader(example)
	gup, err := parseGitUploadPack(ioutil.NopCloser(in))

	if err != nil {
		t.Fatalf("Failed parding git-upload-pack: %v", err)
	}

	// Stringification spits something similar back out...
	if len(example) != len(gup.String()) {
		t.Errorf(
			"Stringified doc should be length %v, got %v",
			len(example), len(gup.String()),
		)
		t.Errorf(
			"Expected String() to return \n%v\nGot:\n%v",
			example, gup.String(),
		)
	}
	/*
		if gup.String() != "001e# service=git-upload-pack\n0000" {
			t.Errorf("Didn't get the right output.")
		}
	*/

	if gup.capabilities != "multi_ack thin-pack side-band side-band-64k ofs-delta shallow no-progress include-tag multi_ack_detailed no-done agent=git/1.8.4" {
		t.Errorf(
			"Expected capabilities to be \n\t%v\nGot:\n\t%v",
			"multi_ack thin-pack side-band side-band-64k ofs-delta shallow no-progress include-tag multi_ack_detailed no-done agent=git/1.8.4",
			gup.capabilities,
		)
	}

	// HEAD
	if gup.refs["HEAD"] != "c7d3d3371baa35587fb66d8a79c6d999a4dafd8e" {
		t.Errorf(
			"Expected refs.HEAD to be \n\t%v\nGot:\n\t%v",
			"c7d3d3371baa35587fb66d8a79c6d999a4dafd8e",
			gup.refs["HEAD"],
		)
	}

	// Can find the tag "update-docs"
	// TODO: Table of stuff we should test here
	var tests = []struct {
		q string
		c string
	}{
		{"update-docs", "48da4910b78e24d8d3a831839cc751700ddc6e10"},
		{"9b36b682", "9b36b682ebbd7bd224b621fb90864821726b11b3"},
		{"1111111111111111111111111111111111111111", "1111111111111111111111111111111111111111"},
	}

	for _, tt := range tests {
		err, commit := gup.findCommitish(tt.q)
		if err != nil || commit != tt.c {
			t.Errorf("Could not find commitish '%v'; got %v and commit %v.", tt.q, err, commit)
		}
	}

	// Set commit at someting bogus blows up
	err = gup.SetHead("does-not-exist")
	if err == nil {
		t.Errorf("Expected SetHead(does-not-exist) to fail. It didn't.")
	}

	err = gup.SetMaster("update-docs")
	if err != nil {
		t.Errorf("Expected SetMaster(update-docs) to work, failed with %v", err)
	} else if gup.refs["refs/heads/master"] != "48da4910b78e24d8d3a831839cc751700ddc6e10" {
		t.Errorf(
			"Expected master to be at %v, but got %v",
			"004448da4910b78e24d8d3a831839cc751700ddc6e10",
			gup.refs["refs/heads/master"],
		)
	}
}
