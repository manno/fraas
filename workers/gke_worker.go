package workers

import (
	"fmt"

	"github.com/gobuffalo/buffalo/worker"
)

var w worker.Worker

func UnsetGKE(args worker.Args) error {
	fmt.Printf("%#v", args)
	return nil
}

func SetGKE(args worker.Args) error {
	fmt.Printf("%#v", args)
	return nil
}
