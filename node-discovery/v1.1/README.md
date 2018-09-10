# Node Discovery v1.1

### Features

* `Bootstrapper` running on `listenPort` has two REST API end points:
	* `/query-p2p-list` : Returns the ListOfPeers.
	* `/join-p2p-net` : Enrolls a peer to the network by adding the peer address to the ListOfPeers.
* `Peer` queries the ListOfPeers and (optionally may join anyone from the list) and enrolls itself in the network. It uses the two endpoints as above.

### Updates

The `Peer` send a `struct` to the `Bootstrapper` instead of a single element. This will be used to transfer `PeerProfile` and `PeerGraph` in future versions.