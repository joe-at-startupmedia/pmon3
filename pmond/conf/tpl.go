package conf

import "time"

type Tpl struct {
	Data                string `yaml:"data"`
	Logs                string `yaml:"logs"`
	HandleInterrupts    bool   `yaml:"handle_interrupts"`
	CmdExecResponseWait int64  `yaml:"cmd_exec_response_wait"`
	Conf                string
}

func (c *Tpl) GetDataDir() string {
	return c.Data
}

func (c *Tpl) GetLogsDir() string {
	return c.Logs
}

func (c *Tpl) ShouldHandleInterrupts() bool {
	return c.HandleInterrupts
}

func (c *Tpl) GetCmdExecResponseWait() time.Duration {
	if c.CmdExecResponseWait >= 0 && c.CmdExecResponseWait <= 10000 {
		return time.Duration(c.CmdExecResponseWait) * time.Millisecond
	} else {
		return 2000 * time.Millisecond
	}
}
