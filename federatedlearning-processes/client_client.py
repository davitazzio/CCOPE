import requests

url = 'http://dtazzioli-processprovider.cloudmmwunibo.it:30001/'

response = requests.get(url)


print(response.text)