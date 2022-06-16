package internal

// ------ message status ------

type MessageStatus string

func (e MessageStatus) String() string {
	return string(e)
}

// ------ message goto action ------

type MessageGoto string

func (e MessageGoto) String() string {
	return string(e)
}
