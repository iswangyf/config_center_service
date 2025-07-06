CREATE DATABASE IF NOT EXISTS configdb DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;

USE configdb;

DROP TABLE IF EXISTS modules;
DROP TABLE IF EXISTS module_groups;

CREATE TABLE module_groups (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE modules (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    group_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(64) NOT NULL,
    content TEXT,
    valid_from DATETIME,
    valid_to DATETIME,
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES module_groups(id)
);

INSERT INTO module_groups (name, description) VALUES ('homepage', 'Homepage config group');

INSERT INTO modules (group_id, name, content, valid_from, valid_to, enabled)
VALUES
(1, 'moduleA', '{"key": "value"}', NOW() - INTERVAL 1 HOUR, NOW() + INTERVAL 1 DAY, TRUE),
(1, 'moduleB', '{"key": "test"}', NOW() - INTERVAL 1 DAY, NOW() + INTERVAL 1 DAY, FALSE);