package main

import (
	"git.iotopen.se/lib/lynx"
	"log"
)

func main() {
	client := lynx.NewClient(&lynx.Options{
		Authenticator: lynx.AuthApiKey{Key: "08d22215ec037f8469d11c3f78c197ee"},
		ApiBase:       "https://lynx-dev.iotopen.se/",
		//ApiBase:       "http://127.0.0.1:8081/",
	})
	installations, err := client.GetInstallations()
	if err != nil {
		log.Fatal("Err:", err)
	}
	for _, v := range installations {
		log.Println(v.ID, v.Name)
	}
	log.Println("---------------------")

	functions, err := client.GetFunctions(19, map[string]string{"name": "*blinder"})
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range functions {
		log.Println(v.Meta["name"])
	}
}
