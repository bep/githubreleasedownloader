package githubreleasedownloader

import (
	"os"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestGetRelease(t *testing.T) {
	if os.Getenv("CI") != "" {
		// TODO(bep) fix this (GH token).
		t.Skip("skipping test in CI")
	}
	c := qt.New(t)

	g, err := New()
	c.Assert(err, qt.IsNil)
	release, err := g.GetRelease("gohugoio", "hugoreleaser", "v0.56.1")
	c.Assert(err, qt.IsNil)
	c.Assert(release.Assets, qt.HasLen, 5)
	linuxTarGz := release.Assets.Filter(func(a Asset) bool {
		return strings.HasSuffix(a.Name, "linux-amd64.tar.gz")
	})[0]
	c.Assert(linuxTarGz.Name, qt.Equals, "hugoreleaser_0.56.1_linux-amd64.tar.gz")
	c.Assert(linuxTarGz.Sha256, qt.Equals, "75e6973223415b15505a585556186c47d08daff8dd2b49c0260205ae660e01ce")
}
