package main

import (
	"sync"
	"time"

	"github.com/iotexproject/Bumblebee/kit/kit"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis"
	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/global"
)

var app = global.App

func main() {
	app.AddCommand("migrate", func(args ...string) {
		global.Migrate()
	})

	app.Execute(func(args ...string) {
		BatchRun(
			func() {
				kit.Run(apis.Root, global.Server())
			},
			global.EventProxy,
			// func() {
			// 	if err := applet_deploy.StartAppletVMs(
			// 		global.WithContext(context.Background()),
			// 	); err != nil {
			// 		panic(err)
			// 	}
			// },
		)
	})
}

func BatchRun(commands ...func()) {
	wg := &sync.WaitGroup{}

	for i := range commands {
		cmd := commands[i]
		wg.Add(1)

		go func() {
			defer wg.Done()
			cmd()
			time.Sleep(200 * time.Millisecond)
		}()
	}
	wg.Wait()
}
