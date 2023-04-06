-- ###############################
--            [Warning]
-- Running this script will delete all data in database.
-- ###############################

DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA IF NOT EXISTS public;
COMMENT ON SCHEMA public IS 'standard public schema';
SET search_path = "public";
SET TIME ZONE 'PRC';

CREATE TABLE "user"
(
    id         serial        NOT NULL,
    union_id   varchar(255),
    name       varchar(255)  NOT NULL UNIQUE,
    email      varchar(255),
    phone      varchar(255),
    password   varchar(255)  NOT NULL,
    created_at timestamp     NOT NULL,
    avatar_url VARCHAR(1024) NOT NULL DEFAULT '/public/user.png',
    PRIMARY KEY ("id")
);

CREATE TABLE problem_type
(
    "id"              serial    NOT NULL,
    "description"     text      NOT NULL,
    "created_at"      timestamp NOT NULL,
    "updated_at"      timestamp NOT NULL,
    "user_id"         integer   NOT NULL,
    "problem_type_id" integer   NOT NULL,
    "is_public"       boolean   NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("user_id") REFERENCES "user" ("id")
);

CREATE TABLE problem_choice
(
    "id"          integer      NOT NULL,
    "choice"      varchar(255) NOT NULL,
    "description" text         NOT NULL,
    "is_correct"  boolean      NOT NULL,
    PRIMARY KEY ("id", "choice"),
    FOREIGN KEY ("id") REFERENCES problem_type ("id") ON DELETE CASCADE
);

CREATE TABLE problem_answer
(
    "id"     integer      NOT NULL,
    "answer" varchar(255) NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("id") REFERENCES "problem_type" ("id") ON DELETE CASCADE
);

CREATE TABLE problemSet
(
    "id"          serial       NOT NULL,
    "name"        varchar(255) NOT NULL,
    "description" text         NOT NULL,
    "created_at"  timestamp    NOT NULL,
    "updated_at"  timestamp    NOT NULL,
    "user_id"     integer      NOT NULL,
    "is_public"   boolean      NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("user_id") REFERENCES "user" ("id")
);

CREATE TABLE problem_in_problemSet
(
    "problem_set_id" integer NOT NULL,
    "problem_id"    integer NOT NULL,
    PRIMARY KEY ("problem_set_id", "problem_id"),
    FOREIGN KEY ("problem_set_id") REFERENCES problemSet ("id") ON DELETE CASCADE,
    FOREIGN KEY ("problem_id") REFERENCES problem_type ("id") ON DELETE CASCADE
);

CREATE TABLE user_favorite_problem
(
    "problem_id" integer   NOT NULL,
    "user_id"    integer   NOT NULL,
    "created_at" timestamp NOT NULL,
    PRIMARY KEY ("problem_id", "user_id"),
    FOREIGN KEY ("problem_id") REFERENCES problem_type ("id") ON DELETE CASCADE,
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);

CREATE TABLE user_favorite_problemSet
(
    "problem_set_id" integer   NOT NULL,
    "user_id"       integer   NOT NULL,
    "created_at"    timestamp NOT NULL,
    PRIMARY KEY ("problem_set_id", "user_id"),
    FOREIGN KEY ("problem_set_id") REFERENCES problemSet ("id") ON DELETE CASCADE,
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);

CREATE TABLE user_wrong_record
(
    "problem_id" integer   NOT NULL,
    "user_id"    integer   NOT NULL,
    "count"      integer   NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    PRIMARY KEY ("problem_id", "user_id"),
    FOREIGN KEY ("problem_id") REFERENCES problem_type ("id") ON DELETE CASCADE,
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);

CREATE TABLE note
(
    "id"         serial       NOT NULL,
    "title"      varchar(255) NOT NULL,
    "content"    text         NOT NULL,
    "created_at" timestamp    NOT NULL,
    "updated_at" timestamp    NOT NULL,
    "user_id"    integer      NOT NULL,
    "is_public"  boolean      NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);

CREATE TABLE note_review
(
    "id"         serial       NOT NULL,
    "title"      varchar(255) NOT NULL,
    "content"    text         NOT NULL,
    "created_at" timestamp    NOT NULL,
    "updated_at" timestamp    NOT NULL,
    "user_id"    integer      NOT NULL,
    "note_id"    integer      NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("note_id") REFERENCES note ("id") ON DELETE CASCADE
);

CREATE TABLE user_like_note_review
(
    "note_review_id" integer   NOT NULL,
    "user_id"        integer   NOT NULL,
    "created_at"     timestamp NOT NULL,
    PRIMARY KEY ("note_review_id", "user_id"),
    FOREIGN KEY ("note_review_id") REFERENCES note_review ("id") ON DELETE CASCADE,
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);

CREATE TABLE user_like_note
(
    "note_id"    integer   NOT NULL,
    "user_id"    integer   NOT NULL,
    "created_at" timestamp NOT NULL,
    PRIMARY KEY ("note_id", "user_id"),
    FOREIGN KEY ("note_id") REFERENCES note ("id") ON DELETE CASCADE,
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);

CREATE TABLE user_favorite_note
(
    "note_id"    integer   NOT NULL,
    "user_id"    integer   NOT NULL,
    "created_at" timestamp NOT NULL,
    PRIMARY KEY ("note_id", "user_id"),
    FOREIGN KEY ("note_id") REFERENCES note ("id") ON DELETE CASCADE,
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);
