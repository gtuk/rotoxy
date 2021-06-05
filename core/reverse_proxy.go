package core

import (
	"context"
	"fmt"
	"github.com/eahydra/socks"
	"github.com/elazarl/goproxy"
	"math/rand"
	"net"
	"net/http"
)

type ReverseProxy struct {
	port            *int
	proxy           *goproxy.ProxyHttpServer
	upstreamProxies *[]TorProxy
}

type upstreamDialer struct {
	forwardDialers []socks.Dialer
}

func (m *ReverseProxy) Start(proxies []TorProxy, port int) error {
	m.port = &port
	m.upstreamProxies = &proxies

	var router socks.Dialer
	proxy := goproxy.NewProxyHttpServer()
	proxy.Tr.DisableKeepAlives = true

	proxyUrls := make([]string, 0)
	for _, proxy := range proxies {
		proxyUrls = append(proxyUrls, fmt.Sprintf("127.0.0.1:%d", *proxy.ProxyPort))
	}

	router, err := buildUpstreamRouter(proxyUrls)
	if err != nil {
		return err
	}

	proxy.ConnectDial = func(network, address string) (net.Conn, error) {
		return router.Dial(network, address)
	}
	proxy.Tr.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		return router.Dial(network, address)
	}

	m.proxy = proxy

	httpListen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	err = http.Serve(httpListen, proxy)
	if err != nil {
		return err
	}

	return nil
}

func newUpstreamDialer(forwardDialers []socks.Dialer) *upstreamDialer {
	return &upstreamDialer{
		forwardDialers: forwardDialers,
	}
}

func (u *upstreamDialer) getRandomDialer() socks.Dialer {
	max := len(u.forwardDialers)
	if max == 0 {
		return u.forwardDialers[0]
	}
	randomDialer := 0 + rand.Intn(max-0)
	return u.forwardDialers[randomDialer]
}

func (u *upstreamDialer) Dial(network, address string) (net.Conn, error) {
	router := u.getRandomDialer()
	conn, err := router.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func buildUpstreamRouter(proxyUrls []string) (socks.Dialer, error) {
	var allForward []socks.Dialer
	for _, proxyUrl := range proxyUrls {
		forward, err := socks.NewSocks5Client("tcp", proxyUrl, "", "", socks.Direct)
		if err != nil {
			return nil, err
		}
		allForward = append(allForward, forward)
	}
	return newUpstreamDialer(allForward), nil
}
