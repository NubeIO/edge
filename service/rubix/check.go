package rubix

import (
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type AppResponse struct {
	Name       string                 `json:"app"`
	Version    string                 `json:"version,omitempty"`
	IsAService bool                   `json:"is_service"`
	AppStatus  *systemctl.SystemState `json:"app_status,omitempty"`
	Error      string                 `json:"error,omitempty"`
}

var DefaultTimeout = 30

var systemOpts = systemctl.Options{
	UserMode: false,
	Timeout:  DefaultTimeout,
}

func (inst *App) ConfirmAppDir(appName string) bool {
	return fileutils.New().DirExists(fmt.Sprintf("%s/%s", DataDir, appName))
}

func (inst *App) ConfirmAppInstallDir(appInstallName string) bool {
	return fileutils.New().DirExists(fmt.Sprintf("%s/%s", AppsInstallDir, appInstallName))
}

func (inst *App) GetAppVersion(appInstallName string) string {
	file := fmt.Sprintf("%s/%s", AppsInstallDir, appInstallName)
	fileInfo, err := os.Stat(file)
	if err != nil {
		return ""
	}
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(file)
		if err != nil {
			return ""
		}
		for _, file := range files {
			fmt.Println(file.Name())
			if checkVersionBool(file.Name()) {
				return file.Name()
			}
		}
	}
	return ""
}

func (inst *App) listFiles(file string) ([]string, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	var dirContent []string
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(file)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			dirContent = append(dirContent, file.Name())
		}
	}
	return dirContent, nil
}

func (inst *App) InstalledApps() ([]AppResponse, error) {
	rootDir := AppsInstallDir
	var files []AppResponse
	app := AppResponse{}
	err := filepath.WalkDir(rootDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.Count(p, string(os.PathSeparator)) == 6 {
			parts := strings.Split(p, "/")
			if len(parts) >= 5 { // app name
				app.Name = parts[5]
			}
			if len(parts) >= 6 { // version
				app.Version = parts[6]
			}
			app.AppStatus = nil
			files = append(files, app)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (inst *App) ConfirmInstalledApps(apps []string) ([]AppResponse, error) {
	var out []AppResponse
	app := AppResponse{}
	for _, ap := range apps {
		installed, err := systemctl.IsInstalled(ap, systemOpts)
		app.Name = ap
		if installed {
			app.IsAService = true
			state, err := systemctl.State(ap, systemOpts)
			if err != nil {
				app.Error = err.Error()
			}
			app.AppStatus = &state
			out = append(out, app)
		} else {
			app.IsAService = false
			app.AppStatus = nil
			app.Error = err.Error()
			out = append(out, app)
		}
	}
	return out, nil
}