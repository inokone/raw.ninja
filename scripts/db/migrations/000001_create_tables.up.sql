CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(36),
    email VARCHAR(255),
    phone VARCHAR(30),
    passhash VARCHAR(100),

    PRIMARY KEY (user_id)
);

CREATE UNIQUE INDEX user__email ON users(email);

CREATE TABLE IF NOT EXISTS descriptors (
    descriptor_id VARCHAR(36),
    user_id VARCHAR(36) references users(user_id),
    filename VARCHAR(1000),
    uploaded TIMESTAMP,
    format VARCHAR(100),
    width INTEGER,
    height INTEGER,
    format VARCHAR(36),

    PRIMARY KEY (descriptor_id)
);

CREATE INDEX descriptor__uploaded ON descriptors(uploaded);
CREATE INDEX descriptor__uploaded_by_user ON descriptors(user_id, uploaded);

CREATE TABLE IF NOT EXISTS tags (
    tag_id VARCHAR(36),
    descriptor_id VARCHAR(36) references descriptors(descriptor_id),
    user_id VARCHAR(36) references users(user_id),
    tag VARCHAR(25),

    PRIMARY KEY (tag_id)
);

CREATE TABLE IF NOT EXISTS photos (
    photo_id VARCHAR(36),
    image_id VARCHAR(36),
    descriptor_id VARCHAR(36) references descriptors(descriptor_id),
    user_id VARCHAR(36) references users(user_id),

    PRIMARY KEY (photo_id)
);