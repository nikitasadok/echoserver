# echoserver
Simple echo server rfc 862 (it doesn't adhere to standard, because it doesn't listen on port 7, but the port and host are parametrized)

It is an echo server which, however, handles some bottlenecks:
1) If connection is idle for more than some time (3 minutes in this implementation), the connection is dropped
2) If there is a number connections equal to maximum number, we find the connection, which was updated last, close it, and open the requested one.

I have also included very simple script in Python directory which does some load testing of the server.
