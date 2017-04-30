package native

import (
	"log"
	"runtime"
	"time"

	"github.com/dertseha/jellui"
)

// Run initializes the environment to run the given application within.
func Run(app jellui.Application, deferrer <-chan func()) {
	runtime.LockOSThread()

	var window *OpenGlWindow
	{
		var err error
		window, err = NewOpenGlWindow()
		if err != nil {
			log.Panicln(err)
		}
	}

	app.Init(window)

	for !window.ShouldClose() {
		select {
		case task := <-deferrer:
			task()
		default:
			time.Sleep(1)
		}
		window.Update()
	}

	window.Close()
}
