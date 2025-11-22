CREATE TABLE IF NOT EXISTS pull_requests(
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    author_id VARCHAR(255) NOT NULL,
    merged_at TIMESTAMP,
    status VARCHAR(255) NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS reviewers(
    pull_request_id VARCHAR(255) NOT NULL,
    reviewer_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (pull_request_id, reviewer_id),
    CONSTRAINT fk_pull_request FOREIGN KEY (pull_request_id) REFERENCES pull_requests(id) ON DELETE CASCADE,
    CONSTRAINT fk_reviewer FOREIGN KEY (reviewer_id) REFERENCES users(id) ON DELETE CASCADE
)