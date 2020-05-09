create database projectA charset utf8mb4;

use projectA;

CREATE USER 'projectAUser'@'localhost' IDENTIFIED BY 'peiqingtao';
grant all privileges on projectA.* to projectAUser@localhost;
flush privileges;

create table if not exists a_categories (
    id int unsigned auto_increment,
    parent_id int unsigned,
    name varchar(255),
    logo varchar(255),
    description varchar(255),
    sort_order int,
    meta_title varchar(255),
    meta_keywords varchar(255),
    meta_description varchar(255),
    created_at timestamp ,
    updated_at timestamp ,
    deleted_at timestamp ,
    primary key (id),
    index (parent_id),
    index (name),
    index (sort_order)
)engine innodb charset utf8mb4;

alter table a_categories add column created_at timestamp, add column updated_at timestamp, add column deleted_at timestamp;