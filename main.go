package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"

	"github.com/Sirupsen/logrus"
)

var (
	Version    string
	CommitHash string
)

type address struct {
	net.IP
}

func (a *address) IsV4() bool {
	return a.To4() != nil
}

func (a *address) IsV6() bool {
	return !a.IsV4()
}

func getAddrOf(ifname string) (addrs []*address, err error) {
	var interfaces []net.Interface
	if ifname == "" {
		interfaces, err = net.Interfaces()
	} else {
		i, err := net.InterfaceByName(ifname)
		if err != nil {
			return nil, err
		}
		interfaces = append(interfaces, *i)
	}
	for _, i := range interfaces {
		rawAddrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, a := range rawAddrs {
			ipn, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			addrs = append(addrs, &address{ipn.IP})
		}
	}
	return addrs, nil
}

func getSuitableAddrs(addrs []*address, v4, v6, linklocal, loopback bool, re *regexp.Regexp, mask *net.IPNet) ([]*address, error) {
	ret := []*address{}
	for _, a := range addrs {
		if a.IsLoopback() && !loopback {
			continue
		}
		if !v6 && a.IsV6() {
			continue
		}
		if !v4 && a.IsV4() {
			continue
		}
		if !linklocal && a.IsLinkLocalUnicast() {
			continue
		}
		if !loopback && a.IsLoopback() {
			continue
		}
		if re != nil {
			if !re.MatchString(a.String()) {
				continue
			}
		}
		if mask != nil {
			if !mask.Contains(a.IP) {
				continue
			}
		}
		ret = append(ret, a)
	}
	if len(ret) == 0 {
		return nil, errors.New("unable to find suitable address")
	}
	return ret, nil
}

func main() {
	opts, err := ParseCmdOptions()
	if err != nil {
		os.Exit(64)
	}
	if opts.VersionFlag {
		fmt.Printf("%s(%s)\n", Version, CommitHash)
		os.Exit(0)
	}
	addrs, err := getAddrOf(opts.IFName)
	if err != nil {
		logrus.Fatal(err)
	}
	var r *regexp.Regexp
	pattern := opts.Positional.Pattern
	if pattern != "" {
		if !opts.EnableRegexp {
			pattern = regexp.QuoteMeta(pattern)
		}
		var err error
		r, err = regexp.Compile(pattern)
		if err != nil {
			logrus.Fatal(err)
		}

	}
	suitables, err := getSuitableAddrs(addrs, opts.NeedIPv4, opts.NeedIPv6, opts.NeedLinkLocal, opts.NeedLoopback, r, opts.IPNet)
	if err != nil {
		logrus.Fatal(err)
	}

	for _, a := range suitables {
		fmt.Println(a.String())
		if opts.First {
			break
		}
	}
}
