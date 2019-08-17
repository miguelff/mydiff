DROP DATABASE IF EXISTS acme_inc;
CREATE DATABASE IF NOT EXISTS acme_inc;
USE acme_inc;

DROP TABLE IF EXISTS  employees;
CREATE TABLE employees (
    birth_date  DATE            NOT NULL,
    first_name  VARCHAR(14)     NOT NULL,
    last_name   VARCHAR(16)     NOT NULL,
    hire_date   DATE            NOT NULL,
    PRIMARY KEY (first_name, last_name)
) ENGINE=INNODB;