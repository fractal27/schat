
-- messages table


CREATE TABLE IF NOT EXISTS messages(
	nick VARCHAR(50),
	contents VARCHAR(1024),
	time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



