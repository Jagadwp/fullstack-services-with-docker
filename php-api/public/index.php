<?php

require_once __DIR__ . '/../vendor/autoload.php';

use App\Database;
use App\UserController;

// Enable error reporting for debugging
error_reporting(E_ALL);
ini_set('display_errors', '1');

// Set headers
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, POST, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type');

// Handle preflight OPTIONS request
if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(200);
    exit;
}

// Initialize database schema
Database::initializeSchema();

// Router
$requestMethod = $_SERVER['REQUEST_METHOD'];
$requestUri = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);

$controller = new UserController();

// Simple routing
if ($requestUri === '/users' || $requestUri === '/') {
    switch ($requestMethod) {
        case 'POST':
            $controller->createUser();
            break;
        case 'GET':
            $controller->getUsers();
            break;
        default:
            http_response_code(405);
            header('Content-Type: application/json');
            echo json_encode(['error' => 'Method not allowed']);
    }
} elseif ($requestUri === '/health') {
    http_response_code(200);
    header('Content-Type: application/json');
    echo json_encode(['status' => 'healthy', 'service' => 'php-api']);
} else {
    http_response_code(404);
    header('Content-Type: application/json');
    echo json_encode(['error' => 'Not found']);
}