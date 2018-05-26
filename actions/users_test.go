package actions

import (
	"manno.name/mm/faas/models"
)

func (as *ActionSuite) Test_Users_New() {
	as.Login()
	res := as.HTML("/users/new").Get()
	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_Users_Create() {
	as.Login()
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(1, count)

	u := &models.User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}

	res := as.HTML("/users").Post(u)
	as.Equal(302, res.Code)

	count, err = as.DB.Count("users")
	as.NoError(err)
	as.Equal(2, count)
}
