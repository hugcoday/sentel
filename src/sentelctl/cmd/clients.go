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
	pb "iothub/api"

	"github.com/spf13/cobra"
)

var clientsCmd = &cobra.Command{
	Use:   "clients",
	Short: "Inquery and control connected client",
	Long:  `Inquery client information and controll client`,
	Run:   clientCmdHandlerFunc,
}

func clientCmdHandlerFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 2 {
		fmt.Println("Usage error, please see help")
		return
	}
	switch args[0] {
	case "list": // Print client list
		reply, err := sentelApi.Clients(&pb.ClientsRequest{Category: args[0]})
		if err != nil {
			fmt.Println("Error:%v", err)
			return
		}
		for _, info := range reply.Clients {
			fmt.Printf("username:%s, cleanSession:%T, peername:%s, connectTime:%s",
				info.UserName, info.CleanSession, info.PeerName, info.ConnectTime)
		}
	case "show":
		if len(args) != 2 {
			fmt.Println("Usage error, please see help")
			return
		}
		reply, err := sentelApi.Clients(&pb.ClientsRequest{Category: args[0], ClientId: args[1]})
		if err != nil {
			fmt.Println("Error:%v", err)
			return
		}
		switch len(reply.Clients) {
		case 0:
			fmt.Printf("No client '%s' information in sentel", args[0])
			return
		case 1:
			info := reply.Clients[0]
			fmt.Printf("username:%s, cleanSession:%T, peername:%s, connectTime:%s",
				info.UserName, info.CleanSession, info.PeerName, info.ConnectTime)
		default:
			fmt.Printf("Error: sentel return multiply user infor for client '%s'", args[0])
			for _, info := range reply.Clients {
				fmt.Printf("username:%s, cleanSession:%T, peername:%s, connectTime:%s",
					info.UserName, info.CleanSession, info.PeerName, info.ConnectTime)
			}
		}

	case "kick":
		if len(args) != 2 {
			fmt.Println("Usage error, please see help")
			return
		}
		reply, err := sentelApi.Clients(&pb.ClientsRequest{Category: args[0], ClientId: args[1]})
		if err != nil {
			fmt.Println("Error:%v", err)
			return
		}
		fmt.Println(reply.Result)

	default:
		fmt.Println("Usage error, please see help")
		return
	}
}
