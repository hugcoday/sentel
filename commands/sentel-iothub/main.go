package main

import (
	"flag"

	"github.com/cloustone/sentel/core"
	"github.com/cloustone/sentel/iothub/api"
	"github.com/cloustone/sentel/iothub/monitor"
	"github.com/cloustone/sentel/iothub/notify"
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
	// Create service manager according to the configuration
	mgr, err := core.NewServiceManager("iothub", config)
	if err != nil {
		glog.Fatal("Failed to launch ServiceManager")
		return
	}
	glog.Error(mgr.Run())
}

func init() {
	core.RegisterService("api", api.Configs, &api.ApiServiceFactory{})
	core.RegisterService("notify", notify.Configs, &notify.NotifyServiceFactory{})
	core.RegisterService("monitor", monitor.Configs, &monitor.MonitorServiceFactory{})
}
