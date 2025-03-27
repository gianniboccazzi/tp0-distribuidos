import logging
import signal
from socket import socket
from communication.protocol import BetProtocol


class ClientHandler:
    def __init__(self, client_sock: socket, lottery_completed, clients_ready, clients_total, file_lock, lottery_lock):
        self.lottery_completed = lottery_completed
        self.clients_ready = clients_ready
        self.clients_total = clients_total
        self.client_sock = client_sock
        self.client_sock.settimeout(5)
        self.file_lock = file_lock
        self.lottery_lock = lottery_lock
        self.protocol = BetProtocol()

    def run(self):
        signal.signal(signal.SIGTERM, self.__signal_handler)
        try:
            client_id = self.protocol.handle_client_connection(self.client_sock, self.lottery_completed, self.file_lock)
            if client_id:
                self.__check_lottery_status(client_id)
        except ConnectionError:
            logging.error("Client disconnected before sending message")
        self.client_sock.close()
    
    def __check_lottery_status(self, client_id):
        with self.lottery_lock:
            self.clients_ready[client_id] = True
            if not self.lottery_completed.value and len(self.clients_ready) == int(self.clients_total):
                self.lottery_completed.value = True
                logging.info("action: sorteo | result: success")

    def __signal_handler(self, signum, frame):
        self.client_sock.close()
