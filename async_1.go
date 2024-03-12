package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/memberlist"
	"gopkg.in/urfave/cli.v1"
)

func action(c *cli.Context) {
	conf := memberlist.DefaultLocalConfig()
	conf.Name = "async_1"
	conf.Events = new(MyEventDelegate)

	list, err := memberlist.Create(conf)
	if err != nil {
		log.Fatal(err)
	}

	local := list.LocalNode()
	log.Printf("async_1 at %s:%d", local.Addr.To4().String(), local.Port)

	//list.Join([]string{
	//	fmt.Sprintf("%s:%d", local.Addr.To4().String(), local.Port),
	//})

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
					log.Printf("node_1 search %d => %s", key, node)
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
	app.Run(os.Args)
}
