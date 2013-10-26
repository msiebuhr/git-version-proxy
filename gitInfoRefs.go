package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type GitUploadPack struct {
	refs         map[string]string
	head         string
	capabilities string
}

func NewGitUploadPack() *GitUploadPack {
	return &GitUploadPack{refs: make(map[string]string)}
}

func parseGitUploadPack(r io.ReadCloser) (*GitUploadPack, error) {
	p := NewGitUploadPack()

	res, e := readPktLine(r)
	if e != nil {
		return nil, e
	}

	// First pack should be "001e# service=git-upload-pack"
	// Second pack is empty...
	// Last pack is empty too...
	refs := res[2:]
	// Thrid shoud have a standard "SHA ref\0capabilities"
	caps := strings.SplitN(refs[0], "\000", 2)
	if len(caps) == 2 {
		p.capabilities = caps[1]
		refs[0] = caps[0]
	}
	// Rest is "SHA ref"
	for _, elem := range refs {
		parts := strings.SplitN(elem, " ", 2)
		if len(parts) == 2 && len(parts[0]) == 40 {
			p.refs[parts[1]] = parts[0]
		} else {
			fmt.Printf("ARGH! Could not figure out '%v'.\n", elem)
		}
	}

	return p, nil
}

func writePktLine(line string) string {
	return fmt.Sprintf("%04x%s", len(line)+4, line)
}

func (p *GitUploadPack) String() string {
	out := []string{
		"001e# service=git-upload-pack\n0000",
		writePktLine(fmt.Sprintf("%s HEAD\000%s\n", p.refs["HEAD"], p.capabilities)),
	}

	// Write everything else
	for ref, commit := range p.refs {
		if ref != "HEAD" {
			out = append(out, writePktLine(fmt.Sprintf("%s %s\n", commit, ref)))
		}
	}

	out = append(out, "0000")

	return strings.Join(out, "")
}

func (p *GitUploadPack) findCommitish(commitish string) (error, string) {
	// If it is a commit-ID, we should just return that
	if len(commitish) == 40 {
		return nil, commitish
	}

	// Look through the refs
	for ref, commit := range p.refs {
		if strings.HasPrefix(commit, commitish) || strings.HasSuffix(ref, commitish) {
			return nil, commit
		}
	}

	return errors.New("Commitish not found"), ""
}

func (p *GitUploadPack) SetMaster(commitish string) error {
	err, commit := p.findCommitish(commitish)
	if err != nil {
		return err
	}
	p.refs["refs/heads/master"] = commit
	return nil
}
