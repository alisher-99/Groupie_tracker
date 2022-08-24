package main

import (
	"log"
	"net/http"
	"os"

	"groupie-tracker/api"
	"groupie-tracker/config"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	appConfig, err := config.InitConfig()
	if err != nil {
		errorLog.Fatal(err)
	}

	infoLog.Printf("Server is listening http://%v:%v/", appConfig.Host, appConfig.Port)
	mux := api.AppMux() // принимаем мультиплексер

	if err := http.ListenAndServe(":"+appConfig.Port, mux); err != nil {
		errorLog.Fatal(err)
	}
}
