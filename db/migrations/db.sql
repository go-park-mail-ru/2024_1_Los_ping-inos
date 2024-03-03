create table if not exists person (
  id serial primary key,
  name text not null default '',
  birthday date not null default current_date,
  description text not null default '',
  location text not null default '', 
  photo text not null default '',
  email text not null default '',
  password text not null default '',
  gender text not null default 'attach helicopter',
  created_at timestamp not null default current_timestamp,
  premium boolean not null default false,
  likes_left integer not null default '0',
  session_id text default ''
);

create table if not exists interest (
  id serial primary key,
  name text not null default ''
);

create table if not exists person_interest (
  person_id serial not null references person(id)
    on delete cascade
    on update cascade,
  interest_id serial not null references interest(id)
    on delete cascade
    on update cascade,
    
  primary key(person_id, interest_id)
);

create table if not exists "like" (
  person_one_id serial not null references person(id)
    on delete cascade
    on update cascade,
  person_two_id serial not null references person(id)
    on delete cascade
    on update cascade,
  
  primary key(person_one_id, person_two_id)
);

create table if not exists dislike (
  person_one_id serial not null references person(id)
    on delete cascade
    on update cascade,
  person_two_id serial not null references person(id)
    on delete cascade
    on update cascade,
  
  primary key(person_one_id, person_two_id)
);
