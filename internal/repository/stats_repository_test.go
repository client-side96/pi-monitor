package repository_test

import (
	"testing"

	"github.com/client-side96/pi-monitor/internal/mocks"
	"github.com/client-side96/pi-monitor/internal/repository"
)

func TestRepository_Stats_ExecuteTemperatureScript(t *testing.T) {
	mockOs := &mocks.MockCommunicator{
		Result: "20.14",
	}
	repo := repository.NewStatsRepository(mockOs)
	temp := repo.ExecuteTemperatureScript()

	if temp != 20.14 {
		t.Errorf("Expected 20.14 but got %f", temp)
	}
}
func TestRepository_Stats_ExecuteCPUScript(t *testing.T) {
	mockOs := &mocks.MockCommunicator{
		Result: "10.01",
	}
	repo := repository.NewStatsRepository(mockOs)
	temp := repo.ExecuteCPUScript()

	if temp != 10.01 {
		t.Errorf("Expected 10.01 but got %f", temp)
	}
}
func TestRepository_Stats_ExecuteMemoryScript(t *testing.T) {
	mockOs := &mocks.MockCommunicator{
		Result: "25.0",
	}
	repo := repository.NewStatsRepository(mockOs)
	temp := repo.ExecuteMemoryScript()

	if temp != 25.0 {
		t.Errorf("Expected 20.14 but got %f", temp)
	}
}
