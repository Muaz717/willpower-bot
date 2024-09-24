CREATE TABLE users
(   
    id SERIAL NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    ch_id INT NOT NULL UNIQUE
);

CREATE TABLE workouts
(
    id SERIAL NOT NULL UNIQUE,
    workout_date VARCHAR(255),
    weight REAL
);

CREATE TABLE users_workouts
(
    id SERIAL NOT NULL UNIQUE,
    chat_id INT REFERENCES users (ch_id) ON DELETE CASCADE NOT NULL,
    workout_id INT REFERENCES workouts (id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE pull_ups
(
    id SERIAL NOT NULL UNIQUE,
    quantity INT,
    date VARCHAR(255)
);

CREATE TABLE users_pullups
(
    id SERIAL NOT NULL UNIQUE,
    chat_id INT REFERENCES users (ch_id) ON DELETE CASCADE NOT NULL,
    pull_up_id INT REFERENCES pull_ups (id) ON DELETE CASCADE NOT NULL
);