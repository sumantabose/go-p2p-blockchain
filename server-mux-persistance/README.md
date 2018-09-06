# Peer-2-Peer Blockchain Network, MUX Server and Persistance

### Getting started

Here we have improved the `server-mux` to added persistence. Previously, for simplicity, we shut down all the nodes if you kill one of them. Here we leave all other nodes running even if one closes.


For the 1st node, `go run defs.go mux.go p2p.go blockchain.go main.go -l 6000 -secio` where you can replace 6000 by your prefered port number. Follow the instructions in the terminal to connect subsequent nodes.

TODO: If one node dies, it should be able to relay the info in its Peer Store to its connected nodes so they can connect to each other.