
import logging
from common.utils import Bet


def parse_message(message: bytes) -> Bet:
    decoded_message = message.decode()
    parts = decoded_message.split("|")
    if len(parts) != 6:
        raise ValueError("Invalid message format")
    return Bet(parts[0], parts[1], parts[2], parts[3], parts[4], parts[5])
    

