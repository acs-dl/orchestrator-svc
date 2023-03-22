-- +migrate Up

create type request_status_enum as enum ('created', 'pending', 'finished', 'failed');

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
    status request_status_enum not null default 'created',
    module_name text not null,
    error text,
    created_at timestamp without time zone not null default current_timestamp,
    constraint fk_modules_name foreign key (module_name) references modules (name) on delete cascade
);

create index requests_payloadaction_ids on requests ((payload ->> 'action'));

-- +migrate Down

drop index if exists requests_payloadaction_ids;
drop table if exists requests;
drop table if exists modules;
drop type if exists request_status_enum;
