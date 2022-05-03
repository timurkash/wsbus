package main

import (
	"flag"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/timurkash/wsbus/internal/conf"
	"github.com/timurkash/wsbus/internal/hub"
	"github.com/timurkash/wsbus/internal/ws"
	"log"
	"net/http"
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
	hub, err := hub.NewHub(bootstrap.Bus, bootstrap.Ws)
	if err != nil {
		panic(err)
	}
	defer hub.CloseBus()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.Serve(hub, w, r)
	})
	if err := http.ListenAndServe(bootstrap.Server.Http.Addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
