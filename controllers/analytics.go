package controllers

import (
	"log"
	"math"
	"strconv"
	"tg_weight_bot/models"
)

func GetAnalyticsBySportLast(sportName string, chatId int) (Message string) {
	s := models.Sport{}
	var records []models.Record
	var diffsign string
	//var percentageDifference float64
	//Показать последний и предпоследний(если есть) результат в этом спорте у этого чата и сравнить.
	s.Name = sportName
	var _, err = s.GetSportByName()

	if err != nil {
		log.Print("error:", err.Error())
		return
	}

	models.DB.Debug().Limit(2).Order("id desc").Where("sport_id = ? AND chat_id >= ?", s.ID, chatId).Preload("Sport").Find(&records)

	if err != nil {
		log.Print("error:", err.Error())
		return
	}
	lastResult := float64(records[0].Result)
	firstResult := float64(records[1].Result)
	//Взять результаты первого и второго
	//Сравнить их и вычислить процент изменения
	//[|(a-b)|/(a+b)/2] × 100
	var resultdif float64 = lastResult - firstResult
	uppv := math.Abs(resultdif)
	ddpv := (lastResult + firstResult) / 2
	var m float64 = uppv / ddpv
	m = m * 100
	log.Print(m)
	if resultdif > 0 {
		diffsign = "\xF0\x9F\x9A\x80"
	} else {
		diffsign = "\xE2\x86\x98"
	}
	log.Print(diffsign)
	//percentageDifference = 10 / 95
	sf := strconv.FormatFloat(m, 'f', 2, 64)
	rMessage := "Последний результат: " + records[0].Oresult + "\n Предпоследний результат: " + records[1].Oresult + "\n Прогресс:" + diffsign + " " + sf + "%\n"

	return rMessage
}
