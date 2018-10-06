<?php

$username = $_REQUEST["username"];
$password = $_REQUEST["password"];

$dbh = new PDO("mysql:dbname=supplychain;host=127.1", "supplyuser", "chainuser");

try {
    $stmt1 = $dbh->prepare("SELECT password FROM login WHERE username=?");
    $stmt1->execute([$username]);
    $password2 = $stmt1->fetch();

    if (empty($password2) || !password_verify($password, $password2['password'])) {
        $msg = 'Wrong username or password!';
        echo json_encode(['status' => 'fail', 'message' => $msg]);
    } else {
        $record = $dbh->prepare("INSERT INTO `session` (`username`, `auth_key`) VALUES(?, ?)");
        $auth_key = md5($username . mt_rand(1, 100000) . time());
        $record->execute([$username, $auth_key]);
        $session_id = $dbh->lastInsertId();

        setcookie('id', $session_id, 0, '/', $_SERVER['HTTP_HOST']);
        setcookie('au', $auth_key, 0, '/', $_SERVER['HTTP_HOST']);


        $url = 'supplychain.html';


        echo json_encode(['status' => 'success', 'message' => $url]);
    }
} catch (PDOException $e) {
    echo json_encode(['status' => 'fail', 'message' => $e->getMessage()]);

}

