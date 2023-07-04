-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS servers
(
    id         UUID PRIMARY KEY NOT NULL   DEFAULT uuid_generate_v1(),
    url        TEXT,
    created_at timestamp without time zone default (now() at time zone 'utc')
);

CREATE TABLE IF NOT EXISTS files
(
    id          UUID PRIMARY KEY NOT NULL   DEFAULT uuid_generate_v1(),
    name        TEXT             NOT NULL,
    size        BIGINT           NOT NULL,
    is_uploaded BOOLEAN          NOT NULL   DEFAULT FALSE,
    created_at  timestamp without time zone default (now() at time zone 'utc')
);

CREATE TABLE IF NOT EXISTS chunks
(
    id          UUID PRIMARY KEY NOT NULL   DEFAULT uuid_generate_v1(),
    file_id     UUID             NOT NULL,
    server_id   UUID             NOT NULL,
    index       INT              NOT NULL,
    size        BIGINT           NOT NULL,
    is_uploaded BOOLEAN          NOT NULL   DEFAULT FALSE,
    hash        TEXT,
    created_at  timestamp without time zone default (now() at time zone 'utc')
);

CREATE INDEX ON chunks (file_id);
CREATE INDEX ON chunks (server_id);

-- +migrate Down
DROP TABLE servers;
DROP TABLE files;
DROP TABLE chunks;
