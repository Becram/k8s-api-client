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
	"flag"
	"fmt"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
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

// func Client(s string) []string {

// 	// deploymentsClient := clientset.AppsV1().Deployments("apps")
// 	deploymentsClient := getClientSet().AppsV1().Deployments(s)
// 	// fmt.Printf("deploy %s\n", deploymentsClient)

// 	deployment, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
// 	// deploymentsClient.
// 	if err != nil {
// 		panic(err)
// 	}
// 	arr_deployment := make([]string, 0)

// 	for _, d := range deployment.Items {
// 		arr_deployment = append(arr_deployment, d.Name)
// 		// fmt.Printf("%s \n", i)
// 	}

// 	// fmt.Printf("string type %T\n", arr_deployment)

// 	return arr_deployment
// }

var home = homedir.HomeDir()
var kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")

func DeploymentRestart(namespace string, deploymentName string) map[string]string {

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	result, getErr := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
	}
	// annotatate := result.Spec.Template.GetAnnotations()
	// for i, d := range annotatate {
	// 	fmt.Printf("annotations at %s: %s\n", i, d)
	// }

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		// result.getErr := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
		annotatate := result.Spec.Template.GetAnnotations()
		annotatate["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
		result.Spec.Template.Annotations = annotatate
		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	return result.GetAnnotations()

}

// func AddSeconds() time.Time {

// 	t := time.Now()
// 	return t.Add(time.Second * 30)

// }

// func int32Ptr(i int32) *int32 { return &i }
