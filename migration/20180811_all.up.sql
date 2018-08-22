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
  point INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT current_timestamp,
  updated_at TIMESTAMPTZ DEFAULT current_timestamp,
  version BIGINT
);
CREATE INDEX assets_created_at_key ON assets (created_at);

CREATE TYPE item_status AS ENUM ('notsold', 'sold', 'soldout');

CREATE TABLE items (
  id CHAR(10) PRIMARY KEY,
  user_id CHAR(10) REFERENCES users(id) NOT NULL,
  label VARCHAR(64) NOT NULL,
  description VARCHAR(256),
  price INT NOT NULL,
  status item_status NOT NULL,
  created_at TIMESTAMPTZ DEFAULT current_timestamp,
  updated_at TIMESTAMPTZ DEFAULT current_timestamp,
  version BIGINT
);
CREATE INDEX items_user_id_created_at_key ON items (user_id, created_at);
CREATE INDEX items_status_created_at_key ON items (status, created_at);
CREATE INDEX items_created_at_key ON items (created_at);

CREATE TABLE deals (
  id CHAR(10) PRIMARY KEY,
  item_id CHAR(10) REFERENCES items(id) UNIQUE NOT NULL,
  seller_id CHAR(10) REFERENCES users(id) NOT NULL,
  buyer_id CHAR(10) REFERENCES users(id) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT current_timestamp
);
CREATE INDEX deals_seller_id_created_at_key ON deals (seller_id, created_at);
CREATE INDEX deals_buyer_id_created_at_key ON deals (buyer_id, created_at);
CREATE INDEX deals_created_at_key ON deals (created_at);
