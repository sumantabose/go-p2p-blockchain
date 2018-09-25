# Peer-2-Peer Blockchain Network, MUX Server, Persistance and Node Discovery

### Getting started

* Step 1: `bash bootstrapper.sh` (This starts the bootstrapper server to enable node-discovery)
* Step 2: In a new terminal, `bash peer.sh` (This starts a peer connecting to the existing network, if any)

### Update
Block structure updated with new transaction fields (interfact).

### ToDo
##### (Future Version)
If one node dies, it should be able to relay the info to the bootstrapper so the rest of the nodes in the network are updates and the connections updated accordingly.
