package util_test

import (
	"fmt"

	"github.com/DemoHn/apm/util"
	// goblin

	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	. "github.com/franela/goblin"
)

func TestPidUsage(t *testing.T) {
	g := Goblin(t)
	var cmd *exec.Cmd
	g.Describe("Util > PidUsage", func() {
		cwd, _ := os.Getwd()
		testHelperPath := filepath.Join(cwd, "../bin/apm-test-helper")

		g.Before(func() {
			cmd = exec.Command(testHelperPath, "normal-with-cost")
			go cmd.Run()
			// wait for a white to start
			time.Sleep(100 * time.Millisecond)
		})

		g.After(func() {
			cmd.Process.Kill()
		})

		g.It("should stat successfully", func() {
			pid := cmd.Process.Pid
			pidUsage := util.NewPidUsage(pid)

			// wait for another 100ms
			time.Sleep(100 * time.Millisecond)
			pidStat := pidUsage.GetStat()
			fmt.Printf("[debug] PID=%d, PPID=%d, CPU=%.2f, Memory=%d KB, Elapsed=%.2f ms\n",
				pidStat.Pid,
				pidStat.PPid,
				pidStat.CPU,
				pidStat.Memory/1000,
				pidStat.Elapsed)

			g.Assert(pidStat.Pid).Eql(pid)
			// pid should not be ppid!
			g.Assert(pidStat.Pid != pidStat.PPid).Eql(true)
			// this cmd will consume CPU, I promise
			g.Assert(pidStat.CPU > 0).Eql(true)
			// since we have wait for 100ms at least
			g.Assert(pidStat.Elapsed > 0.1).Eql(true)
		})

		g.It("should get stat failed /no such PID", func() {
			nonexistPID := -1235
			emptyUsage := util.NewPidUsage(nonexistPID)

			pidStat := emptyUsage.GetStat()
			g.Assert(pidStat == nil).Eql(true)
		})
	})
}
