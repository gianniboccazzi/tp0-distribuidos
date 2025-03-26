import signal
import socket
import logging

from communication.protocol import BetProtocol
from multiprocessing import Process, Manager, Queue



class Server:
    def __init__(self, port, listen_backlog, clients_total):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.active = True
        self.lottery_completed = False
        self.protocol = BetProtocol()
        self.clients_total = clients_total
        self.manager = Manager()
        self.clients_ready = set()
        self.queue = Queue()
        self.processes = []

    def run(self):
        signal.signal(signal.SIGTERM, self.__signal_handler)
        while self.active:
            client_sock = self.__accept_new_connection()
            if client_sock:
                process = Process(target=self.protocol.handle_client_connection, args=(client_sock, self.lottery_completed))
                process.start()
                client_sock.close()
                self.check_lottery_status(client_ready)

    def check_lottery_status(self, client_ready):
        if client_ready:
            self.clients_ready.add(client_ready)
        if not self.lottery_completed and len(self.clients_ready) == int(self.clients_total):
            self.lottery_completed = True
            logging.info("action: sorteo | result: success")



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

    

    
