
-- Create users if not exist

-- Admin user, Password is "password"
INSERT INTO users (first_name, last_name, email, password, access_type, created_at, updated_at, last_login_at)
VALUES 
    ('Admin', 'User', 'admin@example.com', '5f4dcc3b5aa765d61d8327deb882cf99', 'Admin', NOW(), NOW(), NOW())
ON CONFLICT (email) 
DO NOTHING;

-- Customer user, Password is "password"
INSERT INTO users (first_name, last_name, email, password, access_type, created_at, updated_at, last_login_at)
VALUES 
    ('Customer', 'User', 'customer@example.com', '5f4dcc3b5aa765d61d8327deb882cf99', 'Customer', NOW(), NOW(), NOW())
ON CONFLICT (email) 
DO NOTHING;
