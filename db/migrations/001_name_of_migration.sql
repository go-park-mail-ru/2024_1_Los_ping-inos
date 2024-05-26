-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS person (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  birthday DATE NOT NULL DEFAULT CURRENT_DATE,
  description TEXT NOT NULL DEFAULT '',
  location TEXT NOT NULL DEFAULT '',
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  premium BOOLEAN NOT NULL DEFAULT FALSE,
  premium_expires_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  likes_left INTEGER NOT NULL DEFAULT '10',
  gender text not null,
  check(gender = 'male' or gender = 'female')
);

create table if not exists interest (
  id int generated always as IDENTITY PRIMARY KEY,
  name text not null,
  unique(name)
);

create table if not exists person_interest (
  person_id int not null references person(id)
    on delete cascade
    on update cascade,
  interest_id int not null references interest(id)
    on delete cascade
    on update cascade,
    
  primary key(person_id, interest_id)
);

create table if not exists "like" (
  person_one_id int not null references person(id)
    on delete cascade
    on update cascade,
  person_two_id int not null references person(id)
    on delete cascade
    on update cascade,
  
  primary key(person_one_id, person_two_id)
);

create table if not exists dislike (
  person_one_id int not null references person(id)
    on delete cascade
    on update cascade,
  person_two_id int not null references person(id)
    on delete cascade
    on update cascade,
  
  primary key(person_one_id, person_two_id)
);

CREATE TABLE IF NOT EXISTS person_image (
    id int generated always as IDENTITY PRIMARY KEY,
    person_id int not null REFERENCES person(id) 
	on delete cascade
	on update cascade,
    image_url text not null,
    cell_number int not null
    CONSTRAINT person_cell UNIQUE (person_id, cell_number)
);

CREATE TABLE IF NOT EXISTS person_premium (
	id int generated always as IDENTITY PRIMARY KEY,
	person_id int not null references person(id)
    		on delete cascade
    		on update cascade,
	date_started date not null,
	date_ended date
);

create table if not exists message (
   id SERIAL PRIMARY KEY,
   data TEXT NOT NULL DEFAULT '',
   sender_id int not null references person(id)
       on delete cascade
       on update cascade,
   receiver_id int not null references person(id)
       on delete cascade
       on update cascade,
   sent_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

create table if not exists person_claim (
   id SERIAL PRIMARY KEY,
   type int not null references claim(id),
   sender_id int not null references person(id)
       on delete cascade
       on update cascade,
   receiver_id int not null references person(id)
       on delete cascade
       on update cascade
);

create table if not exists claim (
   id SERIAL PRIMARY KEY,
   title TEXT NOT NULL
);
---- create above / drop below ----
-- drop table person;
--
-- drop table interest;
--
-- drop table person_interest;
--
-- drop table "like";
--
-- drop table dislike;
--
-- drop table person_image;
--
-- drop table person_premium;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
