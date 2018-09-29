# Peer-2-Peer Blockchain Network, MUX Server, Persistance, Node Discovery, FileSave and BlockStruct

### Getting started

* Step 1: `bash bootstrapper.sh` (This starts the local bootstrapper server to enable node-discovery). Alternately start cloud bootstrapper in Heroku.
* Step 2: In a new terminal, `bash peer.sh <option>` (This starts a peer connecting to the existing network, if any)
  `<option>` available are:
  * `heroku`: To connect to cloud bootstrapper in Heroku at [http://blockchain-bootstrapper.herokuapp.com](http://blockchain-bootstrapper.herokuapp.com).
  * `local`: To connect to local bootstrapper.
  * `A.B.C.D`: To connect to a bootstrapper in the local network specified using IP.

### Update
1. Block structure updated with new transaction fields (interfact).
2. Correct file saving GOB encoding ensured.
3. Multiple bootstrapper options: Heroku, Local, IP.
4. Query on Blockchain/Transaction data (ongoing)

### ToDo
##### (Future Version)
If one node dies, it should be able to relay the info to the bootstrapper so the rest of the nodes in the network are updates and the connections updated accordingly.
