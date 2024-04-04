create table if not exists person (
  id serial primary key,
  name text not null default '',
  birthday date not null default current_date,
  description text not null default '',
  location text not null default '', 
  photo text not null default '',
  email text not null default '' UNIQUE,
  password text not null default '',
  created_at timestamp not null default current_timestamp,
  premium boolean not null default false,
  likes_left integer not null default '0',
  session_id text default '',
  gender text
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

create table if not exists image (
	url text primary key
);

CREATE TABLE IF NOT EXISTS person_image (
    person_id serial NOT NULL REFERENCES person(id) ON DELETE CASCADE ON UPDATE CASCADE,
    image_url text NOT NULL REFERENCES image(url) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (person_id, image_url),
    UNIQUE(image_url)
);
