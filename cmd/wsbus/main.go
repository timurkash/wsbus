package main

import (
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/timurkash/wsbus/internal/bus"
	"github.com/timurkash/wsbus/internal/conf"
	"log"
)

const (
	Name = "Ws-NATS hub"
)

var (
	flagConf string
)

func init() {
	flag.StringVar(&flagConf, "conf", "./configs", "conf path, eg: -conf conf.yaml")
}

func main() {
	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource(flagConf),
		),
	)
	defer func() {
		err := c.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err := c.Load(); err != nil {
		panic(err)
	}
	var bootstrap conf.Bootstrap
	if err := c.Scan(&bootstrap); err != nil {
		panic(err)
	}
	fmt.Println(bootstrap)
	bus, err := bus.NewBus(bootstrap.Bus)
	if err != nil {
		panic(err)
	}
	bus = bus

}
