# provider-topicprovider

`topicprovider` is a [Crossplane](https://crossplane.io/) Provider that manage the lifecycle of a topic, from the creation, observation and deletion. 

The `MqttTopic` resource manage the external Mqtt topic on a EMQX broker. Each resource requires:
- `Name`: the topic name
- `Host`: the Address of the EMQX broker server
- `Username`: the username for a user on the broker with the configuration privileges. 
- `Password`: the user password.

The resource state reports the topic metrics: 
- `MessagesQos2OutCount`
- `MessagesQos2InCount`
- `MessagesQos1OutCount`
- `MessagesQos1InCount`
- `MessagesQos0OutCount`
- `MessagesQos0InCount`
- `MessagesOutCount`
- `MessagesInCount`
- `MessagesDroppedCount`
- `MessagesQos2OutRate`
- `MessagesQos2InRate`
- `MessagesQos1OutRate`
- `MessagesQos1InRate`
- `MessagesQos0OutRate`
- `MessagesQos0InRate `
- `MessagesOutRate`
- `MessagesInRate`
- `MessagesDroppedRate`

