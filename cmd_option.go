package main

import (
	"errors"
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
	All []bool `short:"a" long:"all"`

	// results
	NeedIPv6      bool
	NeedIPv4      bool
	NeedLinkLocal bool
	NeedLoopback  bool
}

func (c *CmdOption) parseBool(input []bool) bool {
	if len(input) == 0 {
		return false
	}
	return true
}

func (c *CmdOption) prepare() error {
	errstr := []string{}
	incV6 := c.parseBool(c.IncludeIPv6)
	incV4 := c.parseBool(c.IncludeIPv4)
	incLl := c.parseBool(c.IncludeLinkLocal)
	incLb := c.parseBool(c.IncludeLoopback)
	excV6 := c.parseBool(c.ExcludeIPv6)
	excV4 := c.parseBool(c.ExcludeIPv4)
	excLl := c.parseBool(c.ExcludeLinkLocal)
	excLb := c.parseBool(c.ExcludeLoopback)

	if incV6 && excV6 {
		errstr = append(errstr, "conflict: include/exclude ipv6")
	}
	if incV4 && excV4 {
		errstr = append(errstr, "conflict: include/exclude ipv4")
	}
	if incLl && excLl {
		errstr = append(errstr, "conflict: include/exclude linklocal")
	}
	if incLb && excLb {
		errstr = append(errstr, "conflict: include/exclude loopback")
	}
	if len(errstr) > 0 {
		return errors.New(strings.Join(errstr, ", "))
	}

	if c.parseBool(c.All) {
		c.NeedIPv6 = true
		c.NeedIPv4 = true
		c.NeedLinkLocal = true
		c.NeedLoopback = true
		return nil
	}
	if incV6 || (!incV6 && !excV6) {
		c.NeedIPv6 = true
	}
	if incV4 || (!incV4 && !excV4) {
		c.NeedIPv4 = true
	}
	if incLl || (!incLl && !excLl) {
		c.NeedLinkLocal = false
	}
	if incLb || (!incLb && !excLb) {
		c.NeedLoopback = false
	}

	return nil
}

func (c *CmdOption) Pattern() string {
	if len(c.Args) > 1 {
		return c.Args[1]
	}
	return ""
}
