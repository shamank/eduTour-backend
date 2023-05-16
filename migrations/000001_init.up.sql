CREATE TABLE ACCOUNTS
(
    id            serial       not null unique,
    username      varchar(255) not null unique,
    email         varchar(255) not null unique,
    password_hash varchar(255) not null,

    created_at    TIMESTAMP default CURRENT_TIMESTAMP,


    primary key (id)
)

CREATE TABLE REFRESH_TOKENS
(
    id            serial                                         not null unique,
    user_id       int references ACCOUNTS (id) on delete cascade not null,
    refresh_token varchar(255)                                   not null,
    expire_at     TIMESTAMP
        FOREIGN KEY (user_id) REFERENCES ACCOUNT (id)
)