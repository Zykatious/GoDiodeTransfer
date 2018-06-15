# GoDiodeTransfer
A Golang server and client to send files over an air-gapped network via a network diode.

File hashes are checked after the file transfer is complete to make sure the data is intact and will alert on the terminal if a file is corrupt.

Usage is as follows:


### On the receving side:

__go run Server.go__

By default the server uses port 1234 and stores transferred files to ./
This behaviour can be changed by using the argument __-p__ to change the port and __-d__ to change the receiving directory.


### On the client side:

You will need to set a static ARP route so your system knows where to send the file to, this can be done on a Linux terminal like this:
##### \# arp -s 192.168.0.2 aa:bb:cc:dd:ee:ff

Sending a file:

__go run Client.go -f *filename*__

By default the client will send to IP 127.0.0.1 and uses port 1234.
This can be changed by the argument __-l__ (lower case L) to change the receiving IP and __-p__ to change the receiving port.
