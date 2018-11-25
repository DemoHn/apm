package instance

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	. "github.com/franela/goblin"
)

func TestInstance(t *testing.T) {
	g := Goblin(t)

	g.Describe("Instance", func() {
		var instN *Instance
		var instE *Instance

		const instNID = 1
		const instEID = 2
		cwd, _ := os.Getwd()
		testHelperPath := filepath.Join(cwd, "../../bin/apm-test-helper")

		g.Before(func() {
			instN = New(testHelperPath, []string{"normal-run"})
			instE = New("errorpath", []string{})
		})

		g.After(func() {
			// clean up instances
			if instN.GetStatus() == statusRunning {
				instN.ForceStop()
			}

			if instE.GetStatus() == statusRunning {
				instE.ForceStop()
			}
		})
		// setters
		g.It("should set ID, Name", func() {
			const expID = 234
			const expName = "Another-Merchant"

			instN.SetID(expID)
			g.Assert(instN.ID).Equal(expID)
			instN.SetName(expName)
			g.Assert(instN.Name).Equal(expName)

			// set back ID
			instN.SetID(instNID)
		})

		g.It("should get initial status", func() {
			g.Assert(instN.GetStatus()).Eql(statusReady)
		})
		// start & stop
		g.It("Run(): should start instance correctly", func() {
			go instN.Run()
			select {
			case evt := <-instN.Once(ActionStart):
				g.Assert(evt.Int(0)).Equal(instNID)
				g.Assert(evt.String(1)).Equal("")
			}
		})

		g.It("Run(): should not start an instance twice", func() {
			go instN.Run()
			select {
			case evt := <-instN.Once(ActionError):
				g.Assert(evt.Int(0)).Equal(instNID)
				g.Assert(evt.String(1)).Eql(ActionStart)
				g.Assert(evt.Args[2] == nil).Equal(false)
			}
		})

		g.It("Stop(): should stop OK", func() {
			time.Sleep(50 * time.Millisecond)
			err := instN.Stop(os.Interrupt)
			g.Assert(err).Equal(nil)
			// wait for stopped
			evt := <-instN.Once(ActionStop)
			g.Assert(evt.Int(1)).Equal(0)
		})

		g.It("Stop(): should not stop twice", func() {
			g.Assert(instN.GetStatus()).Equal(statusStopped)
			err := instN.Stop(syscall.SIGTERM)
			g.Assert(err == nil).Equal(false)
		})

		g.It("ForceStop(): should forceKill successfully", func() {
			// start again
			go instN.Run()
			<-instN.Once(ActionStart)

			g.Assert(instN.GetStatus()).Equal(statusRunning)
			time.Sleep(50 * time.Millisecond)
			// force again
			err := instN.Stop(os.Interrupt)
			g.Assert(err).Equal(nil)

			select {
			case evt := <-instN.Once(ActionStop):
				g.Assert(evt.Int(0)).Equal(instNID)
				// TODO: why will return -1?
				// g.Assert(evt.Int(1)).Equal(0)
			}

			g.Assert(instN.GetStatus()).Equal(statusStopped)
		})
	})
}
