DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA IF NOT EXISTS public;
COMMENT ON SCHEMA public IS 'standard public schema';
SET search_path = "public";

CREATE TABLE "user"
(
    id         serial        NOT NULL,
    union_id   varchar(255),
    name       varchar(255)  NOT NULL UNIQUE,
    email      varchar(255),
    phone      varchar(255),
    password   varchar(255)  NOT NULL,
    created_at timestamp     NOT NULL,
    avatar_url VARCHAR(1024) NOT NULL DEFAULT '/test/public/user.svg',
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
    FOREIGN KEY ("id") REFERENCES problem_type ("id")
);

CREATE TABLE problem_judge
(
    "id"         serial  NOT NULL,
    "is_correct" boolean NOT NULL,
    FOREIGN KEY ("id") REFERENCES problem_type ("id")
);

CREATE TABLE problem_answer
(
    "id"     serial       NOT NULL,
    "answer" varchar(255) NOT NULL,
    FOREIGN KEY ("id") REFERENCES "problem_type" ("id")
);

CREATE TABLE problemset
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

CREATE TABLE problem_in_problemset
(
    "problemset_id" integer NOT NULL,
    "problem_id"    integer NOT NULL,
    PRIMARY KEY ("problemset_id", "problem_id"),
    FOREIGN KEY ("problemset_id") REFERENCES problemset ("id"),
    FOREIGN KEY ("problem_id") REFERENCES problem_type ("id")
);

CREATE TABLE user_favorite_problemset
(
    "problemset_id" integer NOT NULL,
    "user_id"       integer NOT NULL,
    PRIMARY KEY ("problemset_id", "user_id"),
    FOREIGN KEY ("problemset_id") REFERENCES problemset ("id"),
    FOREIGN KEY ("user_id") REFERENCES "user" ("id")
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
    FOREIGN KEY ("user_id") REFERENCES "user" ("id")
);
