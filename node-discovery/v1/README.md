# Node Discovery v1

### Features

* `Bootstrapper` running on `listenPort` has two REST API end points:
	* `/query-p2p-list` : Returns the ListOfPeers.
	* `/join-p2p-net` : Enrolls a peer to the network by adding the peer address to the ListOfPeers.
* `Peer` queries the ListOfPeers and (optionally may join anyone from the list) and enrolls itself in the network. It uses the two endpoints as above.

### Versions 

For a complete list of versions, check [here](https://github.com/sumantabose/go-p2p-blockchain/tree/master/node-discovery).