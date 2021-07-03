import socket
import unittest

class TestEchoServer(unittest.TestCase):
    def test_load(self):
        TCP_IP = '127.0.0.1'
        TCP_PORT = 3333
        BUFFER_SIZE = 8192
        REQUESTS = 500

        for i in range(REQUESTS):
            want = bytes('Hello from ' + str(i), encoding='utf-8')
            s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            s.connect((TCP_IP, TCP_PORT))
            s.send(want)
            data = s.recv(BUFFER_SIZE)
            self.assertEqual(data, want)
            s.close()