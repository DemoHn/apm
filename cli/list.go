package cli

import (
	"os"
	"strconv"

	"github.com/DemoHn/apm/mod/instance"
	"github.com/DemoHn/apm/mod/master"
	"github.com/DemoHn/tablewriter"
	"github.com/urfave/cli"
)

const stringNA = "N/A"

var listFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "instance ID. Use this to show the information of this one",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "Filter shown instances by name",
	},
}

func listHandler(c *cli.Context) error {
	var resp master.ListInstanceResponse
	// get id
	var id int = c.Int("id")
	var rid *int
	if id != 0 {
		rid = &id
	}
	var name string = c.String("name")
	var rname *string
	if name != "" {
		rname = &name
	}
	req := &master.ListInstanceRequest{
		ID:   rid,
		Name: rname,
	}
	err := sendRequest("Tower.ListInstance", req, &resp)
	if err != nil {
		return err
	}
	// print result with a table
	table := tablewriter.NewWriter(os.Stdout)
	// table settings
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.SetHeader([]string{
		"ID", "name", "status", "↻", "pid", "cpu", "memory", "uptime",
	})

	for _, info := range resp.InstanceInfos {
		table.Append(stringifyInfo(info))
	}
	table.Render()

	return nil
}

func stringifyInfo(info instance.Info) []string {
	sinfo := make([]string, 8)
	// id, name
	sinfo[0] = strconv.Itoa(info.ID)
	sinfo[1] = info.Name
	// status
	sinfo[2] = stringifyStatus(info.Status)
	// restart times
	sinfo[3] = strconv.Itoa(info.RestartTimes)
	// pid, cpu, memory, uptime
	sinfo[4] = stringifyPID(info.PID)
	sinfo[5] = stringifyCPU(info.CPU)
	sinfo[6] = stringifyMemory(info.Memory)
	sinfo[7] = stringifyUptime(info.LaunchTime)
	return sinfo
}

func stringifyStatus(status instance.StatusFlag) string {
	statusMap := map[instance.StatusFlag]string{
		instance.StatusReady:   "ready",
		instance.StatusRunning: "running",
		instance.StatusStopped: "stopeed",
	}

	return statusMap[status]
}

func stringifyPID(pid *int) string {
	if pid == nil {
		return stringNA
	}
	return strconv.Itoa(*pid)
}

func stringifyCPU(cpu *float64) string {
	if cpu == nil {
		return stringNA
	}

	return strconv.FormatFloat((*cpu)*100, 'f', 2, 64) + " %"
}

func stringifyMemory(memory *int64) string {
	if memory == nil {
		return stringNA
	}

	return strconv.FormatInt(*memory/1024, 10) + "K"
}

func stringifyUptime(uptime *float64) string {
	if uptime == nil {
		return stringNA
	}

	return strconv.FormatFloat(*uptime, 'f', 1, 64) + "s"
}
