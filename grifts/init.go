package grifts

import (
	"github.com/gobuffalo/buffalo"
	"manno.name/mm/fraas/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
