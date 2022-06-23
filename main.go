package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type User struct{
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   uint   `json:"age"`
}

type Arguments map[string]string
var args Arguments

func Perform(args Arguments, writer io.Writer) error {
	if _,ok:=args["operation"];!ok || args["operation"]==""{
		return fmt.Errorf("-operation flag has to be specified")
	}
	if _,ok:=args["fileName"];!ok || args["fileName"]==""{
		return fmt.Errorf("-fileName flag has to be specified")
	}
	file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 0755)
	if err != nil{
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.Seek(0,0)
	switch args["operation"] {
	case "list":
		b,err:=listUser(file)
		if err!=nil{
			log.Fatal(err)
		}
		if len(b)>0{
			writer.Write(b)
		}else{
			log.Fatal("file is empty")
		}


	case "add":
		value,ok:=args["item"]
		if !ok || value==""{
			return fmt.Errorf("-item flag has to be specified")
		}

		err=addUser(args["item"],file)
		if err!=nil{
		return err
		}

		break
	case "remove":
		value,ok:=args["id"]
		if value=="" || !ok{
			return fmt.Errorf("-id flag has to be specified")
		}
		err:=removeUser(args["id"],file)
		if err!=nil{
			writer.Write([]byte(err.Error()))
		}

		break
	case "findById":
		value,ok:=args["id"]
		if value=="" || !ok{
			return fmt.Errorf("-id flag has to be specified")
		}
			 u, err := searchUser(args["id"], file)
			if err!=nil{
				return err
			}
			if u!=nil{
				b, _ := json.Marshal(u)
				if b != nil {
					writer.Write(b)
				}
			}else{
				writer.Write([]byte(""))
			}
		break
	default:
		return fmt.Errorf("Operation %s not allowed!", args["operation"])


	}
	return nil
}

func parseArgs() Arguments{
	id:=flag.String("id","","Id user for look for")
	operation:=flag.String("operation","",
		"list - list of users\n"+
			"add - Adding user\n"+
			"findById - Getting user by ID\n"+
			"remove - Delete user")
	fileName:=flag.String("fileName","users.json","File name *.json")
	item:=flag.String("item","","Json filename")

	flag.Parse()
	args = Arguments{
		"id":        *id,
		"operation": *operation,
		"item":      *item,
		"fileName":  *fileName,
	}
	return args
}



func main() {
	arg:=parseArgs()
	err := Perform(arg, os.Stdout)
	if err != nil {
		panic(err)
	}
}


func addUser(item string, file *os.File)error{
		var users []User
		var newUser User
		// check valid json
		if !json.Valid([]byte(item)){
			return fmt.Errorf("-item flag has to be specified")
		}
		err:=json.Unmarshal([]byte(item),&newUser)
		if err!=nil{
			return err
		}
		b, err := ioutil.ReadAll(file)
		if json.Valid(b){
			err=json.Unmarshal(b,&users)
			if err!=nil{
				return err
			}
		}
	u,err:=searchUser(newUser.Id,file)
	// if user with id exist - >
	if u!=nil{
		return fmt.Errorf("Item with id %s already exists",newUser.Id)
	}
	users= append(users,newUser)
	b,err=json.Marshal(&users)
	file.Truncate(0)
	file.Seek(0, 0)
	file.Write(b)
	return nil
}
func searchUser(id string, file *os.File)(u *User ,err error){
	file.Seek(0,0)
	var users []User
	b, err := ioutil.ReadAll(file)
		err=json.Unmarshal(b,&users)
		if err!=nil{
			return nil, fmt.Errorf("%v",err)
		}
		for i:=0;i<len(users);i++{
			if users[i].Id==id{
				return &users[i],nil
			}
		}
	return nil,fmt.Errorf("")
}
func removeUser(id string, file *os.File)error{
	u,_:=searchUser(id,file)
	if u==nil{
		return fmt.Errorf("Item with id %s not found",id)
	}
		var users []User
		b, err := ioutil.ReadAll(file)
		if err!=nil{
			return err
		}
		err=json.Unmarshal(b,&users)
		if err!=nil{
			return err
		}
		for k, u := range users {
			if u.Id == id {
				users = append(users[:k], users[k+1:]...)
			}
		}
		b,err=json.Marshal(users)
		if err!=nil{
			return err
		}
		file.Truncate(0)
		file.Seek(0, 0)
		file.Write(b)
		return nil
}
func listUser(file *os.File) ([]byte, error){
	b, err := ioutil.ReadAll(file)
	if err!=nil{
		return nil,err
	}
	return b,nil
}