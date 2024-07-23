package controllers

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"tg_weight_bot/models"
)

func CreateRecord(sportName, results string, chatId int) {
	s := models.Sport{}
	r := models.Record{}
	sportName = strings.TrimSpace(sportName)
	var assigned_sport_type string
	var final_result float64
	sampleRegexp := regexp.MustCompile(`\d\d`)

	//Проверить тип результатов, разобрать результаты и присвоить тип виду спорта
	result := strings.Split(results, ",")

	for i := 0; i < len(result); i++ {
		log.Print(result[i])
		match := sampleRegexp.MatchString(result[i])

		if strings.Index(result[i], "х") != -1 || strings.Index(result[i], "x") != -1 {
			assigned_sport_type = "countandweight"
			countandweight := strings.Split(result[i], "х")
			r_count, err := strconv.ParseFloat(strings.TrimSpace(countandweight[0]), 64)
			r_weight, err := strconv.ParseFloat(countandweight[1], 64)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			final_result = final_result + r_count*r_weight

		} else if strings.Index(result[i], "-") != -1 {
			assigned_sport_type = "time"
		} else if strings.Index(result[i], ".") != -1 {
			assigned_sport_type = "weight"
		} else if match == true {
			assigned_sport_type = "count"
		}
		//рассчитываем финальный результат исходя из типа спорта
		//10x65 count and weight \Bx\B
		//10 count ^\d\d$
		//108.6 weight [.]
		//12:26 time :

	}
	log.Print(assigned_sport_type)
	log.Print(final_result)
	log.Print(sportName)
	s.Name = sportName
	s.Type = assigned_sport_type

	//Проверяем есть ли такой вид спорта
	//Если нет, то создаем
	var _, err = s.GetSportByName()

	r.Sport = s
	r.Result = int(final_result)
	r.ChatID = chatId

	if err != nil {
		log.Print("error:", err.Error())
		return
	}

	_, err = r.SaveRecord()
	//Записываем результат в рекордс

}
