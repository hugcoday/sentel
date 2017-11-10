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

package iothub

import (
	"flag"
	"path/filepath"
	"time"

	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"

	appsv1beta1 "k8s.io/api/apps/v1beta1"

	apiv1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	BrokerStatusInvalid = 0
	BrokerStatusStarted = 1
	BrokerStatusStoped  = 2
)

type BrokerStatus int

type Broker struct {
	bid         string       // broker identifier
	tid         string       // tenant identifier
	ip          string       // broker ip address
	port        string       // broker port
	status      BrokerStatus // broker status
	createdAt   time.Time    // created time for broker
	lastUpdated time.Time    // last updated time for broker
	pod         *apiv1.Pod   // the attached pod
}

type clusterManager struct {
	clientset *kubernetes.Clientset
}

// newClusterManager retrieve clustermanager instance connected with clustermgr
func newClusterManager(c core.Config) (*clusterManager, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolue path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolue path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &clusterManager{
		clientset: clientset,
	}, nil
}

// createBrokers create a number of brokers for tenant
func (this *clusterManager) createBrokers(tid string, count int32) ([]*Broker, error) {
	podname := "broker-" + tid
	deploymentsClient := this.clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sentel-broker",
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32ptr(count),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "sentel-broker",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  podname,
							Image: "sentel-broker:1.00",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "broker",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		return nil, err
	}
	glog.Infof("broker deployment created:%q.\n", result.GetObjectMeta().GetName())

	// maybe we shoud wait pod to be started
	time.Sleep(5 * time.Second) // TODO

	// get pod list
	pods, err := this.clientset.CoreV1().Pods(podname).List(metav1.ListOptions{})
	if err != nil {
		glog.Fatalf("Failed to get pod list for tenant(%s)", tid)
		return nil, err
	}
	// get all created pods, create broker for each pod
	brokers := []*Broker{}
	for _, pod := range pods.Items {
		broker := &Broker{
			bid:         uuid.NewV4().String(),
			status:      BrokerStatusStarted,
			createdAt:   time.Now(),
			lastUpdated: time.Now(),
			pod:         &pod,
		}
		brokers = append(brokers, broker)
	}

	return brokers, nil
}

// startBroker start specified node
func (this *clusterManager) startBroker(b *Broker) error {
	return nil
}

// stopBroker stop specified node
func (this *clusterManager) stopBroker(b *Broker) error {
	return nil
}

// deleteBrokers stop and delete brokers for tenant
func (this *clusterManager) deleteBrokers(tid string) error {
	podname := "broker-" + tid
	deletePolicy := metav1.DeletePropagationForeground
	deploymentsClient := this.clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	return deploymentsClient.Delete(podname, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

// deleteBroker stop and delete specified broker
func (this *clusterManager) deleteBroker(b *Broker) error {
	return nil
}

// rollbackBrokers rollback tenant's brokers
func (this *clusterManager) rollbackTenantBrokers(oldTenant *Tenant, newTenant *Tenant) error {
	return nil
}

func int32ptr(i int32) *int32 { return &i }
