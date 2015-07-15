package action

import (
	"common"
	"encoding/json"
	"fmt"
	"github.com/dchest/authcookie"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"net/url"
	"time"
	 "gopkg.in/mgo.v2/bson" 
)

type Account struct {
	Acid        int
	User_Name  string
	Password string
	Email       string
	Mobile      string
	Status      int
	Create_time int
}

type StructLogin struct {
	User_Name string
	Password  string
}

const (
	INSERT string = "insert into account_tab (ac_name,ac_password,email,mobile,status,create_time) values (?,?,?,?,?,?)"

	//cookie加密、解密使用
	KEY        string = "QAZWERT4556"
	COOKIENAME string = "MNBVCXZ"
)

func getParams(r *http.Request) (params string) {
	strPostData := r.FormValue("request")
	fmt.Println("strPostData :", strPostData)
	var request common.RequestData

	err := json.Unmarshal([]byte(strPostData), &request)
	if err != nil {
		fmt.Println("json data decode faild :", err)
		return ""
	}
	fmt.Println("request.Params :", request.Params)
	return request.Params
}

func setParams(strMethod string, code int, strMessage string, strData string) (strbody []byte, err error) {
	v1 := common.Response{Method: strMethod, Code: code,Message:strMessage, Data: strData}
	body, err := json.Marshal(v1)
	if err != nil {
		fmt.Println(err)
		return body, err
	}
	return body, nil
}

var logger *log.Logger

func init() {
	if logger == nil {
		logger = common.Log()
	}
}

func register_insert(ac *Account) (ok bool) {
	mydb := common.GetDB()
	if(mydb==nil){
		return false
	}	
	defer common.FreeDB(mydb)

	tx, err := mydb.Begin()
	if err != nil {
		fmt.Println(err)
		return false
	}
	stmt, err := tx.Prepare(INSERT)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(ac.User_Name, ac.Password, ac.Email, ac.Mobile, 0, time.Now().Unix())

	if err != nil {
		fmt.Println(err)
		return false
	}
	tx.Commit()
	return true
}

func register_insert_m(ac *Account) (ok bool) {
	session := common.GetSession()
	if(session==nil){
		return false
	}	
	defer common.FreeSession(session)

	coll := session.DB("at_db").C("user_tab")	
	err:=coll.Insert(ac)
	if(err!=nil){		
		return false
	}
	return true
}

//查询账户是否存在
func isFieldExist_m(name string, value string) bool {
	session := common.GetSession()
	if(session==nil){
		return false
	}	
	defer common.FreeSession(session)

	result := []Account{}
	coll := session.DB("at_db").C("user_tab")

	coll.Find(&bson.M{name:value}).Sort(name).All(&result)
	if len(result)==0 {
		return false;
	}

	return true
}

//帐号注册
func RegisterAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//把前台参数转换成结构体
	strParams := getParams(r)
	var account Account
	err := json.Unmarshal([]byte(strParams), &account)
	var strBody []byte
	if err != nil {
		logger.Println("json data decode faild :", err)
		strBody, _ = setParams("/auth/register", 1, "json data decode faild !", "")
		w.Write(strBody)
		return
	}

	//参数校验
	if account.User_Name == "" {
		logger.Println("action_certification：User_Name can't be empty")
		strBody, _ = setParams("/auth/register", 1, "User_Name can't be empty!", "")
		w.Write(strBody)
		return
	}

	if account.Password == "" {
		logger.Println("action_certification：password can't be empty")
		strBody, _ = setParams("/auth/register", 1, "password can't be empty!", "")
		w.Write(strBody)
		return
	}

	if account.Email == "" {
		logger.Println("action_certification：email can't be empty")
		strBody, _ = setParams("/auth/register", 1, "ac_email can't be empty!", "")
		w.Write(strBody)
		return
	}

	if account.Mobile == "" {
		logger.Println("action_certification：mobile can't be empty")
		strBody, _ = setParams("/auth/register", 1, "mobile can't be empty!", "")
		w.Write(strBody)
		return
	}


	//校验账户、邮箱、手机号码是否已存在
	if true == isFieldExist_m("user_name", account.User_Name) {
		strBody, _ = setParams("/auth/register", 1, "User_Name is already exist!", "")
		w.Write(strBody)
		return
	}
	
	if true == isFieldExist_m("email", account.Email) {
		strBody, _ = setParams("/auth/register", 1, "email is already exist!", "")
		w.Write(strBody)
		return
	}
	if true == isFieldExist_m("mobile", account.Mobile) {
		strBody, _ = setParams("/auth/register", 1, "mobile is already exist!", "")
		w.Write(strBody)
		return
	}

	ok := register_insert_m(&account)
	if ok {
		strBody, _ = setParams("/auth/register", 0, "ok", "")
	} else {
		strBody, _ = setParams("/auth/register", 1, "database error!", "")
	}
	w.Write(strBody)
	return
}

//登录插入
func login_query(login *StructLogin) (ok bool) {
	mydb := common.GetDB()
	if(mydb==nil){
		return false
	}	
	defer common.FreeDB(mydb)

	strSQL := fmt.Sprintf("select count(ac_name) from account_tab where (ac_name='%s' or email='%s' or mobile='%s') and ac_password='%s'", login.User_Name, login.User_Name, login.User_Name, login.Password)
	rows, err := mydb.Query(strSQL)
	if err != nil {
		return false
	} else {
		defer rows.Close()
		var nCount int
		for rows.Next() {
			rows.Scan(&nCount)
		}

		if nCount == 0 {
			return false
		}
	}
	return true
}

//登录插入
func login_query_m(strUser_Name,strPassword string) (ok bool) {
	session := common.GetSession()
	if(session==nil){
		return false
	}	
	defer common.FreeSession(session)

	result := []Account{}
	coll := session.DB("at_db").C("user_tab")

	condition_or:=bson.M{"$or": []bson.M{bson.M{"user_name":strUser_Name},bson.M{"email":strUser_Name},bson.M{"mobile":strUser_Name}}}
	coll.Find(condition_or).Sort(strUser_Name).All(&result)
	// 显示数据
	for _, m := range result {
		if(m.Password==strPassword){
			return true
		}		
	}

	return false;
}
//判断cookie是否存在
func isCookieExist(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie(COOKIENAME)
	if err == nil {
		var cookieValue = cookie.Value
		login := authcookie.Login(cookieValue, []byte(KEY))
		if login != "" {
			strBody, _ := setParams("/auth/login", 0, "ok", "")
			w.Write(strBody)
			return true
		}
	}
	return false
}

//生成cookie，放到reponse对象
func generateCookie(w http.ResponseWriter, r *http.Request, userNmae string, number int) {
	timeLength := 24 * time.Hour
	cookieValue := authcookie.NewSinceNow(userNmae, timeLength, []byte(KEY))
	expire := time.Now().Add(timeLength)
	cookie := http.Cookie{Name: COOKIENAME, Value: cookieValue, Path: "/", Expires: expire, MaxAge: 86400}
	http.SetCookie(w, &cookie)
}

//登录
func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	ok := isCookieExist(w, r)

	if ok {
		strBody, _ := setParams("/auth/login", 0, "ok", "")
		w.Write(strBody)
		return
	}

	//把前台参数转换成结构体
	strParams := getParams(r)
	var login StructLogin
	err := json.Unmarshal([]byte(strParams), &login)
	fmt.Println(login)
	var strBody []byte
	if err != nil {
		logger.Println("json data decode faild :", err)
		strBody, _ = setParams("/auth/login", 1, "json data decode faild !", "")
		w.Write(strBody)
		return
	}

	//参数校验
	if login.User_Name == "" {
		logger.Println("action_certification：User_Name can't be empty")
		strBody, _ = setParams("/auth/login", 1, "User_Name can't be empty!", "")
		w.Write(strBody)
		return
	}
	if login.Password == "" {
		logger.Println("action_certification：password can't be empty")
		strBody, _ = setParams("/auth/login", 1, "password can't be empty!", "")
		w.Write(strBody)
		return
	}

	ok = login_query_m(login.User_Name,login.Password)
	if ok {
		strBody, _ = setParams("/auth/login", 0, "ok", "")
	} else {
		strBody, _ = setParams("/auth/login", 1, "User_Name or pwd not right", "")
	}

	generateCookie(w, r, login.User_Name, 1)

	w.Write(strBody)
	return
}

//注销
func Logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookie := http.Cookie{Name: COOKIENAME, Path: "/", MaxAge: -1}
	http.SetCookie(w, &cookie)
	strBody, _ := setParams("/auth/logout", 0, "ok", "")
	w.Write(strBody)
}

//查询账户是否存在
func isFieldExist(name string, value string) bool {
	mydb := common.GetDB()
	if(mydb==nil){
		return false
	}	
	defer common.FreeDB(mydb)	

	strSQL := fmt.Sprintf("select count(ac_name) from account_tab where %s='%s' ", name, value)
	rows, err := mydb.Query(strSQL)
	if err != nil {
		return false
	} else {
		defer rows.Close()
		var nCount int
		for rows.Next() {
			rows.Scan(&nCount)
		}

		if nCount == 0 {
			return false
		}
	}
	return true
}

// func Logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	strParams := getParams(r)
// 	fmt.Fprint(w, "%s BYE BYE !\n", strParams)
// }

func GetAcidByOpenid(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	strParams := getParams(r)
	var openvalue map[string]interface{}
	err := json.Unmarshal([]byte(strParams), &openvalue)
	var strBody []byte
	if err != nil {
		logger.Println("json data decode faild :", err)
		strBody, _ = setParams("/auth/getacid", 1, "json data decode faild!", "")
		w.Write(strBody)
		return
	}

	strOpenid, ok := openvalue["openid"].(string)
	if !ok {
		strBody, _ = setParams("/auth/getacid", 1, "params error, ac_name miss !", "")
		w.Write(strBody)
		return
	}

	mydb := common.GetDB()
	if(mydb==nil){
		strBody, _ = setParams("/auth/getacid", 1, "database error!!!!", "")
		w.Write(strBody)
		return 
	}	
	defer common.FreeDB(mydb)

	strSQL := fmt.Sprintf("select acid from openid_tab where openid='%s'", strOpenid)
	rows, err := mydb.Query(strSQL)
	if err != nil {
		strBody, _ = setParams("/auth/getacid", 1, "database error !", "")
	} else {
		defer rows.Close()
		var nAcid int
		for rows.Next() {
			rows.Scan(&nAcid)
		}
		if nAcid == 0 {
			strBody, _ = setParams("/auth/getacid", 1, "user acid not exist!", "")
		} else {
			strData := fmt.Sprintf("{\"acid\":\"%d\"}", nAcid)
			strBody, _ = setParams("/auth/getacid", 0, "ok", strData)
		}
	}

	w.Write(strBody)
}

func update_password(strAcName string, strOldPwd string, strNewPwd string) {
	mydb := common.GetDB()
	if(mydb==nil){
		return 
	}	
	defer common.FreeDB(mydb)

	tx, err := mydb.Begin()
	if err != nil {
		fmt.Println(err)
	}
	stmt, err := tx.Prepare(" UPDATE account_tab SET ac_password=? where (ac_name=? or email=? or mobile=?) and ac_password=? ")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(strNewPwd, strAcName, strAcName,strAcName, strOldPwd)

	if err != nil {
		fmt.Println(err)
	}
	tx.Commit()
}

func update_password_m(strAcName string,strNewPwd string) bool {
	session := common.GetSession()
	if(session==nil){
		return false
	}	
	defer common.FreeSession(session)

	coll := session.DB("at_db").C("user_tab")
	// No changes is a no-op and shouldn't return an error.
	condition_or:=bson.M{"$or": []bson.M{bson.M{"user_name":strAcName},bson.M{"email":strAcName},bson.M{"mobile":strAcName}}}
	err := coll.Update(condition_or, bson.M{"$set": bson.M{"password": strNewPwd}})
	if(err==nil){
		return true
	}
	return false
}

func ChangePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	strParams := getParams(r)
	var openvalue map[string]interface{}
	err := json.Unmarshal([]byte(strParams), &openvalue)
	var strBody []byte
	if err != nil {
		logger.Println("json data decode faild :", err)
		strBody, _ = setParams("/auth/changepw", 1, "json data decode faild!", "")
		w.Write(strBody)
		return
	}

	strAcName, ok := openvalue["user_name"].(string)
	if !ok {
		strBody, _ = setParams("/auth/changepw", 1, "params error, user_name miss !", "")
		w.Write(strBody)
		return
	}

	strOldPwd, ok := openvalue["old_password"].(string)
	if !ok {
		strBody, _ = setParams("/auth/changepw", 1, "params error, old_password miss !", "")
		w.Write(strBody)
		return
	}

	strNewPwd, ok := openvalue["new_password"].(string)
	if !ok {
		strBody, _ = setParams("/auth/changepw", 1, "params error, new_password miss !", "")
		w.Write(strBody)
		return
	}

	ok = login_query_m(strAcName,strOldPwd)
	if !ok {
		strBody, _ = setParams("/auth/changepw", 1, "user_name or old pwd not right", "")
		w.Write(strBody)
		return
	}

	ok=update_password_m(strAcName,strNewPwd)
	if ok {
		strBody, _ = setParams("/auth/changepw", 0, "ok", "")
	}else {
		strBody, _ = setParams("/auth/changepw", 1, "change password faild", "")
	}
	
	w.Write(strBody)
	return		

/*
	strSQL := fmt.Sprintf("select count(ac_name) from account_tab where (ac_name='%s' or email='%s' or mobile='%s') and ac_password='%s'",
		strAcName, strAcName, strAcName, strOldPwd)
	fmt.Println("strSQL=", strSQL)

	mydb := common.GetDB()
	if(mydb==nil){
		strBody, _ = setParams("/auth/changepw", 1, "database error!!!!", "")
		w.Write(strBody)
		return 
	}	
	defer common.FreeDB(mydb)

	rows, err := mydb.Query(strSQL)
	if err != nil {
		strBody, _ = setParams("/auth/changepw", 1, "database error !", "")
	} else {
		defer rows.Close()
		var nCount int
		for rows.Next() {
			rows.Scan(&nCount)
		}
		if nCount == 0 {
			strBody, _ = setParams("/auth/changepw", 1, "user not exist or passsword error!", "")
		} else {
			update_password(strAcName, strOldPwd, strNewPwd)
			strBody, _ = setParams("/auth/changepw", 0, "ok", "success")
		}
	}
*/	

}

type AUResponse struct {
	Code int
	Message string
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_name := r.FormValue("user_name")
	token := r.FormValue("token")
	fmt.Println("user_name=",user_name)
	fmt.Println("token=",token)
	
	strCheckURL:="https://connect.funzhou.cn/oauth2/privilige"
	strClientID:="user_center"
	strInterface:="get_user_info"

	value := url.Values{}
	value.Set("token",token)
	value.Set("privilige",strInterface)
	strBody,err:=common.Invoker(common.HTTP_POST,strCheckURL,value)
	if err!=nil {
		fmt.Println(err)
		strBody, _ := setParams("/auth/get_user_info", 1, "submit check  faild !", "")
		w.Write(strBody)
		return
	}
	fmt.Println(strBody);

	var result AUResponse
	err = json.Unmarshal([]byte(strBody), &result)
	if err != nil {
		strBody, _ := setParams("/auth/get_user_info", 1, "json data decode faild!", "")
		w.Write(strBody)
		return 
	}

	if  result.Code!=0 {
		strBody, _ := setParams("/auth/get_user_info", 1, "user check faild beacuse of "+result.Message, "")
		w.Write(strBody)
		return
	}

	session := common.GetSession()
	if(session==nil){
		strBody, _ := setParams("/auth/get_user_info", 1, "get DB faild !!", "")
		w.Write(strBody)
		return 
	}	
	defer common.FreeSession(session)

	user_info := Account{}
	coll := session.DB("at_db").C("user_tab")

	condition_or:=bson.M{"$or": []bson.M{bson.M{"user_name":user_name},bson.M{"email":user_name},bson.M{"mobile":user_name}}}
	coll.Find(condition_or).Sort(user_name).One(&user_info)

	strResult, err := json.Marshal(user_info)
	if err != nil {
		strBody, _ := setParams("/auth/get_user_info", 1, "json data encode faild!", "")
		w.Write(strBody)
		return 
	}

	strResult, _ = setParams("/auth/get_user_info", 0, "ok", string(strResult))
	w.Write(strResult)
	return		
}