CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE profile (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    nik BYTEA NOT NULL,
    nik_bidx VARCHAR(255) NOT NULL,
    name BYTEA NOT NULL,
    name_bidx VARCHAR(255) NOT NULL,
    email BYTEA NOT NULL,
    email_bidx VARCHAR(255) NOT NULL,
    phone BYTEA NOT NULL,
    phone_bidx VARCHAR(255) NOT NULL,
    dob BYTEA NOT NULL
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE nik_text_heap (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    content VARCHAR(100) NOT NULL,
    hash VARCHAR(100) NOT NULL,
);

CREATE TABLE name_text_heap (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    content VARCHAR(100) NOT NULL,
    hash VARCHAR(100) NOT NULL,
);

CREATE TABLE email_text_heap (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    content VARCHAR(100) NOT NULL,
    hash VARCHAR(100) NOT NULL,
);

CREATE TABLE phone_text_heap (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    content VARCHAR(100) NOT NULL,
    hash VARCHAR(100) NOT NULL,
);

