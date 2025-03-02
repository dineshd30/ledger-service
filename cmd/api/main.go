package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dineshd30/ledger-service/internal/ledger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	loadConfig(getEnv())
	configureLogger()

	router := configureRoutes()
	port := getHTTPPort()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	zap.L().Info(fmt.Sprintf("ledger service started at :%s", port))
	err := server.ListenAndServe()
	if err != nil {
		zap.L().Fatal("failed to listen and serve on server", zap.Error(err), zap.String("port", port))
	}
}

// configureRoutes configures service routes
func configureRoutes() *gin.Engine {
	mode := gin.ReleaseMode
	if getEnv() != "prod" {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.Writer.WriteHeader(http.StatusOK)
	})

	uuid := ledger.NewUUIDGenerator()
	store := ledger.NewStore(uuid, initCashLedger(uuid))
	ledgerRoutes := router.Group("/ledger/:ledgerId")
	ledgerRoutes.POST("/transaction", ledger.DoTransaction(store))
	ledgerRoutes.GET("/balance", ledger.ViewBalance(store))
	ledgerRoutes.GET("/statement", ledger.ViewTransactionHistory(store))
	return router
}

// initCashLedger initialises cash ledger
func initCashLedger(uuid ledger.UUIDGenerator) map[string]*ledger.Ledger {
	ledgerId := "304629d2-ba1f-43df-a839-26ceb869645a"
	cashLedger := ledger.Ledger{
		ID:   ledgerId,
		Type: "cash",
		Transactions: []ledger.Transaction{
			{
				ID:             uuid.Generate(),
				Date:           time.Now().UTC().UnixMilli(),
				Type:           ledger.Credit,
				Description:    "Initial transaction",
				Amount:         100,
				RunningBalance: 100.00,
			},
		},
	}

	return map[string]*ledger.Ledger{
		ledgerId: &cashLedger,
	}
}

// configureLogger configures zap logger
func configureLogger() *zap.Logger {
	logLevel := viper.GetString("logs.level")
	conf := zap.NewProductionConfig()
	level := zap.NewAtomicLevel()
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		log.Fatalf("failed to set the log level, defaulting to info level: %s\n", logLevel)
		conf.Level.SetLevel(zap.InfoLevel)
	} else {
		conf.Level = level
	}
	conf.OutputPaths = []string{"stdout"}

	logger, err := conf.Build()
	if err != nil {
		log.Fatalf("failed to build the logger: %s\n", err)
	}

	zap.ReplaceGlobals(logger)
	return logger
}

// getHTTPPort gets http port
func getHTTPPort() string {
	port := viper.GetString("http.port")
	if port == "" {
		port = "8080"
	}
	return port
}
