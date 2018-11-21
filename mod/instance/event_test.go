package instance

import (
	"fmt"
	"testing"

	. "github.com/franela/goblin"
)

func TestEvent(t *testing.T) {
	g := Goblin(t)

	g.Describe("Instance > Event", func() {
		var inst *Instance
		// setup and teardown
		g.Before(func() {
			inst = New("python3", []string{"-m", "http.server", "9090"})
			go inst.Run()
		})

		g.After(func() {
			inst.eventHandle.Close()
		})

		g.It("should receive event: ActionStart", func() {
			// example values
			//var expInstID = 200
			go func() {
				inst.Run()
			}()

			for {
				select {
				case evt := <-inst.Once(ActionError):
					fmt.Println(evt)
				case evt := <-inst.Once(ActionStart):
					fmt.Println(evt)
				}
			}

		})
	})
}
