package app

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	pprint "gthub.com/NubeIO/rubix-cli-app/pkg/helpers/print"
	"testing"
)

func TestAppName_String(t *testing.T) {

	installer, err := New(&App{
		AppName:       "flow",
		Version:       "v0.0.1",
		RubixRootPath: "/data",
		InstallPath:   "aidan",
	})
	if err != nil {
		log.Errorln(err)
	}

	app, err := installer.SelectApp()
	fmt.Println(err)
	if err != nil {
		return
	}

	pprint.PrintJOSN(app)

}