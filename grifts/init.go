package grifts

import (
	"github.com/gobuffalo/buffalo"
	"manno.name/mm/faas/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
