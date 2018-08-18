package models_test

import (
	"manno.name/mm/fraas/models"
)

func (ms *ModelSuite) Test_Deployment() {
	count, err := ms.DB.Count("deployments")
	ms.NoError(err)
	ms.Equal(0, count)

	u := &models.Deployment{
		Name:           "frab3",
		ExternalDomain: "frab3",
	}

	verrs, err := ms.DB.ValidateAndCreate(u)
	ms.NoError(err)
	ms.False(verrs.HasAny())
	ms.NotZero(u.Name)

	count, err = ms.DB.Count("deployments")
	ms.NoError(err)
	ms.Equal(1, count)
}

func (ms *ModelSuite) Test_Deployment_Defaults() {
	count, err := ms.DB.Count("deployments")
	ms.NoError(err)
	ms.Equal(0, count)

	u := &models.Deployment{
		Name:           "frab3",
		ExternalDomain: "frab3",
	}

	u.SetDefaults()

	ms.NotZero(u.Domain)
	ms.NotZero(u.DatabaseID)
	ms.NotZero(u.DatabasePassword)
}
