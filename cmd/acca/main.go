package main

import (
	"context"
	"database/sql"
	"flag"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gebv/acca/engine"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

var (
	VERSION         = "dev"
	pgConnF         = flag.String("pg-conn", "postgres://acca:acca@127.0.0.1:5432/acca?sslmode=disable", "PostgreSQL connection string.")
	grpcAddrsF      = flag.String("grpc-addrs", "127.0.0.1:10001", "gRPC listen address.")
	grpcReflectionF = flag.Bool("grpc-reflection", false, "Enable gRPC reflection.")
)

func main() {
	rand.Seed(time.Now().UnixNano())
	defaultLogger("INFO")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	zap.L().Info("Starting...", zap.String("version", VERSION))
	defer func() { zap.L().Info("Done.") }()

	syncLogger := developLogger(false)
	defer syncLogger()
	handleTerm(cancel)

	sqlDB := setupPostgres(*pgConnF, 0, 5, 5)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(zap.L().Sugar().Debugf))
	_, err := db.Exec("SELECT version();")
	if err != nil {
		zap.L().Panic("Failed to check version to PostgreSQL.", zap.Error(err))
	}

	lis, err := net.Listen("tcp", *grpcAddrsF)
	if err != nil {
		zap.L().Panic("Failed to listen.", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.ChainUnaryServer()),
		grpc.StreamInterceptor(middleware.ChainStreamServer()),
	)

	accountManager := engine.NewAccountManager(db)
	currID, err := accountManager.UpsertCurrency("curr1", nil)
	if err != nil {
		zap.L().Panic("Failed create currency.", zap.Error(err))
	}
	acc1ID, err := accountManager.CreateAccount(currID, "first", nil)
	if err != nil && err != engine.ErrAccountExists {
		zap.L().Panic("Failed create account.", zap.Error(err))
	}
	if acc1ID == 0 {
		acc, err := accountManager.FindAccountByKey(currID, "first")
		if err != nil {
			zap.L().Panic("Failed find account.", zap.Error(err))
		}
		acc1ID = acc.AccountID
	}
	acc2ID, err := accountManager.CreateAccount(currID, "second", nil)
	if err != nil && err != engine.ErrAccountExists {
		zap.L().Panic("Failed create account.", zap.Error(err))
	}
	if acc2ID == 0 {
		acc, err := accountManager.FindAccountByKey(currID, "second")
		if err != nil {
			zap.L().Panic("Failed find account.", zap.Error(err))
		}
		acc2ID = acc.AccountID
	}

	eng := engine.NewSimpleService(db)
	if _, err := eng.InternalTransfer(acc1ID, acc2ID, 100); err != nil {
		zap.L().Panic("Failed to internal transfer.", zap.Error(err))
	}

	// graceful stop
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		go func() {
			<-ctx.Done()
			s.Stop()

		}()
		s.GracefulStop()
	}()

	// TODO: Registry servers

	if *grpcReflectionF {
		// for debug via grpcurl
		reflection.Register(s)
	}

	var wg sync.WaitGroup

	zap.L().Info("gRPC server listen address.", zap.String("address", lis.Addr().String()))
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.Serve(lis); err != nil {
			zap.L().Panic("Failed to serve.", zap.Error(err))
		}
	}()
	wg.Wait()

	// - внутренний grpc АПИ
	// - хандлер для платежек

	/*
		Входящая операция падает в общую очередь
		Колторая обрабатывается в горутине
		Все состояния персистятся в PG
		В случае падения процесса очередь воссоздается из БД (то есть сохранять состояния команд?)
	*/
}

// Configure configure zap logger.
//
// Available values of level:
// - DEBUG
// - INFO
// - WARN
// - ERROR
// - DPANIC
// - PANIC
// - FATAL
func defaultLogger(levelSet string) {
	level := zapcore.InfoLevel
	if err := level.Set(levelSet); err != nil {
		panic(err)
	}
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(level)
	l, err := config.Build(zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(l)
	zap.RedirectStdLog(l.Named("stdlog"))
}

func developLogger(debug bool) func() error {
	zap.L().Sync()

	var config zap.Config
	config = zap.NewDevelopmentConfig()
	config.Development = true
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	if debug {
		config.Level.SetLevel(zap.DebugLevel)
	} else {
		config.Level.SetLevel(zap.InfoLevel)
	}

	l, err := config.Build(
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(l)
	zap.RedirectStdLog(l.Named("stdlog"))

	return l.Sync
}

func productionLogger(debug bool) func() error {
	zap.L().Sync()

	var config zap.Config
	config = zap.NewProductionConfig()
	config.Development = false
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	if debug {
		config.Level.SetLevel(zap.DebugLevel)
	} else {
		config.Level.SetLevel(zap.InfoLevel)
	}

	l, err := config.Build(
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(l)
	zap.RedirectStdLog(l.Named("stdlog"))

	return l.Sync
}

func handleTerm(cancel context.CancelFunc) {
	// handle termination signals: first one gracefully, force exit on the second one
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, unix.SIGTERM, unix.SIGINT)
	go func() {
		s := <-signals
		zap.L().Warn("Shutting down.", zap.String("signal", unix.SignalName(s.(unix.Signal))))
		cancel()

		s = <-signals
		zap.L().Panic("Exiting!", zap.String("signal", unix.SignalName(s.(unix.Signal))))
	}()
}

func setupPostgres(conn string, maxLifetime time.Duration, maxOpen, maxIdle int) *sql.DB {
	sqlDB, err := sql.Open("postgres", conn)
	if err != nil {
		zap.L().Panic("Failed to connect to PostgreSQL.", zap.Error(err))
	}
	sqlDB.SetConnMaxLifetime(maxLifetime)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	if err = sqlDB.Ping(); err != nil {
		zap.L().Panic("Failed to connect ping PostgreSQL.", zap.Error(err))
	}
	zap.L().Info("Postgres - Connected!")

	return sqlDB
}
