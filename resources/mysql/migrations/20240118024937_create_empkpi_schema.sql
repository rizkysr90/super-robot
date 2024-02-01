-- migrate:up
CREATE TABLE IF NOT EXISTS infolmpperformance (
    NIK VARCHAR(20) NOT NULL,
    Name VARCHAR(50) DEFAULT NULL,
    Region VARCHAR(50) DEFAULT NULL,
    JobTitle VARCHAR(50) DEFAULT NULL,
    WorkLocation VARCHAR(255) DEFAULT NULL,
    Period VARCHAR(6) NOT NULL,
    ScoreMonthly decimal(17,2) DEFAULT NULL
);


-- migrate:down

