package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"strings"

	"github.com/armon/go-socks5"
)

// ipAllowed checks if client IP is within allowed CIDRs
func ipAllowed(ip net.IP, networks []*net.IPNet) bool {
	if len(networks) == 0 {
		return true
	}
	for _, n := range networks {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

func main() {
	// Proxy port (default 1080)
	port := os.Getenv("PROXY_PORT")
	if port == "" {
		port = "1080"
	}

	// Parse ALLOWED_IPS (comma-separated CIDRs)
	allowedCIDRs := strings.Split(os.Getenv("ALLOWED_IPS"), ",")
	var networks []*net.IPNet
	for _, cidr := range allowedCIDRs {
		cidr = strings.TrimSpace(cidr)
		if cidr == "" {
			continue
		}
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Printf("[ERR] socks: invalid CIDR in ALLOWED_IPS: %s", cidr)
			continue
		}
		networks = append(networks, network)
	}

	// Optional: authentication
	user := os.Getenv("PROXY_USER")
	pass := os.Getenv("PROXY_PASSWORD")

	conf := &socks5.Config{}

	// Custom Dialer: block IPv6 dests and force IPv4 dialing
	conf.Dial = func(ctx context.Context, network, addr string) (net.Conn, error) {
		// If dest is an IPv6 literal, block with a clean log message
		if host, _, err := net.SplitHostPort(addr); err == nil {
			if ip := net.ParseIP(host); ip != nil && ip.To4() == nil {
				log.Printf("[ERR] socks: blocked IPv6 destination: %s", addr)
				return nil, errors.New("IPv6 not supported")
			}
		}
		// Force IPv4 even for hostnames (avoid v6 resolution/usage)
		switch network {
		case "tcp":
			network = "tcp4"
		case "udp":
			network = "udp4"
		}
		return net.Dial(network, addr)
	}

	if user != "" && pass != "" {
		creds := socks5.StaticCredentials{user: pass}
		cator := socks5.UserPassAuthenticator{Credentials: creds}
		conf.AuthMethods = []socks5.Authenticator{cator}
	}

	server, err := socks5.New(conf)
	if err != nil {
		log.Fatalf("[ERR] socks: failed to create server: %v", err)
	}

	listenAddr := ":" + port
	log.Printf("[INFO] socks: starting proxy on %s", listenAddr)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("[ERR] socks: failed to bind on %s: %v", listenAddr, err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("[ERR] socks: accept failed: %v", err)
			continue
		}

		remoteIP := conn.RemoteAddr().(*net.TCPAddr).IP

		// Block IPv6 clients
		if remoteIP.To4() == nil {
			log.Printf("[ERR] socks: blocked IPv6 client: %s", remoteIP)
			conn.Close()
			continue
		}

		// Check ALLOWED_IPS
		if !ipAllowed(remoteIP, networks) {
			log.Printf("[ERR] socks: connection from disallowed IP: %s", remoteIP)
			conn.Close()
			continue
		}

		log.Printf("[INFO] socks: connection from allowed IP address: %s", remoteIP)
		go server.ServeConn(conn)
	}
}
