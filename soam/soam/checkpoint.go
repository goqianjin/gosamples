package soam

type checkpoint struct {
	checkType           CheckType
	preStatusChecker    preStatusChecker
	postStatusChecker   postStatusChecker
	preReRouterChecker  preReRouterChecker
	postReRouterChecker postReRouterChecker
}
