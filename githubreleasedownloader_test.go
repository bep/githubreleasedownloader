package githubreleasedownloader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestGetRelease(t *testing.T) {
	c := qt.New(t)

	linuxTarGz := getAsset(c, "linux-amd64.tar.gz")

	c.Assert(linuxTarGz.Name, qt.Equals, "helloworld_0.5.0_linux-amd64.tar.gz")
	c.Assert(linuxTarGz.Sha256, qt.Equals, "ba9afe713027e4fc0be39856716d4341d10689c6a3447974296843aa147ea0bd")
}

func TestDownloadAndExtractAssetTarGz(t *testing.T) {
	c := qt.New(t)
	asset := getAsset(c, "linux-amd64.tar.gz")

	tempDir := c.TempDir()
	opts := DownloadAssetOptions{
		TargetBaseDir: tempDir,
	}
	res, err := DownloadAndExtractAsset(asset, opts)
	c.Assert(err, qt.IsNil)

	dirnames, err := os.ReadDir(res.TargetDir)
	c.Assert(err, qt.IsNil)
	// README.md, helloworld, helloworld_0.5.0_linux-amd64.tar.gz.ba9afe713027e4fc0be39856716d4341d10689c6a3447974296843aa147ea0bd
	c.Assert(dirnames, qt.HasLen, 3)
	s := ls(dirnames)
	c.Assert(s, qt.Equals, "-rw-r--r-- 0644 README.md\n-rwxr-xr-x 0755 helloworld\n-rw-r--r-- 0644 helloworld_0_5_0_linux-amd64_tar_gz_ba9afe713027e4fc0be39856716d4341d10689c6a3447974296843aa147ea0bd\n")
}

func TestDownloadAndExtractAssetZip(t *testing.T) {
	c := qt.New(t)
	asset := getAsset(c, ".zip")

	tempDir := c.TempDir()
	opts := DownloadAssetOptions{
		TargetBaseDir: tempDir,
	}
	res, err := DownloadAndExtractAsset(asset, opts)
	c.Assert(err, qt.IsNil)

	dirnames, err := os.ReadDir(res.TargetDir)
	c.Assert(err, qt.IsNil)
	c.Assert(dirnames, qt.HasLen, 3)
	s := ls(dirnames)
	// https://github.com/golang/go/issues/41809
	c.Assert(s, qt.Equals, "-rw-r--r-- 0644 README.md\n-rw-r--r-- 0644 helloworld.exe\n-rw-r--r-- 0644 helloworld_0_5_0_windows-amd64_zip_07e933f0153020c64c8b38d563e47c34ea5c27fe1312ac3b8e9861e9da98e245\n")
}

func TestDownloadAndExtractAssetToSubDirectory(t *testing.T) {
	c := qt.New(t)
	asset := getAsset(c, "linux-amd64.tar.gz")

	tempDir := c.TempDir()
	opts := DownloadAssetOptions{
		TargetBaseDir:         tempDir,
		ExtractToSubDirectory: true,
	}
	res, err := DownloadAndExtractAsset(asset, opts)
	c.Assert(err, qt.IsNil)

	c.Assert(filepath.Base(res.TargetDir), qt.Equals, "helloworld_0_5_0_linux-amd64_tar_gz_ba9afe713027e4fc0be39856716d4341d10689c6a3447974296843aa147ea0bd")

	dirnames, err := os.ReadDir(res.TargetDir)
	c.Assert(err, qt.IsNil)
	// README.md, helloworld
	c.Assert(dirnames, qt.HasLen, 2)

	s := ls(dirnames)
	c.Assert(s, qt.Equals, "-rw-r--r-- 0644 README.md\n-rwxr-xr-x 0755 helloworld\n")
}

func TestDownloadAndExtractAssetWithFilter(t *testing.T) {
	c := qt.New(t)
	asset := getAsset(c, "linux-amd64.tar.gz")

	tempDir := c.TempDir()
	opts := DownloadAssetOptions{
		TargetBaseDir: tempDir,
		ExtractFilter: func(name string, isDir bool) bool {
			return name == "README.md"
		},
	}
	res, err := DownloadAndExtractAsset(asset, opts)
	c.Assert(err, qt.IsNil)

	dirnames, err := os.ReadDir(res.TargetDir)
	c.Assert(err, qt.IsNil)
	c.Assert(dirnames, qt.HasLen, 1)
}

func getAsset(c *qt.C, suffix string) Asset {
	g, err := New()
	c.Assert(err, qt.IsNil)
	release, err := g.GetRelease("bep", "helloworld", "v0.5.0")
	c.Assert(err, qt.IsNil)
	c.Assert(release.Assets, qt.HasLen, 4)
	linuxTarGz := release.Assets.Filter(func(a Asset) bool {
		return strings.HasSuffix(a.Name, suffix)
	})[0]
	return linuxTarGz

}

func ls(dirs []os.DirEntry) string {
	var sb strings.Builder
	for _, d := range dirs {
		fi, err := d.Info()
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&sb, "%s %04o %s\n", fi.Mode(), fi.Mode().Perm(), fi.Name())
	}
	return sb.String()
}
