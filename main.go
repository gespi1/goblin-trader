/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import "goblin-trader/cmd"

func main() {
	cmd.Execute()
}

// import (
// 	"fmt"
// 	"io/ioutil"``
// 	"net/http"
// 	"os"

// 	log "github.com/sirupsen/logrus"
// )

// var ENDPOINT_URL string = "https://api.twelvedata.com"
// var API_TOKEN string = os.Getenv("API_TOKEN")
// var AUTH_PARAM string = fmt.Sprintf("apikey=%s", API_TOKEN)

// func main() {
// 	resp, err := http.Get(ENDPOINT_URL)
// 	if err != nil {
// 		log.Error(err)
// 	}

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Error(err)
// 	}

// 	sb := string(body)
// 	log.Info(sb)
// }
