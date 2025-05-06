CREATE TABLE IF NOT EXISTS users
(
    id     SERIAL PRIMARY KEY,
    name   TEXT        not null,
    phone  TEXT unique not null,

)