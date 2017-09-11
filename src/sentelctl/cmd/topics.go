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

var topicsCmd = &cobra.Command{
	Use:   "topics",
	Short: "List all topics of the broker",
	Long:  `List All topics of the broker and inquery detail topic information`,
	Run:   topicsCmdHandler,
}

func topicsCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 2 {
		fmt.Println("Usage error, please see help")
		return
	}

	req := &pb.TopicsRequest{Category: args[0]}

	switch args[0] {
	case "list": // Print topic list
		reply, err := sentelApi.Topics(req)
		if err != nil {
			fmt.Println("Error:%v", err)
			return
		}
		for _, topic := range reply.Topics {
			fmt.Printf("%s, %s", topic.Topic, topic.Attribute)
		}
	case "show":
		if len(args) != 2 {
			fmt.Println("Usage error, please see help")
			return
		}
		req.Topic = args[1]
		if reply, err := sentelApi.Topics(req); err != nil {
			fmt.Println("Error:%v", err)
			return
		} else if len(reply.Topics) != 1 {
			fmt.Println("Error:sentel server return multiple topics")
			return
		} else {
			topic := reply.Topics[0]
			fmt.Printf("%s, %s", topic.Topic, topic.Attribute)
		}
	default:
		fmt.Println("Usage error, please see help")
		return
	}
}
