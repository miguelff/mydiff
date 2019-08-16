DROP DATABASE IF EXISTS acme_inc;
CREATE DATABASE IF NOT EXISTS acme_inc;
USE acme_inc;

DROP TABLE IF EXISTS  employees;
CREATE TABLE employees (
    id          INT             NOT NULL,
    birth_date  DATE,
    first_name  VARCHAR(14),
    last_name   VARCHAR(16),
    hire_date   DATE,
    PRIMARY KEY (id),
    UNIQUE KEY unique_name(first_name, last_name)
) ENGINE=INNODB CHARACTER SET utf8mb4;