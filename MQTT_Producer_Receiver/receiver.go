package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	kafka_go "github.com/segmentio/kafka-go"
)

// const TOPIC = "$share/gruppoKAFKA/test"
const TOPIC = "$share/gruppoKAFKA/test"
const USERNAME = "297647beea8d0111"
const PASSWORD = "lJI9A9CszbxfpBmXAmNtjnCtVaWLvt187SGcps9C9AAI8MH"
const QOS = 1
const MQTT_HOST = "dtazzioli-edge.cloudmmwunibo.it"
const KAFKA_HOST = "dtazzioli-edge.cloudmmwunibo.it"

var (
	counter_kafka *int
	old_kafka     *int
	counter_mqtt  *int
	old_mqtt      *int
)

var (
	buffer_kafka *[]int64
	buffer_mqtt  *[]int64
)

var done_local1 = false
var done_remote = false

func kafkaConsume(topic string) {

	if strings.HasPrefix(topic, "$") {
		// Remove the beginning of the topic until "/"
		splitTopic := strings.SplitN(topic, "/", 3)
		if len(splitTopic) > 1 {
			topic = splitTopic[2]
		}
	}
	partition := 0
	conn, err := kafka_go.DialLeader(context.Background(), "tcp", KAFKA_HOST, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	log.Default().Printf("Connected to Kafka")

	for {
		msg, err := conn.ReadMessage(1024)
		if err != nil {
			log.Fatal("failed to read message:", err)
		}

		sent_time, err := time.Parse(time.RFC3339Nano, string(msg.Value))
		if err != nil {
			log.Default().Printf("failed to parse time: %s", err.Error())
			continue
		}

		*counter_kafka++

		processing_time := time.Since(sent_time).Nanoseconds()
		*buffer_kafka = append(*buffer_kafka, processing_time)
		file, err := os.OpenFile("kafka_real_time.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("failed to open file:", err)
		}

		if _, err := file.WriteString(fmt.Sprintf("%s;%s;%d\n", sent_time.String(), time.Now().String(), processing_time)); err != nil {
			log.Fatal("failed to write to file:", err)

		}

	}
}

func consumer_newbroker(name string) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://dtazzioli-edge.cloudmmwunibo.it:1883").SetClientID(name)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	c.OptionsReader()
	msgRcvd := func(client mqtt.Client, msg mqtt.Message) {

		sent_time, err := time.Parse(time.RFC3339Nano, string(msg.Payload()))
		if err != nil {
			log.Default().Printf("failed to parse time: %s", err.Error())
		} else {
			*counter_mqtt++
			processing_time := time.Since(sent_time).Nanoseconds()
			*buffer_mqtt = append(*buffer_mqtt, processing_time)

			// se la latenza sale sopra al secondo significa che un nodo è saturo
			if processing_time > 1000000000 && !done_local1 && !done_remote {
				go func() {
					cmd := exec.Command("/bin/bash", "/home/dtazzioli/infrastructure.sh")
					if err := cmd.Run(); err != nil {
						log.Default().Printf("failed to execute ssh command: %s", err.Error())
					}
				}()
				done_remote = true
				done_local1 = true
			}

			file, err := os.OpenFile("mqtt_real_time.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal("failed to open file:", err)
			}
			defer file.Close()

			if _, err := file.WriteString(fmt.Sprintf("%s;%s;%d\n", sent_time.String(), time.Now().String(), processing_time)); err != nil {
				log.Fatal("failed to write to file:", err)

			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	if token := c.Subscribe(TOPIC, QOS, msgRcvd); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}

func main() {

	DeleteTopicMetrics(TOPIC)

	tmp := 0
	tmp2 := 0
	tmp3 := 0
	tmp4 := 0
	counter_kafka = &tmp
	old_kafka = &tmp2
	counter_mqtt = &tmp3
	old_mqtt = &tmp4
	tmp5 := make([]int64, 0)
	tmp6 := make([]int64, 0)
	buffer_kafka = &tmp5
	buffer_mqtt = &tmp6

	go kafkaConsume(TOPIC)
	// go consumer_newbroker("consumer-1")
	// go consumer_newbroker("consumer-2")
	// go consumer_newbroker("consumer-3")
	// go consumer_newbroker("consumer-4")
	go collect_metrics(TOPIC)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigc
		log.Printf("Signal received: %v\n", s)
		messages_received := *counter_kafka
		log.Printf("Messages received: %d\n", messages_received)

		file, err := os.OpenFile("kafka_messages.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("failed to open file:", err)
		}
		defer file.Close()
		for _, value := range *buffer_kafka {
			if _, err := file.WriteString(fmt.Sprintf("%d\n", value)); err != nil {
				log.Fatal("failed to write to file:", err)
			}
		}

		file, err = os.OpenFile("mqtt_messages.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("failed to open file:", err)
		}
		defer file.Close()
		for _, value := range *buffer_mqtt {
			if _, err := file.WriteString(fmt.Sprintf("%d\n", value)); err != nil {
				log.Fatal("failed to write to file:", err)
			}
		}
		os.Exit(0)
	}()

	select {}

}

//#####################################################################################################

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
	ClientConnected      int     `json:"client.subscribe,omitempty"`
}
type MqttTopicObservation struct {
	Metrics    Metrics `json:"metrics,omitempty"`
	Topic      string  `json:"topic,omitempty"`
	CreateTime string  `json:"create_time,omitempty"`
}

func GetTopicMetrics(username string, password string, hostAddress string, topic string) (MqttTopicObservation, int, int, error) {

	if strings.HasPrefix(topic, "$") {
		// Remove the beginning of the topic until "/"
		splitTopic := strings.SplitN(topic, "/", 3)
		if len(splitTopic) > 1 {
			topic = splitTopic[2]

		}

	}

	fmt.Printf("Topic: %s\n", topic)

	url := fmt.Sprintf("http://%s:18083/api/v5/mqtt/topic_metrics/%s", hostAddress, topic)

	// create httprequest

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return MqttTopicObservation{}, -1, -1, err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return MqttTopicObservation{}, -1, -1, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return MqttTopicObservation{}, -1, -1, errors.New("the resource does not exist")
	}

	// read response body
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return MqttTopicObservation{}, -1, -1, err
	}

	// unmarshal body
	var data MqttTopicObservation
	// fmt.Printf("Data: %s\n", buf.String())
	err = json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		return MqttTopicObservation{}, -1, -1, err
	}

	queue := data.Metrics.MessagesInCount - data.Metrics.MessagesOutCount - data.Metrics.MessagesDroppedCount

	return data, queue, data.Metrics.MessagesDroppedCount, nil
}

type AutoGenerated struct {
	DroppedMsgRate              int `json:"dropped_msg_rate"`
	SentMsgRate                 int `json:"sent_msg_rate"`
	PersistedRate               int `json:"persisted_rate"`
	ReceivedMsgRate             int `json:"received_msg_rate"`
	TransformationFailedRate    int `json:"transformation_failed_rate"`
	TransformationSucceededRate int `json:"transformation_succeeded_rate"`
	ValidationFailedRate        int `json:"validation_failed_rate"`
	ValidationSucceededRate     int `json:"validation_succeeded_rate"`
	DisconnectedDurableSessions int `json:"disconnected_durable_sessions"`
	SubscriptionsDurable        int `json:"subscriptions_durable"`
	Subscriptions               int `json:"subscriptions"`
	Topics                      int `json:"topics"`
	Connections                 int `json:"connections"`
	LiveConnections             int `json:"live_connections"`
	RetainedMsgCount            int `json:"retained_msg_count"`
	SharedSubscriptions         int `json:"shared_subscriptions"`
}

func GetMetrics(username string, password string, hostAddress string, topic string) (AutoGenerated, error) {

	// if strings.HasPrefix(topic, "$") {
	// 	// Remove the beginning of the topic until "/"
	// 	splitTopic := strings.SplitN(topic, "/", 3)
	// 	if len(splitTopic) > 1 {
	// 		topic = splitTopic[2]

	// 	}

	// }

	url := fmt.Sprintf("http://%s:18083/api/v5/monitor_current", hostAddress)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return AutoGenerated{}, err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AutoGenerated{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return AutoGenerated{}, errors.New("the resource does not exist")
	}

	// read response body
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return AutoGenerated{}, err
	}

	// unmarshal body
	var data AutoGenerated
	err = json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return AutoGenerated{}, err
	}

	return data, nil
}

func collect_metrics(topic string) {

	if strings.HasPrefix(topic, "$") {
		// Remove the beginning of the topic until "/"
		splitTopic := strings.SplitN(topic, "/", 3)
		if len(splitTopic) > 1 {
			topic = splitTopic[2]

		}

	}
	fmt.Printf("Topic: %s\n", topic)

	for {
		data, _, _, err := GetTopicMetrics(USERNAME, PASSWORD, MQTT_HOST, topic)
		if err != nil {
			log.Default().Printf("Error: %s\n", err.Error())
		}
		data2, err := GetMetrics(USERNAME, PASSWORD, MQTT_HOST, topic)
		if err != nil {
			log.Default().Printf("Error: %s\n", err.Error())
		}
		time.Sleep(1000 * time.Millisecond)

		file, err := os.OpenFile("metrics.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("failed to open file:", err)
		}

		queue := data.Metrics.MessagesInCount - data.Metrics.MessagesOutCount - data.Metrics.MessagesDroppedCount

		if queue > 0 && !done_local1 {
			// if data.Metrics.MessagesDroppedCount > 0 && !done_local {
			// go func() {

			// go consumer_newbroker("consumer2")
			// time.Sleep(5 * time.Second)
			// go consumer_newbroker("consumer3")

			// }()
			done_local1 = true
		}

		if queue > 0 && done_local1 && !done_remote {
			go func() {
				cmd := exec.Command("/bin/bash", "/home/dtazzioli/infrastructure.sh")
				if err := cmd.Run(); err != nil {
					log.Default().Printf("failed to execute ssh command: %s", err.Error())
				}
			}()
			done_remote = true
		}
		file.WriteString(fmt.Sprintf("%d;%d;%d;%d;%d;%f;%f;%f;%d;%d;%d;%d;%t\n",
			data.Metrics.MessagesInCount,
			data.Metrics.MessagesOutCount,
			data.Metrics.MessagesDroppedCount,
			*counter_kafka,
			*counter_mqtt,

			data.Metrics.MessagesInRate,
			data.Metrics.MessagesOutRate,
			data.Metrics.MessagesDroppedRate,
			(*counter_kafka - *old_kafka),
			(*counter_mqtt - *old_mqtt),

			queue,
			data2.Subscriptions,
			done_remote,
		))

		*old_kafka = *counter_kafka
		*old_mqtt = *counter_mqtt

		getMqueueMessages(MQTT_HOST, "java-sub-1")
		getMqueueMessages(MQTT_HOST, "consumer-2")
		getMqueueMessages(MQTT_HOST, "consumer-3")
		getMqueueMessages(MQTT_HOST, "consumer-4")
	}
}

func DeleteTopicMetrics(topic string) {

	if strings.HasPrefix(topic, "$") {
		// Remove the beginning of the topic until "/"
		splitTopic := strings.SplitN(topic, "/", 3)
		if len(splitTopic) > 1 {
			topic = splitTopic[2]
		}
	}

	fmt.Printf("Topic: %s\n", topic)

	url := fmt.Sprintf("http://dtazzioli-edge.cloudmmwunibo.it:18083/api/v5/mqtt/topic_metrics/%s", topic)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(USERNAME, PASSWORD)
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

}

func getMqueueMessages(hostAddress, clientID string) error {
	url := fmt.Sprintf("http://%s:18083/api/v5/clients/%s/mqueue_messages", hostAddress, clientID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(USERNAME, PASSWORD)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("the resource does not exist")
	}

	// read response body
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	// unmarshal body
	var data ClientData
	//
	//
	err = json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return err
	}

	fmt.Printf("Mqueue Messages for %s: %d\n", clientID, len(data.Data))

	return nil

}

type ClientData struct {
	Meta struct {
		Start    string `json:"start"`
		Position string `json:"position"`
		Count    int    `json:"count"`
		Last     string `json:"last"`
	} `json:"meta"`
	Data []struct {
		InsertedAt     string `json:"inserted_at"`
		PublishAt      int64  `json:"publish_at"`
		FromClientid   string `json:"from_clientid"`
		FromUsername   string `json:"from_username"`
		Msgid          string `json:"msgid"`
		MqueuePriority int    `json:"mqueue_priority"`
		Qos            int    `json:"qos"`
		Topic          string `json:"topic"`
		Payload        string `json:"payload"`
	} `json:"data"`
}
