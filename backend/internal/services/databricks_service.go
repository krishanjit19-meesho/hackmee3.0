package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"meesho-clone/internal/models"

	"cloud.google.com/go/storage"
	_ "github.com/databricks/databricks-sql-go"
	"google.golang.org/api/option"
)

// DatabricksService handles connections to Databricks SQL
type DatabricksService struct {
	db            *sql.DB
	storageClient *storage.Client
	bucketName    string
}

// DatabricksConfig holds Databricks connection configuration
type DatabricksConfig struct {
	ServerHostname string
	HTTPPath       string
	AccessToken    string
	Catalog        string
	Schema         string
}

// NewDatabricksService creates a new Databricks service
func NewDatabricksService() (*DatabricksService, error) {
	// Get Databricks configuration from environment variables
	config := &DatabricksConfig{
		ServerHostname: getEnv("DATABRICKS_HOST", ""),
		HTTPPath:       getEnv("DATABRICKS_HTTP_PATH", ""),
		AccessToken:    getEnv("DATABRICKS_TOKEN", ""),
		Catalog:        getEnv("DATABRICKS_CATALOG", "hive_metastore"),
		Schema:         getEnv("DATABRICKS_SCHEMA", "default"),
	}

	// Validate required configuration
	if config.ServerHostname == "" || config.HTTPPath == "" || config.AccessToken == "" {
		return nil, fmt.Errorf("missing required Databricks configuration. Please set DATABRICKS_HOST, DATABRICKS_HTTP_PATH, and DATABRICKS_TOKEN")
	}

	// Create DSN for Databricks SQL
	dsn := fmt.Sprintf("token:%s@%s:443/%s?catalog=%s&schema=%s&timeout=1000&maxRows=10000",
		config.AccessToken,
		config.ServerHostname,
		config.HTTPPath,
		config.Catalog,
		config.Schema,
	)

	// Connect to Databricks SQL
	db, err := sql.Open("databricks", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Databricks SQL: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping Databricks SQL: %w", err)
	}

	// Initialize Google Cloud Storage client
	storageClient, err := storage.NewClient(context.Background(),
		option.WithCredentialsFile(getEnv("GOOGLE_APPLICATION_CREDENTIALS", "")))
	if err != nil {
		log.Printf("Warning: Failed to initialize Google Cloud Storage client: %v", err)
		log.Printf("GCS features will be disabled")
		storageClient = nil
	}

	return &DatabricksService{
		db:            db,
		storageClient: storageClient,
		bucketName:    getEnv("GCS_BUCKET_NAME", "gcs-dsci-adv-ads-dev-prd"),
	}, nil
}

// GetCatalogDataByIDs executes the actual SQL query from the Python script
func (s *DatabricksService) GetCatalogDataByIDs(catalogIDs []string) ([]models.CatalogProduct, error) {
	if len(catalogIDs) == 0 {
		return nil, fmt.Errorf("no catalog IDs provided")
	}

	// Build the IN clause for the query
	placeholders := ""
	args := make([]interface{}, len(catalogIDs))
	for i, id := range catalogIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args[i] = id
	}

	// This is the actual SQL query from the Python script
	// Querying the same tables as the original Databricks notebook
	query := `
		SELECT 
			c.catalog_id,
			c.product_id,
			c.image_url,
			c.category,
			c.sub_category,
			c.price
		FROM catalog_info c
		WHERE c.catalog_id IN (` + placeholders + `)
		ORDER BY c.catalog_id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var products []models.CatalogProduct
	for rows.Next() {
		var product models.CatalogProduct
		err := rows.Scan(
			&product.CatalogID,
			&product.ProductID,
			&product.ImageURL,
			&product.Category,
			&product.SubCategory,
			&product.Price,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		product.Title = fmt.Sprintf("%s - %s", product.Category, product.SubCategory)
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return products, nil
}

// ExecuteComplexQuery executes the full complex query from the Python script
// This queries the actual Databricks tables: gold.product_info, gold.order_master_bi, etc.
func (s *DatabricksService) ExecuteComplexQuery(catalogIDs []string) ([]models.CatalogProduct, error) {
	if len(catalogIDs) == 0 {
		return nil, fmt.Errorf("no catalog IDs provided")
	}

	// Build the IN clause for the query
	placeholders := ""
	args := make([]interface{}, len(catalogIDs))
	for i, id := range catalogIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args[i] = id
	}

	// This is the complex query from the Python script
	// It queries the actual Databricks tables with the same logic
	query := `
		WITH top_products AS (
			SELECT 
				catalog_id, 
				product_id, 
				orders,
				ROW_NUMBER() OVER (PARTITION BY catalog_id ORDER BY orders DESC, product_id) as rnk
			FROM (
				SELECT 
					a.catalog_id, 
					a.product_id, 
					COALESCE(b.orders, 0) as orders
				FROM gold.product_info a
				LEFT JOIN (
					SELECT 
						catalog_id, 
						product_id, 
						COUNT(DISTINCT order_detail_id) as orders
					FROM gold.order_master_bi 
					WHERE test = 0 AND verified = 1
					GROUP BY catalog_id, product_id
				) b ON a.catalog_id = b.catalog_id AND a.product_id = b.product_id
			) ranked_products
		),
		product_images AS (
			SELECT 
				product_id,
				CONCAT('https://images.meesho.com', image_url) as image_url
			FROM silver.supply__products
		),
		taxonomy AS (
			SELECT DISTINCT 
				catalog_id, 
				scat, 
				sscat 
			FROM gold.product_info
		)
		SELECT 
			tp.catalog_id,
			tp.product_id,
			pi.image_url,
			tax.scat as category,
			tax.sscat as sub_category,
			COALESCE(ca.price_sscat_decile, '₹999') as price
		FROM top_products tp
		JOIN product_images pi ON tp.product_id = pi.product_id
		JOIN taxonomy tax ON tp.catalog_id = tax.catalog_id
		LEFT JOIN catalog__attributes_agg ca ON tp.catalog_id = ca.catalog_id
		WHERE tp.rnk = 1 
		AND tp.catalog_id IN (` + placeholders + `)
		ORDER BY tp.catalog_id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		// If the complex query fails, fall back to simple query
		log.Printf("Complex query failed, falling back to simple query: %v", err)
		return s.GetCatalogDataByIDs(catalogIDs)
	}
	defer rows.Close()

	var products []models.CatalogProduct
	for rows.Next() {
		var product models.CatalogProduct
		err := rows.Scan(
			&product.CatalogID,
			&product.ProductID,
			&product.ImageURL,
			&product.Category,
			&product.SubCategory,
			&product.Price,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		product.Title = fmt.Sprintf("%s - %s", product.Category, product.SubCategory)
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return products, nil
}

// GetDataFromGCS retrieves data from Google Cloud Storage (like the Python script)
func (s *DatabricksService) GetDataFromGCS(bucketPath string) ([]byte, error) {
	if s.storageClient == nil {
		return nil, fmt.Errorf("Google Cloud Storage client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bucket := s.storageClient.Bucket(s.bucketName)
	obj := bucket.Object(bucketPath)

	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read from GCS: %w", err)
	}
	defer reader.Close()

	// Read all data
	data := make([]byte, 0)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			data = append(data, buffer[:n]...)
		}
		if err != nil {
			break
		}
	}

	return data, nil
}

// TestConnection tests the Databricks connection
func (s *DatabricksService) TestConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test basic query
	rows, err := s.db.QueryContext(ctx, "SELECT 1 as test")
	if err != nil {
		return fmt.Errorf("failed to execute test query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var result int
		if err := rows.Scan(&result); err != nil {
			return fmt.Errorf("failed to scan test result: %w", err)
		}
		log.Printf("✅ Databricks connection test successful: %d", result)
		return nil
	}

	return fmt.Errorf("no results from test query")
}

// Close closes the database connection
func (s *DatabricksService) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
