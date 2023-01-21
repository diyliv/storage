CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY, 
    user_name VARCHAR(16) NOT NULL,
    user_email VARCHAR(32) NOT NULL UNIQUE, 
    user_hashed_password VARCHAR NOT NULL, 
    user_updated_password TIMESTAMP DEFAULT NOW() NOT NULL, 
    user_created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS users_keys(
    user_id INT NOT NULL, 
    user_public_key VARCHAR NOT NULL,
    user_passphrase VARCHAR(16) NOT NULL
);
