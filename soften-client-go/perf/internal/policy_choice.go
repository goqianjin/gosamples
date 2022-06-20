package internal

type ChoicePolicy interface {
	Next() uint
}
