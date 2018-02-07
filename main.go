package main

import (
	"github.com/zamedic/telegram"
	"github.com/zamedic/dynamodb"
	"github.com/zamedic/multipoolminerbot/user"
	"github.com/zamedic/multipoolminerbot/multipoolminer"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"log"
	"github.com/zamedic/multipoolminerbot/monitor"
)

func main(){

	db := dynamodb.NewConnection()

	telegramStore := telegram.NewDynamoState(db)
	userStore := user.NewDynamoStore(db)


	telegramService := telegram.NewService(telegramStore)
	poolService := multipoolminer.NewService()
	userService := user.NewService(userStore,poolService)
	_ = monitor.NewService(userStore,poolService,telegramService)

	telegramService.RegisterCommand(user.NewStartCommand(telegramStore,telegramService,userStore))
	telegramService.RegisterCommandLet(user.NewAddTokenCommandlet(userService,telegramService))

	log.Println("All systems go!")
	errs := make(chan error, 2)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Println("terminated", <-errs)



}
