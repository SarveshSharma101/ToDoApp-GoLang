package main

import (
	controller "ToDoApp/Controller"
	service "ToDoApp/Service"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	env := "Local"
	if len(os.Args) > 1 {
		env = os.Args[1]
	}
	username, password, url, dbname := readConfigs("../Resources/Config/AppConfig.yaml", env)

	//Initialize the DB
	service.InitDbConnection(dbname, username, password, url)
	//Get router
	router := controller.GetRouter()
	fmt.Println("Server is listening at 7000...")
	//Start the server
	http.ListenAndServe(":7000", router)
}

func readConfigs(configFilePath string, env string) (uname, password, url, dbname string) {
	//Read the yaml file for configs
	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}
	//convert the byte array to map
	configMap := map[string]map[string]map[string]string{}
	if err := yaml.Unmarshal(yamlFile, &configMap); err != nil {
		panic(err)
	}
	//fetch the value of env
	dbConfigs := configMap["Environment"][env]
	if len(dbConfigs) <= 0 {
		fmt.Println("No such env: ", env)
		panic("Please provide valid env")
	}

	//validate required values are present in the config files else panic
	if len(dbConfigs["dbName"]) <= 0 {
		fmt.Println("DB name not found on config file, for env: ", env)
		panic("Please provide DB name in configs")
	} else if len(dbConfigs["url"]) <= 0 {
		fmt.Println("DB url not found on config file, for env: ", env)
		panic("Please provide DB url in configs")
	} else if len(dbConfigs["password"]) <= 0 {
		fmt.Println("DB password not found on config file, for env: ", env)
		panic("Please provide DB password in configs")
	} else if len(dbConfigs["username"]) <= 0 {
		fmt.Println("DB username not found on config file, for env: ", env)
		panic("Please provide DB username in configs")
	} else {
		dbname = dbConfigs["dbName"]
		url = dbConfigs["url"]
		password = dbConfigs["password"]
		uname = dbConfigs["username"]
	}

	return uname, password, url, dbname
}
