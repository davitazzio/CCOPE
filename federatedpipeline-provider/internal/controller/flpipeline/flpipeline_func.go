package flpipeline

import (
	"context"
	"fmt"

	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

func (c *external) connect_kube(address string) (*dynamic.DynamicClient, error) {
	//for real implementation download kubeconfig file from gitlab given address

	config := &rest.Config{
		Host: address,
	}

	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		c.logger.Debug("Error in creating clientset")
		c.logger.Debug(err.Error())
		return nil, err
	}
	return clientset, nil

}

func (c *external) get_flclient_data(clientset *dynamic.DynamicClient, name string) (client_data, error) {
	result, err := clientset.Resource(schema.GroupVersionResource{
		Group:    "federatedlearning.providerflclient.crossplane.io",
		Version:  "v1alpha1",
		Resource: "flclients",
	}).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		c.logger.Debug(fmt.Sprintf("Error in getting resource %s", name))
		c.logger.Debug(err.Error())
		return client_data{}, err
	}
	result_byte, _ := result.MarshalJSON()
	c.logger.Debug(string(result_byte))

	obj := result.Object["status"]
	status, _ := obj.(map[string]interface{})["atProvider"].(map[string]interface{})
	c.logger.Debug(fmt.Sprintf("status: %+v", status))
	active_str, _ := status["active"].(string)
	epochs_str, _ := status["epochs"].(string)
	batch_size_str, _ := status["batch_size"].(string)
	//convert to bool
	active, _ := strconv.ParseBool(active_str)
	//convert to int
	epochs, _ := strconv.Atoi(epochs_str)
	batch_size, _ := strconv.Atoi(batch_size_str)

	if active {
		c.logger.Debug("flclient è attivo")
	} else {
		c.logger.Debug("flclient NON è attivo")
	}
	c.logger.Debug(fmt.Sprintf("epochs: %d , batch_size: %d", epochs, batch_size))
	return client_data{epochs: epochs, batch_size: batch_size}, nil
}

func (c *external) get_topic_queue(clientset *dynamic.DynamicClient, name string) (bool, error) {
	result, err := clientset.Resource(schema.GroupVersionResource{
		Group:    "topic.topicprovider.crossplane.io",
		Version:  "v1alpha1",
		Resource: "mqtttopics",
	}).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		c.logger.Debug(fmt.Sprintf("Error in getting resource %s", name))
		c.logger.Debug(err.Error())
		return true, err
	}

	result_bite, _ := result.MarshalJSON()
	c.logger.Debug(string(result_bite))

	obj := result.Object["status"]
	status, _ := obj.(map[string]interface{})["atResource"].(map[string]interface{})
	degraded_str, _ := status["degraded"].(string)

	//convert to bool
	degraded, _ := strconv.ParseBool(degraded_str)

	return degraded, nil
}
