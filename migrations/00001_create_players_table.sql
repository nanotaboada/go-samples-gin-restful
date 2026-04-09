-- +goose Up
CREATE TABLE IF NOT EXISTS players (
    id           TEXT        PRIMARY KEY,
    firstName    VARCHAR(100),
    middleName   VARCHAR(100),
    lastName     VARCHAR(100),
    dateOfBirth  TEXT,
    squadNumber  INTEGER     UNIQUE NOT NULL,
    position     VARCHAR(50),
    abbrPosition VARCHAR(10),
    team         VARCHAR(100),
    league       VARCHAR(100),
    starting11   BOOLEAN
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_players_squad_number ON players (squadNumber);

-- +goose Down
DROP INDEX IF EXISTS idx_players_squad_number;
DROP TABLE IF EXISTS players;
