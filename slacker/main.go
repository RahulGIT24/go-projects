package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/shomali11/slacker"
	"github.com/joho/godotenv"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
	error := godotenv.Load(".env")
	if error!=nil{
		log.Fatal(error)
	}

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.Command("My yob is <year>", &slacker.CommandDefinition{
		Description: "Age Calculator",
		Examples:    []string{"My yob is 2020"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				fmt.Println(err)
			}
			currentYear := time.Now().Year()
			age := 2024 - yob
			if yob > currentYear{
				response.Reply("You were not born yet!")
				return
			}
			r := fmt.Sprintf("Age is %d", age)
			response.Reply(r)
		},
	})

	go printCommandEvents(bot.CommandEvents())

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
