package scalabilitymanager

import (
	"context"
	"fmt"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateNewBroker(address string, nodeport string, remoteuser string, name string, logger logging.Logger) (bool, error) {
	// var kubeconfig string
	// if home := homedir.HomeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }

	// kubeconfig := flag.Lookup("kubeconfig").Value.(flag.Getter).Get().(string)

	// flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return false, err
	}
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return false, err
	}

	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "mqtt.mqttprovider.crossplane.io/v1alpha1",
			"kind":       "MqttBroker",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"forProvider": map[string]interface{}{
					"nodeAddress": address,
					"nodePort":    nodeport,
					"remoteUser":  remoteuser,
				},
			},
		},
	}

	result, err := clientset.Resource(schema.GroupVersionResource{
		Group:    "mqtt.mqttprovider.crossplane.io",
		Version:  "v1alpha1",
		Resource: "mqttbrokers",
	}).Create(context.TODO(), obj, metav1.CreateOptions{})

	if err != nil {
		logger.Debug(fmt.Sprintf("Error creating custom resource: %v", err))
		return false, err
	} else {
		logger.Debug(fmt.Sprintf("Custom resource created successfully ", result))
		return true, nil
	}

}

func ObserveBroker(name string, logger logging.Logger) (bool, int64, error) {

	// var kubeconfig string
	// if home := homedir.HomeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }

	// kubeconfig := flag.Lookup("kubeconfig").Value.(flag.Getter).Get().(string)
	// flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err)
	}
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	result, err := clientset.Resource(schema.GroupVersionResource{
		Group:    "mqtt.mqttprovider.crossplane.io",
		Version:  "v1alpha1",
		Resource: "mqttbrokers",
	}).Get(context.TODO(), "example-3", metav1.GetOptions{})

	if err != nil {
		logger.Debug(err.Error())
		return false, -1, err
	}

	result_bite, _ := result.MarshalJSON()
	logger.Debug(fmt.Sprintf(string(result_bite)))

	obj := result.Object["status"]
	status, _ := obj.(map[string]interface{})["atProvider"].(map[string]interface{})
	active, _ := status["active"].(bool)
	queueState, _ := status["queueState"].(int64)

	if active {
		logger.Debug("Il broker è attivo")
	} else {
		logger.Debug("Il broker NON è attivo")
	}

	logger.Debug(fmt.Sprintf("numero messaggi in coda: %d", queueState))

	return active, queueState, nil

}

func CreateNewProcess(address string, nodeport string, remoteuser string, programPath string, name string, logger logging.Logger) (bool, error) {
	// var kubeconfig string
	// if home := homedir.HomeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }

	// kubeconfig := flag.Lookup("kubeconfig").Value.(flag.Getter).Get().(string)
	// flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return false, err
	}
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return false, err
	}

	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "process.processprovider.crossplane.io/v1alpha1",
			"kind":       "Process",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"forProvider": map[string]interface{}{
					"nodeAddress": address,
					"nodePort":    nodeport,
					"remoteUser":  remoteuser,
					"programPath": programPath,
				},
				"providerConfigRef": map[string]interface{}{
					"name": "processprovider-config",
				},
			},
		},
	}

	result, err := clientset.Resource(schema.GroupVersionResource{
		Group:    "process.processprovider.crossplane.io",
		Version:  "v1alpha1",
		Resource: "processes",
	}).Create(context.TODO(), obj, metav1.CreateOptions{})

	if err != nil {
		logger.Debug(fmt.Sprintf("Error creating custom resource: %v", err))
		return false, err
	} else {
		logger.Debug(fmt.Sprintf("Custom resource created successfully ", result))
		return true, nil
	}

}

func CreateNewProcessConsumer(address string, nodeport string, remoteuser string, name string, logger logging.Logger) (bool, error) {

	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return false, err
	}
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return false, err
	}

	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "process.processprovider.crossplane.io/v1alpha1",
			"kind":       "Process",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"forProvider": map[string]interface{}{
					"nodeAddress": address,
					"nodePort":    nodeport,
					"remoteUser":  remoteuser,
				},
				"providerConfigRef": map[string]interface{}{
					"name": "processprovider-config",
				},
			},
		},
	}

	result, err := clientset.Resource(schema.GroupVersionResource{
		Group:    "process.processprovider.crossplane.io",
		Version:  "v1alpha1",
		Resource: "processes",
	}).Create(context.TODO(), obj, metav1.CreateOptions{})

	if err != nil {
		logger.Debug(fmt.Sprintf("Error creating custom resource: %v", err))
		return false, err
	} else {
		logger.Debug(fmt.Sprintf("Custom resource created successfully ", result))
		return true, nil
	}

}

// func CreateClusterRoleBinding(logger logging.Logger) (bool, error) {
// 	/*
// 		apiVersion: rbac.authorization.k8s.io/v1
// 		kind: ClusterRoleBinding
// 		metadata:
// 			labels:
// 				app.kubernetes.io/component: exporter
// 				app.kubernetes.io/name: kube-state-metrics
// 				app.kubernetes.io/version: 2.12.0
// 			name: kube-state-metrics
// 		roleRef:
// 			apiGroup: rbac.authorization.k8s.io
// 			kind: ClusterRole
// 			name: kube-state-metrics
// 		subjects:
// 			- kind: ServiceAccount
// 			name: kube-state-metrics
// 			namespace: kube-system
// 	*/

// 	obj := &unstructured.Unstructured{
// 		Object: map[string]interface{}{
// 			"apiVersion": "rbac.authorization.k8s.io/v1",
// 			"kind":       "ClusterRoleBinding",
// 			"metadata": map[string]interface{}{

// 				"labels": map[string]interface{}{
// 					"app.kubernetes.io/component": "exporter",
// 					"app.kubernetes.io/name":      "kube-state-metrics",
// 					"app.kubernetes.io/version":   "2.12.0",
// 				},
// 				"name": "prova",
// 			},

// 			"roleRef": map[string]interface{}{
// 				"apiGroup": "rbac.authorization.k8s.io",
// 				"kind":     "ClusterRole",
// 				"name":     "cluster-admin",
// 			},
// 			"subjects": map[string]interface{}{
// 				"kind":      "ServiceAccount",
// 				"name":      "scalability-provider-account",
// 				"namespace": "crossplane-system",
// 			},
// 		},
// 	}
// 	config, err := clientcmd.BuildConfigFromFlags("", "")
// 	if err != nil {
// 		return false, err
// 	}
// 	clientset, err := dynamic.NewForConfig(config)
// 	if err != nil {
// 		return false, err
// 	}

// 	result, err := clientset.Resource(schema.GroupVersionResource{
// 		Group:    "rbac.authorization.k8s.io",
// 		Version:  "v1",
// 		Resource: "rolebindings",
// 	}).Create(context.TODO(), obj, metav1.CreateOptions{})

// 	if err != nil {
// 		logger.Debug(fmt.Sprintf("Error creating custom resource: %v", err))
// 		return false, err
// 	} else {
// 		logger.Debug(fmt.Sprintf("Custom resource created successfully ", result))
// 		return true, nil
// 	}

// }
