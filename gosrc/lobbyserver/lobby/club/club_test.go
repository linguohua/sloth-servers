package club

import(
	"testing"
	"log"
	"net/http"
	"time"
	"lobbyserver/lobby"
	"io/ioutil"
	"github.com/golang/protobuf/proto"
	"bytes"
)

// TestSomething 测试用例
func TestSomething(t *testing.T) {
	log.Println("TestSomething")

	// testCreateClub("10000002")
	// testLoadMyClubs("10000002")
	// testDeleteClub("10000002")
	// testLoadClubMembers("10000002")
	// testJoinClub("10000003")
	// testLoadClubEvent("10000002")
	// testJoinApproval("10000002", "10000003", "yes", "5")
	// testClubQuit("10000003")
	// testLoadMyClubs("10000003")
	testCreateClubRoom("10000002")
	// testLoadClubRooms("10000002")
	// testDeleteClubRoom("10000002")


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
	var url = "http://localhost:3002/lobby/uuid/disbandClub?tk="+ tk + "&clubID=9949cd58-7e97-11e9-a192-107b445225b6"
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
	var url = "http://localhost:3002/lobby/uuid/loadClubMembers?tk="+ tk + "&clubID=6b512ef0-7b77-11e9-a192-107b445225b6"
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

func testJoinClub(id string) {
	tk := lobby.GenTK(id)
	// tk := "vpequ8ELk8xCTPN-heLzghqikggNF85xeH1AyElDSHY="
	var url = "http://localhost:3002/lobby/uuid/joinClub?tk="+ tk + "&clubNumber=24367"
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

	// reply := &MsgClubLoadMembersReply{}
	// buf := msgClubReply.GetContent()

	// err = proto.Unmarshal(buf, reply)
	// if err != nil {
	// 	log.Println("err:", err)
	// }
	buf := msgClubReply.GetContent()

	log.Println("reply:", string(buf))
}

func testLoadClubEvent(id string) {
	tk := lobby.GenTK(id)
	// tk := "vpequ8ELk8xCTPN-heLzghqikggNF85xeH1AyElDSHY="
	var url = "http://localhost:3002/lobby/uuid/loadClubEvents?tk="+ tk + "&clubID=6b512ef0-7b77-11e9-a192-107b445225b6"
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

	reply := &MsgClubLoadEventsReply{}
	buf := msgClubReply.GetContent()

	err = proto.Unmarshal(buf, reply)
	if err != nil {
		log.Println("err:", err)
	}

	log.Println("reply:", reply)
}

func testJoinApproval(id string, applicantID string, agree string, eID string) {
	tk := lobby.GenTK(id)
	// tk := "vpequ8ELk8xCTPN-heLzghqikggNF85xeH1AyElDSHY="

	var url = "http://localhost:3002/lobby/uuid/joinApproval?tk="+ tk +
	"&clubID=3ca992e8-7e99-11e9-a192-107b445225b6&applicantID="+applicantID+"&agree="+ agree +"&eID=" + eID
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

	buf := msgClubReply.GetContent()

	log.Println("reply:", string(buf))
}

func testClubQuit(id string) {
	tk := lobby.GenTK(id)
	// tk := "vpequ8ELk8xCTPN-heLzghqikggNF85xeH1AyElDSHY="

	var url = "http://localhost:3002/lobby/uuid/quitClub?tk="+ tk + "&clubID=6b512ef0-7b77-11e9-a192-107b445225b6"
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

	buf := msgClubReply.GetContent()
	if len(buf) == 0 {
		log.Println("len(buf) == 0")
		return
	}

	reply := &MsgClubLoadMyClubsReply{}
	err = proto.Unmarshal(buf, reply)
	if err != nil {
		log.Println("err:", err)
	}

	log.Println("reply:", reply)
}

func testCreateClubRoom(id string) {
	tk := lobby.GenTK(id)
	// tk := "vpequ8ELk8xCTPN-heLzghqikggNF85xeH1AyElDSHY="
	config := `{"playerNumAcquired":4, "payNum":0, "payType":0, "handNum":4, "roomType":1, "modName":"game1"}`
	createRoomReq := &lobby.MsgCreateRoomReq{}
	createRoomReq.Config = &config

	buf, err := proto.Marshal(createRoomReq)
	if err != nil {
		log.Println("testCreateClubRoom, error:", err)
		return
	}

	var url = "http://localhost:3002/lobby/uuid/createClubRoom?tk="+ tk + "&clubID=9949cd58-7e97-11e9-a192-107b445225b6"
	client := &http.Client{Timeout: time.Second * 60}
	req, err := http.NewRequest("OPTIONS", url,  bytes.NewBuffer(buf))

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

	reply := &lobby.MsgCreateRoomRsp{}
	err = proto.Unmarshal(body, reply)
	if err != nil {
		log.Println("testCreateClubRoom, err:", err)
		return
	}

	if reply.GetResult() != 0 {
		log.Printf("reply errCode:%d, retMsg:%s", reply.GetResult(), reply.GetRetMsg())
	} else {
			log.Println("reply:", reply)
	}
	// log.Println("reply:", reply)
}

func testLoadClubRooms(id string) {
	tk := lobby.GenTK(id)

	var url = "http://localhost:3002/lobby/uuid/loadClubRooms?tk="+ tk + "&clubID=9949cd58-7e97-11e9-a192-107b445225b6"
	client := &http.Client{Timeout: time.Second * 60}
	req, err := http.NewRequest("POST", url,  nil)

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

	reply := &lobby.MsgLoadRoomListRsp{}
	err = proto.Unmarshal(body, reply)
	if err != nil {
		log.Println("testCreateClubRoom, err:", err)
		return
	}

	if reply.GetResult() != 0 {
		log.Printf("reply errCode:%d, retMsg:%s", reply.GetResult(), reply.GetRetMsg())
	} else {
			log.Println("reply:", reply)
	}
	// log.Println("reply:", reply)
}

func testDeleteClubRoom(id string) {
	tk := lobby.GenTK(id)

	var url = "http://localhost:3002/lobby/uuid/deleteClubRoom?tk="+ tk + "&clubID=5fab2ce0-7e06-11e9-a192-107b445225b6&roomID=2d4958eb-5162-429d-beb8-0d81509ad89c"
	client := &http.Client{Timeout: time.Second * 60}
	req, err := http.NewRequest("POST", url,  nil)

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

	reply := &lobby.MsgDeleteRoomReply{}
	err = proto.Unmarshal(body, reply)
	if err != nil {
		log.Println("testDeleteClubRoom, err:", err)
		return
	}

	if reply.GetResult() != 0 {
		log.Printf("reply errCode:%d", reply.GetResult())
	} else {
			log.Println("reply:", reply)
	}
	// log.Println("reply:", reply)
}
