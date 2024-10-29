CREATE TABLE IF NOT EXISTS banners (
    id SERIAL PRIMARY KEY,
    deskbannerurl VARCHAR(255),
    mobbannerurl VARCHAR(255),
    title VARCHAR(255),
    subtitle VARCHAR(255),
    ctatext VARCHAR(255),
    ctaaction VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS  products (
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    category       JSONB,
    catalog        TEXT,
    features       JSONB,
    specifications JSONB,
    techinfo       JSONB,
    tags           JSONB,
    images         JSONB
);

CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    CONSTRAINT unique_tag_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    imgurl VARCHAR(255),
    description TEXT,
    CONSTRAINT unique_cat_name UNIQUE (name)
);