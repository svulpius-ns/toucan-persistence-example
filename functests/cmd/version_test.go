/*
Copyright © 2019 Netskope
*/

package cmd

import (
	"fmt"
	"testing"

	"github.com/onsi/gomega"

	"github.com/netskope/piratetreasure/internal/build"
)

var (
	expectedVersionShort = fmt.Sprintf("%s version %s (%s)",
		build.AppName,
		build.Version,
		build.GitSha,
	)
	expectedVersionLong = fmt.Sprintf("%s:\n Version: %s\n GitSha: %s\n Built: %s by %s",
		build.AppName,
		build.Version,
		build.GitSha,
		build.BuildTime,
		build.BuiltBy,
	)
)

func TestAppVersionFlags(t *testing.T) {
	stdout, err := execApp("--version")

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(stdout).Should(gomega.ContainSubstring(expectedVersionShort))
}

func TestAppVersionCmd(t *testing.T) {
	stdout, err := execApp("version")

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(stdout).Should(gomega.ContainSubstring(expectedVersionLong))
}
