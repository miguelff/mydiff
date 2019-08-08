--  Sample employee database adapted from https://raw.githubusercontent.com/datacharmer/test_db/master/employees.sql
--

DROP DATABASE IF EXISTS employees;
CREATE DATABASE IF NOT EXISTS employees;
USE employees;

SELECT 'CREATING DATABASE STRUCTURE' as 'INFO';

DROP TABLE IF EXISTS dept_emp,
                     dept_manager,
                     titles,
                     salaries,
                     employees,
                     departments;

/*!50503 set default_storage_engine = InnoDB */;
/*!50503 select CONCAT('storage engine: ', @@default_storage_engine) as INFO */;

CREATE TABLE employees (
    emp_id      INT             NOT NULL,
    birth_date  DATE            NOT NULL,
    first_name  VARCHAR(14)     NOT NULL,
    last_name   VARCHAR(16)     NOT NULL,
    gender      ENUM ('M','F')  NOT NULL,
    hire_date   DATE            NOT NULL,
    PRIMARY KEY (emp_id)
);

CREATE TABLE departments (
    dept_id     CHAR(4)         NOT NULL,
    dept_name   VARCHAR(40)     NOT NULL,
    PRIMARY KEY (dept_id),
    UNIQUE  KEY (dept_name)
);

CREATE TABLE dept_manager (
   emp_id       INT             NOT NULL,
   dept_id      CHAR(4)         NOT NULL,
   from_date    DATE            NOT NULL,
   to_date      DATE            NOT NULL,
   FOREIGN KEY (emp_id)  REFERENCES employees (emp_id)    ON DELETE CASCADE,
   FOREIGN KEY (dept_id) REFERENCES departments (dept_id) ON DELETE CASCADE,
   PRIMARY KEY (emp_id,dept_id)
);

CREATE TABLE dept_emp (
    emp_id      INT             NOT NULL,
    dept_id     CHAR(4)         NOT NULL,
    from_date   DATE            NOT NULL,
    to_date     DATE            NOT NULL,
    FOREIGN KEY (emp_id)  REFERENCES employees   (emp_id)  ON DELETE CASCADE,
    FOREIGN KEY (dept_id) REFERENCES departments (dept_id) ON DELETE CASCADE,
    PRIMARY KEY (emp_id,dept_id)
);

CREATE TABLE titles (
    emp_id      INT             NOT NULL,
    title       VARCHAR(50)     NOT NULL,
    from_date   DATE            NOT NULL,
    to_date     DATE,
    FOREIGN KEY (emp_id) REFERENCES employees (emp_id) ON DELETE CASCADE,
    PRIMARY KEY (emp_id,title, from_date)
)
;

CREATE TABLE salaries (
    emp_id      INT             NOT NULL,
    salary      INT             NOT NULL,
    from_date   DATE            NOT NULL,
    to_date     DATE            NOT NULL,
    FOREIGN KEY (emp_id) REFERENCES employees (emp_id) ON DELETE CASCADE,
    PRIMARY KEY (emp_id, from_date)
)
;

CREATE OR REPLACE VIEW dept_emp_latest_date AS
    SELECT emp_id, MAX(from_date) AS from_date, MAX(to_date) AS to_date
    FROM dept_emp
    GROUP BY emp_id;

# shows only the current department for each employee
CREATE OR REPLACE VIEW current_dept_emp AS
    SELECT l.emp_id, dept_id, l.from_date, l.to_date
    FROM dept_emp d
        INNER JOIN dept_emp_latest_date l
        ON d.emp_id=l.emp_id AND d.from_date=l.from_date AND l.to_date = d.to_date;

flush /*!50503 binary */ logs;