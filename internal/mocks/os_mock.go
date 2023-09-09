package mocks

type MockCommunicator struct {
	Result string
}

func (c *MockCommunicator) ExecuteScript(script string) string {
	return c.Result
}
