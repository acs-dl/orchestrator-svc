-- +migrate Up

create type request_status_enum as enum ('pending', 'in progress', 'success', 'invited', 'failed');
create type transaction_action_enum as enum ('single','delete_user');

create table if not exists modules (
    name text primary key,
    title text not null,
    topic text not null,
    prefix text not null,
    link text not null,
    is_module boolean not null default false -- in table `modules` except `gitlab`, `github` and other can be `unverified-svc`, `identity-svc`, but in some cases i need only modules such `telegram`, `gitlab`
);

create table if not exists requests (
    id uuid primary key,
    from_user_id bigint not null,
    to_user_id bigint not null,
    payload jsonb not null,
    status request_status_enum not null default 'pending',
    module_name text not null,
    error text,
    description text,
    created_at timestamp without time zone not null default current_timestamp,
    constraint fk_modules_name foreign key (module_name) references modules (name) on delete cascade
);

create index requests_payload_action_ids on requests ((payload ->> 'action'));

create table if not exists request_transactions(
    id uuid primary key,
    action transaction_action_enum not null,
    requests jsonb not null  --map[uuid]:bool - to check whether request was handled
);

-- +migrate Down

drop table if exists request_transactions;
drop index if exists requests_payload_action_ids;
drop table if exists requests;
drop table if exists modules;
drop type if exists request_status_enum;
drop type if exists transaction_action_enum;
