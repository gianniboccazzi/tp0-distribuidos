import signal
import socket
import logging

from communication.protocol import BetProtocol
from multiprocessing import Process, Manager, Queue

from common.client_handler import ClientHandler



class Server:
    def __init__(self, port, listen_backlog, clients_total):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.active = True
        self.clients_total = clients_total
        self.manager = Manager()
        self.lottery_completed = self.manager.Value('b', False)
        self.clients_ready = self.manager.dict()
        self.file_lock = self.manager.Lock()
        self.lottery_lock = self.manager.Lock()
        self.processes = []

    def run(self):
        signal.signal(signal.SIGTERM, self.__signal_handler)
        while self.active:
            client_sock = self.__accept_new_connection()
            if client_sock:
                client_handler = ClientHandler(client_sock, self.lottery_completed, self.clients_ready,
                                                self.clients_total, self.file_lock, self.lottery_lock)
                process = Process(target=client_handler.run)
                process.start()
                self.processes.append(process)
                client_sock.close()
        self.__shutdown()

    def __shutdown(self):
        """
        Shutdown server

        Function that closes the server socket and sets the active flag to False
        """
        self.active = False
        self._server_socket.close()
        for process in self.processes:
            if process.is_alive():
                process.terminate()
        for process in self.processes:
            process.join(timeout=1)
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
        self.__shutdown()

    