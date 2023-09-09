package mocks

type MockCommunicator struct {
	Result string
}

func (c *MockCommunicator) ExecuteScript(_ string) string {
	return c.Result
}
