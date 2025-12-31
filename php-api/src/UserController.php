<?php

namespace App;

use PDO;
use PDOException;

class UserController
{
    private PDO $db;

    public function __construct()
    {
        $this->db = Database::getConnection();
    }

    public function createUser(): void
    {
        // Get JSON input
        $input = file_get_contents('php://input');
        $data = json_decode($input, true);

        // Validate input
        if (json_last_error() !== JSON_ERROR_NONE) {
            $this->jsonResponse(['error' => 'Invalid JSON'], 400);
            return;
        }

        if (empty($data['name']) || empty($data['email'])) {
            $this->jsonResponse(['error' => 'Name and email are required'], 400);
            return;
        }

        // Validate email format
        if (!filter_var($data['email'], FILTER_VALIDATE_EMAIL)) {
            $this->jsonResponse(['error' => 'Invalid email format'], 400);
            return;
        }

        // Insert user
        try {
            $stmt = $this->db->prepare(
                'INSERT INTO users (name, email) VALUES (:name, :email)'
            );
            
            $stmt->execute([
                'name' => trim($data['name']),
                'email' => trim($data['email'])
            ]);

            $userId = $this->db->lastInsertId();

            // Fetch the created user
            $stmt = $this->db->prepare('SELECT id, name, email FROM users WHERE id = :id');
            $stmt->execute(['id' => $userId]);
            $user = $stmt->fetch();

            $this->jsonResponse($user, 201);
        } catch (PDOException $e) {
            // Check for duplicate email
            if ($e->getCode() === '23000') {
                $this->jsonResponse(['error' => 'Email already exists'], 409);
            } else {
                error_log("Database error: " . $e->getMessage());
                $this->jsonResponse(['error' => 'Failed to create user'], 500);
            }
        }
    }

    public function getUsers(): void
    {
        try {
            $stmt = $this->db->query('SELECT id, name, email FROM users ORDER BY id DESC');
            $users = $stmt->fetchAll();
            $this->jsonResponse($users);
        } catch (PDOException $e) {
            error_log("Database error: " . $e->getMessage());
            $this->jsonResponse(['error' => 'Failed to fetch users'], 500);
        }
    }

    private function jsonResponse(mixed $data, int $statusCode = 200): void
    {
        http_response_code($statusCode);
        header('Content-Type: application/json');
        echo json_encode($data, JSON_PRETTY_PRINT);
    }
}