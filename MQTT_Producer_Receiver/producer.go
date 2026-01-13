package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// const TOPIC = "$share/gruppoKAFKA/mqtt/test"
const TOPIC = "test"
const USERNAME = "297647beea8d0111"
const PASSWORD = "lJI9A9CszbxfpBmXAmNtjnCtVaWLvt187SGcps9C9AAI8MH"
const QOS = 1

func DeleteTopicMetrics() {

	var topic string
	if strings.Contains(TOPIC, "$share") {
		splitTopic := strings.SplitN(TOPIC, "/", 3)

		if len(splitTopic) > 1 {
			topic = splitTopic[2]
		} else {
			topic = TOPIC
		}
	} else {
		topic = TOPIC
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

func CreateTopicMetric() {

	url := "http://dtazzioli-edge.cloudmmwunibo.it:18083/api/v5/mqtt/topic_metrics"
	var topic string
	if strings.Contains(TOPIC, "$share") {
		splitTopic := strings.SplitN(TOPIC, "/", 3)

		if len(splitTopic) > 1 {
			topic = splitTopic[2]
		} else {
			topic = TOPIC
		}
	} else {
		topic = TOPIC
	}

	fmt.Printf("Topic: %s\n", topic)

	requestJSON := fmt.Sprintf(`{"topic":"%s"}`, topic)

	ioreader := bytes.NewBuffer([]byte(requestJSON))

	req, err := http.NewRequest("POST", url, ioreader)
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

	fmt.Printf("%d", resp.StatusCode)
	var data interface{}
	json.Unmarshal(buf.Bytes(), &data)
	fmt.Println(data)
}

func GetTopicMetrics() {
	var topic string
	if strings.Contains(TOPIC, "$share") {
		splitTopic := strings.SplitN(TOPIC, "/", 3)

		if len(splitTopic) > 1 {
			topic = splitTopic[2]
		} else {
			topic = TOPIC
		}
	} else {
		topic = TOPIC
	}
	fmt.Printf("Topic: %s\n", topic)

	url := fmt.Sprintf("http://dtazzioli-edge.cloudmmwunibo.it:18083/api/v5/mqtt/topic_metrics/%s", topic)

	req, err := http.NewRequest("GET", url, nil)
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

	fmt.Println(resp.StatusCode)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		panic(err)
	}

	var data interface{}
	json.Unmarshal(buf.Bytes(), &data)
	fmt.Println(data)
}

func publish(clientID string) {

	index := 0
	sleep_times := []int{1000, 666, 500, 400, 333, 285, 250, 222, 200, 182, 167, 154, 142, 133, 125, 118, 111, 105, 100}

	clientNum, err := strconv.Atoi(strings.TrimPrefix(clientID, "pub-"))
	if err != nil {
		panic(err)
	}

	opts := mqtt.NewClientOptions().AddBroker("tcp://dtazzioli-edge.cloudmmwunibo.it:1883").SetClientID(clientID)
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	go func() {
		for range time.Tick(3 * time.Second) {
			index++
			fmt.Printf("%f messages per second (sleep: %d)\n", 1+(float64(index)*0.5)+float64((clientNum*10)), sleep_times[index])
			if index >= len(sleep_times)-1 {
				return
			}
		}
	}()

	for {
		var token mqtt.Token
		if token = c.Publish(TOPIC, QOS, false, time.Now().Format(time.RFC3339Nano)); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}
		time.Sleep(time.Duration(sleep_times[index]) * time.Millisecond)

	}
}

func publish2(clientID string, speed chan int) {

	sleep := 200
	opts := mqtt.NewClientOptions().AddBroker("tcp://dtazzioli-edge.cloudmmwunibo.it:1883").SetClientID(clientID)
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		select {
		case <-time.After(time.Duration(sleep) * time.Millisecond):
			{
				var token mqtt.Token
				if token = c.Publish(TOPIC, QOS, false, time.Now().Format(time.RFC3339Nano)); token.Wait() && token.Error() != nil {
					fmt.Println(token.Error())
				}
			}
		case s := <-speed:
			{
				if s < 0 {
					return
				}
				sleep = s
			}
		}
	}
}

func main() {
	DeleteTopicMetrics()
	// getMqttConfiguratio("dtazzioli-edge.cloudmmwunibo.it", USERNAME, PASSWORD)
	// configureBroker("dtazzioli-edge.cloudmmwunibo.it", USERNAME, PASSWORD)
	CreateTopicMetric()

	// fmt.Println("Produzione 5msg/sec")

	speed_channel := make(chan int)
	go publish2("pub-0", speed_channel)

	// // speed_channel2 := make(chan int)
	// // go publish2("pub-1", speed_channel2)

	// // speed_channel3 := make(chan int)
	// // go publish2("pub-2", speed_channel3)

	// // speed_channel4 := make(chan int)
	// // go publish2("pub-3", speed_channel4)

	// // speed_channel5 := make(chan int)
	// // go publish2("pub-4", speed_channel5)

	// // speed_channell6 := make(chan int)
	// // go publish2("pub-5", speed_channell6)

	time.Sleep(30 * time.Second)

	// fmt.Println("Produzione 15msg/sec")
	// speed_channel <- 70
	// time.Sleep(1 * time.Second)
	// // fmt.Println("Produzione 28msg/sec")
	// // speed_channel2 <- 200
	// // time.Sleep(1 * time.Second)
	// // // fmt.Println("Produzione 36msg/sec")
	// // speed_channel3 <- 200
	// // time.Sleep(1 * time.Second)
	// // // fmt.Println("Produzione 60msg/sec")
	// // speed_channel4 <- -1
	// // speed_channel5 <- -1
	// // speed_channell6 <- 80

	// time.Sleep(360 * time.Second)

	// fmt.Println("END")

}

func configureBroker(hostAddress, username, password string) error {

	url := fmt.Sprintf("http://%s:18083/api/v5/configs/global_zone", hostAddress)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")
	byteBody, err := json.Marshal(fillRequest())
	if err != nil {
		fmt.Printf("error in marshalling body: %s", err.Error())
		return err
	}
	req.Body = io.NopCloser(bytes.NewReader(byteBody))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("error in configuring broker: %d", resp.StatusCode)
		// return errors.New("error in configuring broker")
	}

	// read response body
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		fmt.Printf("error in reading response body: %s", err.Error())
		return err
	} else {
		fmt.Println(buf.String())
	}

	// unmarshal body

	return nil
}

type ConfigBrokerRequest struct {
	Mqtt            Mqtt            `json:"mqtt"`
	FlappingDetect  FlappingDetect  `json:"flapping_detect"`
	ForceShutdown   ForceShutdown   `json:"force_shutdown"`
	ForceGc         ForceGc         `json:"force_gc"`
	DurableSessions DurableSessions `json:"durable_sessions"`
}
type Mqtt struct {
	IdleTimeout                string `json:"idle_timeout"`
	MaxPacketSize              string `json:"max_packet_size"`
	MaxClientidLen             int    `json:"max_clientid_len"`
	MaxTopicLevels             int    `json:"max_topic_levels"`
	MaxTopicAlias              int    `json:"max_topic_alias"`
	RetainAvailable            bool   `json:"retain_available"`
	WildcardSubscription       bool   `json:"wildcard_subscription"`
	SharedSubscription         bool   `json:"shared_subscription"`
	SharedSubscriptionStrategy string `json:"shared_subscription_strategy"`
	// SharedSubscriptionInitialStickyPick string  `json:"shared_subscription_initial_sticky_pick"`
	ExclusiveSubscription  bool    `json:"exclusive_subscription"`
	IgnoreLoopDeliver      bool    `json:"ignore_loop_deliver"`
	StrictMode             bool    `json:"strict_mode"`
	ResponseInformation    string  `json:"response_information"`
	ServerKeepalive        string  `json:"server_keepalive"`
	KeepaliveMultiplier    float64 `json:"keepalive_multiplier"`
	KeepaliveCheckInterval string  `json:"keepalive_check_interval"`
	RetryInterval          string  `json:"retry_interval"`
	UseUsernameAsClientid  bool    `json:"use_username_as_clientid"`
	PeerCertAsUsername     string  `json:"peer_cert_as_username"`
	PeerCertAsClientid     string  `json:"peer_cert_as_clientid"`
	ClientAttrsInit        []any   `json:"client_attrs_init"`
	SessionExpiryInterval  string  `json:"session_expiry_interval"`
	MessageExpiryInterval  string  `json:"message_expiry_interval"`
	MaxAwaitingRel         int     `json:"max_awaiting_rel"`
	MaxQosAllowed          int     `json:"max_qos_allowed"`
	MqueuePriorities       string  `json:"mqueue_priorities"`
	MqueueDefaultPriority  string  `json:"mqueue_default_priority"`
	MqueueStoreQos0        bool    `json:"mqueue_store_qos0"`
	MaxMqueueLen           int     `json:"max_mqueue_len"`
	MaxInflight            int     `json:"max_inflight"`
	MaxSubscriptions       string  `json:"max_subscriptions"`
	UpgradeQos             bool    `json:"upgrade_qos"`
	AwaitRelTimeout        string  `json:"await_rel_timeout"`
}
type FlappingDetect struct {
	Enable     bool   `json:"enable"`
	WindowTime string `json:"window_time"`
	MaxCount   int    `json:"max_count"`
	BanTime    string `json:"ban_time"`
}
type ForceShutdown struct {
	Enable         bool   `json:"enable"`
	MaxMailboxSize int    `json:"max_mailbox_size"`
	MaxHeapSize    string `json:"max_heap_size"`
}
type ForceGc struct {
	Enable bool   `json:"enable"`
	Count  int    `json:"count"`
	Bytes  string `json:"bytes"`
}
type DurableSessions struct {
	Enable                 bool   `json:"enable"`
	BatchSize              int    `json:"batch_size"`
	IdlePollInterval       string `json:"idle_poll_interval"`
	HeartbeatInterval      string `json:"heartbeat_interval"`
	SessionGcInterval      string `json:"session_gc_interval"`
	SessionGcBatchSize     int    `json:"session_gc_batch_size"`
	MessageRetentionPeriod string `json:"message_retention_period"`
}

func fillRequest() *ConfigBrokerRequest {
	return &ConfigBrokerRequest{
		Mqtt: Mqtt{
			IdleTimeout:                "15s",
			MaxPacketSize:              "32MB",
			MaxClientidLen:             65535,
			MaxTopicLevels:             128,
			MaxTopicAlias:              65535,
			RetainAvailable:            true,
			WildcardSubscription:       true,
			SharedSubscription:         true,
			SharedSubscriptionStrategy: "random",
			// SharedSubscriptionInitialStickyPick: "random",
			ExclusiveSubscription:  false,
			IgnoreLoopDeliver:      false,
			StrictMode:             false,
			ResponseInformation:    "",
			ServerKeepalive:        "disabled",
			KeepaliveMultiplier:    1.5,
			KeepaliveCheckInterval: "12m",
			RetryInterval:          "1s",
			UseUsernameAsClientid:  false,
			PeerCertAsUsername:     "disabled",
			PeerCertAsClientid:     "disabled",
			ClientAttrsInit:        []any{},
			SessionExpiryInterval:  "12m",
			MessageExpiryInterval:  "infinity",
			MaxAwaitingRel:         100,
			MaxQosAllowed:          2,
			MqueuePriorities:       "disabled",
			MqueueDefaultPriority:  "highest",
			MqueueStoreQos0:        true,
			MaxMqueueLen:           100,
			MaxInflight:            65535,
			MaxSubscriptions:       "infinity",
			UpgradeQos:             false,
			AwaitRelTimeout:        "12m",
		},
		FlappingDetect: FlappingDetect{
			Enable:     false,
			WindowTime: "12m",
			MaxCount:   15,
			BanTime:    "12m",
		},
		ForceShutdown: ForceShutdown{
			Enable:         true,
			MaxMailboxSize: 1000,
			MaxHeapSize:    "1024KB",
		},
		ForceGc: ForceGc{
			Enable: true,
			Count:  16000,
			Bytes:  "32MB",
		},
		DurableSessions: DurableSessions{
			Enable:                 false,
			BatchSize:              100,
			IdlePollInterval:       "12m",
			HeartbeatInterval:      "12m",
			SessionGcInterval:      "12m",
			SessionGcBatchSize:     100,
			MessageRetentionPeriod: "12m",
		},
	}
}

func getMqttConfiguratio(hostAddress, username, password string) error {
	url := fmt.Sprintf("http://%s:18083/api/v5/configs/global_zone/", hostAddress)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		fmt.Printf("error in reading broker configuration: %d", resp.StatusCode)
		// return errors.New("error in configuring broker")
	}

	// read response body
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		fmt.Printf("error in reading response body: %s", err.Error())
		return err
	} else {
		fmt.Println(buf.String())
	}

	// unmarshal body

	var config ConfigBrokerRequest
	err = json.Unmarshal(buf.Bytes(), &config)
	if err != nil {
		fmt.Printf("error in unmarshalling response body: %s", err.Error())
		return err
	}
	fmt.Printf("Broker Configuration: %+v\n", config)
	fmt.Printf("Broker Configuration: %+v\n", config.Mqtt)
	return nil
}
