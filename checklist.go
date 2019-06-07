package checklist

const (
	StatusDone = iota
	StatusWait
	StatusFail
)

type Checklist struct {
	elements []ChecklistElement
}

type ChecklistElement struct {
	text      string
	script    string
	completed bool
}

func (c ChecklistElement) doStuff() {
}
