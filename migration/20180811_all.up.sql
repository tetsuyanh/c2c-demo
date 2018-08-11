CREATE TABLE users (
  id CHAR(10) PRIMARY KEY,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX users_created_at_key ON users (created_at);

CREATE TABLE sessions (
  id CHAR(10) PRIMARY KEY,
  user_id CHAR(10) REFERENCES users(id) NOT NULL,
  token CHAR(32) UNIQUE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT current_timestamp
);
CREATE INDEX sessions_user_id_key ON sessions (user_id);
CREATE INDEX sessions_created_at_key ON sessions (created_at);

CREATE TABLE authentications (
  id CHAR(10) PRIMARY KEY,
  user_id CHAR(10) REFERENCES users(id) UNIQUE NOT NULL,
  email VARCHAR(256) UNIQUE NOT NULL,
  password CHAR(60) NOT NULL,
  token CHAR(32) UNIQUE NOT NULL,
  enabled BOOL DEFAULT false,
  created_at TIMESTAMPTZ DEFAULT current_timestamp,
  updated_at TIMESTAMPTZ DEFAULT current_timestamp
);
CREATE INDEX authentications_created_at_key ON authentications (created_at);

CREATE TABLE assets (
  id CHAR(10) PRIMARY KEY,
  user_id CHAR(10) REFERENCES users(id) UNIQUE NOT NULL,
  point CHAR(32) DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT current_timestamp,
  updated_at TIMESTAMPTZ DEFAULT current_timestamp
);
CREATE INDEX assets_created_at_key ON assets (created_at);

