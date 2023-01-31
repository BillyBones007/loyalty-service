CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users(uuid UUID NOT NULL PRIMARY KEY, login VARCHAR(50) UNIQUE NOT NULL, password VARCHAR(255) NOT NULL);
CREATE TABLE IF NOT EXISTS balance(user_id UUID NOT NULL PRIMARY KEY, balance INT DEFAULT 0);
CREATE TABLE IF NOT EXISTS orders( num_order VARCHAR(11) UNIQUE NOT NULL, time_order DATE NOT NULL, status VARCHAR(12) NOT NULL, debet_points INT DEFAULT 0, credit_points INT DEFAULT 0, user_id UUID NOT NULL REFERENCES ballance (user_id));
