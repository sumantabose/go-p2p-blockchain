<?php
//phpinfo();die;
$username = $_REQUEST["username"];
$password = password_hash($_REQUEST["password"], PASSWORD_DEFAULT);

$dbh = new PDO("mysql:dbname=supplychain;host=127.1", "supplyuser", "chainuser");

try {
    $stmt1 = $dbh->prepare("INSERT INTO `login` (`username`, `password`) VALUES(?,?)");
    $stmt1->execute([$username, $password]);

    echo json_encode(['status' => 'success', 'message' => 'User registration successful!']);
} catch (PDOException $e) {
    echo json_encode(['status' => 'fail', 'message' => $e->getMessage()]);

}
