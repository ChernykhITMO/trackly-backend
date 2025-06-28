CREATE TABLE plans
(
    id       SERIAL PRIMARY KEY,
    habit_id INT,
    plan_unit TEXT,
    goal INT,
    start_time TIMESTAMP,
    close_time TIMESTAMP
);
CREATE TABLE habits
(
    id SERIAL PRIMARY KEY,
    habit_name TEXT NOT NULL,
    description TEXT,
    user_id int,
    constraint us_id foreign key (user_id) references users(id),
    start_date TIMESTAMP,
    notifications BOOLEAN

);
ALTER TABLE plans
    ADD CONSTRAINT fk_plans_habit_id
        FOREIGN KEY (habit_id)
            REFERENCES habits(id);

CREATE TABLE habit_scores(
    id SERIAL PRIMARY KEY ,
    habit_id int,
    constraint hab_id foreign key (habit_id) references habits(id),
    plan_id int,
    constraint pl_id foreign key (plan_id) references plans(id),
    date_time TIMESTAMP,
    value int
);