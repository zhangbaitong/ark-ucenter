package main

import (
	"fmt"
	_"strings"
	"common"
	"gopkg.in/mgo.v2/bson" 
)
type Response struct {
	Code int
	Message string
}

type ATUserInfo struct {
      Id bson.ObjectId "_id"
	Ac_id   int
	Create_time   int
	Info map[string] string
}

type UserJoinLocus struct {
	Id bson.ObjectId "_id"
	Mid string
	Mid_c string
	Show_id  string
	Create_time   int
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

func  SaveLocus(show_id string) {
	session := common.GetSession()
	if(session==nil){
		return 
	}
	defer common.FreeSession(session)

	var result []ATUserInfo
	//result := ATUserInfo[]
	coll := session.DB("at_db").C("user_tab")

	err:=coll.Find(&bson.M{}).All(&result)
	if err!=nil {
		fmt.Println(err)	
		return 
	}

	Locus:=UserJoinLocus{}
	coll = session.DB("at_db").C("locus_tab")
	for i := 0; i < len(result); i++ {

		Locus.Id=bson.NewObjectId()
		Locus.Mid=result[i].Id.Hex()
		Locus.Mid_c=result[i].Id.Hex()
		Locus.Show_id=show_id
		Locus.Create_time=result[i].Create_time
		Locus.Info=result[i].Info
		err=coll.Insert(&Locus)
		if(err!=nil){		
			fmt.Println(err)
			break;
		}
	}	
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
//List:=strings.Split("qq,email,mobile,weibo", ",")
//InsertDictionary("check_list",List)

//GetDictionary("check_list")
	SaveLocus("old");
}
