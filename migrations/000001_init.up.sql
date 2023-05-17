CREATE TABLE ACCOUNTS
(
    id            serial       not null unique,
    username      varchar(255) not null unique,
    email         varchar(255) not null unique,
    password_hash varchar(255) not null,

    created_at    TIMESTAMP default CURRENT_TIMESTAMP
);

CREATE TABLE REFRESH_TOKENS
(
    id            serial                                         not null unique,
    user_id       int references ACCOUNTS (id) on delete cascade not null,
    refresh_token varchar(255)                                   not null,
    expire_at     TIMESTAMP,
    black_list    bool default false                             not null,
    FOREIGN KEY (user_id) REFERENCES ACCOUNTS (id)
);

CREATE TABLE EVENTS
(
    id          serial       not null unique,
    name        varchar(255) not null,
    description text         not null,
    start_date  TIMESTAMP,
    end_date    TIMESTAMP,
    cover       varchar(255),
    content     text,
    created_at  TIMESTAMP default CURRENT_TIMESTAMP
);

CREATE TABLE CATEGORIES
(
    id   serial       not null unique,
    name varchar(255) not null unique,
    slug varchar(255) not null unique

);

CREATE TABLE EVENTS_CATEGORIES
(
    id          serial not null unique,
    event_id    int,
    category_id int,

    FOREIGN KEY (event_id) REFERENCES EVENTS (id),
    FOREIGN KEY (category_id) REFERENCES CATEGORIES (id)
);