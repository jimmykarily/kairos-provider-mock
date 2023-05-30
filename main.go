package main

import (
	"encoding/json"
	"os"

	"github.com/kairos-io/kairos-sdk/bus"
	kairos "github.com/kairos-io/kairos/pkg/config"
	"github.com/mudler/go-pluggable"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type CloudInit struct {
	Install kairos.Install `yaml:"install,omitempty" json:"install,omitempty"`
}

func main() {
	plugins := pluggable.NewPluginFactory(
		pluggable.FactoryPlugin{
			EventType:     bus.EventInstall,
			PluginHandler: HandleInstall,
		},
	)

	if err := plugins.Run(pluggable.EventType(os.Args[1]), os.Stdin, os.Stdout); err != nil {
		logrus.Fatal(err)
	}
}

func HandleInstall(event *pluggable.Event) pluggable.EventResponse {
	var cloudInit CloudInit

	cloudInit.Install.GrubOptions["saved_entry"] = "registration"

	ccBytes, err := yaml.Marshal(cloudInit)
	if err != nil {
		return pluggable.EventResponse{
			Error: err.Error(),
		}
	}

	rpayload := map[string]string{
		"cc":       string(ccBytes),
		"device":   "auto",
		"reboot":   "false",
		"poweroff": "false",
	}

	data, _ := json.Marshal(rpayload)
	return pluggable.EventResponse{
		Data: string(data),
	}
}
