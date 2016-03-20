// ClueGetter - Does things with mail
//
// Copyright 2016 Dolf Schimmel, Freeaqingme.
//
// This Source Code Form is subject to the terms of the two-clause BSD license.
// For its contents, please refer to the LICENSE file.
//
package lua

import (
	"github.com/miekg/dns"
	"github.com/yuin/gopher-lua"

	"errors"
	"fmt"
	"net"
	"strings"
)

func DnsLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), dnsExports)
	L.Push(mod)
	return 1
}

var dnsExports = map[string]lua.LGFunction{
	"queryTxt": dnsQueryTxt,
}

func dnsQuery(L *lua.LState, query string, qType uint16) ([]dns.RR, error) {
	config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
	c := new(dns.Client)

	m := new(dns.Msg)
	m.RecursionDesired = true
	m.SetQuestion(dns.Fqdn(query), qType)

	r, _, err := c.Exchange(m, net.JoinHostPort(config.Servers[0], config.Port))
	if r == nil {
		return nil, err
	}

	if r.Rcode != dns.RcodeSuccess {
		return nil,
			errors.New(fmt.Sprintf("invalid answer name %s after MX query for '%s'", query, query))
	}

	return r.Answer, nil
}

func dnsQueryTxt(L *lua.LState) int {
	query := L.ToString(1)

	res, err := dnsQuery(L, query, dns.TypeTXT)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	ret := L.NewTable()
	for _, a := range res {
		ret.Append(lua.LString(strings.Join(a.(*dns.TXT).Txt, "")))
	}

	L.Push(ret)
	return 1
}
