create table urls
(
    id         integer
        constraint urls_pk
            primary key autoincrement,
    code text    default '',
    url       text    default '',
    shortUrl       text    default '',
    domain       text    default '',
    counter    integer default 0
);

create unique index urls_id_uindex
    on urls (id);