#!/usr/bin/env python

# code:
# set up the webhook and wait for a connection
import socket

webhook_host='localhost'
webhook_port=3333

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
print "binding on host", webhook_host, "port", webhook_port
s.bind((webhook_host, webhook_port))
print "listening"
s.listen(1)
print "accepting connections"

while True:
    conn, (addr, port) = s.accept()
    print "connection ", conn, " from ", addr, " port ", port

# output like:
# binding on host level02-3.stripe-ctf.com port 9894
# listening
# accepting connections