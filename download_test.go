package vspb

import "testing"

func TestGitClone(t *testing.T) {
	cases := []struct {
		Name    string
		Addr    string
		Version string
		Path    string
	}{
		{
			"clone master",
			"https://github.com/AkarinVS/L-SMASH-Works.git",
			"",
			"./download_tests/clone-master",
		},
		{
			"clone with tag",
			"https://github.com/AkarinVS/L-SMASH-Works.git",
			"vA.3j",
			"./download_tests/clone-with-tag",
		},
		{
			"clone with hash",
			"https://github.com/AkarinVS/L-SMASH-Works.git",
			"dffecae8ea055ae69a43894ce77960fc39c869d0",
			"./download_tests/clone-with-hash",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			err := gitClone(c.Addr, c.Version, c.Path)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGrabDownload(t *testing.T) {
	cases := []struct {
		Name string
		Addr string
		Path string
	}{
		{
			Name: "normal download",
			Addr: "https://www.python.org/ftp/python/3.11.1/Python-3.11.1.tar.xz",
			Path: "./download_tests",
		},
		{
			Name: "cover download",
			Addr: "https://www.python.org/ftp/python/3.11.1/Python-3.11.1.tar.xz",
			Path: "./download_tests",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			err := grabDownload(c.Addr, c.Path)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestIsGitHash(t *testing.T) {
	cases := []struct {
		Name     string
		Version  string
		Expected bool
	}{
		{
			"7",
			"d73d2dc",
			false,
		},
		{
			"40",
			"d73d2dc1675d2b72a632382db641897f50edbede",
			true,
		},
		{
			"8",
			"d73d2dc1",
			false,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if isGitHash(c.Version) != c.Expected {
				t.Errorf("expected %v, got %v", c.Expected, isGitHash(c.Version))
			}
		})
	}
}
