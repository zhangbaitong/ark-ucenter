package main

import (
	"fmt"
	_"strings"
	"common"
	"gopkg.in/mgo.v2/bson" 
	"os/exec"
)
type Response struct {
	Code int
	Message string
}

type ATUserInfo struct {
      Id bson.ObjectId "_id"
	Ac_id   int
	Info map[string] string
}

type Field_List struct {
	Name string
	List []string
}

func InsertDictionary(Name string,List []string) (ok bool){
	session := common.GetSession()
	if(session==nil){
		return false
	}	
	defer common.FreeSession(session)

	coll := session.DB("at_db").C("dictionary_tab")	
	DicList:=Field_List{Name:Name,List:List}
	err:=coll.Insert(DicList)	
	if(err!=nil){		
		return false
	}
	return true	
}

func  GetDictionary(Name string) {
	session := common.GetSession()
	if(session==nil){
		return 
	}
	defer common.FreeSession(session)

	result := Field_List{}
	coll := session.DB("at_db").C("dictionary_tab")

	err:=coll.Find(&bson.M{"name": Name}).One(&result)
	if err!=nil {
	fmt.Println(err)	
		return 
	}
	fmt.Println(result)	
	return 
}

/*
sh := exec.Command("/bin/echo '"+ SubstrAfter(container.Name,0)+"' >> /etc/dnsmasq.d/dnsmasq.hosts", "service dnsmasq restart")
out, err := sh.CombinedOutput()
fmt.Println("out=", string(out), "err=", err)
if err != nil {
fmt.Println(err, ":", string(out))
} 
*/
func main() {
/*	
	out, err := exec.Command("date").Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("The date is %s\n", out)
*/	
	//List:=strings.Split("qq,email,mobile,weibo", ",")
	//InsertDictionary("check_list",List)
	
	sh := exec.Command("/usr/bin/echo","test",">>","test.dat")
	out, err := sh.CombinedOutput()
	fmt.Println("out=", string(out), "err=", err)
	if err != nil {
		fmt.Println(err, ":", string(out))
	} 
	
	//GetDictionary("check_list")
}
