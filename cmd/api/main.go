package main

import (
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"

	"github.com/CodingFervor/data-middle-platform/internal/database"
)

func main() {
	r := gin.Default()
	r.Use(CORS())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Format(time.RFC3339)})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", Login)
		auth := api.Group("/")
		auth.Use(AuthMiddleware())
		{
			// Data sources
			auth.GET("/datasources", ListDataSources)
			auth.POST("/datasources", CreateDataSource)
			auth.PUT("/datasources/:id", UpdateDataSource)
			auth.DELETE("/datasources/:id", DeleteDataSource)
			auth.POST("/datasources/:id/test", TestDataSource)

			// Metadata
			auth.GET("/metadata/tables", ListTables)
			auth.GET("/metadata/tables/:id", GetTableDetail)
			auth.GET("/metadata/tables/:id/columns", ListColumns)
			auth.POST("/metadata/sync/:datasource_id", SyncMetadata)
			auth.PUT("/metadata/tables/:id/tags", UpdateTableTags)

			// Data lineage
			auth.GET("/lineage/table/:id", GetTableLineage)
			auth.GET("/lineage/column/:id", GetColumnLineage)

			// ETL Pipelines
			auth.GET("/pipelines", ListPipelines)
			auth.POST("/pipelines", CreatePipeline)
			auth.GET("/pipelines/:id", GetPipeline)
			auth.PUT("/pipelines/:id", UpdatePipeline)
			auth.DELETE("/pipelines/:id", DeletePipeline)
			auth.POST("/pipelines/:id/run", RunPipeline)
			auth.GET("/pipelines/:id/history", PipelineHistory)
			auth.POST("/pipelines/:id/schedule", SchedulePipeline)

			// Data warehouse models
			auth.GET("/models", ListModels)
			auth.POST("/models", CreateModel)
			auth.PUT("/models/:id", UpdateModel)
			auth.POST("/models/:id/build", BuildModel)
			auth.GET("/models/:id/data", QueryModelData)

			// Data quality
			auth.GET("/quality/rules", ListQualityRules)
			auth.POST("/quality/rules", CreateQualityRule)
			auth.POST("/quality/check", RunQualityCheck)
			auth.GET("/quality/results", QualityResults)
			auth.GET("/quality/dashboard", QualityDashboard)

			// Data API (DaaS)
			auth.GET("/data-apis", ListDataAPIs)
			auth.POST("/data-apis", CreateDataAPI)
			auth.PUT("/data-apis/:id", UpdateDataAPI)
			auth.POST("/data-apis/:id/publish", PublishDataAPI)
			auth.POST("/data-apis/:id/unpublish", UnpublishDataAPI)

			// Data queries
			auth.POST("/query/sql", ExecuteSQL)
			auth.POST("/query/preview/:table_id", PreviewData)

			// Data governance
			auth.GET("/governance/classifications", ListClassifications)
			auth.POST("/governance/classifications", CreateClassification)
			auth.GET("/governance/access-logs", AccessLogs)

			// Dashboard
			auth.GET("/dashboard/overview", DashboardOverview)
			auth.GET("/dashboard/storage", StorageStats)
			auth.GET("/dashboard/usage", UsageStats)
		}
	}
	log.Println("Data Middle Platform starting on :8080")
	addr := ":" + strconv.Itoa(8080)
	srv := &http.Server{Addr: addr, Handler: r}
	go func() {
		logger.Info("server listening", "port", 8080)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("forced shutdown", "error", err)
	}
	logger.Info("server exited")
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" { c.AbortWithStatus(http.StatusNoContent); return }
		c.Next()
	}
}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"}); return }
		c.Next()
	}
}

func Login(c *gin.Context)                  { c.JSON(http.StatusOK, gin.H{"message": "login"}) }
func ListDataSources(c *gin.Context)        { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func CreateDataSource(c *gin.Context)       { c.JSON(http.StatusCreated, gin.H{"message": "datasource created"}) }
func UpdateDataSource(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{"message": "datasource updated"}) }
func DeleteDataSource(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{"message": "datasource deleted"}) }
func TestDataSource(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{"data": gin.H{"connected": true}}) }
func ListTables(c *gin.Context)             { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func GetTableDetail(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{"data": gin.H{}}) }
func ListColumns(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func SyncMetadata(c *gin.Context)           { c.JSON(http.StatusOK, gin.H{"message": "metadata synced"}) }
func UpdateTableTags(c *gin.Context)        { c.JSON(http.StatusOK, gin.H{"message": "tags updated"}) }
func GetTableLineage(c *gin.Context)        { c.JSON(http.StatusOK, gin.H{"data": gin.H{}}) }
func GetColumnLineage(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{"data": gin.H{}}) }
func ListPipelines(c *gin.Context)          { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func CreatePipeline(c *gin.Context)         { c.JSON(http.StatusCreated, gin.H{"message": "pipeline created"}) }
func GetPipeline(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{"data": gin.H{}}) }
func UpdatePipeline(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{"message": "pipeline updated"}) }
func DeletePipeline(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{"message": "pipeline deleted"}) }
func RunPipeline(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{"message": "pipeline triggered"}) }
func PipelineHistory(c *gin.Context)        { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func SchedulePipeline(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{"message": "pipeline scheduled"}) }
func ListModels(c *gin.Context)             { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func CreateModel(c *gin.Context)            { c.JSON(http.StatusCreated, gin.H{"message": "model created"}) }
func UpdateModel(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{"message": "model updated"}) }
func BuildModel(c *gin.Context)             { c.JSON(http.StatusOK, gin.H{"message": "model build triggered"}) }
func QueryModelData(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func ListQualityRules(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func CreateQualityRule(c *gin.Context)      { c.JSON(http.StatusCreated, gin.H{"message": "rule created"}) }
func RunQualityCheck(c *gin.Context)        { c.JSON(http.StatusOK, gin.H{"message": "quality check triggered"}) }
func QualityResults(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func QualityDashboard(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{"data": gin.H{}}) }
func ListDataAPIs(c *gin.Context)           { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func CreateDataAPI(c *gin.Context)          { c.JSON(http.StatusCreated, gin.H{"message": "data API created"}) }
func UpdateDataAPI(c *gin.Context)          { c.JSON(http.StatusOK, gin.H{"message": "data API updated"}) }
func PublishDataAPI(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{"message": "data API published"}) }
func UnpublishDataAPI(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{"message": "data API unpublished"}) }
func ExecuteSQL(c *gin.Context)             { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func PreviewData(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func ListClassifications(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func CreateClassification(c *gin.Context)   { c.JSON(http.StatusCreated, gin.H{"message": "classification created"}) }
func AccessLogs(c *gin.Context)             { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
func DashboardOverview(c *gin.Context)      { c.JSON(http.StatusOK, gin.H{"data": gin.H{}}) }
func StorageStats(c *gin.Context)           { c.JSON(http.StatusOK, gin.H{"data": gin.H{}}) }
func UsageStats(c *gin.Context)             { c.JSON(http.StatusOK, gin.H{"data": []gin.H{}}) }
