package main

type Option struct {
	Name string
}

type IAction interface {
	Do(index int, options ...Option) int
}

type MyAction struct {
}

func (a *MyAction) Do(index int, options ...Option) int {
	return index * 2
}

func main() {
	var action IAction
	action = &MyAction{}
	action.Do(1)
}
