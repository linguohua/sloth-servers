package auth

import (
	"testing"
	"log"
	"net/http"
	"time"
	"lobbyserver/lobby"
	"io/ioutil"
	"github.com/golang/protobuf/proto"
	"bytes"
)

func TestAuth(t *testing.T) {
	testWxLogin()
}

func testWxLogin() {

	var code = "033AOtuJ1cpAF20PqEtJ1MTAuJ1AOtul"
	var encrypteddata = "PQE/LuDxKMl69olmaZr0IlfwIRwlxOCSdSHLuKcHUV+JUp1Mr4tJeSJd1arKlt2UJYzD3tho3IQGtAF9q0hb8yv95oydQuAcIPcitYrf6gtnrYIG4dFghHrZATld9ACTVVE6fQlAKzSymKPRbhH82iMgUXj6D7sXYUDHQnS585vf0ykD9Vn95YyhXwSSz8FGkIyss3HtWPF+ydtCvbsb03ab9SqfAe001GZZw1A7Dt4F2ct7I1FPMbF887fUaOSZLJjQJkphXqgBvf3Gi3MOPQi7mZ8qUlK9gDsaC5aUZvHN7KvK0h8edlv1YJxP+5aecQGehNO9DJPOzXEXJXsUCWKm/WxJQZrI3Gqe4JyffEl++DsHkkoe/zaJ3O490aRLYvLGGuUaUtBEzklJsIKRv1cmOKbqb8oGoEmP/2oC8HvxkKcoqpgp3cScoSrVxjAmPG20mxFuTNWJioMkHk9vKhOH5W5jrN0ATijfirsuh10="
	var iv = "f4IrtZDL8ht8lDz6NIqTlQ=="

	var url = "http://localhost:3004/lobby/uuid/wxLogin"

	wxLogin := &lobby.MsgWxLogin{}
	wxLogin.Code = &code
	wxLogin.Encrypteddata = &encrypteddata
	wxLogin.Iv = &iv

	buf, err := proto.Marshal(wxLogin)
	if err != nil {
		log.Println("err:", err)
		return
	}

	client := &http.Client{Timeout: time.Second * 60}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))

	resp, err := client.Do(req)
	if err != nil {
		log.Println("err: ", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Println("resp.StatusCode != 200, resp.StatusCode:", resp.StatusCode)
		return
	}

	errcode := resp.Header.Get("error")
	if errcode != "" {
		log.Println("errorcode: ", errcode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("handlerChat error:", err)
		return
	}

	msgLoginReply := &lobby.MsgLoginReply{}
	err = proto.Unmarshal(body, msgLoginReply)
	if err != nil {
		log.Println("err:", err)
	}


	log.Println("msgLoginReply:", msgLoginReply)
}