package process

import (
	"fmt"

	// goblin

	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/franela/goblin"
)

func TestPidUsage(t *testing.T) {
	g := goblin.Goblin(t)
	var cmd *exec.Cmd

	g.Describe("Util > PidUsage", func() {
		cwd, _ := os.Getwd()
		testHelperPath := filepath.Join(cwd, "../../bin/apm-test-helper")

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
			pidUsage := NewPidUsage()
			pidUsage.SetPID(pid)
			// wait for another 100ms
			time.Sleep(100 * time.Millisecond)
			// skip tests if OS is not *ix
			if !pidUsage.supportedOS() {
				fmt.Printf("Skip this test since current OS doesn't support yet\n")
				return
			}

			pidStat, err := pidUsage.GetStat()
			if err != nil {
				g.Fail(err)
			}
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
			emptyUsage := NewPidUsage()
			emptyUsage.SetPID(nonexistPID)
			pidStat, err := emptyUsage.GetStat()

			g.Assert(err != nil).Equal(true)
			g.Assert(pidStat == nil).Eql(true)
		})
	})
}
