package main

import "context"

type AccountProcessor struct {
}
type IssueProcessor struct {
}
type SwaggerAccount struct {
	ID       string `json:"ID"`
	Company  string `json:"company"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Status   string `json:"status"`
}

type SwaggerIssue struct {
	Company string `json:"company"`
	Token   string `json:"token"`
}

// @Summary 查询所有用户帐号
// @Description 查询所有用户帐号
// @Tags Account
// @Accept  json
// @Produce  json
// @Param X-Namespace header string true "命名空间"
// @Success 200 {anrry} SwaggerAccount "帐号的数组"
// @Router /v/1/accounts [get]
func (p *AccountProcessor) ProcessRetrieveMany(ctx context.Context) error {
	return nil
}

// @Summary 添加用户帐号(注册)
// @Description 添加用户帐号、用户注册
// @Tags Account
// @Accept  json
// @Produce  json
// @Param X-Namespace header string true "命名空间"
// @Param name body string true "用户名" default(admin)
// @Param password formData string true "密码"
// @Param email formData string true "邮箱帐号"
// @Param company formData string true "命名空间"
// @Param status formData string true "用户帐号的状态" Enums(Active, Disabled)
// @Success 204 "No Content"
// @Router /v/1/accounts [post]
func (p *AccountProcessor) ProcessCreate(ctx context.Context) error {
	return nil
}

// @Summary 用户登录
// @Description 用户登录
// @Tags Issue
// @Accept  json
// @Produce  json
// @Param X-Namespace header string true "命名空间"
// @Param realm body string true "The authentication realm." default(Vince)
// @Param name body string true "用户名" default(admin)
// @Param password body string true "密码"
// @Success 200 {object} SwaggerIssue "命名空间和token"
// @Router /v/1/issue [post]
func (p *IssueProcessor) ProcessCreate(ctx context.Context) error {
	return nil
}
