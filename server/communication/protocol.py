
import logging
from common.utils import Bet


def parse_message(message: bytes) -> Bet:
    decoded_message = message.decode()
    parts = decoded_message.split("|")
    if len(parts) != 6:
        raise ValueError("Invalid message format")
    return Bet(parts[0], parts[1], parts[2], parts[3], parts[4], parts[5])
    

def parse_batch(message: bytes) -> tuple[list[Bet], bool]:
    decoded_message = message.decode()
    eof = False
    if decoded_message.endswith("|||"):
        eof = True
        decoded_message = decoded_message[:-3]
    parts = decoded_message.split("||")
    bets = []
    for part in parts:
        bet_parts = part.split("|")
        if len(bet_parts) != 6:
            raise ValueError("Invalid message format")
        bets.append(Bet(bet_parts[0], bet_parts[1], bet_parts[2], bet_parts[3], bet_parts[4], bet_parts[5]))
    return bets, eof