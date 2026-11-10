-- Initial schema for api_logs and audits
CREATE TABLE IF NOT EXISTS api_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    hash_id VARCHAR(255) NULL,
    method VARCHAR(10) NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    request_body LONGTEXT,
    response_body LONGTEXT,
    status_code INT NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    employee_id VARCHAR(50),
    execution_time_ms DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_api_logs_endpoint (endpoint),
    INDEX idx_api_logs_created_at (created_at)
);

CREATE TABLE IF NOT EXISTS audits (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    table_name VARCHAR(100) NOT NULL,
    record_id BIGINT NOT NULL,
    action ENUM('INSERT', 'UPDATE', 'DELETE') NOT NULL,
    old_data JSON NULL,
    new_data JSON NULL,
    changed_fields JSON NULL,
    created_by BIGINT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_audits_table_record (table_name, record_id),
    INDEX idx_audits_created_at (created_at)
);
