-- ###############################
--            [Warning]
-- Running this script will delete all data in database.
-- ###############################

DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA IF NOT EXISTS public;
COMMENT ON SCHEMA public IS 'standard public schema';
SET search_path = "public";
SET TIME ZONE 'PRC';

create table if not exists "user"
(
    id         serial
        primary key,
    open_id    varchar(255),
    name       varchar(255)                                         not null
        unique,
    email      varchar(255)                                         not null,
    phone      varchar(255),
    password   varchar(255)                                         not null,
    created_at timestamp                                            not null,
    avatar_url varchar(1024) default '/user.png'::character varying not null,
    nick_name  varchar(255)                                         not null
);

alter table "user"
    owner to postgres;

create table if not exists problem_type
(
    id              serial
        primary key,
    description     text      not null,
    created_at      timestamp not null,
    updated_at      timestamp not null,
    user_id         integer   not null
        references "user",
    problem_type_id integer   not null,
    is_public       boolean   not null,
    analysis        text
);

alter table problem_type
    owner to postgres;

create table if not exists problem_choice
(
    id          integer      not null
        references problem_type
            on delete cascade,
    choice      varchar(255) not null,
    description text         not null,
    is_correct  boolean      not null,
    primary key (id, choice)
);

alter table problem_choice
    owner to postgres;

create table if not exists problem_answer
(
    id     integer      not null
        primary key
        references problem_type
            on delete cascade,
    answer varchar(255) not null
);

alter table problem_answer
    owner to postgres;

create table if not exists problem_set
(
    id          serial
        primary key,
    name        varchar(255)        not null,
    description text                not null,
    created_at  timestamp           not null,
    updated_at  timestamp           not null,
    user_id     integer             not null
        references "user",
    is_public   boolean             not null,
    group_id    integer default 0   not null,
    area_id     integer default 100 not null
);

alter table problem_set
    owner to postgres;

create table if not exists problem_in_problem_set
(
    problem_set_id integer not null
        references problem_set
            on delete cascade,
    problem_id     integer not null
        references problem_type
            on delete cascade,
    primary key (problem_set_id, problem_id)
);

alter table problem_in_problem_set
    owner to postgres;

create table if not exists user_favorite_problem
(
    problem_id integer   not null
        references problem_type
            on delete cascade,
    user_id    integer   not null
        references "user"
            on delete cascade,
    created_at timestamp not null,
    primary key (problem_id, user_id)
);

alter table user_favorite_problem
    owner to postgres;

create table if not exists user_favorite_problem_set
(
    problem_set_id integer   not null
        references problem_set
            on delete cascade,
    user_id        integer   not null
        references "user"
            on delete cascade,
    created_at     timestamp not null,
    primary key (problem_set_id, user_id)
);

alter table user_favorite_problem_set
    owner to postgres;

create table if not exists user_wrong_record
(
    problem_id integer   not null
        references problem_type
            on delete cascade,
    user_id    integer   not null
        references "user"
            on delete cascade,
    count      integer   not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    primary key (problem_id, user_id)
);

alter table user_wrong_record
    owner to postgres;

create table if not exists note
(
    id         serial
        primary key,
    title      varchar(255) not null,
    content    text         not null,
    created_at timestamp    not null,
    updated_at timestamp    not null,
    user_id    integer      not null
        references "user"
            on delete cascade,
    is_public  boolean      not null
);

alter table note
    owner to postgres;

create table if not exists note_review
(
    id         serial
        primary key,
    title      varchar(255) not null,
    content    text         not null,
    created_at timestamp    not null,
    updated_at timestamp    not null,
    user_id    integer      not null
        references "user"
            on delete cascade,
    note_id    integer      not null
        references note
            on delete cascade
);

alter table note_review
    owner to postgres;

create table if not exists user_like_note_review
(
    note_review_id integer   not null
        references note_review
            on delete cascade,
    user_id        integer   not null
        references "user"
            on delete cascade,
    created_at     timestamp not null,
    primary key (note_review_id, user_id)
);

alter table user_like_note_review
    owner to postgres;

create table if not exists user_like_note
(
    note_id    integer   not null
        references note
            on delete cascade,
    user_id    integer   not null
        references "user"
            on delete cascade,
    created_at timestamp not null,
    primary key (note_id, user_id)
);

alter table user_like_note
    owner to postgres;

create table if not exists user_favorite_note
(
    note_id    integer   not null
        references note
            on delete cascade,
    user_id    integer   not null
        references "user"
            on delete cascade,
    created_at timestamp not null,
    primary key (note_id, user_id)
);

alter table user_favorite_note
    owner to postgres;

create table if not exists problem_judge
(
    id         integer not null
        primary key
        references problem_type
            on delete cascade,
    is_correct boolean not null
);

alter table problem_judge
    owner to postgres;

create table if not exists "group"
(
    id          serial
        primary key,
    name        varchar(255)              not null,
    description text                      not null,
    invitation  varchar(255)              not null,
    created_at  timestamp                 not null,
    user_id     integer                   not null
        references "user"
            on delete cascade,
    area_id     integer       default 100 not null,
    avatar_url  varchar(1024) default '/group.png'::character varying
);

alter table "group"
    owner to postgres;

create table if not exists group_member
(
    group_id   integer               not null
        references "group"
            on delete cascade,
    user_id    integer               not null
        references "user"
            on delete cascade,
    created_at timestamp             not null,
    is_admin   boolean default false not null,
    is_owner   boolean default false,
    primary key (group_id, user_id)
);

alter table group_member
    owner to postgres;

create table if not exists discussion
(
    id             serial
        primary key,
    title          varchar(255) not null,
    content        text         not null,
    created_at     timestamp    not null,
    updated_at     timestamp    not null,
    user_id        integer      not null
        references "user"
            on delete cascade,
    group_id       integer      not null
        references "group"
            on delete cascade,
    is_public      boolean      not null,
    like_count     integer default 0,
    favorite_count integer default 0
);

alter table discussion
    owner to postgres;

create table if not exists user_favorite_discussion
(
    discussion_id integer   not null
        references discussion
            on delete cascade,
    user_id       integer   not null
        references "user"
            on delete cascade,
    created_at    timestamp not null,
    primary key (discussion_id, user_id)
);
alter table user_favorite_discussion
    owner to postgres;

create table if not exists discussion_review
(
    id            serial
        primary key,
    title         varchar(255) not null,
    content       text         not null,
    created_at    timestamp    not null,
    updated_at    timestamp    not null,
    user_id       integer      not null
        references "user"
            on delete cascade,
    discussion_id integer      not null
        references discussion
            on delete cascade,
    like_count    integer default 0
);

alter table discussion_review
    owner to postgres;

create table if not exists note_problem
(
    note_id    integer not null
        references note
            on delete cascade,
    problem_id integer not null
        references problem_type
            on delete cascade,
    created_at timestamp default now(),
    primary key (note_id, problem_id)
);

alter table note_problem
    owner to postgres;

create table if not exists user_like_discussion_review
(
    discussion_review_id integer   not null
        references discussion_review
            on delete cascade,
    user_id              integer   not null
        references "user"
            on delete cascade,
    created_at           timestamp not null,
    primary key (discussion_review_id, user_id)
);

alter table user_like_discussion_review
    owner to postgres;

create table if not exists user_like_discussion
(
    discussion_id integer   not null
        references discussion
            on delete cascade,
    user_id       integer   not null
        references "user"
            on delete cascade,
    created_at    timestamp not null,
    primary key (discussion_id, user_id)
);

alter table user_like_discussion
    owner to postgres;

create table if not exists group_application
(
    id         serial
        primary key,
    user_id    integer   not null
        references "user"
            on delete cascade,
    group_id   integer   not null
        references "group"
            on delete cascade,
    created_at timestamp not null,
    status     integer   not null,
    message    text
);

alter table group_application
    owner to postgres;

create table if not exists area
(
    id   integer primary key,
    name varchar(255) not null
);

alter table area
    owner to postgres;

INSERT INTO area (id, name)
VALUES (1, '综合');
INSERT INTO area (id, name)
VALUES (2, '计算机');
INSERT INTO area (id, name)
VALUES (3, '经济金融');
INSERT INTO area (id, name)
VALUES (4, '电子信息');
INSERT INTO area (id, name)
VALUES (5, '数学');
INSERT INTO area (id, name)
VALUES (6, '生物');
INSERT INTO area (id, name)
VALUES (7, '医学');
INSERT INTO area (id, name)
VALUES (8, '物理');
INSERT INTO area (id, name)
VALUES (9, '化学');
INSERT INTO area (id, name)
VALUES (10, '历史');
INSERT INTO area (id, name)
VALUES (11, '建筑');
INSERT INTO area (id, name)
VALUES (12, '交通');
INSERT INTO area (id, name)
VALUES (13, '人文社科');
INSERT INTO area (id, name)
VALUES (14, '外语');
INSERT INTO area (id, name)
VALUES (15, '体育健康');
INSERT INTO area (id, name)
VALUES (16, '公务员');
INSERT INTO area (id, name)
VALUES (17, '教师');
INSERT INTO area (id, name)
VALUES (18, '天文学');
INSERT INTO area (id, name)
VALUES (19, '地理');
INSERT INTO area (id, name)
VALUES (20, '政治');
INSERT INTO area (id, name)
VALUES (100, '其他');
CREATE EXTENSION pg_trgm;
alter function to_tsvector(text) immutable;
CREATE INDEX problem_set_search_idx ON problem_set USING GIN (to_tsvector(problem_set.name || problem_set.description));
CREATE INDEX group_search_idx ON "group" USING GIN (to_tsvector("group".name || "group".description));
CREATE INDEX note_search_idx ON note USING GIN (to_tsvector(note.title || note.content));