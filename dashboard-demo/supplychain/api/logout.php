<?php



setcookie('id', 0, 0, '/', $_SERVER['HTTP_HOST']);
setcookie('au', 0, 0, '/', $_SERVER['HTTP_HOST']);


$url = 'http://localhost/supplychain/index.html';


echo json_encode(['status' => 'success', 'message' => $url]);


