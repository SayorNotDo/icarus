CREATE DATABASE IF NOT EXISTS icarus DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE icarus;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS department;
CREATE TABLE department (
    did     int(11) NOT NULL AUTO_INCREMENT,
    name    varchar(255) NOT NULL,
    center  varchar(255) NOT NULL,
    company varchar(255) NOT NULL,
    PRIMARY KEY (did)
);

DROP TABLE IF EXISTS user;
CREATE TABLE user (
    uid             int(11) NOT NULL AUTO_INCREMENT,
    username        varchar(255) NOT NULL,
    password        varchar(255) NOT NULL,
    chinese_name    varchar(255),
    department_id   int,
    role_id         int(11) NOT NULL,
    employee_id     varchar(255) NOT NULL,
    position        varchar(255) NOT NULL,
    email           varchar(255),
    phone           varchar(255),
    status          int(11) NOT NULL,
    join_date       date         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login_time timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (uid),
    FOREIGN KEY (department_id) REFERENCES department (did)
);

DROP TABLE IF EXISTS project;
CREATE TABLE project (
    pid              int(11) NOT NULL AUTO_INCREMENT,
    create_time      timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_update_time timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    name             varchar(255) NOT NULL,
    manager          varchar(255) NOT NULL,
    description      text,
    status           int(11) NOT NULL,
    start_time       timestamp,
    finish_time      timestamp,
    reference        text,
    PRIMARY KEY (pid)
);

DROP TABLE IF EXISTS test_plan;
CREATE TABLE test_plan (
    tpid             int(11) NOT NULL AUTO_INCREMENT,
    create_time      timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    creator          timestamp    NOT NULL,
    last_update_time timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    start_time       timestamp,
    finish_time      timestamp,
    deadline         date,
    name             varchar(255) NOT NULL,
    project          varchar(255) NOT NULL,
    test_case        text,
    description      text,
    status           int(11) NOT NULL,
    reference        text,
    PRIMARY KEY (tpid)
);

DROP TABLE IF EXISTS test_case;
CREATE TABLE test_case (
    tcid             int (11) NOT NULL AUTO_INCREMENT,
    name             varchar(255)   NOT NULL,
    create_time      timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    creator          varchar(255)   NOT NULL,
    priority         int(11) NOT NULL,
    status           int(11) NOT NULL,
    object           int(11) NOT NULL,
    reference        text,
    description      text,
    content          text,
    PRIMARY KEY (tcid)
);

DROP TABLE IF EXISTS task;
CREATE TABLE task (
    tid              int(11) NOT NULL AUTO_INCREMENT,
    create_time      timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    start_time       timestamp,
    finish_time      timestamp,
    name             varchar(255),
    description      text,
    status           int(11) NOT NULL,
    priority         int(11) NOT NULL,
    PRIMARY KEY (tid)
);

DROP TABLE IF EXISTS failure;
CREATE TABLE failure (
    fid              int(11) NOT NULL AUTO_INCREMENT,
    create_time      timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    task             varchar(255),
    type             int(11) NOT NULL,
    PRIMARY KEY (fid)
);

DROP TABLE IF EXISTS bug;
CREATE TABLE bug (
    bid              int(11) NOT NULL AUTO_INCREMENT,
    create_time      timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_update_time timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    fixed_time       timestamp,
    failure_id       int,
    name             text,
    priority         int(11) NOT NULL,
    test_case        varchar(255) NOT NULL,
    status           int(11) NOT NULL,
    actual_result    text,
    analysis         text,
    reference        text,
    PRIMARY KEY (bid),
    FOREIGN KEY (failure_id) REFERENCES failure (fid)
);
SET FOREIGN_KEY_CHECKS = 1;