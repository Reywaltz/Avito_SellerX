
INSERT INTO users (username, created_at) values ('Reywaltz', NOW());
INSERT INTO users (username, created_at) values ('Wol4ara', NOW());
INSERT INTO users (username, created_at) values ('Hitchpock', NOW());
INSERT INTO users (username, created_at) values ('Yunus', NOW());


insert into chats (name, created_at) values ('Vagu', NOW());
insert into chats (name, created_at) values ('Avito Tech', NOW());

insert into users_chats(chat_id, user_id) values (1, 1);
insert into users_chats(chat_id, user_id) values (1, 2);
insert into users_chats(chat_id, user_id) values (2, 4);
insert into users_chats(chat_id, user_id) values (2, 3);

insert into messages(chat, author, text, created_at) values (1, 2, 'Hello to Vagu chat', NOW());
insert into messages(chat, author, text, created_at) values (2, 3, 'Hello to Avito chat', NOW());
insert into messages(chat, author, text, created_at) values (2, 3, 'Hello to Avito chat #1', NOW());
insert into messages(chat, author, text, created_at) values (2, 3, 'Hello to Avito chat #2', NOW());
