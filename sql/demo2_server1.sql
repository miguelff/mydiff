DROP DATABASE IF EXISTS acme_inc;
CREATE DATABASE IF NOT EXISTS acme_inc CHARACTER SET UTF8mb4 COLLATE utf8mb4_bin;
USE acme_inc;

DROP TABLE IF EXISTS  employees;
CREATE TABLE employees (
    id          INT             NOT NULL,
    birth_date  DATE            NOT NULL,
    first_name  VARCHAR(14)     NOT NULL,
    last_name   VARCHAR(16)     NOT NULL,
    hire_date   DATE,
    PRIMARY KEY (id),
    UNIQUE KEY unique_name(first_name, last_name)
);

CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) NOT NULL,
	UNIQUE KEY version_key(version)
);

INSERT INTO schema_migrations values (20190815193300);
INSERT INTO schema_migrations values (20190817000000);