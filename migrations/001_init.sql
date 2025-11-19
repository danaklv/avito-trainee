CREATE TABLE IF NOT EXISTS teams {
    team_id SERIAL NOT NULL PRIMARY KEY,
    teamname VARCHAR(50) NOT NULL UNIQUE,
};

CREATE TABLE IF NOT EXISTS users {
    user_id SERIAL NOT NULL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    team_id BIGINT NOT NULL REFERENCES teams(team_id) ON DELETE CASCADE
};


CREATE TABLE IF NOT EXISTS pull_requests {
    pull_request_id SERIAL NOT NULL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    author BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    status VARCHAR(6) NOT NULL CHECK (status IN ('OPEN', 'MERGED'))
}

CREATE TABLE IF NOT EXISTS reviewers {
    user_id BIGINT REFERENCES users(user_id) ON DELETE CASCADE,
    pull_request_id BIGINT REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE
    PRIMARY KEY (pull_request_id, user_id)
}