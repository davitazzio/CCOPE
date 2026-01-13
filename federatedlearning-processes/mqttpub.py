import time
import paho.mqtt.client as mqtt
import numpy as np
import argparse

def syntetic_data():
    a = [1,2,3,3,5,2,7,8,6,1]
    e = [2,2,5,3,1,4,2,1,2,3]
    features = np.random.rand(10)
    #polinomial
    label = np.dot(a,np.power(features,e))
    return ",".join(map(str, features))+","+str(label)


def on_publish(client, userdata, mid, reason_code, properties):
    # reason_code and properties will only be present in MQTTv5. It's always unset in MQTTv3
    try:
        userdata.remove(mid)
    except KeyError:
        print("on_publish() is called with a mid not present in unacked_publish")
        print("This is due to an unavoidable race-condition:")
        print("* publish() return the mid of the message sent.")
        print("* mid from publish() is added to unacked_publish by the main thread")
        print("* on_publish() is called by the loop_start thread")
        print("While unlikely (because on_publish() will be called after a network round-trip),")
        print(" this is a race-condition that COULD happen")
        print("")
        print("The best solution to avoid race-condition is using the msg_info from publish()")
        print("We could also try using a list of acknowledged mid rather than removing from pending list,")
        print("but remember that mid could be re-used !")



def main():
    #get args from command line
    parser = argparse.ArgumentParser()
    parser.add_argument("--broker", help="broker address")
    parser.add_argument("--topic", help="topic name")
    args = parser.parse_args()

    unacked_publish = set()
    mqttc = mqtt.Client(mqtt.CallbackAPIVersion.VERSION2)
    mqttc.on_publish = on_publish

    mqttc.user_data_set(unacked_publish)
    mqttc.connect(args.broker) #192.168.17.48:18083 server address  mqtt.eclipseprojects.io, as-sensiblecity1.cloudmmwunibo.it
    #set client identifier
    mqttc.loop_start()

    # Our application produce some messages
    for i in range(10000):
        msg_info = mqttc.publish(args.topic, syntetic_data(), qos=1)
        unacked_publish.add(msg_info.mid)
        time.sleep(0.1)

    # for i in range(10):
    #     msg_info2 = mqttc.publish("prova_topic", "message2", qos=1)
    #     unacked_publish.add(msg_info2.mid)

    # Wait for all message to be published
    while len(unacked_publish):
        time.sleep(0.1)

    # Due to race-condition described above, the following way to wait for all publish is safer
    msg_info.wait_for_publish()
    # msg_info2.wait_for_publish()

    mqttc.disconnect()
    mqttc.loop_stop()

if __name__ == "__main__":
    main()