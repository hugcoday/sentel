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

var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "List all MQTT session of the broker",
	Long:  `List all MQTT session of the broker, or specified session`,
	Run:   sessionsCmdHandler,
}

func sessionsCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 2 {
		fmt.Println("Usage error, please see help")
		return
	}

	req := &pb.SessionsRequest{Category: args[0], Conditions: make(map[string]bool)}

	switch args[0] {
	case "list": // Print client list
		// Format conditions
		if len(args) == 2 {
			switch args[1] {
			case "persistent":
				req.Conditions["persistent"] = true
			case "transient":
				req.Conditions["transient"] = true
			}
		}
		reply, err := sentelApi.Sessions(req)
		if err != nil {
			fmt.Println("Error:%v", err)
			return
		}
		for _, session := range reply.Sessions {
			fmt.Printf(`clientid=%s, created_at=%s, clean_session=%s, 
				max_inflight=%s,infliaht=%s,inqueue=%s,droped=%s,awaiting_rel=%s,
				awaiting_comp=%s,awaiting_ack=%s`,
				session.ClientId, session.CreatedAt, session.CleanSession,
				session.MessageMaxInflight,
				session.MessageInflight,
				session.MessageInQueue,
				session.MessageDropped,
				session.AwaitingRel,
				session.AwaitingComp,
				session.AwaitingAck)
		}
	case "show":
		if len(args) != 2 {
			fmt.Println("Usage error, please see help")
			return
		}
		req.ClientId = args[1]
		if reply, err := sentelApi.Sessions(req); err != nil {
			fmt.Println("Error:%v", err)
			return
		} else if len(reply.Sessions) != 1 {
			fmt.Println("Error:sentel server return multiple sessions")
			return
		} else {
			session := reply.Sessions[0]
			fmt.Printf(`clientid=%s, created_at=%s, clean_session=%s, 
				max_inflight=%s,infliaht=%s,inqueue=%s,droped=%s,awaiting_rel=%s,
				awaiting_comp=%s,awaiting_ack=%s`,
				session.ClientId, session.CreatedAt, session.CleanSession,
				session.MessageMaxInflight,
				session.MessageInflight,
				session.MessageInQueue,
				session.MessageDropped,
				session.AwaitingRel,
				session.AwaitingComp,
				session.AwaitingAck)

		}
	default:
		fmt.Println("Usage error, please see help")
		return
	}
}
