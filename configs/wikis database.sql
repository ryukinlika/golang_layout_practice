create database wikis;
use wikis;

drop table if exists pages;

create table pages (
    id      int auto_increment not null,
    title   varchar(255) not null,
    body    varchar(255) not null,
    primary key (`id`)
);

insert into pages     
    (title, body)
values
    ("Hello World", "Hello world! First entry of the database!"),
    ("Golang", "Go (Golang) is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. Go is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency"),
    ("Python", "Python is an interpreted high-level general-purpose programming language. Its design philosophy emphasizes code readability with its use of significant indentation.");

