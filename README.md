# GoDiodeTransfer
A Golang server and client to send data over an air-gapped network via a network diode.

Usage is as follows:


On the receving side:

__go run Server.go__

By default the server uses port 1234 and stores transferred files to ./
This behaviour can be changed by using the argument -p to change the port and -d to change the receiving directory.


On the client side:

__go run Client.go -f *filename*__

By default the client will send to IP 127.0.0.1 and uses port 1234.
This can be changed by the argument -l to change the receiving IP and -p to change the receiving port.
