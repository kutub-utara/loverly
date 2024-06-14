BEGIN;

-- Create the table users
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,

    -- Utility columns
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    email VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TYPE GENDER AS ENUM ('male', 'female');

-- Create the table profiles
CREATE TABLE profiles(
    id BIGSERIAL PRIMARY KEY,   

    -- Utility columns
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    user_id BIGINT NOT NULL UNIQUE,
    name VARCHAR NOT NULL,
    birthday DATE,
    gender GENDER,
    location VARCHAR,
    bio TEXT,
    profile_picture VARCHAR,
    interests TEXT
);

-- Create the table photos
CREATE TABLE photos(
    id BIGSERIAL PRIMARY KEY,   

    -- Utility columns
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    user_id BIGINT NOT NULL,
    photo VARCHAR
);

CREATE TYPE DIRECTION AS ENUM ('left', 'right');

-- Create the table swipes
CREATE TABLE swipes(
    id BIGSERIAL PRIMARY KEY,

    -- Utility columns
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    swiper_id BIGINT NOT NULL,
    swiped_id BIGINT NOT NULL,
    direction DIRECTION
);

-- Create the table matchs
CREATE table matchs(
    id BIGSERIAL PRIMARY KEY,

     -- Utility columns
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    user_id_1 BIGINT NOT NULL,
    user_id_2 BIGINT NOT NULL
);

CREATE TYPE PLAN AS ENUM ('unlimited', 'verified');

-- Create the table subscriptions
CREATE table subscriptions(
    id BIGSERIAL PRIMARY KEY,

     -- Utility columns
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    user_id BIGINT NOT NULL,
    plan PLAN,
    start_date DATE,
    end_date DATE
);

ALTER TABLE ONLY profiles
    ADD CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES users(id) NOT VALID;

ALTER TABLE ONLY photos
    ADD CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES users(id) NOT VALID;

ALTER TABLE ONLY swipes
    ADD CONSTRAINT swiper_id FOREIGN KEY (swiper_id) REFERENCES users(id) NOT VALID;

ALTER TABLE ONLY swipes
    ADD CONSTRAINT swiped_id FOREIGN KEY (swiped_id) REFERENCES users(id) NOT VALID;

ALTER TABLE ONLY swipes
    ADD CONSTRAINT unique_swipes_id UNIQUE (swiper_id, swiped_id);

ALTER TABLE ONLY matchs
    ADD CONSTRAINT user_id_1 FOREIGN KEY (user_id_1) REFERENCES users(id) NOT VALID;

ALTER TABLE ONLY matchs
    ADD CONSTRAINT user_id_2 FOREIGN KEY (user_id_2) REFERENCES users(id) NOT VALID;

ALTER TABLE ONLY matchs
    ADD CONSTRAINT unique_matchs_id UNIQUE (user_id_1, user_id_2);

ALTER TABLE ONLY subscriptions
    ADD CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES users(id) NOT VALID;

COMMIT;