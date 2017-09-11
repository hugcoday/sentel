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

var subscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "List all subscriptions of the broker",
	Long:  `All software has versions. This is Hugo's`,
	Run:   subscriptionsCmdHandler,
}

func subscriptionsCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 2 {
		fmt.Println("Usage error, please see help")
		return
	}

	req := &pb.SubscriptionsRequest{Category: args[0]}

	switch args[0] {
	case "list": // Print topic list
		reply, err := sentelApi.Subscriptions(req)
		if err != nil {
			fmt.Println("Error:%v", err)
			return
		}
		for _, sub := range reply.Subscriptions {
			fmt.Printf("clientid:%s, topic:%s, attribute:%s",
				sub.ClientId, sub.Topic, sub.Attribute)
		}
	case "show":
		if len(args) != 2 {
			fmt.Println("Usage error, please see help")
			return
		}
		req.Subscription = args[1]
		if reply, err := sentelApi.Subscriptions(req); err != nil {
			fmt.Println("Error:%v", err)
			return
		} else if len(reply.Subscriptions) != 1 {
			fmt.Println("Error:sentel server return multiple subscriptions")
			return
		} else {
			sub := reply.Subscriptions[0]
			fmt.Printf("clientid:%s, topic:%s, attribute:%s",
				sub.ClientId, sub.Topic, sub.Attribute)

		}
	default:
		fmt.Println("Usage error, please see help")
		return
	}
}
