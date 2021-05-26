package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/go-connections/tlsconfig"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func main() {
	fmt.Println("邮件通知价格阈值：", os.Args[1])
	set, err := strconv.ParseFloat(os.Args[1], 32)
	if err != nil {
		fmt.Println("参数错误")
		return
	}
	tls, _ := tlsconfig.Client(tlsconfig.Options{CAFile: "sina.cn", InsecureSkipVerify: true})
	tr := &http.Transport{
		TLSClientConfig: tls,
	}
	client := &http.Client{Transport: tr}

	c := time.NewTicker(10 * time.Second)
	sendMailTicker := time.NewTicker(30 * time.Minute)
	sendMailTickerFlag := true
	for {
		select {
		case <-c.C:
			t := time.Now()
			if t.Weekday() >= 1 && t.Weekday() <= 5 && ((t.Local().Hour() == 9 && t.Minute() >= 30) || (t.Local().Hour() == 11 && t.Minute() <= 30) ||
				t.Local().Hour() == 10 || (t.Local().Hour() >= 13 && t.Local().Hour() <= 14)) {
				res, err := client.Get("https://hq.sinajs.cn/?_=0.4173090047767789&list=sz000586")
				if err != nil {
					fmt.Println("error occurred!:", err)
					continue
				}
				data, _ := ioutil.ReadAll(transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder()))
				res.Body.Close()
				result := string(data)
				sp := strings.Split(result, ",")
				fmt.Println("汇源通信：", sp[3], " ", sp[len(sp)-3], " ", sp[len(sp)-2])
				cur, _ := strconv.ParseFloat(sp[3], 32)
				if cur > set && sendMailTickerFlag {
					sendMailTickerFlag = false
					mail("汇源通信股票提醒", sp[3])
				}
			}
		case <-sendMailTicker.C:
			sendMailTickerFlag = true
		}
	}
}

func mail(subject string, price string) {
	user := "********@126.com"
	password := "*******"
	host := "smtp.126.com:25"
	to := "******@qq.com"

	// subject := "使用Golang发送邮件"

	body := `
		<html>
		<body>
		<h3>
		汇源通信股票价格超过预设值，` + "当前价格：" + price +
		`</h3>
		</body>
		</html>
		`
	fmt.Println("send email")

	err := SendToMail(user, password, host, to, subject, body, "html")
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}
}

func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}
