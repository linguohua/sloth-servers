package club

// func sendMail(myMail *webdata.Mail) {
// 	buf, err := json.Marshal(myMail)
// 	if err != nil {
// 		log.Println("marshal mail error")
// 		return
// 	}
// 	// 注意签名要用到signKey
// 	var sign = genSign(fmt.Sprintf("%d", myMail.GameID), fmt.Sprintf("%d", myMail.PlayerID), fmt.Sprintf("%d", myMail.Type), myMail.Subject, myMail.Title, myMail.ExpirationTime)
// 	var url = config.MailServer + "/mail/addmail?sign=" + sign

// 	client := &http.Client{Timeout: time.Second * 60}
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Println("err: ", err)
// 		return
// 	}

// 	if resp.StatusCode != 200 {
// 		log.Println("resp.StatusCode != 200")
// 		return
// 	}

// 	errcode := resp.Header.Get("error")
// 	if errcode != "" {
// 		log.Println("errorcode: ", errcode)
// 		return
// 	}

// 	body := make([]byte, 1024)
// 	read, err := resp.Body.Read(body)
// 	if err != nil && read < 1 {
// 		log.Println("read message body err: ", err)
// 		return
// 	}
// 	body = body[:read]
// 	log.Println("msg: ", string(body))
// }

func sendClubMail(msg string, userID string) {
	// TODO: llwant mysql
	// userIDInt64, _ := strconv.ParseInt(userID, 10, 32)
	// var myMail = &webdata.Mail{}
	// myMail.GameID = config.LobbyID
	// //用户ID
	// myMail.Type = 1
	// myMail.Subject = "俱乐部邮件"
	// myMail.Title = "俱乐部"
	// myMail.Text = msg
	// myMail.ExpirationTime = "2018-04-04 14:25:22"
	// myMail.PlayerID = userIDInt64

	// sendMail(myMail)
}
