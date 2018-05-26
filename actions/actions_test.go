package actions

import (
	"testing"

	"github.com/gobuffalo/suite"
	"manno.name/mm/faas/models"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	as := &ActionSuite{suite.NewAction(App())}
	suite.Run(t, as)
}

func (as *ActionSuite) Login() {
	u := &models.User{
		Email:                "admin@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}
	verrs, err := u.Create(as.DB)
	as.NoError(err)
	as.False(verrs.HasAny())
	as.Session.Set("current_user_id", u.ID)
}
