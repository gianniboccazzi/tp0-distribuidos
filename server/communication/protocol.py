
import logging
import signal
import socket
from common.utils import Bet, has_won, load_bets, store_bets


ERROR_RES = "ERR"
ACK_RES = "ACK"
NONE_RES = "NONE"


def parse_batch(message: str, client_id) -> list[Bet]:
    parts = message.split("||")
    bets = []
    for part in parts:
        bet_parts = part.split("|")
        if len(bet_parts) != 5:
            raise ValueError("Invalid message format")
        bets.append(Bet(client_id, bet_parts[0], bet_parts[1], bet_parts[2], bet_parts[3], bet_parts[4]))
    return bets

class BetProtocol:
    def receive_batches(self, client_sock: socket.socket, client_id, file_lock):
        eof = False
        bets = []
        while not eof:
            try:
                buffer = b""
                buffer = self.receive_until_delimiter(client_sock, buffer)
                message_length_bytes, remaining_data = buffer.split(b"|", 1)
                message_length = int(message_length_bytes.decode())
                remaining_data += self.__recv_all(client_sock, message_length - len(remaining_data))
                decoded_message = remaining_data.decode()
                if decoded_message == "EOF":
                    eof = True
                    break
                bets = parse_batch(decoded_message, client_id)
                with file_lock:
                    store_bets(bets)
                self.__send_response(client_sock, ACK_RES)
                logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")
            except ValueError as e:
                logging.error(f"action: apuesta_recibida | result: fail | error: {e}")
                self.__send_response(client_sock, ERROR_RES)
                return
            except (socket.timeout, ConnectionError, OSError) as e:
                logging.error(f"action: apuesta_recibida | result: fail | cantidad: {len(bets)} | error: {e}")
                return
        return client_id

    def receive_until_delimiter(self, client_sock, buffer):
        while b"|" not in buffer:
            chunk = client_sock.recv(2)
            if not chunk:
                raise ConnectionError("Client disconnected before sending message")
            buffer += chunk
        return buffer
        

    def handle_client_connection(self, client_sock: socket.socket, lottery_ready, file_lock):
        client_sock.settimeout(5)
        buffer = b""
        buffer = self.receive_until_delimiter(client_sock, buffer)
        message_length_bytes, remaining_data = buffer.split(b"|", 1)
        message_length = int(message_length_bytes.decode())
        remaining_data += self.__recv_all(client_sock, message_length - len(remaining_data))
        payload = remaining_data.decode()
        client_id, action = payload.split("|")
        if action == "BETS":
            return self.receive_batches(client_sock, client_id, file_lock)
        if action == "WINNERS":
            return self.send_winners(client_sock,client_id, lottery_ready, file_lock)
        return 

    
    def send_winners(self, client_sock: socket.socket, client_id, lottery_ready, file_lock):
        if not lottery_ready.value:
            self.__send_response(client_sock, ERROR_RES)
            return
        winners = []
        with file_lock:
            bets = load_bets()
        for bet in bets:
            if has_won(bet) and bet.agency == int(client_id):
                winners.append(bet.document)
        if len(winners) == 0:
            payload = NONE_RES
        else:
            payload = "|".join(winners)
        self.__send_response(client_sock,payload)
        return
        
    def __send_response(self, client_sock: socket.socket, response: str):
        response_length = len(response)
        response_bytes = f"{response_length}|{response}".encode()
        bytesSent = 0
        while bytesSent < len(response_bytes):
            bytesSent += client_sock.send(response_bytes[bytesSent:])
        return
    
    
    def __recv_all(self, sock, length):
        data = b""
        while len(data) < length:
            chunk = sock.recv(length - len(data))
            if not chunk:
                raise ConnectionError("Client disconnected while receiving message")
            data += chunk
        return data

    