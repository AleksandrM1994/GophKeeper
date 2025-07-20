-- +goose Up
-- +goose StatementBegin
CREATE
TYPE private_data_type AS ENUM ('UNKNOWN', 'TEXT', 'FILE', 'AUTH', 'BANK');

create table if not exists private_data
(
    id         varchar(50) primary key,
    data       json        not null,
    type private_data_type default 'UNKNOWN',
    created_at TIMESTAMP   not null,
    updated_at TIMESTAMP   not null,
    user_id    varchar(50) not null
);

CREATE INDEX private_data_on_user_id_type_idx ON private_data (user_id, type);
-- +goose StatementEnd