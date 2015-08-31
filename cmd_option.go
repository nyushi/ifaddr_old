package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

func ParseCmdOptions() (*CmdOption, error) {
	opts := &CmdOption{}
	args, err := flags.ParseArgs(opts, os.Args)
	if err != nil {
		return nil, err
	}
	if err := opts.prepare(); err != nil {
		return nil, err
	}
	opts.Args = args
	return opts, nil
}

type CmdOption struct {
	Args []string

	IFName           string `short:"i" long:"interface"`
	Netmask          string `short:"m" long:"netmask"`
	First            bool   `short:"1" long:"first"`
	EnableRegexp     bool   `short:"E" long:"regexp"`
	IncludeIPv6      []bool `long:"include-ipv6"`
	ExcludeIPv6      []bool `long:"exclude-ipv6"`
	IncludeIPv4      []bool `long:"include-ipv4"`
	ExcludeIPv4      []bool `long:"exclude-ipv4"`
	IncludeLinkLocal []bool `long:"include-linklocal"`
	ExcludeLinkLocal []bool `long:"exclude-linklocal"`
	IncludeLoopback  []bool `long:"include-loopback"`
	ExcludeLoopback  []bool `long:"exclude-loopback"`

	// shortcuts
	All      []bool `short:"a" long:"all"`
	OnlyIPv6 []bool `short:"6" long:"only-ipv6"`
	OnlyIPv4 []bool `short:"4" long:"only-ipv4"`

	// results
	NeedIPv6      bool
	NeedIPv4      bool
	NeedLinkLocal bool
	NeedLoopback  bool
	IPNet         *net.IPNet
}

func (c *CmdOption) parseBool(input []bool) bool {
	if len(input) == 0 {
		return false
	}
	return true
}

func (c *CmdOption) prepare() error {
	values := map[string]bool{
		"include-ipv6":      c.parseBool(c.IncludeIPv6),
		"include-ipv4":      c.parseBool(c.IncludeIPv4),
		"include-linklocal": c.parseBool(c.IncludeLinkLocal),
		"include-loopback":  c.parseBool(c.IncludeLoopback),
		"exclude-ipv6":      c.parseBool(c.ExcludeIPv6),
		"exclude-ipv4":      c.parseBool(c.ExcludeIPv4),
		"exclude-linklocal": c.parseBool(c.ExcludeLinkLocal),
		"exclude-loopback":  c.parseBool(c.ExcludeLoopback),
		"only-ipv6":         c.parseBool(c.OnlyIPv6),
		"only-ipv4":         c.parseBool(c.OnlyIPv4),
		"all":               c.parseBool(c.All),
	}
	conflicts := [][]string{
		{"only-ipv6", "only-ipv4"},
		{"only-ipv6", "include-ipv4"},
		{"only-ipv4", "include-ipv6"},
		{"only-ipv6", "exclude-ipv6"},
		{"only-ipv4", "exclude-ipv4"},
		{"include-ipv6", "exclude-ipv6"},
		{"include-ipv4", "exclude-ipv4"},
		{"include-linklocal", "exclude-linklocal"},
		{"include-loopback", "exclude-loopback"},
		{"all", "exclude-ipv6"},
		{"all", "exclude-ipv4"},
		{"all", "exclude-linklocal"},
		{"all", "exclude-loopback"},
	}
	errstr := []string{}

	for _, keys := range conflicts {
		b := true
		for _, key := range keys {
			b = b && values[key]
		}
		if b {
			errstr = append(errstr, fmt.Sprintf("conflict: %s", strings.Join(keys, "/")))
		}
	}
	if c.Netmask != "" {
		var err error
		_, c.IPNet, err = net.ParseCIDR(c.Netmask)
		if err != nil {
			errstr = append(errstr, fmt.Sprintf("invalid netmask: %s", err))
		}
	} else {
		c.IPNet = nil
	}
	if len(errstr) > 0 {
		return errors.New(strings.Join(errstr, ", "))
	}

	// default
	c.NeedIPv6 = true
	c.NeedIPv4 = true
	c.NeedLinkLocal = false
	c.NeedLoopback = false

	if values["only-ipv6"] {
		c.NeedIPv6 = true
		c.NeedIPv4 = false
	}
	if values["only-ipv4"] {
		c.NeedIPv6 = false
		c.NeedIPv4 = true
	}
	if values["all"] {
		c.NeedIPv6 = true
		c.NeedIPv4 = true
		c.NeedLinkLocal = true
		c.NeedLoopback = true
	}
	if values["include-linklocal"] {
		c.NeedLinkLocal = true
	}
	if values["include-loopback"] {
		c.NeedLoopback = true
	}
	return nil
}

func (c *CmdOption) Pattern() string {
	if len(c.Args) > 1 {
		return c.Args[1]
	}
	return ""
}
