<?php

$response='[
  {
      "Name": "VrqC",
    "SerialNo": 1,
    "CodeNo": 6343,
    "Location": 1
  },
  {
      "Name": "LSKx",
    "SerialNo": 2,
    "CodeNo": 6123,
    "Location": 2
  },
  {
      "Name": "mNms",
    "SerialNo": 3,
    "CodeNo": 6809,
    "Location": 1
  },
  {
      "Name": "Pehh",
    "SerialNo": 4,
    "CodeNo": 6360,
    "Location": 2
  },
  {
      "Name": "ShUU",
    "SerialNo": 5,
    "CodeNo": 6470,
    "Location": 1
  },
  {
      "Name": "mDpk",
    "SerialNo": 6,
    "CodeNo": 6195,
    "Location": 1
  }
]';




$response = json_decode($response, true);

foreach ($response as $item){
     if($item['Location']==1){
         print "hi";
         $item['Location']="factory";
     }
}




echo json_encode($response);