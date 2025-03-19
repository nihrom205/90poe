CREATE TABLE ports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    key varchar(255),
    name varchar(255),
    city varchar(255),
    country varchar(255),
    alias text,
    regions text,
    coordinates text,
    province varchar(255),
    timezone varchar(255),
    unlocs text,
    code varchar(255)
);