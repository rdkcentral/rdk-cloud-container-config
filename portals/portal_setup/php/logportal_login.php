<?php
session_name("APP_LOG_SESSION");

session_start();

if (isset($_SESSION['logged_in']) && $_SESSION['logged_in'] === true) {
    header("Location: logportal.php");
    exit;
}

$error = "";

$credFile = '##############################';

if (!is_readable($credFile)) {
    die("Credential file missing or not readable");
}

$creds = parse_ini_file($credFile, false, INI_SCANNER_RAW);
$validUser = $creds['username'] ?? '';
$token =$creds['token'] ?? '';
$url = "############################################";

$payload = [
    "secretName" => $validUser
];

$data = json_encode($payload);

$options = [
    "http" => [
        "method" => "POST",
        "header" => "Authorization: Bearer $token\r\n" .
        "Content-Type: application/json\r\n",
        "Content-Length: " . strlen($data) . "\r\n",
        "content" => $data,
        "ignore_errors" => true,
        "timeout" => 30
    ]
];

$context = stream_context_create($options);
$response = file_get_contents($url, false, $context);

$secretData = json_decode($response, true);

$validPass = $secretData['secretValue'] ?? '';

/*  BOTH checks included */
if (
    $_SERVER['REQUEST_METHOD'] === 'POST' &&
    isset($_POST['username'], $_POST['password'])
) {
    $username = trim($_POST['username']);
    $password = $_POST['password'];

    if ($username === $validUser && $password === $validPass) {
        $_SESSION['logged_in'] = true;
        $_SESSION['username'] = $username;

        header("Location: logportal.php");
        exit;
    } else {
        $error = "Invalid username or password";
    }
}
?>


<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
    <style>
        body {
            margin: 0;
            height: 100vh;
            font-family: Arial, Helvetica, sans-serif;
            background: #ffffff;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .login-box {
            width: 520px;
            background: #ffffff;
            border: 1px solid #e5e7eb;
            border-radius: 8px;
            padding: 30px 35px;
            box-shadow: 0 4px 10px rgba(0,0,0,0.06);
        }

        .login-title {
            text-align: center;
            font-size: 20px;
            font-weight: 600;
            color: #111827;
            margin-bottom: 10px;
        }

        hr {
            border: none;
            border-top: 1px solid #e5e7eb;
            margin: 15px 0 25px 0;
        }

        .form-row {
            display: grid;
            grid-template-columns: 140px 1fr;
            align-items: center;
            margin-bottom: 18px;
        }

        .form-row label {
            font-size: 14px;
            color: #374151;
        }

        .form-row input {
            padding: 10px;
            font-size: 14px;
            border: 1px solid #d1d5db;
            border-radius: 6px;
            width: 100%;
        }

        .form-row input:focus {
            outline: none;
            border-color: #2563eb;
        }

        .actions {
            text-align: right;
            margin-top: 10px;
        }

        .actions button {
            padding: 10px 28px;
            font-size: 14px;
            background: #2563eb;
            color: #ffffff;
            border: none;
            border-radius: 6px;
            cursor: pointer;
        }

        .actions button:hover {
            background: #1d4ed8;
        }

        .error {
            background: #fee2e2;
            color: #b91c1c;
            padding: 10px;
            border-radius: 6px;
            margin-bottom: 20px;
            text-align: center;
            font-size: 14px;
        }
    </style>
</head>
<body>

<div class="login-box">
    <div class="login-title">Enter credentials</div>
    <hr>

    <?php if ($error): ?>
        <div class="error"><?php echo $error; ?></div>
    <?php endif; ?>

    <form method="post">
        <div class="form-row">
            <label for="username">Username:</label>
            <input type="text" id="username" name="username" required>
        </div>

        <div class="form-row">
            <label for="password">Password:</label>
            <input type="password" id="password" name="password" required>
        </div>

        <div class="actions">
            <button type="submit">Login</button>
        </div>
    </form>
</div>

</body>
</html>

