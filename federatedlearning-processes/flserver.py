from flask import Flask, request
import pickle as pkl
from AI import *
def syntetic_data():
    a = [1,2,3,3,5,2,7,8,6,1]
    e = [2,2,5,3,1,4,2,1,2,3]
    features = np.random.rand(10)
    #polinomial
    label = np.dot(a,np.power(features,e))
    return features, label

app = Flask(__name__)
server_model = RegressionModel()

models_received = []

number_models_to_train = 2
def update_model():
    mean_weights = []
    for i in range(len(models_received[0])):
        mean_weights.append(np.mean([layer[i] for layer in models_received], axis=0))
    server_model.load_weights(mean_weights)
    models_received.clear()
    print("Model updated")


# create an endpoint at http://localhost:/3000/
@app.route('/', methods=['POST'])
def home():
    data = request.get_data()
    data = pkl.loads(data)
    models_received.append(data["weights"])
    if len(models_received) == number_models_to_train:
        update_model()

    #testing model
    features, label = syntetic_data()
    predicted = server_model.predict(features)
    print(f"error {(label-predicted.numpy()[0])**2}")
    return 'Data received successfully!'

@app.route('/', methods=['GET'])
def get():
    return pkl.dumps(server_model.get_weights())


if __name__ == '__main__':
    app.run("0.0.0.0",port=3000,debug=True)