package prolabs

type Prolabs struct {
	Name        string
	Flags       int
	Machines    int
	Difficulty  string
	Progression int
}

type model struct {
	Prolabs []Prolabs
	width   int
	height  int
}
