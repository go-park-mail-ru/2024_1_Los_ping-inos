CREATE TABLE IF NOT EXISTS questions (
    id INT generated always as IDENTITY PRIMARY KEY,
    tittle TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS csat (
    id INT generated always as IDENTITY PRIMARY KEY,
    q1 INT NOT NULL,
    tittle_id INT NOT NULL REFERENCES questions(id)
    check(q1 >= 0 and q1 <= 5)
);