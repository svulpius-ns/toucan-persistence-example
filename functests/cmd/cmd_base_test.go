/*
Copyright © 2019-2020 Netskope
*/

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/netskope/piratetreasure/functests"
)

var (
	compiledPath string
	packagePath  = "github.com/netskope/piratetreasure/cmd/piratetreasure"
)

// execApp invokes the app with the args specified.
func execApp(args ...string) (string, error) {
	command := exec.Command(compiledPath, args...)
	session, err := gexec.Start(command, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	session.Wait()

	stdout := string(session.Out.Contents())

	return stdout, err
}

// TestMain runs all of the tests defined.
func TestMain(m *testing.M) {
	if functests.Functional() {
		fmt.Println(functests.EnabledMsg)
	} else {
		fmt.Println(functests.DisabledMsg)
		os.Exit(0)
	}

	before()
	rc := m.Run()
	after()
	os.Exit(rc)
}

func before() {
	gomega.RegisterFailHandler(ginkgo.Fail)
	var err error
	compiledPath, err = gexec.Build(packagePath)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
}

func after() {
	gexec.CleanupBuildArtifacts()
}
