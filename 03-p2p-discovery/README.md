# Peer-2-Peer Blockchain Network, MUX Server, Persistance and Node Discovery

### Getting started

* Step 1: `go run bootstrapper/bootstrapper.go` (This starts the bootstrapper server to enable node-discovery)
* Step 2: In a new terminal, `go run *.go` or `bash run-node.sh` (This starts a node connecting to the existing network, if any)

### Update

Node discovery feature added (facilitated by Bootstrapper).

### ToDo
##### (Future Version)
If one node dies, it should be able to relay the info to the bootstrapper so the rest of the nodes in the network are updates and the connections updated accordingly.
