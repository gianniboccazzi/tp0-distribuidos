import signal
import socket
import logging

from .utils import store_bets
from communication.protocol import parse_message


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.active = True

    def run(self):
        signal.signal(signal.SIGTERM, self.__signal_handler)
        while self.active:
            client_sock = self.__accept_new_connection()
            if client_sock:
                self.__handle_client_connection(client_sock)

    def __handle_client_connection(self, client_sock: socket.socket):
        """
        Read the message from the client and send an ack if the message was received correctly
        Then, store the bet in the storage file

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            client_sock.settimeout(5)
            buffer = b""
            while b"|" not in buffer:
                chunk = client_sock.recv(3)
                if not chunk:
                    raise ConnectionError("Client disconnected before sending message")
                buffer += chunk
            message_length_str, remaining_data = buffer.split(b"|", 1)
            message_length = int(message_length_str.decode())
            remaining_data += self.recv_all(client_sock, message_length - len(remaining_data))
            bet = parse_message(remaining_data)
            self.__send_flag(client_sock, "ACK\n")
            logging.info(f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}")
            store_bets([bet])  
        except ValueError as e:
            logging.error(f"action: apuesta_almacenada | result: fail | error: {e}")
            self.__send_flag(client_sock, "ERR\n") 
        except (socket.timeout, ConnectionError, OSError) as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")


    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        try:
            c, addr = self._server_socket.accept()            
        except OSError as e:
            logging.error(f'action: accept_connections | result: fail | error: {e}')
            return
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

    def __signal_handler(self, signum, frame):
        """
        Signal handler for SIGTERM signal

        Function that is called when a SIGTERM signal is received
        """
        self.active = False
        self._server_socket.close()

    def recv_all(self, sock, length):
        data = b""
        while len(data) < length:
            chunk = sock.recv(length - len(data))
            if not chunk:
                raise ConnectionError("Client disconnected while receiving message")
            data += chunk
        return data

    def __send_flag(self, client_sock: socket.socket, flag: str):
        """
        Send ack to client

        Function that sends an ack message to the client
        """
        res = client_sock.sendall(flag.encode())
        if res == 0:
            raise ConnectionError("Client disconnected before sending ack message")
