# Peer-2-Peer Blockchain Network, MUX Server and Persistance

### Getting started

Here we have improved the `server-mux` to added persistence. Previously, for simplicity, we shut down all the nodes if you kill one of them. Here we leave all other nodes running even if one closes.


TODO: If one node dies, it should be able to relay the info in its Peer Store to its connected nodes so they can connect to each other.