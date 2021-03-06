package dns

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/datawire/teleproxy/internal/pkg/tpu"
)

type searchDomains struct {
	Interface string
	Domains   string
}

func OverrideSearchDomains(domains string) func() {
	if runtime.GOOS != "darwin" {
		return func() {}
	}

	ifaces, _ := getIfaces()
	previous := []searchDomains{}

	for _, iface := range ifaces {
		// setup dns search path
		domain, _ := getSearchDomains(iface)
		setSearchDomains(iface, domains)
		previous = append(previous, searchDomains{iface, domain})
	}

	// return function to restore dns search paths
	return func() {
		for _, prev := range previous {
			setSearchDomains(prev.Interface, prev.Domains)
		}
	}
}

func getIfaces() (ifaces []string, err error) {
	lines, err := tpu.ShellLogf("networksetup -listallnetworkservices | fgrep -v '*'", log)
	if err != nil {
		return
	}
	for _, line := range strings.Split(lines, "\n") {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			ifaces = append(ifaces, line)
		}
	}
	return
}

func getSearchDomains(iface string) (domains string, err error) {
	domains, err = tpu.ShellLogf(fmt.Sprintf("networksetup -getsearchdomains '%s'", iface), log)
	domains = strings.TrimSpace(domains)
	return
}

func setSearchDomains(iface, domains string) (err error) {
	_, err = tpu.ShellLogf(fmt.Sprintf("networksetup -setsearchdomains '%s' '%s'", iface, domains), log)
	return
}
