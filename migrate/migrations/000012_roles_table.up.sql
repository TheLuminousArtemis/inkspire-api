CREATE TABLE if not exists roles (
    id bigserial primary key,
    name varchar(255) not null unique,
    level int not null default 0,
    description text
);

INSERT INTO roles (name, level, description) VALUES
('user', 1, 'User role with limited access'), 
('moderator',2, 'Moderator can update other users posts'),
('admin',3, 'Administrator role with full access');