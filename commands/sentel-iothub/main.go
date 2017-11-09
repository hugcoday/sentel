package main

import (
	"flag"

	"github.com/cloustone/sentel/core"
	"github.com/cloustone/sentel/iothub"
	"github.com/golang/glog"
)

var (
	configFileFullPath = flag.String("c", "stentel-iothub.conf", "config file")
)

func main() {
	flag.Parse()
	glog.Info("Starting iothub ...")

	// Get configuration
	config, err := core.NewWithConfigFile(*configFileFullPath)
	if err != nil {
		glog.Fatal(err)
		flag.PrintDefaults()
		return
	}

	// Initialize iothub at startup
	if err := iothub.InitializeIothub(config); err != nil {
		glog.Fatal(err)
		return
	}

	// Create service manager according to the configuration
	mgr, err := core.NewServiceManager("iothub", config)
	if err != nil {
		glog.Fatal("Failed to launch ServiceManager")
		return
	}
	glog.Error(mgr.Run())
}

func init() {
	for group, values := range iothub.Configs {
		core.RegisterConfig(group, values)
	}
	core.RegisterService("api", nil, &iothub.ApiServiceFactory{})
	core.RegisterService("notify", nil, &iothub.NotifyServiceFactory{})
}
