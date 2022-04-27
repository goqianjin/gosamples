package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Users struct {
	Phone  string `form:"phone" json:"phone" validate:"required"`
	passwd string `form:"passwd" json:"passwd" validate:"required,max=20,min=6"`
	Code   string `form:"code" json:"code" validate:""`
}

func main() {
	users := &Users{
		Phone:  "1326654487",
		passwd: "123a",
		Code:   "1234567",
	}
	validate := validator.New()
	err := validate.Struct(users)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err) //Key: 'Users.Passwd' Error:Field validation for 'Passwd' failed on the 'min' tag
			//return
		}
	}

	var boolTest bool
	err = validate.Var(boolTest, "required")
	if err != nil {
		fmt.Println(err)
	}
	var stringTest string = ""
	err = validate.Var(stringTest, "required")
	if err != nil {
		fmt.Println(err)
	}
	var stringTest1 = "12"
	err = validate.Var(stringTest1, "required,numeric,max=20,min=6")
	if err != nil {
		fmt.Println(err)
	}
}
