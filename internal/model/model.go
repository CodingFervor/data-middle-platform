package model

import "time"

type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"`
	Name      string    `json:"name" db:"name"`
	Role      string    `json:"role" db:"role"` // admin, developer, analyst, viewer
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type DataSource struct {
	ID         int64     `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Type       string    `json:"type" db:"type"` // mysql, postgresql, mongodb, api, csv, kafka
	Host       string    `json:"host" db:"host"`
	Port       int       `json:"port" db:"port"`
	Database   string    `json:"database" db:"database"`
	Username   string    `json:"username" db:"username"`
	Password   string    `json:"-" db:"password"`
	Params     string    `json:"params" db:"params"` // JSON extra config
	Status     string    `json:"status" db:"status"` // active, inactive, error
	LastSyncAt *time.Time `json:"last_sync_at" db:"last_sync_at"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type MetaTable struct {
	ID           int64     `json:"id" db:"id"`
	DataSourceID int64     `json:"data_source_id" db:"data_source_id"`
	Schema       string    `json:"schema" db:"schema"`
	Name         string    `json:"name" db:"name"`
	Comment      string    `json:"comment" db:"comment"`
	RowCount     int64     `json:"row_count" db:"row_count"`
	SizeBytes    int64     `json:"size_bytes" db:"size_bytes"`
	Tags         string    `json:"tags" db:"tags"` // JSON array
	Classification string `json:"classification" db:"classification"`
	SyncedAt     time.Time `json:"synced_at" db:"synced_at"`
}

type MetaColumn struct {
	ID          int64  `json:"id" db:"id"`
	TableID     int64  `json:"table_id" db:"table_id"`
	Name        string `json:"name" db:"name"`
	DataType    string `json:"data_type" db:"data_type"`
	Comment     string `json:"comment" db:"comment"`
	Nullable    bool   `json:"nullable" db:"nullable"`
	PrimaryKey  bool   `json:"primary_key" db:"primary_key"`
	Sensitive   bool   `json:"sensitive" db:"sensitive"`
}

type Pipeline struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Type        string    `json:"type" db:"type"` // batch, streaming
	Config      string    `json:"config" db:"config"` // JSON pipeline definition
	Schedule    string    `json:"schedule" db:"schedule"` // cron expression
	Status      string    `json:"status" db:"status"` // active, paused, error
	LastRunAt   *time.Time `json:"last_run_at" db:"last_run_at"`
	CreatedBy   int64     `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type PipelineRun struct {
	ID         int64     `json:"id" db:"id"`
	PipelineID int64     `json:"pipeline_id" db:"pipeline_id"`
	Status     string    `json:"status" db:"status"` // running, success, failed
	StartTime  time.Time `json:"start_time" db:"start_time"`
	EndTime    *time.Time `json:"end_time" db:"end_time"`
	Duration   int       `json:"duration" db:"duration"` // seconds
	RowsRead   int64     `json:"rows_read" db:"rows_read"`
	RowsWritten int64    `json:"rows_written" db:"rows_written"`
	Error      string    `json:"error" db:"error"`
}

type DWModel struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Layer       string    `json:"layer" db:"layer"` // ods, dwd, dws, ads
	TableName   string    `json:"table_name" db:"table_name"`
	SQL         string    `json:"sql" db:"sql"`
	Columns     string    `json:"columns" db:"columns"` // JSON
	Status      string    `json:"status" db:"status"` // draft, active, deprecated
	LastBuiltAt *time.Time `json:"last_built_at" db:"last_built_at"`
	CreatedBy   int64     `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type QualityRule struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	TableID     int64     `json:"table_id" db:"table_id"`
	ColumnID    *int64    `json:"column_id" db:"column_id"`
	Type        string    `json:"type" db:"type"` // null_check, unique, range, regex, custom
	Config      string    `json:"config" db:"config"` // JSON
	Severity    string    `json:"severity" db:"severity"` // low, medium, high
	Enabled     bool      `json:"enabled" db:"enabled"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type QualityResult struct {
	ID         int64     `json:"id" db:"id"`
	RuleID     int64     `json:"rule_id" db:"rule_id"`
	TableID    int64     `json:"table_id" db:"table_id"`
	PassRate   float64   `json:"pass_rate" db:"pass_rate"` // 0-100
	TotalRows  int64     `json:"total_rows" db:"total_rows"`
	FailedRows int64     `json:"failed_rows" db:"failed_rows"`
	Status     string    `json:"status" db:"status"` // pass, fail
	CheckedAt  time.Time `json:"checked_at" db:"checked_at"`
}

type DataAPI struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Path        string    `json:"path" db:"path"`
	Method      string    `json:"method" db:"method"` // GET, POST
	Description string    `json:"description" db:"description"`
	SQL         string    `json:"sql" db:"sql"`
	Params      string    `json:"params" db:"params"` // JSON
	CacheTTL    int       `json:"cache_ttl" db:"cache_ttl"` // seconds
	QPSLimit    int       `json:"qps_limit" db:"qps_limit"`
	Published   bool      `json:"published" db:"published"`
	CreatedBy   int64     `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type DataLineage struct {
	ID            int64 `json:"id" db:"id"`
	SourceType    string `json:"source_type" db:"source_type"` // table, column
	SourceID      int64 `json:"source_id" db:"source_id"`
	TargetType    string `json:"target_type" db:"target_type"`
	TargetID      int64 `json:"target_id" db:"target_id"`
	TransformDesc string `json:"transform_desc" db:"transform_desc"`
	PipelineID    *int64 `json:"pipeline_id" db:"pipeline_id"`
}

type Classification struct {
	ID     int64  `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Level  string `json:"level" db:"level"` // public, internal, confidential, secret
}

type AccessLog struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Action    string    `json:"action" db:"action"` // query, export, api_call
	Resource  string    `json:"resource" db:"resource"`
	Detail    string    `json:"detail" db:"detail"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
