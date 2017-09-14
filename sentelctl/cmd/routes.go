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

var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "List routes informations of broker",
	Long:  `List all routes ro routes for specific topic`,
	Run:   routesCmdHandler,
}

func routesCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 2 {
		fmt.Println("Usage error, please see help")
		return
	}
	req := &pb.RoutesRequest{Category: args[0], Service: ""}
	switch args[0] {
	case "list":
		if len(args) == 2 {
			req.Service = args[1]
		}
		if reply, err := sentelApi.Routes(req); err != nil {
			fmt.Printf("Error:%v", err)
			return
		} else {
			for _, info := range reply.Routes {
				fmt.Printf("%s ->%v", info.Topic, info.Route)
			}
		}
	case "show":
		if len(args) != 2 {
			fmt.Println("Usage error, please see help")
			return
		}
		req.Topic = args[1]
		if reply, err := sentelApi.Routes(req); err != nil {
			fmt.Printf("Error:%v", err)
			return
		} else {
			for _, info := range reply.Routes {
				fmt.Printf("%s ->%v", info.Topic, info.Route)
			}
		}
	}
}
