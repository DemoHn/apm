package util

import (
	// goblin

	"testing"

	. "github.com/franela/goblin"
)

func TestPidUsage(t *testing.T) {
	g := Goblin(t)

	g.Describe("Util > PidUsage", func() {
		g.It("could read data", func() {
			/*u := New(1933)
			for i := 0; i < 7; i++ {
				fmt.Println(u.GetStat())
				time.Sleep(500 * time.Millisecond)
			}
			*/
		})
	})
}
