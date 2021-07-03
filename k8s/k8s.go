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

	"github.com/Becram/k8s-api-client/notifier"
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

func DeploymentRestart(namespace string, deploymentName string) map[string]string {

	config, err := rest.InClusterConfig()
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
	clientset, err := kubernetes.NewForConfig(config)
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

		// result.getErr := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
		annotatate := result.Spec.Template.GetAnnotations()
		fmt.Printf("annotate: %s", annotatate)
		annotatate["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
		result.Spec.Template.Annotations = annotatate
		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("update failed: %v", retryErr))
	}
	return result.Spec.Template.GetAnnotations()

}

func RestartDeployment(w http.ResponseWriter, r *http.Request) {
	deployment := r.PostFormValue("Name")
	namespace := r.PostFormValue("NS")
	fmt.Printf("Restarting: %s", deployment)

	statuses := Statuses{
		Status{Deployment: deployment, RestartedAt: DeploymentRestart(namespace, deployment)["kubectl.kubernetes.io/restartedAt"]},
	}
	fmt.Printf("SStatus: %s", statuses)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		panic(fmt.Errorf("failed to get status: %v", err))
	}
	notifier.SendSlackNotification(deployment, "Restarted at: "+DeploymentRestart(namespace, deployment)["kubectl.kubernetes.io/restartedAt"])

}
