import numpy as np


class Metrics:
    def __init__(self):
        self.loss = []
        self.average_loss = 0
        self.throughput = []
        self.average_throughput = 0

    def add_loss(self, loss):
        self.loss.append(loss)
        if len(self.loss) > 10:
            self.loss.pop(0)

    def get_metrics(self):
        if len(self.loss) == 0:
            return {
                "loss": "0",
                "average_loss_5": "0",
                "average_loss_10": "0",
            }
        return {
            "loss": str(np.float64(self.loss[-1])),
            "average_loss_5": str(np.mean(self.loss[-5:],dtype=np.float64)),
            "average_loss_10": str(np.mean(self.loss,dtype=np.float64)),
        }
