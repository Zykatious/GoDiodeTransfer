# GoDiodeTransfer
A Golang server and client to send files over an air-gapped network via a network diode.

Usage is as follows:


On the receving side:

__go run Server.go__

By default the server uses port 1234 and stores transferred files to ./
This behaviour can be changed by using the argument __-p__ to change the port and __-d__ to change the receiving directory.


On the client side:

__go run Client.go -f *filename*__

By default the client will send to IP 127.0.0.1 and uses port 1234.
This can be changed by the argument __-l__ (lower case L) to change the receiving IP and __-p__ to change the receiving port.
