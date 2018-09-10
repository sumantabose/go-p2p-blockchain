# Node Discovery v1.1

### Features

For basic features, see [Node Discovery v1](https://github.com/sumantabose/go-p2p-blockchain/tree/master/node-discovery/v1).

### Updates

The `Peer` send a `struct` to the `Bootstrapper` instead of a single element. This will be used to transfer `PeerProfile` and `PeerGraph` in future versions.