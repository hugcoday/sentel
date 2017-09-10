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

var brokerCmd = &cobra.Command{
	Use:   "broker",
	Short: "Inquery broker status, such as metrics, clients",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Invalid usage, please see help")
			return
		}
		switch args[0] {
		case "stats":
			getBrokerStats()
		case "metrics":
			getBrokerMetrics()
		default:
			fmt.Println("Invalid usage, please see help")
		}
	},
}

func brokerHelpFunc(cmd *cobra.Command, args []string) {
	fmt.Println("broker stats|metrics")
}

func getBrokerStats() {
	reply, err := sentelApi.Broker(&pb.BrokerRequest{Category: "stats"})
	if err != nil {
		fmt.Println("Error:%v", err)
		return
	}
	for key, val := range reply.Stats {
		fmt.Println("%10s:%10s", key, val)
	}
}

func getBrokerMetrics() {
	reply, err := sentelApi.Broker(&pb.BrokerRequest{Category: "metrics"})
	if err != nil {
		fmt.Println("Error:%v", err)
		return
	}
	for key, val := range reply.Metrics {
		fmt.Println("%10s:%10s", key, val)
	}

}
