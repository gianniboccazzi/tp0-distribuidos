
import logging
import signal
import socket
from common.utils import Bet, store_bets


ERROR_RES = "ERR"
ACK_RES = "ACK"


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

class BetProtocol:
    def handle_client_connection_batches(self, client_sock: socket.socket):
        eof = False
        client_sock.settimeout(5)
        bets = []
        while not eof:
            try:
                buffer = b""
                while b"|" not in buffer:
                    chunk = client_sock.recv(2)
                    if not chunk:
                        raise ConnectionError("Client disconnected before sending message")
                    buffer += chunk

                message_length_bytes, remaining_data = buffer.split(b"|", 1)
                message_length = int(message_length_bytes.decode())
                remaining_data += self.__recv_all(client_sock, message_length - len(remaining_data))
                bets, eof = parse_batch(remaining_data)
                self.__send_response(client_sock, ACK_RES)
                logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")
                store_bets(bets)
            except ValueError as e:
                logging.error(f"action: apuesta_recibida | result: fail | error: {e}")
                self.__send_response(client_sock, ERROR_RES)
                break
            except (socket.timeout, ConnectionError, OSError) as e:
                logging.error(f"action: apuesta_recibida | result: fail | cantidad: {len(bets)} | error: {e}")
                break
    
    def __send_response(self, client_sock: socket.socket, response: str):
        """
        Send ack to client

        Function that sends an ack message to the client
        """
        res = client_sock.sendall(response.encode())
        if res == 0:
            raise ConnectionError("Client disconnected before sending ack message")
    
    def __recv_all(self, sock, length):
        data = b""
        while len(data) < length:
            chunk = sock.recv(length - len(data))
            if not chunk:
                raise ConnectionError("Client disconnected while receiving message")
            data += chunk
        return data

    