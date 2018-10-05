<?php
include 'common.php';

$Name = $_REQUEST["name"];
$CodeNo = $_REQUEST["code"];

if($username=="factory" || $username=="admin"){
    $postaction="http://localhost:8080/post/"."$CodeNo"."/"."$Name";
    $response = file_get_contents($postaction);
    $msg='You have sent the record successfully!';
    echo json_encode(['status' => 'success', 'message' => $msg]);
}else{
    $msg="Access deny";
    echo json_encode(['status' => 'fail', 'message' => $msg]);
}




