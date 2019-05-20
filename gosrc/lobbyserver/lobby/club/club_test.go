package club

import(
	"testing"
	"log"
	"net/http"
	"time"
	"lobbyserver/lobby"
	"io/ioutil"
	"github.com/golang/protobuf/proto"
)

// TestSomething 测试用例
func TestSomething(t *testing.T) {
	log.Println("TestSomething")

	testCreateClub("10000002")
	// testLoadMyClubs("10000002")
	// testDeleteClub("10000002")
	// testLoadClubMembers("10000002")
}



func testCreateClub(id string) {
	tk := lobby.GenTK(id)
	// tk := "vpequ8ELk8xCTPN-heLzghqikggNF85xeH1AyElDSHY="
	var url = "http://localhost:3002/lobby/uuid/createClub?tk="+ tk + "&clname=mytest"
	client := &http.Client{Timeout: time.Second * 60}
	req, err := http.NewRequest("GET", url, nil)

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

	msgClubReply := &MsgClubReply{}
	err = proto.Unmarshal(body, msgClubReply)
	if err != nil {
		log.Println("err:", err)
	}


	createClubReply := &MsgCreateClubReply{}
	buf := msgClubReply.GetContent()

	err = proto.Unmarshal(buf, createClubReply)
	if err != nil {
		log.Println("err:", err)
	}

	log.Println("createClubReply:", createClubReply)
}

func testLoadMyClubs(id string) {
	tk := lobby.GenTK(id)
	// tk := "vpequ8ELk8xCTPN-heLzghqikggNF85xeH1AyElDSHY="
	var url = "http://localhost:3002/lobby/uuid/loadMyClubs?tk="+ tk
	client := &http.Client{Timeout: time.Second * 60}
	req, err := http.NewRequest("GET", url, nil)

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

	msgClubReply := &MsgClubReply{}
	err = proto.Unmarshal(body, msgClubReply)
	if err != nil {
		log.Println("err:", err)
	}


	reply := &MsgClubLoadMyClubsReply{}
	buf := msgClubReply.GetContent()

	err = proto.Unmarshal(buf, reply)
	if err != nil {
		log.Println("err:", err)
	}

	log.Println("reply:", reply)
}

func testDeleteClub(id string) {
	tk := lobby.GenTK(id)
	// tk := "vpequ8ELk8xCTPN-heLzghqikggNF85xeH1AyElDSHY="
	var url = "http://localhost:3002/lobby/uuid/disbandClub?tk="+ tk + "&clubID=7d8904b8-7ac7-11e9-a192-107b445225b6"
	client := &http.Client{Timeout: time.Second * 60}
	req, err := http.NewRequest("GET", url, nil)

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

	msgClubReply := &MsgClubReply{}
	err = proto.Unmarshal(body, msgClubReply)
	if err != nil {
		log.Println("err:", err)
	}

	if msgClubReply.GetReplyCode() == int32(ClubReplyCode_RCError) {
		genericRely := &MsgCubOperGenericReply{}
		err = proto.Unmarshal(msgClubReply.GetContent(), genericRely)
		if err != nil {
			log.Println("parse error:", err)
		}

		log.Println("errCode:", genericRely.GetErrorCode())
		return
	}

	reply := &MsgClubLoadMyClubsReply{}
	buf := msgClubReply.GetContent()

	err = proto.Unmarshal(buf, reply)
	if err != nil {
		log.Println("err:", err)
	}

	log.Println("reply:", reply)
}

func testLoadClubMembers(id string) {
	tk := lobby.GenTK(id)
	// tk := "vpequ8ELk8xCTPN-heLzghqikggNF85xeH1AyElDSHY="
	var url = "http://localhost:3002/lobby/uuid/loadClubMembers?tk="+ tk + "&clubID=7d8904b8-7ac7-11e9-a192-107b445225b6"
	client := &http.Client{Timeout: time.Second * 60}
	req, err := http.NewRequest("GET", url, nil)

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

	msgClubReply := &MsgClubReply{}
	err = proto.Unmarshal(body, msgClubReply)
	if err != nil {
		log.Println("err:", err)
	}

	if msgClubReply.GetReplyCode() == int32(ClubReplyCode_RCError) {
		genericRely := &MsgCubOperGenericReply{}
		err = proto.Unmarshal(msgClubReply.GetContent(), genericRely)
		if err != nil {
			log.Println("parse error:", err)
		}

		log.Println("errCode:", genericRely.GetErrorCode())
		return
	}

	reply := &MsgClubLoadMembersReply{}
	buf := msgClubReply.GetContent()

	err = proto.Unmarshal(buf, reply)
	if err != nil {
		log.Println("err:", err)
	}

	log.Println("reply:", reply)
}