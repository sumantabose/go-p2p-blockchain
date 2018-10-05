<?php



$postaction="http://localhost:8080/next";
$response = file_get_contents($postaction);
$msg = "It succeed!";

echo json_encode(['status' => 'success', 'message' => $msg]);