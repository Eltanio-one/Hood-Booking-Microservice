DROP TABLE IF EXISTS users, hoods, bookings;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    passhash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    emergency_telephone INT NOT NULL,
    research_group VARCHAR(255) NOT NULL
);

CREATE TABLE hoods (
    id SERIAL PRIMARY KEY,
    hood_number INT NOT NULL,
    room VARCHAR(255) NOT NULL
);

CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    hoodnumber INT NOT NULL,
    booking_date TIMESTAMP WITH TIME ZONE
);
