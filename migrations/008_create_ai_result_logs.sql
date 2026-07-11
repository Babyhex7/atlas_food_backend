CREATE TABLE ai_result_logs (
    id CHAR(36) PRIMARY KEY,
    submission_id CHAR(36) NOT NULL UNIQUE,
    input_payload JSON NOT NULL,
    raw_response JSON NOT NULL,
    overall_status ENUM('good','less','excess') NOT NULL,
    model_used VARCHAR(50) NOT NULL DEFAULT 'llama3-8b-8192',
    token_used INT NULL,
    latency_ms INT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_submission (submission_id),
    INDEX idx_status (overall_status),
    CONSTRAINT fk_ai_result_logs_submission FOREIGN KEY (submission_id) REFERENCES survey_submissions(id) ON DELETE CASCADE
);