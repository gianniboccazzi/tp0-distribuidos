import signal
import socket
import logging

from communication.protocol import BetProtocol



class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.active = True
        self.protocol = BetProtocol()

    def run(self):
        signal.signal(signal.SIGTERM, self.__signal_handler)
        while self.active:
            client_sock = self.__accept_new_connection()
            if client_sock:
                self.protocol.handle_client_connection(client_sock)
                client_sock.close()

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
        self._server_socket.close()

    

    
