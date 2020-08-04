create table users
(
    id         char(32)                               not null
        constraint users_pk
            primary key,
    username   varchar(255)                           not null,
    created_at timestamp with time zone default now() not null
);

create unique index users_id_uindex
    on users (id);

create unique index users_username_uindex
    on users (username);

create table chats
(
    id         serial                                 not null
        constraint chats_pk
            primary key,
    name       varchar(255)                           not null,
    created_at timestamp with time zone default now() not null
);

create unique index chats_id_uindex
    on chats (id);

create unique index chats_name_uindex
    on chats (name);

create table messages
(
    id         serial                                 not null
        constraint messages_pk
            primary key,
    chat       integer                                not null
        constraint messages_chats_id_fk
            references chats,
    author     char(32)                               not null
        constraint messages_users_id_fk
            references users,
    text       text                                   not null,
    created_at timestamp with time zone default now() not null
);

create unique index messages_id_uindex
    on messages (id);

create index messages_author_index
    on messages (author);

create index messages_chat_index
    on messages (chat);

create table chats_users
(
    chat_id integer  not null
        constraint chats_users_chats_id_fk
            references chats,
    user_id char(32) not null
        constraint chats_users_users_id_fk
            references users,
    constraint chats_users_pk
        primary key (chat_id, user_id)
);

