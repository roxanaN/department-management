CREATE DATABASE IF NOT EXISTS department_management;

USE department_management;

CREATE TABLE roles (
    id CHAR(36) PRIMARY KEY,
    name ENUM('admin', 'manager', 'user') UNIQUE NOT NULL
);

CREATE TABLE departments (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    flags INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE departments_hierarchy (
    parent_id CHAR(36) NOT NULL,
    child_id CHAR(36) NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES departments(id) ON DELETE CASCADE,
    FOREIGN KEY (child_id) REFERENCES departments(id) ON DELETE CASCADE,
    PRIMARY KEY (parent_id, child_id)
);

CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) DEFAULT '',
    is_password_set BOOLEAN DEFAULT FALSE,
    activated BOOLEAN NOT NULL DEFAULT TRUE,
    role_id CHAR(36),
    department_id CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE SET NULL,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL
);

INSERT INTO roles (id, name)
VALUES (UUID(), 'admin'), (UUID(), 'manager'), (UUID(), 'user');

INSERT INTO departments (id, name, flags)
VALUES (UUID(), 'HR', 1), (UUID(), 'IT', 1), (UUID(), 'Marketing', 1), (UUID(), 'BE', 1), (UUID(), 'FE', 1), (UUID(), 'QA', 1);

INSERT INTO departments_hierarchy (parent_id, child_id)
VALUES ((SELECT id FROM departments WHERE name = 'HR'), (SELECT id FROM departments WHERE name = 'IT')),
       ((SELECT id FROM departments WHERE name = 'IT'), (SELECT id FROM departments WHERE name = 'BE')),
       ((SELECT id FROM departments WHERE name = 'IT'), (SELECT id FROM departments WHERE name = 'FE')),
       ((SELECT id FROM departments WHERE name = 'IT'), (SELECT id FROM departments WHERE name = 'QA')),
       ((SELECT id FROM departments WHERE name = 'HR'), (SELECT id FROM departments WHERE name = 'Marketing')),
       ((SELECT id FROM departments WHERE name = 'IT'), (SELECT id FROM departments WHERE name = 'Marketing'));

-- Add admin user to database
INSERT INTO users (id, email, password, is_password_set, activated, role_id, department_id)
VALUES (
    UUID(), 
    'roxana.nemulescu@gmail.com', 
    '$2a$10$4tFYxSGwZmyQGwKR3ugVk.fKO/YtWz1u9xo8vSDL4edZvtqVyPnyG', 
    TRUE, 
    TRUE, 
    (SELECT id FROM roles WHERE name = 'admin'), 
    (SELECT id FROM departments WHERE name = 'HR')
);

GRANT ALL PRIVILEGES ON * . * TO 'root'@'localhost';