<?php

$id=$_GET["id"];

$postaction="http://localhost:8080/query/product/"."$id";

$response = file_get_contents($postaction);

echo $response;