CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,  -- Store hashed passwords
    role VARCHAR(50) NOT NULL
);

-- Insert a predefined admin user (password should be hashed):
-- AdminPass@123
INSERT INTO users (username, password, role) VALUES ('admin', '01dff78c1dd4a0f69be8c79b2eb8a175', 'admin');
