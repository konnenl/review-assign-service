ALTER TABLE users ADD COLUMN team_name VARCHAR(255);

UPDATE users u
SET team_name = ut.team_name
FROM users_teams ut
WHERE u.id = ut.user_id;

ALTER TABLE users
  ADD CONSTRAINT fk_team FOREIGN KEY (team_name) REFERENCES teams(name) ON DELETE CASCADE;

DROP TABLE IF EXISTS users_teams;