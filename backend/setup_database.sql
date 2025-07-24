-- Create database
CREATE DATABASE IF NOT EXISTS meesho_clone;
USE meesho_clone;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    INDEX idx_user_id (user_id),
    INDEX idx_phone_number (phone_number),
    INDEX idx_deleted_at (deleted_at)
);

-- Insert sample data (optional)
INSERT IGNORE INTO users (user_id, phone_number, name) VALUES 
('user_sample123', '9876543210', 'Sample User'),
('user_test456', '8765432109', 'Test User');

-- Show tables and data
SHOW TABLES;
SELECT * FROM users; 