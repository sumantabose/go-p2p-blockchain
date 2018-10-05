<?php
include 'common.php';

$party=$_GET["field"];

if($username!=$party && $username != 'admin') {
    $msg = 'Not authorized';
    echo json_encode(['status' => 'fail', 'message' => $msg]);
}else{

    if ($party == "factory") {
        $field = 1;
    } elseif ($party == "shipment1") {
        $field = 2;
    } elseif ($party == "warehouse") {
        $field = 3;
    } elseif ($party == "shipment2") {
        $field = 4;
    } elseif ($party == "merchant") {
        $field = 5;
    } elseif ($party == "customer") {
        $field = 6;
    }
    $postaction = "http://localhost:8080/info/" . "$field";

    $response = file_get_contents($postaction);
    $response = json_decode($response, TRUE);

    if (empty($response)) {
        $response = [];
    }

    foreach ($response as $k => $v) {

        if ($response[$k]['Location'] == 1) {
            $response[$k]['Location'] = 'factory';
        } elseif ($response[$k]['Location'] == 2) {
            $response[$k]['Location'] = 'shipment 1';
        } elseif ($response[$k]['Location'] == 3) {
            $response[$k]['Location'] = 'warehouse';
        } elseif ($response[$k]['Location'] == 4) {
            $response[$k]['Location'] = 'shipment 2';
        } elseif ($response[$k]['Location'] == 5) {
            $response[$k]['Location'] = 'merchant';
        } elseif ($response[$k]['Location'] == 6) {
            $response[$k]['Location'] = 'customer';
        }
    }


    echo json_encode($response);
}