package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	spider_base_url     string = "http://www.jikexueyuan.com/path/docker/"
	login_url           string = "http://yjs.whu.edu.cn/ssfw/login.jsp"
	verify_code_url     string = "http://yjs.whu.edu.cn//ssfw/captcha.do"
	post_login_info_url string = "http://yjs.whu.edu.cn/ssfw/j_spring_ids_security_check"
	index_url           string = "http://yjs.whu.edu.cn/ssfw/pygl/xkgl/xskb.do?timetip="
	username            string = "2020202090029"
	password            string = "213843"
)

func login() []byte {
	////获取登陆界面的cookie
	//c := &http.Client{}
	//req, _ := http.NewRequest("GET", login_url, nil)
	//res, _ := c.Do(req)
	//req.URL, _ = url.Parse(verify_code_url)
	//var temp_cookies = res.Cookies()
	//for _, v := range res.Cookies() {
	//	req.AddCookie(v)
	//}
	//// 获取验证码
	//var verify_code string
	//for {
	//	res, _ = c.Do(req)
	//	file, _ := os.Create("verify.gif")
	//	io.Copy(file, res.Body)
	//	fmt.Println("请查看verify.gif， 然后输入验证码， 看不清输入0重新获取验证码")
	//	fmt.Scanf("%s", &verify_code)
	//	if verify_code != "0" {
	//		break
	//	}
	//	res.Body.Close()
	//}
	//
	//{
	//	value := "mode=db&j_username=" + username + "&j_password=" + password + "&validateCode=" + verify_code
	//	req, _ = http.NewRequest("POST", post_login_info_url, strings.NewReader(value))
	//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	req.Header.Set("Connection", "keep-alive")
	//	req.Header.Set("Host", "yjs.whu.edu.cn")
	//	req.Header.Set("Referer", login_url)
	//	for _, v := range temp_cookies {
	//		req.AddCookie(v)
	//	}
	//	resp, _ := c.Do(req)
	//	data, _ := ioutil.ReadAll(res.Body)
	//	fmt.Println(string(data))
	//	resp.Body.Close()
	//}

	{
		c1 := &http.Client{}
		timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		req1, _ := http.NewRequest("GET", index_url+timestamp, nil)
		req1.Header.Set("Accept", "text/html; charset=UTF-8")
		req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req1.Header.Set("Connection", "keep-alive")
		req1.Header.Set("Host", "yjs.whu.edu.cn")
		req1.Header.Set("Referer", "http://yjs.whu.edu.cn/ssfw/index.do")
		req1.Header.Set("Cookie", "JSESSIONID=0000H5Pna-7Ir6SxnzfbDsJBO1P:18u2sfi06; iPlanetDirectoryPro=izRMff1vkjsSXk1CQ2fiSe")
		//for _, v := range temp_cookies {
		//	req1.AddCookie(v)
		//}
		//req1.AddCookie(&http.Cookie{Raw: "iPlanetDirectoryPro=Hu1uAK6IIdoIZVOM2Y1dY5; JSESSIONID=0000RpHqiHQr1RLY4Av7SameSSo:18u2sfi06"})
		res1, _ := c1.Do(req1)
		data, _ := ioutil.ReadAll(res1.Body)

		res1.Body.Close()

		return data
	}

}

type ClassInfo struct {
	ClassName    string
	ClassTime    string
	ClassPos     string
	ClassTeacher string
}

func main() {
	htmlData := login()
	//fmt.Println(string(htmlData))
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(htmlData))
	tempH, _ := doc.Html()
	str := strings.Replace(tempH, "<br/>", "\n", -1)
	doc, _ = goquery.NewDocumentFromReader(strings.NewReader(str))

	var classList string
	doc.Find(".t_con").Children().Each(func(i int, selection *goquery.Selection) {
		//if !selection.Is("br") {
		txt := selection.Text()
		if txt != " " && !(strings.Contains(txt, "第 ") && strings.Contains(txt, " 节")) &&
			txt != "上午" && txt != "中午" && txt != "下午" && txt != "晚上" {
			classList += "\n" + selection.Text()
		}
		//}
	})
	clist := strings.Split(classList, "\n")

	i := 0
	t := &ClassInfo{}
	classInfoList := make([]*ClassInfo, 0)
	for _, c := range clist {
		if c == "" {
			continue
		}
		switch i % 3 {
		case 0:
			t = &ClassInfo{}
			t.ClassName = c
		case 1:
			t.ClassTime = c
		case 2:
			tt := strings.Split(c, " ")
			t.ClassPos = tt[0]
			t.ClassTeacher = tt[2]
			classInfoList = append(classInfoList, t)
		}
		i++
	}
	fmt.Println(classInfoList)

}
