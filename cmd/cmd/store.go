package cmd

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	dbase "gthub.com/NubeIO/rubix-cli-app/database"
	pprint "gthub.com/NubeIO/rubix-cli-app/pkg/helpers/print"
	"gthub.com/NubeIO/rubix-cli-app/service/apps"
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "manage rubix service apps",
	Long:  `do things like install an app, the device must have internet access to download the apps`,
	Run:   runApps,
}

type InstallResp struct {
	RespBuilder *apps.RespBuilder `json:"response_builder"`
}

func runApps(cmd *cobra.Command, args []string) {
	db := initDB()
	if flgApp.addStore {
		products, _ := json.Marshal([]string{"RubixCompute", "AllLinux"}) // product.ProductType
		db.DropAppStores()
		store := &apps.Store{
			Name:                    flgApp.appName,
			AppTypeName:             "Go",
			AllowableProducts:       products,
			DownloadPath:            flgApp.downloadPath,
			RubixRootPath:           flgApp.rubixRootPath,
			AppPath:                 "/data/flow-framework",
			Repo:                    "flow-framework",
			ServiceName:             "nubeio-flow-framework",
			RunAsUser:               "root",
			Port:                    1660,
			AppsPath:                "/data/rubix-apps/installed",
			ServiceDescription:      "nubeio-app rubix flow-framework",
			ServiceWorkingDirectory: "/data/rubix-apps/installed/flow-framework",
			ServiceExecStart:        "/data/rubix-apps/installed/flow-framework/app-amd64 -p 1660 -g /data/flow-framework -d data -prod",
		}

		app, err := db.CreateAppStore(store)
		if err != nil {
			log.Errorln(err)
			return
		}
		pprint.PrintJOSN(app)
	}

	if flgApp.installApp {

		app, err := db.InstallApp(&dbase.App{
			AppName: flgApp.appName,
			Version: flgApp.version,
			Token:   flgApp.token,
		})
		pprint.PrintJOSN(app)
		if err != nil {
			log.Println("install app err", err)
			return
		}

	}

}

var flgApp struct {
	token         string
	owner         string
	appName       string
	arch          string
	version       string
	downloadPath  string
	rubixRootPath string
	addStore      bool
	installApp    bool
}

func init() {
	RootCmd.AddCommand(appsCmd)
	flagSet := appsCmd.Flags()
	flagSet.StringVar(&flgApp.token, "token", "", "github oauth2 token value (optional)")
	flagSet.StringVarP(&flgApp.appName, "app", "", "", "rubix-wires, wires or RubixWires")
	flagSet.StringVar(&flgApp.version, "version", "latest", "version of build")
	flagSet.StringVar(&flgApp.downloadPath, "download-path", "", "download path")
	flagSet.StringVar(&flgApp.rubixRootPath, "rubix-path", "", "rubix main path")
	flagSet.BoolVarP(&flgApp.addStore, "store-add", "", false, "add a new app to the store")
	flagSet.BoolVarP(&flgApp.installApp, "install", "", false, "install an app")

}
