package internal

// ------ checkpoint ------

type Checkpoint struct {
	CheckType          CheckType
	PreStatusChecker   PreStatusChecker
	PostStatusChecker  PostStatusChecker
	PreRerouteChecker  PreRerouteChecker
	PostRerouteChecker PostRerouteChecker
}
