/*
Copyright 2017 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Becram/k8s-api-client/pkg/notifier"
	util "github.com/Becram/k8s-api-client/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

type Status struct {
	Deployment  string `json:"Name"`
	RestartedAt string `json:"RestartedAt"`
}

type Statuses []Status

func appInit() *rest.Config {

	config, err := rest.InClusterConfig()
	fmt.Printf("incluster error %v\n", err)
	if err != nil {
		// fallback to kubeconfig
		home := homedir.HomeDir()
		kubeconfig := filepath.Join(home, ".kube", "config")
		if envvar := os.Getenv("KUBECONFIG"); len(envvar) > 0 {
			kubeconfig = envvar
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Printf("The kubeconfig cannot be loaded: %v\n", err)
			os.Exit(1)
		}

	}
	return config
}

func DeploymentUpdate(namespace string, deploymentName string) map[string]string {

	clientset, err := kubernetes.NewForConfig(appInit())
	if err != nil {
		panic(err)
	}
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	result, getErr := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("failed to get latest version of deployment: %v", getErr))
	}
	// annotatate := result.Spec.Template.GetAnnotations()
	// for i, d := range annotatate {
	// 	fmt.Printf("annotations at %s: %s\n", i, d)
	// }

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		if len(result.Spec.Template.GetAnnotations()) != 0 {
			annotate := result.Spec.Template.GetAnnotations()
			annotate["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
			result.Spec.Template.Annotations = annotate
			for i, d := range annotate {
				fmt.Printf("annotations at %s: %s\n", i, d)
			}

		} else {
			fmt.Print("No annoattions")
			annotate := make(map[string]string)
			annotate["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
			result.Spec.Template.Annotations = annotate
			for i, d := range annotate {
				fmt.Printf("annotations from empty at %s: %s\n", i, d)
			}
		}

		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("update failed: %v", retryErr))
	}
	return result.Spec.Template.GetAnnotations()

}

func DeploymentGet(namespace string) []string {
	var deploymentList []string
	clientset, err := kubernetes.NewForConfig(appInit())
	if err != nil {
		panic(err)
	}
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	result, getErr := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if getErr != nil {
		panic(fmt.Errorf("failed to get list of deployments: %v", getErr))
	}
	for _, d := range result.Items {
		deploymentList = append(deploymentList, d.Name)
		// fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
	return deploymentList
}

func RestartDeployment(w http.ResponseWriter, r *http.Request) {
	deployment := r.PostFormValue("Name")
	namespace := r.PostFormValue("NS")
	if util.FindString(DeploymentGet(namespace), deployment) {
		statuses := Statuses{
			Status{Deployment: deployment, RestartedAt: DeploymentUpdate(namespace, deployment)["kubectl.kubernetes.io/restartedAt"]},
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(statuses); err != nil {
			panic(fmt.Errorf("failed to get status: %v", err))
		}

		notifier.SendSlackNotification(deployment, "Restarted at: "+DeploymentUpdate(namespace, deployment)["kubectl.kubernetes.io/restartedAt"])
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode("Deployment not found"); err != nil {
			panic(fmt.Errorf("failed to get status: %v", err))
		}

	}

}

func ListDeployment(w http.ResponseWriter, r *http.Request) {
	namespace := r.PostFormValue("NS")

	// return json.Marshal(DeploymentList(namespace))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(DeploymentGet(namespace)); err != nil {
		panic(fmt.Errorf("failed to get status: %v", err))
	}

}
