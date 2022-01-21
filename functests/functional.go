/*
Copyright © 2019-2020 Netskope
*/

package functests

import (
	"flag"
)

const (
	// DisabledMsg occurs when functional testing is not enabled via the -functests flag
	DisabledMsg = "functional tests are not enabled"
	// EnabledMsg occurs when function testing is enabled via the -functests flag
	EnabledMsg = "functional tests are enabled"
)

var (
	functional = flag.Bool("functests", false, "run functional tests")
)

// Functional returns true when functional tests are enabled
func Functional() bool {
	flag.Parse()
	return *functional
}
