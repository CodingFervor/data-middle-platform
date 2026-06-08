CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY, username VARCHAR(50) UNIQUE NOT NULL, password VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL, role VARCHAR(20) DEFAULT 'analyst' CHECK (role IN ('admin','developer','analyst','viewer')),
    status VARCHAR(20) DEFAULT 'active', created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE data_sources (
    id BIGSERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('mysql','postgresql','mongodb','api','csv','kafka','elasticsearch')),
    host VARCHAR(200), port INT, database_name VARCHAR(100),
    username VARCHAR(100), password VARCHAR(255),
    params JSONB DEFAULT '{}', status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active','inactive','error')),
    last_sync_at TIMESTAMP, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE meta_tables (
    id BIGSERIAL PRIMARY KEY, data_source_id BIGINT NOT NULL REFERENCES data_sources(id),
    schema_name VARCHAR(100), name VARCHAR(200) NOT NULL, comment TEXT,
    row_count BIGINT DEFAULT 0, size_bytes BIGINT DEFAULT 0,
    tags JSONB DEFAULT '[]', classification VARCHAR(50),
    synced_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(data_source_id, schema_name, name)
);

CREATE TABLE meta_columns (
    id BIGSERIAL PRIMARY KEY, table_id BIGINT NOT NULL REFERENCES meta_tables(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL, data_type VARCHAR(50),
    comment TEXT, nullable BOOLEAN DEFAULT TRUE, primary_key BOOLEAN DEFAULT FALSE,
    sensitive BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_columns_table ON meta_columns(table_id);

CREATE TABLE pipelines (
    id BIGSERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL, description TEXT,
    type VARCHAR(20) NOT NULL CHECK (type IN ('batch','streaming')),
    config JSONB DEFAULT '{}', schedule VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active','paused','error')),
    last_run_at TIMESTAMP, created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE pipeline_runs (
    id BIGSERIAL PRIMARY KEY, pipeline_id BIGINT NOT NULL REFERENCES pipelines(id),
    status VARCHAR(20) DEFAULT 'running' CHECK (status IN ('running','success','failed','killed')),
    start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP, end_time TIMESTAMP,
    duration_seconds INT DEFAULT 0, rows_read BIGINT DEFAULT 0, rows_written BIGINT DEFAULT 0,
    error TEXT
);

CREATE INDEX idx_runs_pipeline ON pipeline_runs(pipeline_id);

CREATE TABLE dw_models (
    id BIGSERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL, description TEXT,
    layer VARCHAR(10) NOT NULL CHECK (layer IN ('ods','dwd','dws','ads')),
    table_name VARCHAR(200) NOT NULL, sql_text TEXT,
    columns JSONB DEFAULT '[]',
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft','active','deprecated')),
    last_built_at TIMESTAMP, created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE quality_rules (
    id BIGSERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL,
    table_id BIGINT NOT NULL REFERENCES meta_tables(id), column_id BIGINT REFERENCES meta_columns(id),
    type VARCHAR(20) NOT NULL CHECK (type IN ('null_check','unique','range','regex','custom_sql')),
    config JSONB DEFAULT '{}', severity VARCHAR(10) DEFAULT 'medium' CHECK (severity IN ('low','medium','high')),
    enabled BOOLEAN DEFAULT TRUE, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE quality_results (
    id BIGSERIAL PRIMARY KEY, rule_id BIGINT NOT NULL REFERENCES quality_rules(id),
    table_id BIGINT NOT NULL REFERENCES meta_tables(id),
    pass_rate DECIMAL(5,2) DEFAULT 100, total_rows BIGINT DEFAULT 0, failed_rows BIGINT DEFAULT 0,
    status VARCHAR(10) DEFAULT 'pass' CHECK (status IN ('pass','fail')),
    checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_quality_table ON quality_results(table_id);

CREATE TABLE data_apis (
    id BIGSERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL,
    path VARCHAR(200) UNIQUE NOT NULL, method VARCHAR(10) DEFAULT 'GET',
    description TEXT, sql_text TEXT NOT NULL, params JSONB DEFAULT '[]',
    cache_ttl INT DEFAULT 0, qps_limit INT DEFAULT 100,
    published BOOLEAN DEFAULT FALSE, created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE data_lineage (
    id BIGSERIAL PRIMARY KEY,
    source_type VARCHAR(10) NOT NULL CHECK (source_type IN ('table','column')),
    source_id BIGINT NOT NULL,
    target_type VARCHAR(10) NOT NULL CHECK (target_type IN ('table','column')),
    target_id BIGINT NOT NULL,
    transform_desc TEXT, pipeline_id BIGINT REFERENCES pipelines(id)
);

CREATE TABLE classifications (
    id BIGSERIAL PRIMARY KEY, name VARCHAR(50) NOT NULL,
    level VARCHAR(20) NOT NULL CHECK (level IN ('public','internal','confidential','secret'))
);

INSERT INTO classifications (name, level) VALUES
('Public Data', 'public'), ('Internal Data', 'internal'),
('Confidential Data', 'confidential'), ('Secret Data', 'secret');

CREATE TABLE access_logs (
    id BIGSERIAL PRIMARY KEY, user_id BIGINT NOT NULL REFERENCES users(id),
    action VARCHAR(20) NOT NULL CHECK (action IN ('query','export','api_call','manage')),
    resource VARCHAR(200), detail TEXT, ip_address VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_access_user ON access_logs(user_id);
CREATE INDEX idx_access_time ON access_logs(created_at);

INSERT INTO users (username, password, name, role) VALUES
('admin', '$2a$10$dummyhash', 'Admin', 'admin');
