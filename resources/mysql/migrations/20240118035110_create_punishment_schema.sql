-- migrate:up
CREATE TABLE IF NOT EXISTS punishment (
    Reg_No VARCHAR(21) NOT NULL,
    Reg_Date datetime DEFAULT NULL,
    ID VARCHAR(20) DEFAULT NULL,
    Join_Date datetime DEFAULT NULL,
    Resign_Date datetime DEFAULT NULL,
    Directorate VARCHAR(30) DEFAULT NULL,
    Division VARCHAR(50) DEFAULT NULL,
    Department VARCHAR(50) DEFAULT NULL,
    Job_Title VARCHAR(50) DEFAULT NULL,
    JG VARCHAR(5) DEFAULT NULL,
    Punishment_Type VARCHAR(3) DEFAULT NULL,
    Start DATETIME DEFAULT NULL,
    End DATETIME DEFAULT NULL,
    Description VARCHAR(750) DEFAULT NULL,
    Remarks VARCHAR(850) NOT NULL,
    Status VARCHAR(11) NOT NULL,
    Classification VARCHAR(255) DEFAULT NULL
)

-- migrate:down

