package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type User struct{
	id 		string		`json: id`
	email 	string		`json: email`
	age 	int			`json: age`
}

type Arguments map[string]string


func Perform(args Arguments, writer io.Writer) error {
	return nil
}

func main() {
	var users []User

	//err := Perform(parseArgs(), os.Stdout)
	//if err != nil {
	//	panic(err)
	//}
	file, err := os.OpenFile("users.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil{
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	var b []byte
	b, err = ioutil.ReadAll(file)
	fmt.Println(string(b))
	checkErr(err)
	if json.Valid(b){
		err=json.Unmarshal(b,&users)
		checkErr(err)
	}
	for k,v:=range users{
		fmt.Println(k,v)
	}

}


func checkErr (err error){
	if err!=nil {
		fmt.Println("Error : ---->", err)
	}
}