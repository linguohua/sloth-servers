package share

import (

)

// TestSomething 测试用例
func TestSomething(t *testing.T) {




}


func testGetShareInfo(id string) {
	tk := lobby.GenTK(id)
	// tk := "vpequ8ELk8xCTPN-heLzghqikggNF85xeH1AyElDSHY="
	var url = "http://localhost:3002/lobby/uuid/getShareInfo?tk="+ tk + "&sence=1&mediaType=1&"
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