
CREATE TABLE IF NOT EXISTS csat (
    id INT generated always as IDENTITY PRIMARY KEY,
    q1 INT NOT NULL
    check(q1 >= 0 and q1 <= 5)
);