package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/envy"

	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/packr"

	"manno.name/mm/fraas/models"
	"manno.name/mm/fraas/workers"

	fh "manno.name/mm/fraas/fraas-helpers"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_fraas_session",
		})

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Load complex site config
		if err := fh.ConfigFromEnv(); err != nil {
			app.Stop(err)
		}

		// Register jobs
		w := app.Worker
		if err := w.Register("unset_gke", workers.UnsetGKE); err != nil {
			app.Stop(err)
		}
		if err := w.Register("set_gke", workers.SetGKE); err != nil {
			app.Stop(err)
		}
		if err := w.Register("unset_db", workers.UnsetDB); err != nil {
			app.Stop(err)
		}
		if err := w.Register("set_db", workers.SetDB); err != nil {
			app.Stop(err)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
			app.Stop(err)
		}
		app.Use(T.Middleware())

		app.GET("/", HomeHandler)

		app.Use(SetCurrentUser)
		app.Use(Authorize)
		app.GET("/users/new", UsersNew)
		app.POST("/users", UsersCreate)
		app.GET("/signin", AuthNew)
		app.POST("/signin", AuthCreate)
		app.DELETE("/signout", AuthDestroy)
		app.Resource("/deployments", DeploymentsResource{})
		app.POST("/deployments/{deployment_id}/set", DeploymentsSet)
		app.POST("/deployments/{deployment_id}/unset", DeploymentsUnset)

		app.Middleware.Skip(Authorize, HomeHandler, AuthNew, AuthCreate)
		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}
