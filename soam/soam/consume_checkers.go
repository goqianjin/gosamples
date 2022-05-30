package soam

type consumeCheckers struct {
	preReRouteChecker  preReRouterChecker  //
	postReRouteChecker postReRouterChecker //
	reRouter           *multiTopicsRouter  // 配置路由

	preBlockingChecker  preStatusChecker  //
	postBlockingChecker postStatusChecker //
	blockingReRouter    *statusRouter     // 状态路由
	prePendingChecker   preStatusChecker  //
	postPendingChecker  postStatusChecker //
	pendingReRouter     *statusRouter     // 状态路由
	preRetryingChecker  preStatusChecker  //
	postRetryingChecker postStatusChecker //
	retryingReRouter    *statusRouter     // 状态路由
	preDeadChecker      preStatusChecker  //
	postDeadChecker     postStatusChecker //
	deadReRouter        *statusRouter     // 状态路由
	preDiscardChecker   preStatusChecker  //
	postDiscardChecker  preStatusChecker  //

	preUpgradeChecker  preStatusChecker  //
	postUpgradeChecker postStatusChecker //
	upgradeReRouter    *gradeRouter      // 升降路由: 升级为NewReady
	preDegradeChecker  preStatusChecker  //
	postDegradeChecker postStatusChecker //
	degradeReRouter    *gradeRouter      // 升降路由: 升级为NewReady
}
