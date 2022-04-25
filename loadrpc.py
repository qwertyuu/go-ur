from pygoridge import RPC, SocketRelay
import numpy as np
import json
from pygoridge.constants import PayloadType
import csv
from numpy import genfromtxt
import shap
my_data = genfromtxt('dataset_AI_Random_7.csv',delimiter=',', skip_header=1)

rpc = RPC(SocketRelay("127.0.0.1", 6001))

#def go_ur_ai(payload: dict) -> dict:
#    return json.loads(rpc("GoUr.Infer", json.dumps(payload)))
#
#ret = go_ur_ai({
#    'pawn_per_player': 3,
#    'ai_pawn_out': 0,
#    'enemy_pawn_out': 0,
#    'dice': 3,
#    'ai_pawn_positions': [1, 2],
#    'enemy_pawn_positions': [1, 2],
#})
#print(ret)

class memfile:
    def __init__(self):
        self.b = []
    def write(self, b):
        self.b += b
    def asbytes(self):
        return bytes(self.b)


def go_ur_ai_numpy(x):
    mm = memfile()
    np.save(mm, x)
    return np.array(json.loads(rpc("GoUr.InferNumpy", mm.asbytes(), PayloadType.PAYLOAD_RAW)))

with open('dataset_AI_Random_7.csv', 'r') as f:
    csv_reader = csv.reader(f)

print(my_data)
print(go_ur_ai_numpy(my_data))

explainer = shap.KernelExplainer(go_ur_ai_numpy, my_data[:20], link="logit")
shap_values = explainer.shap_values(my_data[:20], nsamples=100)
print(shap_values)