package actions

import "manno.name/mm/fraas/models"

func createDeployment(as *ActionSuite, name string) (string, *models.Deployment) {
	d := &models.Deployment{
		Name:           name,
		ExternalDomain: name,
	}

	verrs, err := as.DB.ValidateAndCreate(d)
	as.NoError(err)
	as.False(verrs.HasAny())
	return d.ID.String(), d
}

func (as *ActionSuite) Test_DeploymentsResource_List() {
	as.Login()
	createDeployment(as, "frab3")

	res := as.HTML("/deployments").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "frab3")
}

func (as *ActionSuite) Test_DeploymentsResource_Show() {
	as.Login()
	id, _ := createDeployment(as, "frab3")

	res := as.HTML("/deployments/" + id).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), id)
}

func (as *ActionSuite) Test_DeploymentsResource_New() {
	as.Login()
	res := as.HTML("/deployments/new").Get()
	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_DeploymentsResource_Create() {
	as.Login()
	d := &models.Deployment{
		Name:           "frab4",
		ExternalDomain: "frab4",
	}

	res := as.HTML("/deployments").Post(d)
	as.Equal(302, res.Code)
	err := as.DB.First(d)
	as.NoError(err)
	as.Equal("ffrab4", d.DatabaseID)
}

func (as *ActionSuite) Test_DeploymentsResource_Edit() {
	as.Login()
	id, _ := createDeployment(as, "frab3")

	res := as.HTML("/deployments/" + id + "/edit").Get()
	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_DeploymentsResource_Update() {
	as.Login()
	id, d := createDeployment(as, "frab3")

	d.Name = "frab4"
	res := as.HTML("/deployments/" + id).Put(d)
	as.Equal(302, res.Code)

	err := as.DB.Find(d, id)
	as.NoError(err)
	as.Equal("frab4", d.Name)
}

func (as *ActionSuite) Test_DeploymentsResource_Destroy() {
	as.Login()
	id, _ := createDeployment(as, "frab3")

	count, err := as.DB.Count("deployments")
	as.NoError(err)
	as.Equal(1, count)

	res := as.HTML("/deployments/" + id).Delete()
	as.Equal(302, res.Code)

	count, err = as.DB.Count("deployments")
	as.NoError(err)
	as.Equal(0, count)
}

func (as *ActionSuite) Test_Deployments_Set() {
	as.Login()
	id, d := createDeployment(as, "frab3")

	res := as.HTML("/deployments/" + id + "/set").Post(d)
	as.Equal(302, res.Code)
}

func (as *ActionSuite) Test_Deployments_Unset() {
	as.Login()
	id, d := createDeployment(as, "frab3")

	res := as.HTML("/deployments/" + id + "/unset").Post(d)
	as.Equal(302, res.Code)
}
