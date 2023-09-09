package repository

import (
	"github.com/client-side96/pi-monitor/internal/os"
	"github.com/client-side96/pi-monitor/internal/util"
)

type IStatsRepository interface {
	ExecuteTemperatureScript() float64
	ExecuteCPUScript() float64
	ExecuteMemoryScript() float64
}

type StatsRepository struct {
	os os.Communicator
}

func NewStatsRepository(os os.Communicator) *StatsRepository {
	return &StatsRepository{
		os: os,
	}
}

func (repo *StatsRepository) ExecuteTemperatureScript() float64 {
	var tempScript = "temperature.sh"
	return util.ToFloat(repo.os.ExecuteScript(tempScript))
}

func (repo *StatsRepository) ExecuteCPUScript() float64 {
	var cpuScript = "cpu.sh"
	return util.ToFloat(repo.os.ExecuteScript(cpuScript))
}

func (repo *StatsRepository) ExecuteMemoryScript() float64 {
	var memoryScript = "memory.sh"
	return util.ToFloat(repo.os.ExecuteScript(memoryScript))
}
