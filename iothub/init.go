package iothub

import (
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

func RunWithConfigFile(fileName string) error {
	glog.Info("Starting iothub ...")

	// Get configuration
	config, err := core.NewWithConfigFile(fileName)
	if err != nil {
		return err
	}

	// Initialize iothub at startup
	glog.Info("Initializing iothub...")
	if err := InitializeIothub(config); err != nil {
		return err
	}

	// Create service manager according to the configuration
	mgr, err := core.NewServiceManager("iothub", config)
	if err != nil {
		return err
	}
	return mgr.Run()
}

func init() {
	for group, values := range defaultConfigs {
		core.RegisterConfig(group, values)
	}
	core.RegisterService("api", nil, &ApiServiceFactory{})
	core.RegisterService("notify", nil, &NotifyServiceFactory{})
}
