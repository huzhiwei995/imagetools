package main

import (
	"fmt"
	"os"
	"log"
	"strconv"
	"time"
	"github.com/gocolly/colly"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"io/ioutil"
	"io"
	"bytes"
	"flag"
)
var CountNum int
var logger *log.Logger
var logFile *os.File
var imgPath string
var l int
func init() {
	filePath := "./imageLog.txt"
	l = 0
	logFile,err := os.OpenFile(filePath,os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("创建文件错误 : %v\n", err)
	}
	logger = log.New(logFile,"[INFO]",log.Ldate)
	logger.SetFlags(log.Lshortfile)
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
/**
 * 自定义下载页数
 * 自定义下载地址
 * 自定义下载条件
 */
func main() {
	i := 0
	flag.String("f", "", "使用示例 -i 1 -u D:/image/")
	var Count= flag.Int("i", 1, "指定最大下载页数")
	var imgPathurl= flag.String("u", "./image", "指定图片下载路径，默认当前目录创建image")
	flag.Bool("t", false, "条件搜索，暂不提供")
	flag.Parse()
	CountNum = *Count
	imgPath = *imgPathurl
	exist, err := PathExists(imgPath)
	if err != nil {
		fmt.Println(err)
	}
	if !exist {
		err := os.Mkdir(imgPath, os.ModePerm)
		if err != nil {
			fmt.Printf("目录创建失败![%v]\n", err)
		} else {
			fmt.Printf("目录创建成功!\n")
		}
	}
	for i <= CountNum{
		err := work("https://anime-pictures.net/pictures/view_posts/"+strconv.Itoa(i)+"?aspect=16%3A9&order_by=date&ldate=0&lang=en")
		if err != nil {
			fmt.Println("连接错误：",err)
			break
		}
		i++
		time.Sleep(2*time.Second)
	}
	defer logFile.Close()
}
func work(url string) (err error) {
	Info("url是：",url)
	fmt.Println("url是：",url)
	c := colly.NewCollector(colly.MaxDepth(2))
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36 Edge/16.16299"
	c.OnHTML("*", func(e *colly.HTMLElement) {
		e.DOM.Find(".img_block_big").Each(func(i int, cs *goquery.Selection) {
			a,_ := cs.Find("a").Eq(0).Attr("href")
			c.Visit(e.Request.AbsoluteURL(a))
		})
	})
	c.OnHTML(".download_icon", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		link = "https://anime-pictures.net"+link
		name,err := getImg(link)
		if err != nil {
			Info(name,"下载失败,原因：",err)
		}
		fmt.Println(name,"下载成功")
		l++
	})
	/*c.OnResponse(func(r *colly.Response) {
		CountNum++
	})*/
	c.OnError(func(r *colly.Response, err error) {
		Info("返回的URL:", r.Request.URL, "错误返回:", r, "\nError:", err)
		fmt.Println("返回的URL:", r.Request.URL, "错误返回:", r, "\nError:", err)
	})
	c.Visit(url)
	return nil
}
func getImg(urlimg string) (string, error) {
	imgPathend := imgPath[len(imgPath)-1:]
	if imgPathend == "/"{
		imgPath = imgPath[0:len(imgPath)-1]
	}
	name := imgPath+"/H_"+strconv.Itoa(l)+".jpg"
	out, err := os.Create(name)
	if err != nil{
		Info("创建文件失败",err)
		fmt.Println("创建文件失败",err)
		return name,err
	}
	defer out.Close()
	resp, err := http.Get(urlimg)
	if err != nil{
		Info("请求图片url失败",err)
		fmt.Println("请求图片url失败",err)
		return name,err
	}
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		Info("读取文件失败",err)
		fmt.Println("读取文件失败",err)
		return name,err
	}
	_, err = io.Copy(out, bytes.NewReader(pix))
	if err != nil {
		Info("下载文件失败",err)
		fmt.Println("下载文件失败",err)
		return name,err
	}
	return name,err
}
func Info(v ...interface{}) {
	logger.SetPrefix("[Info]")
	logger.Println(v)
}