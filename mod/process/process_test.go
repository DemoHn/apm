package process_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/DemoHn/apm/mock"
	"github.com/DemoHn/apm/mod/process"
	"github.com/franela/goblin"
	"github.com/golang/mock/gomock"
)

func TestProcess(t *testing.T) {
	g := goblin.Goblin(t)

	var proc *process.Process
	var cmd *exec.Cmd
	var mockPidUsage *mock.MockIPidUsage
	g.Describe("Process", func() {

		cwd, _ := os.Getwd()
		testHelperPath := filepath.Join(cwd, "../../bin/apm-test-helper")

		g.Before(func() {
			cmd = exec.Command(testHelperPath, "normal-with-cost")
			// costruct pidusage mock
			ctrl := gomock.NewController(t)
			mockPidUsage = mock.NewMockIPidUsage(ctrl)

			proc = &process.Process{
				Cmd:      cmd,
				PidUsage: mockPidUsage,
			}
		})

		g.After(func() {
			cmd.Process.Kill()
		})

		g.It("should not getPID /not started (beforeAll)", func() {
			fakePid := proc.GetPID()
			g.Assert(fakePid).Equal(0)
		})

		g.It("should start process correctly", func() {
			// stub
			mockPidUsage.EXPECT().SetPID(gomock.Any())

			if err := proc.Start(); err != nil {
				g.Fail(err)
			}

			newPID := proc.GetPID()
			g.Assert(newPID > 0).Equal(true)
		})

		g.It("should start fail /invalid command", func() {
			// stub
			ctrl2 := gomock.NewController(t)
			mockPidUsage2 := mock.NewMockIPidUsage(ctrl2)
			mockPidUsage2.EXPECT().SetPID(gomock.Any())

			wrongProc := &process.Process{
				Cmd:      exec.Command("non-existing-command"),
				PidUsage: mockPidUsage2,
			}

			err := wrongProc.Start()
			g.Assert(err != nil).Equal(true)
		})

		g.It("should stop success", func() {
			err := proc.Stop(syscall.SIGTERM)
			if err != nil {
				g.Fail(err)
			}

			<-time.After(500 * time.Millisecond)
			// stop an stopped process again still returns no error
			err2 := proc.Stop(syscall.SIGTERM)
			g.Assert(err2 != nil).Equal(false)
		})

		g.It("should GetUsage() with error", func() {
			e := fmt.Errorf("E")
			mockPidUsage.EXPECT().GetStat().Return(nil, e).Times(1)

			stat := proc.GetUsage()
			g.Assert(stat == nil).Equal(true)
		})

		g.It("should GetUsage() with mock return data", func() {
			data := &process.PidStat{}
			mockPidUsage.EXPECT().GetStat().Return(data, nil).Times(1)

			stat := proc.GetUsage()
			g.Assert(stat).Equal(data)
		})
	})
}
