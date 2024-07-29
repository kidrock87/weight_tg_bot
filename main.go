package main

import (
	//"github.com/gin-contrib/cors"

	"log"
	"net/http"
	"os"
	"strings"
	"tg_weight_bot/controllers"
	"tg_weight_bot/middlewares"
	"tg_weight_bot/models"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	models.ConnectDatabase()

	err := godotenv.Load(".env")
	telegram_apitoken := os.Getenv("TELEGRAM_APITOKEN")
	log.Printf(telegram_apitoken)

	bot, err := tgbotapi.NewBotAPI(telegram_apitoken)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil { // ignore any non-Message updates
			if update.Message.Text == "/status" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неправильный формат данных или команды")
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages

			var textMessage string = update.Message.Text

			FormatedTextMessage := strings.Replace(textMessage, "\n", "\r\n", -1)
			//Разделяем текст по строкам
			s := strings.Split(strings.Replace(FormatedTextMessage, "\r\n", "\n", -1), "\n")
			for i := 0; i < len(s); i++ {
				//Разделяем строки по двоеточию, слева название вида спорта, справа результаты.
				result := strings.Split(s[i], ":")
				//Если строка раздвоилась то дальше обрабатываем результаты
				//Если нет, то ругаем пользователя за формат.
				if len(result) > 1 {
					sporttype := result[0]
					results := result[1]
					res := controllers.CreateRecord(sporttype, results, int(update.Message.Chat.ID))
					if res == "Success" {
						returnMessage := controllers.GetAnalyticsBySportLast(strings.TrimSpace(sporttype), int(update.Message.Chat.ID))
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, returnMessage)
						if _, err := bot.Send(msg); err != nil {
							log.Panic(err)
						}

					}
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неправильный формат данных или команды")
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
				}
			}
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /sayhi and /status."
		case "chart":
			aaa := strings.Split(update.Message.Text, " ")
			if len(aaa) > 1 {
				returnMessage := controllers.GetAnalyticsChartBySport(strings.TrimSpace(aaa[1]), int(update.Message.Chat.ID))
				photo := tgbotapi.NewPhoto(update.Message.From.ID, tgbotapi.FileURL(returnMessage))
				if _, err = bot.Send(photo); err != nil {
					log.Panic(err)
				}
				msg.Text = "returnMessage"

			} else {

			}
		case "lastresult":
			aaa := strings.Split(update.Message.Text, " ")
			if len(aaa) > 1 {
				returnMessage := controllers.GetAnalyticsBySportLast(strings.TrimSpace(aaa[1]), int(update.Message.Chat.ID))
				msg.Text = returnMessage

			} else {
				msg.Text = "Нет аналитики по этому виду спорта"
			}

		case "status":
			msg.Text = "I'm ok."
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}

	}

	models.ConnectDatabase()

	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())

	r.GET("/api/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.Static("/assets", "./assets")

	//r.POST("/api/auth/signup", controllers.Register)
	//r.POST("/api/auth/signin", controllers.Login)

	//r.GET("/api/books/:id", controllers.FindBook)
	//r.GET("/api/books", controllers.FindBooks)

	//r.GET("/api/chapters", controllers.FindChapters)
	//r.GET("/api/chapters/:id", controllers.FindChapter)

	//r.GET("/api/categories/:id", controllers.FindCategory)

	//r.GET("/api/posts/:id", controllers.FindPost)

	//api := r.Group("/api")
	//api.Use(middlewares.JwtAuthMiddleware())

	//api.POST("/auth/verifyphone", controllers.VerifyPhone)
	//api.POST("/auth/verifyphonecode", controllers.VerifyPhoneCode)
	//api.POST("/auth/createaccount", controllers.CreateAccount)
	//api.POST("/auth/createavatar", controllers.CreateAvatar)

	//api.POST("/books", controllers.CreateBook)
	//api.PATCH("/books/:id", controllers.UpdateBook)
	//api.DELETE("/books/:id", controllers.DeleteBook)

	//api.POST("/posts", controllers.CreatePost)
	//api.PATCH("/posts/:id", controllers.UpdatePost)

	//api.GET("/tags/:name", controllers.FindBooksByTag)
	//api.POST("/tags/:name", controllers.FindBooksByTag)

	//api.POST("/chapters", controllers.CreateChapter)
	//api.PATCH("/chapters/:id", controllers.UpdateChapter)
	//api.DELETE("/chapter/:id", controllers.DeleteChapter)

	//api.POST("/categories", controllers.CreateCategory)
	//api.PATCH("/categories/:id", controllers.UpdateCategory)
	//api.DELETE("/categories/:id", controllers.DeleteCategory)

	//protected := r.Group("/admin")
	//protected.Use(middlewares.JwtAuthMiddleware())
	//protected.GET("/user", controllers.CurrentUser)

	r.Run()
}
