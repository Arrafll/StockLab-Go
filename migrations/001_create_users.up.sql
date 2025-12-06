CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(100) NULL,
    role VARCHAR(50) DEFAULT 'staff',
    avatar BYTEA NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (email, name, phone, role, password)
VALUES ('admin@gmail.com', 'admin', '081412412412', 'admin', '$2a$10$fJetUatvJZjhAphuzg90ju6CQB/WIs2pZgT6faLOWTuCAwWpTANYG');