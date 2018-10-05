<?php




$id = $_COOKIE["id"];
$auth_key = $_COOKIE["au"];

$dbh = new PDO("mysql:dbname=supplychain;host=127.1", "supplyuser", "chainuser");

$stmt1 = $dbh->prepare("SELECT username FROM session WHERE id=? AND auth_key=?");
$stmt1->execute([$id, $auth_key]);
$username = $stmt1->fetch();

if (empty($username)) {
    header('http://Localhost/login/login.html');
    exit('ACCESS DENIED, login required');
} else {
    $username = $username['username'];
}
