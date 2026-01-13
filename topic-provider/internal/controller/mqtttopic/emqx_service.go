package mqtttopic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/provider-topicprovider/apis/topic/v1alpha1"
)

type MqttTopicObservation struct {
	Metrics    Metrics `json:"metrics,omitempty"`
	Topic      string  `json:"topic,omitempty"`
	CreateTime string  `json:"create_time,omitempty"`
}

type Metrics struct {
	MessagesQos2OutCount int     `json:"messages.qos2.out.count,omitempty"`
	MessagesQos2InCount  int     `json:"messages.qos2.in.count,omitempty"`
	MessagesQos1OutCount int     `json:"messages.qos1.out.count,omitempty"`
	MessagesQos1InCount  int     `json:"messages.qos1.in.count,omitempty"`
	MessagesQos0OutCount int     `json:"messages.qos0.out.count,omitempty"`
	MessagesQos0InCount  int     `json:"messages.qos0.in.count,omitempty"`
	MessagesOutCount     int     `json:"messages.out.count,omitempty"`
	MessagesInCount      int     `json:"messages.in.count,omitempty"`
	MessagesDroppedCount int     `json:"messages.dropped.count,omitempty"`
	MessagesQos2OutRate  float64 `json:"messages.qos2.out.rate,omitempty"`
	MessagesQos2InRate   float64 `json:"messages.qos2.in.rate,omitempty"`
	MessagesQos1OutRate  float64 `json:"messages.qos1.out.rate,omitempty"`
	MessagesQos1InRate   float64 `json:"messages.qos1.in.rate,omitempty"`
	MessagesQos0OutRate  float64 `json:"messages.qos0.out.rate,omitempty"`
	MessagesQos0InRate   float64 `json:"messages.qos0.in.rate,omitempty"`
	MessagesOutRate      float64 `json:"messages.out.rate,omitempty"`
	MessagesInRate       float64 `json:"messages.in.rate,omitempty"`
	MessagesDroppedRate  float64 `json:"messages.dropped.rate,omitempty"`
}

func GetTopics(username string, password string, hostAddress string, logger logging.Logger) error {

	url := fmt.Sprintf("http://%s:18083/api/v5/topics", hostAddress)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Info(fmt.Sprintf("[GetTopics]: error while creating the HTTP request: %s", err.Error()))
		return err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Info(fmt.Sprintf("[GetTopics]: error while getting topics: %s", err.Error()))
		return err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		logger.Info(fmt.Sprintf("[GetTopics]: error while getting topics: %s", err.Error()))
		return err
	}

	var data interface{}
	json.Unmarshal(buf.Bytes(), &data)
	return nil
}

func GetTopicMetrics(username string, password string, hostAddress string, topicName string, logger logging.Logger) (v1alpha1.MqttTopicObservation, int, int, error) {

	url := fmt.Sprintf("http://%s:18083/api/v5/mqtt/topic_metrics/%s", hostAddress, topicName)

	// create httprequest
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Info(fmt.Sprintf("[GetTopicMetrics]: error while creating the HTTP request: %s", err.Error()))
		return v1alpha1.MqttTopicObservation{}, -1, -1, err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Info(fmt.Sprintf("[GetTopicMetrics]: error while sending the HTTP request: %s", err.Error()))
		return v1alpha1.MqttTopicObservation{}, -1, -1, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logger.Info(fmt.Sprintf("[GetTopicMetrics]: status code: %d: %s", resp.StatusCode, err.Error()))
		return v1alpha1.MqttTopicObservation{}, -1, -1, err
	}

	// read response body
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		logger.Info(fmt.Sprintf("[GetTopicMetrics]: error while reading the response: %s", err.Error()))
		return v1alpha1.MqttTopicObservation{}, -1, -1, err
	}

	// unmarshal body
	var data MqttTopicObservation
	err = json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		logger.Info(fmt.Sprintf("[GetTopicMetrics]: error while unmarshalling the response: %s", err.Error()))
		return v1alpha1.MqttTopicObservation{}, -1, -1, err
	}
	result := v1alpha1.MqttTopicObservation{
		Topic:      data.Topic,
		CreateTime: data.CreateTime,
		Metrics: v1alpha1.Metrics{
			MessagesQos2OutCount: strconv.Itoa(data.Metrics.MessagesQos2OutCount),
			MessagesQos2InCount:  strconv.Itoa(data.Metrics.MessagesQos2InCount),
			MessagesQos1OutCount: strconv.Itoa(data.Metrics.MessagesQos1OutCount),
			MessagesQos1InCount:  strconv.Itoa(data.Metrics.MessagesQos1InCount),
			MessagesQos0OutCount: strconv.Itoa(data.Metrics.MessagesQos0OutCount),
			MessagesQos0InCount:  strconv.Itoa(data.Metrics.MessagesQos0InCount),
			MessagesOutCount:     strconv.Itoa(data.Metrics.MessagesOutCount),
			MessagesInCount:      strconv.Itoa(data.Metrics.MessagesInCount),
			MessagesDroppedCount: strconv.Itoa(data.Metrics.MessagesDroppedCount),
			MessagesQos2OutRate:  strconv.FormatFloat(data.Metrics.MessagesQos2OutRate, 'f', 5, 64),
			MessagesQos2InRate:   strconv.FormatFloat(data.Metrics.MessagesQos2InRate, 'f', 5, 64),
			MessagesQos1OutRate:  strconv.FormatFloat(data.Metrics.MessagesQos1OutRate, 'f', 5, 64),
			MessagesQos1InRate:   strconv.FormatFloat(data.Metrics.MessagesQos1InRate, 'f', 5, 64),
			MessagesQos0OutRate:  strconv.FormatFloat(data.Metrics.MessagesQos0OutRate, 'f', 5, 64),
			MessagesQos0InRate:   strconv.FormatFloat(data.Metrics.MessagesQos0InRate, 'f', 5, 64),
			MessagesOutRate:      strconv.FormatFloat(data.Metrics.MessagesOutRate, 'f', 5, 64),
			MessagesInRate:       strconv.FormatFloat(data.Metrics.MessagesInRate, 'f', 5, 64),
			MessagesDroppedRate:  strconv.FormatFloat(data.Metrics.MessagesDroppedRate, 'f', 5, 64),
		},
	}

	queue := data.Metrics.MessagesInCount - data.Metrics.MessagesOutCount

	return result, queue, data.Metrics.MessagesDroppedCount, nil
}

func CreateTopicMetric(username string, password string, hostAddress string, topicName string, logger logging.Logger) error {

	url := fmt.Sprintf("http://%s:18083/api/v5/mqtt/topic_metrics", hostAddress)

	ioreader := bytes.NewBuffer([]byte(fmt.Sprintf(`{"topic":"%s"}`, topicName)))

	req, err := http.NewRequest("POST", url, ioreader)
	if err != nil {
		logger.Info(fmt.Sprintf("[CreateTopicMetric]: error while creating the request: %s", err.Error()))
		return err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Info(fmt.Sprintf("[CreateTopicMetric]: error while sending the request: %s", err.Error()))
		return err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		logger.Info(fmt.Sprintf("[CreateTopicMetric]: error while reading the response: %s", err.Error()))
		return err
	}

	var data interface{}
	json.Unmarshal(buf.Bytes(), &data)

	return nil

}

func DeleteTopicMetrics(username string, password string, hostAddress string, topicName string, logger logging.Logger) error {

	url := fmt.Sprintf("http://%s:18083/api/v5/mqtt/topic_metrics/%s", hostAddress, topicName)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		logger.Info(fmt.Sprintf("[DeleteTopicMetrics]: error while creating the request: %s", err.Error()))
		return err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Info(fmt.Sprintf("[DeleteTopicMetrics]: error while sending the request: %s", err.Error()))
		return err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		logger.Info(fmt.Sprintf("[DeleteTopicMetrics]: error while reading the response: %s", err.Error()))
		return err
	}

	var data interface{}
	json.Unmarshal(buf.Bytes(), &data)

	return nil

}

func ResetTopicMetric(username string, password string, hostAddress string, topicName string, logger logging.Logger) error {

	url := fmt.Sprintf("http://%s:18083/api/v5/mqtt/topic_metrics", hostAddress)

	ioreader := bytes.NewBuffer([]byte(`{"action":"reset"}`))

	req, err := http.NewRequest("POST", url, ioreader)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		panic(err)
	}

	var data interface{}
	json.Unmarshal(buf.Bytes(), &data)
	fmt.Println(data)
	return nil
}

func LookupTopicInfoByName(username string, password string, hostAddress string, topicName string, logger logging.Logger) error {

	url := fmt.Sprintf("http://%s:18083/api/v5/topics/%s", hostAddress, topicName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		panic(err)
	}

	var data interface{}
	json.Unmarshal(buf.Bytes(), &data)
	fmt.Println(data)
	return nil

}
