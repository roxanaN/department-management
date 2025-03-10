DELIMITER //

CREATE PROCEDURE CreateDepartment(IN dept_id CHAR(36), IN dept_name VARCHAR(255))
BEGIN
    INSERT INTO departments (id, name, flags) VALUES (dept_id, dept_name, 1); -- Set bit 1 (Active)
END //

CREATE PROCEDURE UpdateDepartment(IN dept_id CHAR(36), IN dept_name VARCHAR(255), IN dept_flags INT)
BEGIN
    UPDATE departments SET name = dept_name, flags = dept_flags WHERE id = dept_id;
END //

CREATE PROCEDURE DeleteDepartment(
    IN dept_id CHAR(36)
)
BEGIN
    UPDATE departments SET flags = (flags & ~1) | 2 WHERE id = dept_id;
END //

CREATE PROCEDURE GetAllDepartments()
BEGIN
    SELECT id, name, flags FROM departments;
END //

CREATE PROCEDURE GetHierarchy(IN dept_name VARCHAR(255))
BEGIN
    SELECT d.id, d.name, d.flags
    FROM departments d
    JOIN departments_hierarchy dh ON d.id = dh.child_id
    WHERE dh.parent_id = (SELECT id FROM departments WHERE name = dept_name);
END //

CREATE PROCEDURE GetRoleIDByName(IN role_name VARCHAR(255))
BEGIN
    SELECT id FROM roles WHERE name = role_name;
END //

CREATE PROCEDURE GetDepartmentIDByName(IN department_name VARCHAR(255))
BEGIN
    SELECT id FROM departments WHERE name = department_name;
END //

CREATE PROCEDURE AddDepartmentToHierarchy(IN parent_id CHAR(36), IN child_id CHAR(36))
BEGIN
    INSERT INTO departments_hierarchy (parent_id, child_id) VALUES (parent_id, child_id);
END //

CREATE PROCEDURE AddUser(
    IN user_id CHAR(36),
    IN user_email VARCHAR(255),
    IN role_id CHAR(36),
    IN department_id CHAR(36)
)
BEGIN
    INSERT INTO users (id, email, is_password_set, activated, role_id, department_id)
    VALUES (user_id, user_email, FALSE, FALSE, role_id, department_id);
END //

CREATE PROCEDURE ChangeUserDepartment(
    IN user_id CHAR(36),
    IN new_department_id CHAR(36)
)
BEGIN
    UPDATE users
    SET department_id = new_department_id
    WHERE id = user_id;
END //

CREATE PROCEDURE ChangeUserRole(
    IN user_id CHAR(36),
    IN new_role_id CHAR(36)
)
BEGIN
    UPDATE users
    SET role_id = new_role_id
    WHERE id = user_id;
END //

CREATE PROCEDURE GetUserByEmail(IN user_email VARCHAR(255))
BEGIN
    SELECT 
        users.id,
        users.email,
        users.password,
        users.is_password_set,
        users.activated,
        roles.name AS role,
        departments.name AS department
    FROM 
        users
    LEFT JOIN 
        roles ON users.role_id = roles.id
    LEFT JOIN 
        departments ON users.department_id = departments.id
    WHERE 
        users.email = user_email;
END //

CREATE PROCEDURE SetUserPassword(IN user_email VARCHAR(255), IN hashed_password VARCHAR(255))
BEGIN
    UPDATE users
    SET password = hashed_password, is_password_set = TRUE
    WHERE email = user_email;
END //

CREATE PROCEDURE UpdateUserPasswordStatus(IN p_email VARCHAR(255), IN p_status BOOLEAN)
BEGIN
    UPDATE users SET is_password_set = p_status WHERE email = p_email;
END //

CREATE PROCEDURE GetAllUsers()
BEGIN
    SELECT u.id, u.email, r.name AS role, d.name AS department
    FROM users u
    LEFT JOIN roles r ON u.role_id = r.id
    LEFT JOIN departments d ON u.department_id = d.id;
END //

CREATE PROCEDURE DeleteUser(IN user_id CHAR(36))
BEGIN
    DELETE FROM users WHERE id = user_id;
END //

CREATE PROCEDURE GetUsersByDepartment(IN dept_id CHAR(36))
BEGIN
    SELECT u.id, u.email, r.name AS role
    FROM users u
    JOIN roles r ON u.role_id = r.id
    WHERE u.department_id = dept_id;
END //

DELIMITER ;