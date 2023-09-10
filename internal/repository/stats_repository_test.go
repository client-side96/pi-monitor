package repository_test

import (
	"testing"

	"github.com/client-side96/pi-monitor/internal/repository"
	"github.com/client-side96/pi-monitor/internal/testutil"
	"github.com/client-side96/pi-monitor/internal/testutil/mocks"
)

func TestRepository_Stats_ExecuteTemperatureScript(t *testing.T) {
	mockOs := &mocks.MockCommunicator{
		Result: "20.14",
	}
	repo := repository.NewStatsRepository(mockOs)
	result := repo.ExecuteTemperatureScript()

	testutil.AssertEqual(t, 20.14, result)
}
func TestRepository_Stats_ExecuteCPUScript(t *testing.T) {
	mockOs := &mocks.MockCommunicator{
		Result: "10.01",
	}
	repo := repository.NewStatsRepository(mockOs)
	result := repo.ExecuteCPUScript()

	testutil.AssertEqual(t, 10.01, result)
}
func TestRepository_Stats_ExecuteMemoryScript(t *testing.T) {
	mockOs := &mocks.MockCommunicator{
		Result: "25.0",
	}
	repo := repository.NewStatsRepository(mockOs)
	result := repo.ExecuteMemoryScript()

	testutil.AssertEqual(t, 25.0, result)
}
