-- Write your migrate up statements here
CREATE TYPE genders AS ENUM (
	'male',
	'female');

CREATE TABLE IF NOT EXISTS person (
  id INT generated always as IDENTITY PRIMARY KEY,
  name TEXT NOT NULL,
  birthday DATE NOT NULL DEFAULT CURRENT_DATE,
  description TEXT NOT NULL,
  location TEXT NOT NULL, 
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  premium BOOLEAN NOT NULL DEFAULT FALSE,
  likes_left INTEGER NOT NULL DEFAULT '10',
  gender genders NOT NULL
);

create table if not exists interest (
  id int generated always as IDENTITY PRIMARY KEY,
  name text not null
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
);

CREATE TABLE IF NOT EXISTS person_premium (
	id int generated always as IDENTITY PRIMARY KEY,
	person_id int not null references person(id)
    		on delete cascade
    		on update cascade,
	date_started date not null,
	date_ended date
);
---- create above / drop below ----
drop table person;

drop table interest;

drop table person_interest;

drop table "like";

drop table dislike;

drop table person_image;

drop table person_premium;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
