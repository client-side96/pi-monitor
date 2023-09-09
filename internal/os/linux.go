package os

import (
	"log"
	"os/exec"

	"github.com/client-side96/pi-monitor/internal/config"
)

type LinuxCommunicator struct {
	env config.Environment
}

func NewLinuxCommunicator(env config.Environment) *LinuxCommunicator {
	return &LinuxCommunicator{
		env: env,
	}
}

func (lc *LinuxCommunicator) ExecuteScript(script string) string {
	result, err := exec.Command(lc.env.ScriptDir + script).Output()
	if err != nil {
		log.Fatalf("Script failed: %s", err)
	}
	return string(result)
}
