package main

import (
	"context"
	"github.com/hashicorp/memberlist"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"strconv"
	"time"
)

func action(c *cli.Context) {
	conf := memberlist.DefaultLocalConfig()
	conf.Name = "async_2"
	conf.BindPort = 7947
	conf.AdvertisePort = conf.BindPort
	conf.Events = new(MyEventDelegate)

	list, err := memberlist.Create(conf)
	if err != nil {
		log.Fatal(err)
	}

	local := list.LocalNode()
	log.Printf("async_2 at %s:%d", local.Addr.To4().String(), local.Port)

	//list.Join([]string{
	//	fmt.Sprintf("%s:%d", local.Addr.To4().String(), local.Port),
	//})

	join := c.String("join")
	log.Printf("async_2 join to %s", join)

	if _, err := list.Join([]string{join}); err != nil {
		log.Fatal(err)
	}

	// -------
	stopCtx, cancel := context.WithCancel(context.TODO())
	go wait_signal(cancel)

	tick := time.NewTicker(3 * time.Second)
	run := true
	for run {
		select {
		case <-tick.C:
			devt := conf.Events.(*MyEventDelegate)
			if devt == nil {
				log.Printf("consistent isnt initialized")
				continue
			}
			log.Printf("current node size: %d", devt.consistent.Size())

			for key := 1; key <= 10; key++ {
				node, ok := devt.consistent.GetNode(strconv.Itoa(key))
				if ok == true {
					log.Printf("node_2 search %d => %s", key, node)
				} else {
					log.Printf("no node available")
				}
			}
		case <-stopCtx.Done():
			log.Printf("stop called")
			run = false
		}
	}
	tick.Stop()
	log.Printf("bye.")
}

func main() {
	app := cli.NewApp()
	app.Action = action
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "join, j",
			Usage:  "cluster join address",
			Value:  "127.0.0.1:xxxx",
			EnvVar: "JOIN_ADDR",
		},
	}
	app.Run(os.Args)
}
