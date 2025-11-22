CREATE TABLE users_teams(
    user_id VARCHAR(255) NOT NULL,
    team_name VARCHAR(255) NOT NULL,
    PRIMARY KEY (user_id, team_name),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_team FOREIGN KEY (team_name) REFERENCES teams(name) ON DELETE CASCADE
);

ALTER TABLE users DROP COLUMN IF EXISTS team_name;