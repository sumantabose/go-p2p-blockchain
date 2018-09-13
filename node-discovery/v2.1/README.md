# Node Discovery v2.1

### Features

For history of features, see [Node Discovery](https://github.com/sumantabose/go-p2p-blockchain/tree/master/node-discovery) v1.0, v1.1, v2.0.

### Updates in v2.1

The `Peer` asks for `PeerPort` to the `Bootstrapper`, which is used to create the `PeerProfile` and later used to update the `PeerGraph`.

### ToDo

* Error Handling and Exceptions
* Race condition in assigning `PeerPort` over the network
