-- +migrate Up
CREATE TABLE workout_sessions (
    session_id INT AUTO_INCREMENT PRIMARY KEY,
    training_date DATE NOT NULL,
    user_id INT
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE exercises (
    exercise_id INT AUTO_INCREMENT PRIMARY KEY,
    session_id INT NOT NULL,
    exercise_name VARCHAR(255) NOT NULL,
    FOREIGN KEY (session_id) REFERENCES workout_sessions(session_id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE sets (
    set_id INT AUTO_INCREMENT PRIMARY KEY,
    exercise_id INT NOT NULL,
    set_number INT NOT NULL,
    weight DECIMAL(5,2) NOT NULL,
    reps INT NOT NULL,
    FOREIGN KEY (exercise_id) REFERENCES exercises(exercise_id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
