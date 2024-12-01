CREATE TABLE IF NOT EXISTS banners (
    id CHAR(6) DEFAULT substring(gen_random_uuid()::text from 1 for 6) PRIMARY KEY,
    deskbannerurl VARCHAR(255),
    mobbannerurl VARCHAR(255),
    title VARCHAR(255),
    subtitle VARCHAR(255),
    ctatext VARCHAR(255),
    ctaaction VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS products (
    id CHAR(6) DEFAULT substring(gen_random_uuid()::text from 1 for 6) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category TEXT[],
    catalog json,
    feature_desc TEXT,
    feature_list json,
    specifications json,
    techinfo json,
    tags TEXT[],
    image_gallery json,
    CONSTRAINT unique_prod_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS tags (
    id CHAR(6) DEFAULT substring(gen_random_uuid()::text from 1 for 6) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    CONSTRAINT unique_tag_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS categories (
    id CHAR(6) DEFAULT substring(gen_random_uuid()::text from 1 for 6) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    imgurl VARCHAR(255),
    description TEXT,
    CONSTRAINT unique_cat_name UNIQUE (name)
);


CREATE TABLE IF NOT EXISTS industries (
    id CHAR(6) DEFAULT substring(gen_random_uuid()::text from 1 for 6) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    imgurl VARCHAR(255),
    title TEXT,
    subtitle TEXT,
    content TEXT,
    CONSTRAINT unique_industry_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS awards (
    id CHAR(6) DEFAULT substring(gen_random_uuid()::text from 1 for 6) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    awardurl VARCHAR(255),
    type VARCHAR(255),
    CONSTRAINT unique_awards_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS clients (
    id CHAR(6) DEFAULT substring(gen_random_uuid()::text from 1 for 6) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    clienturl VARCHAR(255),
    CONSTRAINT unique_clients_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS orders (
    id CHAR(10) DEFAULT substring(regexp_replace(gen_random_uuid()::text, '-', '', 'g') from 1 for 10) PRIMARY KEY,
    invoice_number VARCHAR(255) NOT NULL,
    transporter_name VARCHAR(255),
    lr_number VARCHAR(255),
    batch_number VARCHAR(255),
    order_status VARCHAR(255),
    remarks VARCHAR(255),
    CONSTRAINT unique_order_number UNIQUE (invoice_number)
);


CREATE TABLE IF NOT EXISTS contact (
    id CHAR(10) DEFAULT substring(regexp_replace(gen_random_uuid()::text, '-', '', 'g') from 1 for 10) PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL,
    city VARCHAR(255),
    country VARCHAR(255),
    companyname VARCHAR(255),
    emailid VARCHAR(255),
    phone VARCHAR(255),
    requirement VARCHAR(255)
);
