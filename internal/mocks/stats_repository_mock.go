package mocks

type MockStatsRepository struct {
	Temp float64
	CPU  float64
	Mem  float64
}

func (m *MockStatsRepository) ExecuteTemperatureScript() float64 {
	return m.Temp
}
func (m *MockStatsRepository) ExecuteCPUScript() float64 {
	return m.CPU
}
func (m *MockStatsRepository) ExecuteMemoryScript() float64 {
	return m.Mem
}
