package c2

type Task struct {
	Recipient Agent
	Command   string
}

type Agent struct {
	UUID     string
	Hostname string
	Username string
	OS       string
	Alias    string
}

type ServerState struct {
	RegisteredAgents []Agent
	SelectedAgent    Agent
	Tasks            []Task
	IsServing        bool
}
