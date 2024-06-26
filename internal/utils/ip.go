package utils

import (
	"net"
	"strings"
)

// IsLocalIP 判断是否为本地IP
func IsLocalIP(ipString string) bool {
	var ip = net.ParseIP(ipString)
	if ip == nil {
		return false
	}

	// IPv6
	if strings.Contains(ipString, ":") {
		return ip.String() == "::1"
	}

	// IPv4
	ip = ip.To4()
	if ip == nil {
		return false
	}
	if ip[0] == 127 ||
		ip[0] == 10 ||
		(ip[0] == 172 && ip[1]&0xf0 == 16) ||
		(ip[0] == 192 && ip[1] == 168) {
		return true
	}

	return false
}

// IsIPv4 是否为IPv4
func IsIPv4(ip string) bool {
	var data = net.ParseIP(ip)
	if data == nil {
		return false
	}
	if strings.Contains(ip, ":") {
		return false
	}
	return data.To4() != nil
}

// IsIPv6 是否为IPv6
func IsIPv6(ip string) bool {
	var data = net.ParseIP(ip)
	if data == nil {
		return false
	}
	return !IsIPv4(ip)
}

// IsWildIP 宽泛地判断一个数据是否为IP
func IsWildIP(v string) bool {
	var l = len(v)
	if l == 0 {
		return false
	}

	// for [IPv6]
	if v[0] == '[' && v[l-1] == ']' {
		return IsWildIP(v[1 : l-1])
	}

	return net.ParseIP(v) != nil
}
