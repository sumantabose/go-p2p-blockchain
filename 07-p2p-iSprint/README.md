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
4. Query on both RawMaterial and Delivery Transaction data

### Transaction Structure
#### RawMaterial Transaction
```
  {
        "SerialNo": 0,
        "ProductCode": "",
        "ProductName": "",
        "ProductBatchNo": "",
        "Quantity": 0,
        "RawMaterial": {
            "RawMaterialBatchNo": "",
            "RawMaterialsID": "",
            "RawMaterialName": "",
            "RawMaterialQuantity": 0,
            "RawMaterialMeasurementUnit": ""
        }
    }
```

#### Delivery Transaction
```
  {
        "SerialNo": 0,
        "RecGenerator": "",
        "ShipmentID": "",
        "Timestamp": "",
        "Longitude": "",
        "Latitude": "",
        "ShippedFromCompID": "",
        "ShippedToCompID": "",
        "LocationID": "",
        "DeliveryStatus": "",
        "DeliveryType": "",
        "Product": {
            "ProductCode": "",
            "ProductName": "",
            "ProductBatch": {
                "ProductBatchNo": "",
                "ProductBatchQuantity": 0
            }
        },
        "Document": {
            "DocumentURL": "",
            "DocumentType": "",
            "DocumentHash": "",
            "DocumentSign": ""
        }
    }
```

### Block Structure
```
{
    "Index": 4,
    "Timestamp": "2018-10-02 12:03:16.59392278 +0800 SGT m=+203.143341201",
    "TxnType": 0,
    "TxnPayload": "",
    "Comment": "0",
    "Proposer": "/ip4/10.27.119.37/tcp/5001/ipfs/QmTsZbgjT1NWE7kX3WH1oye8uTfJAAQ1MuWFZuPz8meapF",
    "PrevHash": "9fa31b704e70aba20c4649fb21e2472354f4bb4e8d7dec2bc1d76c3c401b5f99",
    "ThisHash": "0a8ba60f43d5699f9866e0113aa44245bcb3a821f40a0b6685f4024fca8a0895"
}
```

### ToDo
##### (Future Version)
1. If one node dies, it should be able to relay the info to the bootstrapper so the rest of the nodes in the network are updates and the connections updated accordingly.
2. Reliable and robust consensus protocol.
