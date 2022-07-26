package weibo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
)

const HEADER = "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"

func getTID() (string, error) {
	c := http.Client{}
	r, _ := http.NewRequest(http.MethodGet, "https://passport.weibo.com/visitor/genvisitor?cb=gen_callback&fp={\"os\":\"1\",\"browser\":\"Chrome70,0,3538,25\",\"fonts\":\"undefined\",\"screenInfo\":\"1920*1080*24\",\"plugins\":\"\"}", nil)
	r.Header.Set("User-Agent", HEADER)
	resp, _ := c.Do(r)
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("request error")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	bodyText := string(body)
	tidIndex := strings.Index(bodyText, "tid")
	if tidIndex == -1 {
		return "", errors.New("err")
	}
	newTidIndex := strings.Index(bodyText, "new_tid")
	return bodyText[tidIndex+6 : newTidIndex-3], nil
}

func getCookie(tid string) (string, string, error) {
	c := http.Client{}
	r, _ := http.NewRequest(http.MethodGet, "https://passport.weibo.com/visitor/visitor?a=incarnate&t="+tid+"&w=3&c&cb=restore_back&from=weibo", nil)
	r.Header.Set("User-Agent", HEADER)
	resp, _ := c.Do(r)
	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New("err")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	bodyText := string(body)
	subIndex := strings.Index(bodyText, "sub")
	subpIndex := strings.Index(bodyText, "subp")
	if subIndex == -1 {
		return "", "", errors.New("err")
	}
	if len(bodyText[subIndex+6:subpIndex-3]) < 5 {
		return "", "", errors.New("err")
	}
	return bodyText[subIndex+6 : subpIndex-3], bodyText[subpIndex+7 : len(bodyText)-5], nil
}

func GetTrending() string {
	var sub string
	var subp string
	for i := 0; i < 10; i++ {
		tid, err := getTID()
		if err != nil {
			time.Sleep(time.Microsecond * 100)
			continue
		}
		sub, subp, err = getCookie(tid)
		if err != nil {
			time.Sleep(time.Microsecond * 100)
			continue
		}
		break
	}

	// fmt.Println(sub, subp)
	c := http.Client{}
	r, _ := http.NewRequest(http.MethodGet, "https://s.weibo.com/top/summary/summary", nil)
	r.Header.Set("User-Agent", HEADER)
	r.AddCookie(&http.Cookie{Name: "SUB", Value: sub})
	r.AddCookie(&http.Cookie{Name: "SUBP", Value: subp})
	var resp *http.Response
	for i := 0; i < 10; i++ {
		resp, _ = c.Do(r)
		if resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(time.Microsecond * 100)
	}

	if resp.StatusCode != http.StatusOK {
		return "获取微博热搜失败"
	}

	// bo, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(bo))
	doc, _ := htmlquery.Parse(resp.Body)
	// fmt.Println(doc)
	list := htmlquery.Find(doc, "//ul/li")
	for i, n := range list {
		if i < 5 {
			continue
		}
		removeSpace := strings.TrimSpace(htmlquery.InnerText(n))
		if strings.Contains(removeSpace, "•") {
			continue
		}
		fmt.Println(removeSpace)
	}
	// fmt.Println(list)
	return ""
}
