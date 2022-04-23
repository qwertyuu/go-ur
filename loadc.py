import ctypes
import json

so = ctypes.cdll.LoadLibrary('./go-ur_infer.so')
infer = so.infer
infer.argtypes = [ctypes.c_char_p]
infer.restype = ctypes.c_void_p

def go_ur_ai(a: dict) -> dict:
    ptr = infer(json.dumps(a).encode('utf-8'))
    out = ctypes.string_at(ptr)
    return json.loads(out.decode('utf-8'))

ret = go_ur_ai({
    'pawn_per_player': 3,
    'ai_pawn_out': 0,
    'enemy_pawn_out': 0,
    'dice': 3,
    'ai_pawn_positions': [1, 2],
    'enemy_pawn_positions': [1, 2],
})
print(ret)
