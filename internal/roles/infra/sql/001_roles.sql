/*
    This table represents the roles of the users
*/
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    label ENUM('prof', 'admin') UNIQUE NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Migrate the former default role before constraining the column to the
-- supported role labels. The role ID is preserved, so existing assignments
-- continue to work as prof assignments.
UPDATE roles SET label = 'prof' WHERE CAST(label AS CHAR) = 'teacher';

ALTER TABLE roles
    MODIFY COLUMN label ENUM('prof', 'admin') NOT NULL;

/*
    This table represents the many-to-many relationship between users and roles
*/
CREATE TABLE IF NOT EXISTS users_roles (
    id BIGINT UNSIGNED AUTO_INCREMENT NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    role_id BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY(id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- insert the default roles in roles table
INSERT IGNORE INTO roles (label)
VALUES ('prof'),
    ('admin');
