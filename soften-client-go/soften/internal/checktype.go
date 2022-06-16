package internal

// ------ check type ------

type CheckType string

func (e *CheckType) String() string {
	return string(*e)
}
