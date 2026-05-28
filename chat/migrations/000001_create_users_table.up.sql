CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username varchar(50) NOT NULL,
  phone_number varchar(20) NOT NULL UNIQUE
);
