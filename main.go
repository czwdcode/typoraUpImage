package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

type config struct { //json结构体成员需要大写开头,原因大写外部可见json是另一个包
	Message, Branch, Token, UserName, Repositorie, Folder, BucketDomain string
}

// Gitee Post提交数据结构体.
type Gitee struct {
	Access_token string `json:"access_token"`
	Content      string `json:"content"`
	Message      string `json:"message"`
	Branch       string `json:"branch"`
}

type reJsoncc struct {
	//ReJson返回参数结构体解析
	Content struct {
		Download_url string `json:"download_url"`
	} `json:"content"`
}

func main() {
	rand.Seed(time.Now().UnixNano())      //置随机种子
	picPathSet := os.Args[1:]             //获得命令行参数,初始化文件路径,args第一个是本身路径不需要
	ch := make(chan int, len(picPathSet)) // 创建通道容量为传进来的参数数量
	//fmt.Println(picPathSet)
	//picPathSet = append(picPathSet, "D:\\goStudy\\picUp\\1611217821.png")
	//picPathSet = append(picPathSet, "D:\\goStudy\\picUp\\1611217821.png")
	jsonData := new(config) //创建结构体指针
	//conPath, _ := os.Getwd() // 获取当前运行路径
	//conPath = conPath + "/config.json" // 组合配置文件路径
	//fmt.Println(conPath)
	configData, err := ioutil.ReadFile("./config.json") //读出8字节切片数据
	if err != nil {
		//defer fmt.Println("configErr")
		//defer fmt.Println(conPath)
		panic(err) //读取不到文件直接退出
	}
	err = json.Unmarshal(configData, jsonData) //json反序列化转换成结构体(源,目标)
	if err != nil {
		panic(err) //反序列化转换失败
	}
	contentType := "application/json;charset=UTF-8" //构建contentType参数
	urlPost := jsonData.BucketDomain + jsonData.UserName + "/" +
		jsonData.Repositorie + "/" + "contents/" + jsonData.Folder + "/" //url拼接
	//fmt.Println("urlpost--", urlPost)
	// if len(picPathSet) > 0 { //typora 第一行必须打印 Upload Success:,通过判断 picpathset 文件判断是否有文件传入
	// 	fmt.Println("Upload Success:")
	// }
	for i, Path := range picPathSet { //启动go程
		go upPic(Path, &urlPost, &contentType, jsonData, i, ch) //(文件路径,地址指针,结构体指针)
	}
	// 用Channel判断,是否 运行完毕
	for range picPathSet {
		<-ch
		//select {
		//case <-ch:
		//case <-time.After(time.Second * 10):// 避免超时
		//
		//}
	}
	close(ch) //关闭通道
	//time.Sleep(3 * time.Second) //延时命令防止主go程结束
}

func upPic(picPath string, urlPost, contentType *string, jsonData *config, i int, ch chan int) {
	//fmt.Println(*urlPost)
	// 保证通道有数据写入避免死锁用defer是因为执行完毕才写入,defer后面必须跟函数
	defer func() {
		ch <- i
	}()
	strsp := picPath[:4]
	if strsp == "http" {
		fmt.Println("Already uploaded")
		runtime.Goexit()
	}
	base64Pic := imagesToBase64(picPath) //image转base64后返回的数据
	//fmt.Println(base64Pic)
	pathSuffix := filepath.Ext(picPath) //文件后缀名,用filepath因为path只能用于反斜杠路径
	urlPostCompelet := *urlPost + time.Now().Format("20060102") +
		strconv.Itoa(rand.Intn(1000000)) + pathSuffix //完整的url
	//fmt.Println(urlPostCompelet)

	data := new(Gitee) //构造结构体变量为一个指针
	//初始化结构体
	data.Access_token = jsonData.Token
	data.Content = base64Pic
	data.Message = jsonData.Message
	data.Branch = jsonData.Branch
	//初始化结束
	//fmt.Println(*data)
	rawUrl := Post(&urlPostCompelet, data, contentType) //调用post发送数据(post地址指针,post数据,post结构头)
	//fmt.Println(result)
	fmt.Println(rawUrl) //打印返回的图片源地址
}

func imagesToBase64(path string) string { //base64转换

	picDate, err := ioutil.ReadFile(path) //读取文件数据
	if err != nil {                       //判断是否出错,出错记录错误
		fmt.Println("ReadFile", err)
		runtime.Goexit()
	}
	base64Pic := base64.StdEncoding.EncodeToString(picDate)
	return base64Pic
}

// Post 函数
func Post(url *string, data interface{}, contentType *string) string {
	// 超时时间：5秒.
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(data) //结构体转换json数据
	//fmt.Println(*url)
	resp, err := client.Post(*url, *contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("postErr", err)
		runtime.Goexit()
		//panic(err)
	}
	result, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close() // 读完需要关闭防止内存泄露
	reJsonzz := new(reJsoncc)
	_ = json.Unmarshal(result, reJsonzz)
	return reJsonzz.Content.Download_url // 返回图片的源地址
}
