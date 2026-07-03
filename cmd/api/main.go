package main

import (
	"fmt"
	"gestao/config"
	"gestao/db"
	"gestao/internal/controller"
	"gestao/internal/repository"
	"gestao/internal/router"
	"gestao/internal/service"
	"os"
	"time"
)

func main() {
	app := &application{
		API_port:           fmt.Sprintf(":%s", os.Getenv("PORT")),
		API_maxHeaderBytes: config.GetInt("API_MAX_HEADER_BYTES", 1<<20), // 1 MB
		API_readTimeout:    config.GetDuration("API_READ_TIMEOUT", 10*time.Second),
		API_writeTimeout:   config.GetDuration("API_WRITE_TIMEOUT", 10*time.Second),
	}

	database := &db.BancoDados{
		ConnectionString: os.Getenv("DATABASE_URL"),
		Driver:           config.GetString("DB_DRIVER", "postgres"),
		MaxOpenConns:     config.GetInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:     config.GetInt("DB_MAX_IDLE_CONNS", 25),
		MaxIdleTime:      config.GetDuration("DB_MAX_IDLE_TIME", 15*time.Minute),
	}

	dbConn, err := database.Conectar()
	if err != nil {
		fmt.Printf("Erro ao conectar ao banco de dados: %v\n", err)
		return
	}
	defer dbConn.Close()

	/*
		fmt.Println("Sincronizando tabelas com o banco de dados...")
		if err := db.IniciarTabelas(dbConn); err != nil {
			fmt.Printf("Erro crítico ao sincronizar banco de dados: %v\n", err)
			return
		}
	*/
	repositorio := repository.NewRepository(dbConn)
	service := service.NewService(repositorio, dbConn)
	controller := controller.NewController(service)

	r := router.CarregarRotas(controller)

	if err := app.iniciarApp(r); err != nil {
		fmt.Printf("Erro ao iniciar o servidor: %v\n", err)
	}
}
