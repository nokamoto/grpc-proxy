package yaml

type Observe struct {
	Observe struct {
		Logs []Log
	}
}

type Log struct {
	Name string
	File string
}
