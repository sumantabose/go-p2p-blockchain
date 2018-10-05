<?php
include 'common.php';

$Name = $_REQUEST["name"];
$CodeNo = $_REQUEST["code"];



if($username=="factory"){
    $data = array(
        'CodeNo'=>intval($CodeNo),
        'Name'=>$Name
    );

    $query = json_encode($data);


    $options['http'] = array(
        'timeout'=>60,
        'method' => 'POST',
        'header' => 'Content-type:application/json',
        'content' => $query
    );

    $url = "http://localhost:8080/post";
    $context = stream_context_create($options);
    $result = file_get_contents($url, false, $context);

    $msg='You have sent the record successfully!';
    echo json_encode(['status' => 'success', 'message' => $msg]);
}else{
    $msg="Access deny";
    echo json_encode(['status' => 'fail', 'message' => $msg]);
}




