CREATE TABLE USERS
(
    id            serial                              not null unique,
    username      varchar(255)                        not null unique,
    email         varchar(255)                        not null unique,
    password_hash varchar(255)                        not null,
    phone         varchar(255) unique,
    avatar        varchar(255),

    first_name    varchar(255),
    last_name     varchar(255),
    middle_name   varchar(255),

    created_at    TIMESTAMP default CURRENT_TIMESTAMP not null
--     last_visit_at TIMESTAMP

);

CREATE TABLE ROLES
(
    id   serial not null unique,
    name varchar unique
);

INSERT INTO ROLES
VALUES (0, 'user'),
       (1, 'admin');

CREATE TABLE USERS_ROLES
(
    user_id int not null,
    role_id int not null,
    PRIMARY KEY (user_id, role_id),

    FOREIGN KEY (user_id) REFERENCES USERS (id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES ROLES (id) ON DELETE CASCADE
);

CREATE TABLE REFRESH_TOKENS
(
    id            serial                    not null unique,
    user_id       int references USERS (id) not null,
    refresh_token varchar(255)              not null,
    expire_at     TIMESTAMP                 not null,
    black_list    bool default false        not null,
    FOREIGN KEY (user_id) REFERENCES USERS (id) ON DELETE CASCADE
);

-- CREATE TABLE EVENTS
-- (
--     id          serial       not null unique,
--     name        varchar(255) not null,
--     description text         not null,
--     start_date  TIMESTAMP,
--     end_date    TIMESTAMP,
--     cover       varchar(255),
--     content     text,
--     created_at  TIMESTAMP default CURRENT_TIMESTAMP
-- );
--
-- CREATE TABLE CATEGORIES
-- (
--     id   serial       not null unique,
--     name varchar(255) not null unique,
--     slug varchar(255) not null unique
--
-- );
--
-- CREATE TABLE EVENTS_CATEGORIES
-- (
--     id          serial not null unique,
--     event_id    int,
--     category_id int,
--
--     FOREIGN KEY (event_id) REFERENCES EVENTS (id),
--     FOREIGN KEY (category_id) REFERENCES CATEGORIES (id)
-- );