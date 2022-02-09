CREATE TABLE IF NOT EXISTS shortened(
   id serial PRIMARY KEY,
   user_id VARCHAR (128) NOT NULL,
   short_url VARCHAR (128) UNIQUE NOT NULL,
   original_url VARCHAR (1024) NOT NULL
);