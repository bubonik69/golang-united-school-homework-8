package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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
	operate:=args["operation"]
	switch operate {
	case "list":
		listUser(args["fileName"],writer)
	case "add":
		if args["item"]==""{
			return fmt.Errorf("-item flag has to be specified")
		}
		if _,ok:=args["item"];ok{
			err:=addUser(args["item"],file,writer)
			if err!=nil{
			writer.Write([]byte(err.Error()))
			}
		}
		break
	case "remove":
		if args["id"]==""{
			return fmt.Errorf("-id flag has to be specified")
		}
		if _,ok:=args["id"];ok{
			err:=removeUser(args["id"],file)
			if err!=nil{
				writer.Write([]byte(err.Error()))
			}
		}
		break
	case "findById":
		if args["id"]==""{
			return fmt.Errorf("-id flag has to be specified")
		}
		if _,ok:=args["id"];ok {
			err, u := searchUser(args["id"], file)
			checkErr(err, writer)
			if u!=nil{
				b, _ := json.Marshal(u)
				if b != nil {
					writer.Write(b)
				}
			}else{
				writer.Write([]byte(""))
			}

		}
		break
	default:
		return fmt.Errorf("Operation %s not allowed!", operate)


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
	//var buffer bytes.Buffer
	//err:=Perform(args,&buffer)
	//if err!=nil{
	//	fmt.Println(err)
	//}
	arg:=parseArgs()
	err := Perform(arg, os.Stdout)
	if err != nil {
		panic(err)
	}


}


func checkErr (err error,writer io.Writer){
	if err!=nil {
		writer.Write([]byte(error.Error(err)))
	}
}

func addUser(item string, file *os.File, writer io.Writer)error{
		var users []User
		var newUser User
		// check valid json
		defer file.Seek(0,0)
		file.Seek(0,0)
		if !json.Valid([]byte(item)){
			return fmt.Errorf("-item flag has to be specified")
		}
		err:=json.Unmarshal([]byte(item),&newUser)
		checkErr(err, writer)
		//check id
		b, err := ioutil.ReadAll(file)
		checkErr(err, writer)

		if json.Valid(b){
			err=json.Unmarshal(b,&users)
			if err!=nil{
				checkErr(err,writer)
			}
		}
	err,u:=searchUser(newUser.Id,file)
	// if user with id exist - >
	if u!=nil{
		return fmt.Errorf("Item with id %s already exists",newUser.Id)
	}
	users= append(users,newUser)
	b,err=json.Marshal(&users)
	checkErr(err,writer)
	file.Write(b)
	return nil
}
func searchUser(id string, file *os.File)(err error, u *User){
	defer file.Seek(0,0)
	var users []User
	file.Seek(0,0)
	b, err := ioutil.ReadAll(file)
		err=json.Unmarshal(b,&users)
		if err!=nil{
			return fmt.Errorf("%v",err),nil
		}
		for i:=0;i<len(users);i++{
			if users[i].Id==id{
				return nil,&users[i]
			}
		}
	return nil,u
}
func removeUser(id string, file *os.File)error{
	defer file.Seek(0,0)
	err,u:=searchUser(id,file)
	if u!=nil{
		return fmt.Errorf("Item with id %s not found",id)
	}
	var users []User
	file.Seek(0,0)
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
func listUser(fileName string, writer io.Writer){

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil{
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.Seek(0,0)
	var b []byte
	b, err = ioutil.ReadAll(file)
	checkErr(err, writer)
	writer.Write(b)
}