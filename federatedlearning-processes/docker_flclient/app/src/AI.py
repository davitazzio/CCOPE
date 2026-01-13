import tensorflow as tf
import numpy as np
import time

tf.random.set_seed(42)
np.random.seed(42)

class RegressionModel:
    def __init__(self):
        input_size = 10
        self.model = tf.keras.models.Sequential(
            [
                tf.keras.layers.Input(shape=(input_size,)),
                tf.keras.layers.Dense(64, activation="sigmoid"),
                tf.keras.layers.Dense(64, activation="sigmoid"),
                tf.keras.layers.Dense(1, activation="linear"),
            ]
        )
        self.optimizer = tf.keras.optimizers.Adam()

    def train_step(self, data, label):
        with tf.GradientTape() as tape:
            t_start = time.time_ns()
            # data = np.expand_dims(data, axis=0)
            prediction = self.model(data, training=True)
            time_inference = time.time_ns() - t_start
            t = time.time_ns()
            loss = tf.reduce_mean(tf.square(label - prediction))
            grads = tape.gradient(loss, self.model.trainable_variables)
            self.optimizer.apply_gradients(zip(grads, self.model.trainable_variables))
            time_train = time.time_ns() - t
            # print("{},{},{}".format(t_start, time_train, time_inference))
        return prediction, loss

    def predict(self, data):
        t = time.time()
        data = np.expand_dims(data, axis=0)
        output = self.model(data)
        # print("Inference time: ", time.time() - t)
        return output
    
    def get_weights(self):
        return self.model.get_weights()
    

    def load_weights(self, weights):
        self.model.set_weights(weights)