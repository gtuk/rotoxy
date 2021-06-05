package main

import (
	"context"
	"fmt"
	"github.com/gtuk/rotating-tor-proxy/core"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
)

func main() {
	var numberTorInstances int
	var port int
	var circuitInterval int

	app := &cli.App{
		Name:  "rotoxy",
		Usage: "run a rotating Tor proxy server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "tors",
				Value:       1,
				Usage:       "number of Tor proxies that should run",
				Destination: &numberTorInstances,
			},
			&cli.IntFlag{
				Name:        "port",
				Value:       8080,
				Usage:       "port where the reverse proxy should listen on",
				Destination: &port,
			},
			&cli.IntFlag{
				Name:        "circuitInterval",
				Value:       30,
				Usage:       "number in seconds after a new circuit should be requested",
				Destination: &circuitInterval,
			},
		},
		Action: func(c *cli.Context) error {
			return run(port, numberTorInstances, circuitInterval)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(port int, numberTorInstances int, circuitInterval int) error {
	log.Println(fmt.Sprintf("Starting tor proxies"))

	proxies := make([]core.TorProxy, 0)
	ch := make(chan core.TorProxy, numberTorInstances)

	g, _ := errgroup.WithContext(context.Background())

	for i := 1; i <= numberTorInstances; i++ {
		g.Go(func() error {
			torProxy, err := core.CreateTorProxy(circuitInterval)
			if err != nil {
				if torProxy != nil {
					torProxy.Close()
				}

				return err
			}

			ch <- *torProxy

			return nil
		})
	}

	err := g.Wait()
	close(ch)

	for proxy := range ch {
		proxies = append(proxies, proxy)
	}
	defer core.CloseProxies(proxies)

	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Started %d tor proxies", len(proxies)))
	log.Println(fmt.Sprintf("Start reverse proxy on port %d", port))

	reverseProxy := &core.ReverseProxy{}
	err = reverseProxy.Start(proxies, port)
	if err != nil {
		return err
	}

	return nil
}
