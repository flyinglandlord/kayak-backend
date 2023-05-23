-- ###############################
--            [Warning]
-- Running this script will delete all data in database.
-- ###############################

DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA IF NOT EXISTS public;
COMMENT ON SCHEMA public IS 'standard public schema';
SET search_path = "public";
SET TIME ZONE 'PRC';

create table "user"
(
    id         serial
        primary key,
    union_id   varchar(255),
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

create table problem_type
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

create table problem_choice
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

create table problem_answer
(
    id     integer      not null
        primary key
        references problem_type
            on delete cascade,
    answer varchar(255) not null
);

alter table problem_answer
    owner to postgres;

create table problem_set
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

create table problem_in_problem_set
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

create table user_favorite_problem
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

create table user_favorite_problem_set
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

create table user_wrong_record
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

create table note
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

create table note_review
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

create table user_like_note_review
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

create table user_like_note
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

create table user_favorite_note
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

create table problem_judge
(
    id         integer not null
        primary key
        references problem_type
            on delete cascade,
    is_correct boolean not null
);

alter table problem_judge
    owner to postgres;

create table "group"
(
    id          serial
        primary key,
    name        varchar(255)        not null,
    description text                not null,
    invitation  varchar(255)        not null,
    created_at  timestamp           not null,
    user_id     integer             not null
        references "user"
            on delete cascade,
    area_id     integer default 100 not null
);

alter table "group"
    owner to postgres;

create table group_member
(
    group_id   integer   not null
        references "group"
            on delete cascade,
    user_id    integer   not null
        references "user"
            on delete cascade,
    created_at timestamp not null,
    primary key (group_id, user_id)
);

alter table group_member
    owner to postgres;

create table discussion
(
    id         serial
        primary key,
    title      varchar(255)      not null,
    content    text              not null,
    created_at timestamp         not null,
    updated_at timestamp         not null,
    user_id    integer           not null
        references "user"
            on delete cascade,
    group_id   integer           not null
        references "group"
            on delete cascade,
    is_public  boolean           not null,
    like_count integer default 0 not null
);

alter table discussion
    owner to postgres;

create table discussion_review
(
    id            serial
        primary key,
    title         varchar(255)      not null,
    content       text              not null,
    created_at    timestamp         not null,
    updated_at    timestamp         not null,
    user_id       integer           not null
        references "user"
            on delete cascade,
    discussion_id integer           not null
        references discussion
            on delete cascade,
    like_count    integer default 0 not null
);

alter table discussion_review
    owner to postgres;

create table user_like_discussion_review
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

create table user_like_discussion
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

create table note_problem
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

