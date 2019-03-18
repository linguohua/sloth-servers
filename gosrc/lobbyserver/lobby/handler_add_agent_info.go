package lobby

import (
	"encoding/json"
	"net/http"
	"unicode"

	log "github.com/sirupsen/logrus"
)

func replyAddAgentInfo(w http.ResponseWriter, msg string, isSuccessed bool) {

	type Rsp struct {
		Successed bool   `json:"Successed"`
		Message   string `json:"Message"`
	}

	var rsp = &Rsp{}

	rsp.Successed = isSuccessed
	rsp.Message = msg
	rspjson, _ := json.Marshal(rsp)
	w.WriteHeader(200)
	w.Write(rspjson)
	return
}

func isPhoneNumber(s string) bool {
	if len(s) != 11 {
		return false
	}

	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func isValidCode(s string) bool {
	if s == "" {
		return false
	}

	if len(s) > 6 {
		return false
	}

	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// OnAddAgentInfo 请求认证
func OnAddAgentInfo(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("OnAddAutoInfo, userID:", userID)
	// TODO: llwant mysql
	// r.ParseForm()
	// // GameID := r.FormValue("GameID")
	// // UserID := r.FormValue("UserID")
	// UserID := userID
	// phone := r.FormValue("phone")
	// code := r.FormValue("code")

	// if !isPhoneNumber(phone) {
	// 	replyAddAgentInfo(w, "请输入正确的手机号", false)
	// 	return
	// }

	// if code == "" {
	// 	replyAddAgentInfo(w, "请输入验证码", false)
	// 	return
	// }

	// if !isValidCode(code) {
	// 	replyAddAgentInfo(w, "请输正确的验证码", false)
	// 	return
	// }

	// // 验证验证码是否正确
	// err := webdata.CheckPhoneAuthCode(UserID, phone, code)
	// if err != nil {
	// 	log.Printf("CheckPhoneAuthCode err, reason:%s", err)
	// 	replyAddAgentInfo(w, fmt.Sprintf("%v", err), false)
	// 	return
	// }

	// // 检查是否已经写入到库
	// conn := pool.Get()
	// defer conn.Close()

	// exist, err := redis.Int(conn.Do("EXISTS", gconst.AgentInfo+phone))
	// if err != nil {
	// 	replyAddAgentInfo(w, "信息提交失败", false)
	// 	return
	// }

	// if exist == 1 {
	// 	replyAddAgentInfo(w, "该号码已提交成功，请勿反复提交", false)
	// 	return
	// }

	// 验证手机号码是否正确
	// type 1:手机号 2:身份证
	// result, err := webdata.CheckRealNameAuthPlus(phone, "1")
	// if err != nil || result != 1 {
	// 	log.Printf("CheckRealNameAuthPlus type:1 err, reason:%s", err)
	// 	rsp.Successed = false
	// 	rsp.Message = "该号码已提交成功，请勿反复提交"
	// 	rspjson, _ := json.Marshal(rsp)
	// 	w.WriteHeader(200)
	// 	w.Write(rspjson)
	// 	return
	// }

	// // AgentInfo

	// // 写入数据库
	// err = webdata.SetRealNameAuthInfo(UserID, phone, "", "")
	// if err != nil {
	// 	log.Printf("SetPhoneAuthCode err, reason:%s", err)
	// 	rsp.Successed = false
	// 	rsp.Message = "身份信息写入数据库失败"
	// 	rspjson, _ := json.Marshal(rsp)
	// 	w.WriteHeader(200)
	// 	w.Write(rspjson)
	// 	return
	// }

	// conn.Do("HSET", gconst.AgentInfo+phone, "userID", userID)

	// replyAddAgentInfo(w, "信息提交成功", true)
}
