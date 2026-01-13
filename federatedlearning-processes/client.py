#send a request to local host
import requests

# url = 'http://localhost:3001/'
# data = {
#     "port": "3001",
#     "epochs": "1",
#     "batch_size": "32", 
#     "client_id": "0", 
#     "update_period": "100",
#     "topic_name": "sensors_topic",
#     "server_address": "http://localhost:3000/",
#     "broker_address": "mqtt.eclipseprojects.io"#"as-sensiblecity1.cloudmmwunibo.it"
#     }
# response = requests.post(url, json=data)
# print(response.text)


# url = 'http://nodeserf.cloudmmwunibo.it:30001/'
url = "http://dtazzioli-edge.cloudmmwunibo.it:30001/"


data = {
    "port": "3001",
    "epochs": "1",
    "batch_size": "10", 
    "client_id": "0", 
    "update_period": "100.0",
    "topic_name": "provatopicprovider",
    "server_address": "http://nodeserf.cloudmmwunibo.it:3000/",
    "broker_address": "as-sensiblecity1.cloudmmwunibo.it"#"as-sensiblecity1.cloudmmwunibo.it", mqtt.eclipseprojects.io
    }
response = requests.post(url, json=data)
print(response.text)