import paho.mqtt.client as mqtt
from flask import Flask, request, jsonify
from flask.logging import default_handler
import requests
from AI import *
from metrics import *
import pickle as pkl
from flask_executor import Executor
import argparse
import logging

app = Flask(__name__)
executor = Executor(app)
root = logging.getLogger()
root.addHandler(default_handler)
root.setLevel(logging.DEBUG)


getting_data = False

parameters = {
    "port": "3001",
    "epochs": 1,
    "batch_size": 1, 
    "update_period": 1, 
    "client_id": "0",
    "topic_name": "prova",
    "server_address": "http://localhost:3000/",
    "broker_address": "as-sensiblecity1.cloudmmwunibo.it"
    }

metrics = Metrics()

def on_subscribe(client, userdata, mid, reason_code_list, properties):
    # Since we subscribed only for a single channel, reason_code_list contains
    # a single entry
    if reason_code_list[0].is_failure:
        logging.debug(f"Broker rejected you subscription: {reason_code_list[0]}")
    else:
        logging.debug(f"Broker granted the following QoS: {reason_code_list[0].value}")

def on_unsubscribe(client, userdata, mid, reason_code_list, properties):
    # Be careful, the reason_code_list is only present in MQTTv5.
    # In MQTTv3 it will always be empty
    if len(reason_code_list) == 0 or not reason_code_list[0].is_failure:
        logging.debug("unsubscribe succeeded (if SUBACK is received in MQTTv3 it success)")
    else:
        logging.debug(f"Broker replied with failure: {reason_code_list[0]}")
    client.disconnect()

def on_message(client, userdata, message):
    # userdata is the structure we choose to provide, here it's a list()
    userdata.append(message.payload)
    
    # We only want to process 10 messages
    if len(userdata) >= parameters["batch_size"]:
        client.unsubscribe(parameters["topic_name"])

def on_connect(client, userdata, flags, reason_code, properties):
    if reason_code.is_failure:
        logging.debug(f"Failed to connect: {reason_code}. loop_forever() will retry connection")
    else:
        # we should always subscribe from on_connect callback to be sure
        # our subscribed is persisted across reconnections.
        client.subscribe(parameters["topic_name"], qos=1)

mqttc = mqtt.Client(mqtt.CallbackAPIVersion.VERSION2)
mqttc.on_connect = on_connect
mqttc.on_message = on_message
mqttc.on_subscribe = on_subscribe
mqttc.on_unsubscribe = on_unsubscribe

#network
batch_size = 100
model = RegressionModel()

def data_preprocessing(data):
    data = [x.decode("utf-8").split(",") for x in data]
    data = np.array(data, dtype=np.float32)
    # We are assuming that the first 10 values are the input and the last one is the label
    return data[:, :-1], data[:, -1]

def get_data_and_train():
    logging.debug("start getting data loop")
    while True:
        mqttc.user_data_set([])
        # logging.debug(f"connecting to {parameters['broker_address']}")
        #connect with username and password
        # mqttc.username_pw_set(username="admin", password="public")
        mqttc.connect(parameters['broker_address'])
        mqttc.loop_forever()
        logging.debug("Data received")
        x, y = data_preprocessing(mqttc.user_data_get())
        logging.debug(f"Data preprocessed, x shape: {x.shape}, y shape: {y.shape}")
        for epoch in range(parameters["epochs"]):
            prediction, loss = model.train_step(x,y)

        send_to_server = {}
        send_to_server["prediction"] = prediction.numpy().tolist()
        send_to_server["loss"] = loss.numpy().tolist()
        send_to_server["weights"] = model.get_weights()
        send_to_server["id"] = parameters["client_id"]
        #has to send also the number of datapoints used for training
        #possibility of sending only gradients?

        logging.debug(f"Loss:{np.mean(loss.numpy())} ")

        metrics.add_loss(np.mean(loss.numpy()))

        compressed_data = pkl.dumps(send_to_server)
        logging.debug("Sending data to server")
        response = requests.post(parameters["server_address"], data=compressed_data,timeout=5)
        if response.status_code != 200:
            logging.debug(f"Error in sending data to server, status code: {response.status_code}")
        else:
            logging.debug(response.text)

#must be called to start training
@app.route('/', methods=['POST'])
def home():
    global parameters
    global getting_data
    data = request.get_json()
    logging.debug(data)
    parameters["port"] = int(data["port"])
    parameters["epochs"] = int(data["epochs"])
    parameters["batch_size"] = int(data["batch_size"])
    parameters["client_id"] = str(data["client_id"])
    parameters["topic_name"] = str(data["topic_name"])
    parameters["update_period"] = float(data["update_period"])
    parameters["broker_address"] = str(data["broker_address"])
    parameters["server_address"] = str(data["server_address"])
    #if get_data_and_train is already running, don't start another one
    if not getting_data:
        executor.submit(get_data_and_train)
        getting_data = True

    return f"parameters received, parameters: {parameters}"


@app.route('/', methods=['GET'])
def get():
    global parameters
    #return loss, and average loss 5, average loss 10 and state of the client
    response = metrics.get_metrics()
    response.update({k: str(v) for k, v in parameters.items()})
    return jsonify(response)



if __name__ == '__main__':  
    argparser = argparse.ArgumentParser()
    argparser.add_argument("--port", type=int, default=3001)
    args = argparser.parse_args()
    parameters["port"] = args.port
    app.run(host='0.0.0.0', port=args.port, debug=True)


    

    