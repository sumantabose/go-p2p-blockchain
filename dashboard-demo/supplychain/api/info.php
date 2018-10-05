<?php



$postaction="http://localhost:8080/info";

$response = file_get_contents($postaction);

$response = json_decode($response, TRUE);

foreach ($response as &$item){
    if($item['Location']==1){
        $item['Location']="factory";
    }elseif($item['Location']==2){
        $item['Location']="shipment1";
    }elseif($item['Location']==3){
        $item['Location']="warehouse";
    }elseif($item['Location']==4){
        $item['Location']="shipment2";
    }elseif($item['Location']==5){
        $item['Location']="shop";
    }elseif($item['Location']==6){
        $item['Location']="customer";
    }
}
unset($item);




$response = json_encode($response);

echo $response;