create table users
(
    id           uuid not null
        primary key,
    first_name   text,
    last_name    text,
    phone_number text
        constraint uni_users_phone_number
            unique,
    address      text,
    pin          text,
    created_date timestamp with time zone,
    balance      numeric
);

alter table users
    owner to postgres;

create table top_ups
(
    id             uuid not null
        primary key,
    user_id        uuid
        constraint fk_users_top_ups
            references users,
    amount         numeric,
    balance_before numeric,
    balance_after  numeric,
    created_date   timestamp with time zone
);

alter table top_ups
    owner to postgres;

create index idx_top_ups_user_id
    on top_ups (user_id);

create table payments
(
    id             uuid not null
        primary key,
    user_id        uuid
        constraint fk_payments_user
            references users,
    amount         numeric,
    remarks        text,
    balance_before numeric,
    balance_after  numeric,
    created_date   timestamp with time zone
);

alter table payments
    owner to postgres;

create index idx_payments_user_id
    on payments (user_id);

create table transfers
(
    id             uuid not null
        primary key,
    from_user_id   uuid
        constraint fk_transfers_from_user
            references users,
    to_user_id     uuid
        constraint fk_transfers_to_user
            references users,
    amount         numeric,
    remarks        text,
    balance_before numeric,
    balance_after  numeric,
    created_date   timestamp with time zone
);

alter table transfers
    owner to postgres;

create index idx_transfers_to_user_id
    on transfers (to_user_id);

create index idx_transfers_from_user_id
    on transfers (from_user_id);

