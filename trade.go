package dfTrade

import (
	"dfTrade/util"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const (
	pubKey = "-----BEGIN PUBLIC KEY-----\n" +
		"MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDHdsyxT66pDG4p73yope7jxA92\n" +
		"c0AT4qIJ/xtbBcHkFPK77upnsfDTJiVEuQDH+MiMeb+XhCLNKZGp0yaUU6GlxZdp\n" +
		"+nLW8b7Kmijr3iepaDhcbVTsYBWchaWUXauj9Lrhz58/6AE/NF0aMolxIGpsi+ST\n" +
		"2hSHPu3GSXMdhPCkWQIDAQAB\n" +
		"-----END PUBLIC KEY-----"
)
const webUrl = "https://jywg.18.cn"

type Trade struct {
	userId       string
	password     string
	identifyCode string
	randNum      string
	validateKey  string
	cookies      []*http.Cookie
}

func (this *Trade) perUrl(p string) string {
	uri, _ := url.Parse(webUrl)
	v := url.Values{}
	v.Add("validatekey", this.validateKey)
	uri.RawQuery = v.Encode()
	uri.Path = path.Join(uri.Path, p)
	return uri.String()
}

func (this *Trade) get(url string, queryData url.Values) (string, error) {

	req, err := http.NewRequest("GET", url, strings.NewReader(queryData.Encode()))
	if err != nil {
		return "", err
	}
	req.URL.RawQuery = queryData.Encode()

	fmt.Println(req.URL.String())
	//req.Header.Set("User-Agent", "同花顺/7.0.10 CFNetwork/1333.0.4 Darwin/21.5.0")
	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), err

}

func (this *Trade) post(url string, postData url.Values) (string, error) {

	req, err := http.NewRequest("POST", url, strings.NewReader(postData.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "jywg.18.cn")
	req.Header.Set("Origin", "https://jywg.18.cn")

	for _, cookie := range this.cookies {
		req.AddCookie(cookie)
	}

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if len(resp.Cookies()) != 0 {
		this.cookies = append(this.cookies, resp.Cookies()...)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func (this *Trade) yzmOcr() {

	randNum := "0." + util.RandDigitString(15)
	yzmUrl := fmt.Sprintf("https://jywg.18.cn/Login/YZM?randNum=%s", randNum)

	postData := url.Values{}
	postData.Set("url", yzmUrl)
	yzm, err := this.get("http://192.168.1.3:8868/ocr", postData)
	if err != nil {
		log.Println(err)
	}
	log.Println("验证码:", yzm)

	this.randNum = randNum
	this.identifyCode = yzm
}

func (this *Trade) htmlValidateKey() error {
	u := "https://jywg.18.cn/Search/Position"
	postData := url.Values{}
	//postData.Set("revokes", "20220714_589264")
	resp, err := this.post(u, postData)
	if err != nil {
		return err
	}
	index := strings.Index(resp, "<input id=\"em_validatekey\" type=\"hidden\" value=\"")
	if index < 0 {
		return errors.New("无数据")
	}
	if len(resp) < index+36+48 {
		return errors.New("解析失败")
	}
	validateKey := resp[index+48 : index+36+48]
	log.Println("validateKey:", validateKey)
	this.validateKey = validateKey
	return nil
}

func (this *Trade) Login(userId string, password string) error {
	this.userId = userId
	this.password = password
	this.yzmOcr()

	pwd := util.RsaEncode([]byte(this.password), []byte(pubKey))

	postData := url.Values{}
	postData.Set("userId", this.userId)
	postData.Set("password", base64.StdEncoding.EncodeToString(pwd))
	postData.Set("randNumber", this.randNum)
	postData.Set("identifyCode", this.identifyCode)
	postData.Set("duration", "1800")
	postData.Set("authCode", "")
	postData.Set("type", "Z")

	resp, err := this.post("https://jywg.18.cn/Login/Authentication", postData)
	if err != nil {
		return err
	}
	status := gjson.Get(resp, "Status")
	if status.String() != "0" {
		message := gjson.Get(resp, "Message")
		return errors.New(message.String())
	}
	if err = this.htmlValidateKey(); err != nil {
		return err
	}
	log.Println("登录成功")
	return nil
}

func (this *Trade) QueryAssetAndPositionV1() {
	u := "https://jywg.18.cn/Com/queryAssetAndPositionV1"
	postData := url.Values{}
	postData.Set("moneyType", "RMB")
	resp, err := this.post(u, postData)
	if err != nil {
		log.Println(err)

	}
	//fmt.Println(resp)

	fmt.Println("总金额:", gjson.Get(resp, "Data").Array()[0].Get("Zzc").Float())
	fmt.Println("总市值:", gjson.Get(resp, "Data").Array()[0].Get("Zxsz").Float())
	fmt.Println("可用资金:", gjson.Get(resp, "Data").Array()[0].Get("Kyzj").Float())
	fmt.Println("可取资金:", gjson.Get(resp, "Data").Array()[0].Get("Kqzj").Float())
	fmt.Println("资金余额:", gjson.Get(resp, "Data").Array()[0].Get("Zjye").Float())
	fmt.Println("冻结资金:", gjson.Get(resp, "Data").Array()[0].Get("Djzj").Float())
	fmt.Println("冻结资金:", gjson.Get(resp, "Data").Array()[0].Get("Dryk").Float())
	fmt.Println("持仓盈亏:", gjson.Get(resp, "Data").Array()[0].Get("Ljyk").Float())
	fmt.Println("持仓:", gjson.Get(resp, "Data").Array()[0].Get("positions").Array())
	//fmt.Println(resp, err)
}

/*
SendOrder
amount: 单位股 最低100股
tradeType: 买卖 B/S
*/
func (this *Trade) SendOrder(code string, market string, price decimal.Decimal, amount int, tradeType string) (string, error) {
	u := this.perUrl("Trade/SubmitTradeV2")
	postData := url.Values{}
	postData.Set("stockCode", code)
	postData.Set("price", price.String())
	postData.Set("amount", fmt.Sprintf("%d", amount))
	postData.Set("tradeType", tradeType)
	//postData.Set("zqmc", "中概互联网ETF")
	//postData.Set("gddm", "")
	postData.Set("market", market)
	resp, err := this.post(u, postData)

	return resp, err

}

func (this *Trade) GetRevokeList() {
	u := this.perUrl("Trade/GetRevokeList")
	postData := url.Values{}
	resp, err := this.post(u, postData)
	if err != nil {
		log.Println(err)
	}
	/*
		{"Status":0,"Count":1,"Data":[{"Htxh":"B250022392","Wtbh":"1412433"}],"Errcode":0}
		Wtbh:委托编号
	*/
	fmt.Println(resp)

}

/*
CancelOrder
撤单
*/
func (this *Trade) CancelOrder(revokes string) (string, error) {
	u := this.perUrl("Trade/RevokeOrders")
	postData := url.Values{}
	//postData.Set("revokes", "20220714_589204,20220714_589205")
	postData.Set("revokes", revokes)
	return this.post(u, postData)
}

func (this *Trade) Positions() {

}

func NewDFTrade() *Trade {
	return &Trade{}
}
