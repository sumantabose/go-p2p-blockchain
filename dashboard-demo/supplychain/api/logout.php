<?php



setcookie('id', 0, 0, '/', $_SERVER['HTTP_HOST']);
setcookie('au', 0, 0, '/', $_SERVER['HTTP_HOST']);


$url = 'index.html';


echo json_encode(['status' => 'success', 'message' => $url]);


