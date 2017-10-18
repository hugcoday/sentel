//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package cmd

import (
	"fmt"

	pb "github.com/cloustone/sentel/iothub/api"

	"github.com/spf13/cobra"
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "List all running services and start or stop service",
	Long:  `List all running service, start or stop service`,
	Run:   servicesCmdHandler,
}

func servicesCmdHandler(cmd *cobra.Command, args []string) {
	switch len(args) {
	case 0:
		// List all services status
		listAllServicesStatus(cmd, args)
	case 1, 2, 3:
		handleServiceCommand(cmd, args[1:])
	default:
		fmt.Println("Error:Invalid usage, please sell help")
		return
	}
}

func listAllServicesStatus(cmd *cobra.Command, args []string) {
	req := &pb.ServicesRequest{
		Category: "list",
	}
	reply, err := sentelApi.Services(req)
	if err != nil {
		fmt.Printf("Error:%v", err)
		return
	}
	for _, service := range reply.Services {
		fmt.Printf("service '%s' is listening on:'%s'", service.ServiceName, service.Listen)
		fmt.Printf("\tacceptors:%d", service.Acceptors)
		fmt.Printf("\tmax_clients:%d", service.MaxClients)
		fmt.Printf("\tcurrent_clients:%d", service.CurrentClients)
		fmt.Printf("\tshutdown_count:%d", service.Acceptors)
	}
	return

}

// handleServiceCommand hanle service specific commands, such as start, stop
// for example,
// sentelctl services start mqtt:tcp 127.0.0.1:8081(optional)
// sentelctl services stop mmqtt:tcp
func handleServiceCommand(cmd *cobra.Command, args []string) {
	req := &pb.ServicesRequest{
		Category: args[0],
	}

	switch args[0] {
	case "start":
		if len(args) == 2 {
			req.ServiceName = args[1]
		} else if len(args) == 3 {
			req.ServiceName = args[1]
			req.Listen = args[2]
		} else {
			fmt.Println("Error:Invalid usage, please see help")
			return
		}
		if _, err := sentelApi.Services(req); err != nil {
			fmt.Println("Error:Invalid usage, please sell help")
			return
		}
		fmt.Println("Service '%s' is successfuly started", args[1])

	case "stop":
		if len(args) == 2 {
			req.ServiceName = args[1]
		} else {
			fmt.Println("Error:Invalid usage, please see help")
			return
		}
		if _, err := sentelApi.Services(req); err != nil {
			fmt.Println("Error:Invalid usage, please sell help")
			return
		}
		fmt.Println("Service '%s' is successfuly stoped", args[1])

	default:
		fmt.Println("Error:Invalid usage, please sell help")
		return
	}
}
