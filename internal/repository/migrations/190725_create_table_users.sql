-- +goose Up
-- +goose StatementBegin
create table if not exists users
(
    id       varchar(50) primary key,
    login    varchar(50) not null unique,
    password text        not null,
    jwt      text        not null
);

CREATE INDEX users_on_login_password_idx ON users (login, password);
-- +goose StatementEnd