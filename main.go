package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/rxxcy/image-download/skk"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// Pwd 获取当前路径
func Pwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln("获取基础路径失败：", err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// Mkdir 创建文件夹
func Mkdir(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, 0766)
		if err != nil {
			return false
		}
		return true
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0766)
		if err != nil {
			return false
		}
		return true
	}
	return true
}

func Downloads(title string, url string) {
	var index = 1
	x := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"))
	_ = x.Limit(&colly.LimitRule{DomainGlob: "*.xiezhen.*", Parallelism: 5})
	x.OnHTML(".article-content img", func(e *colly.HTMLElement) {
		src := e.Attr("src")
		path := fmt.Sprintf("%s/%s", baseDir, title)
		Mkdir(path)
		extras := regExt.FindStringSubmatch(src)
		if len(extras) < 1 {
			skk.RedOnly("未知的文件格式")
			return
		}
		ext := extras[1]
		f := fmt.Sprintf("%d", index)
		if index < 10 {
			f = fmt.Sprintf("0%s", f)
		}
		fileName := fmt.Sprintf("%s/%s.%s", path, f, ext)
		fmt.Printf("\r下载第 %s 张", f)
		image := Image{Src: src, Path: fileName}
		CreateDownloadTask(image, ch)
		index++
	})

	x.OnError(func(r *colly.Response, err error) {
		log.Fatal("err: ", err.Error())
	})

	x.OnRequest(func(r *colly.Request) {
		skk.Blue("\n下载图集", title)
	})

	err := x.Visit(url)
	if err != nil {
		log.Println(err.Error())
		return
	}

}

func Download(url string, path string) {
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("referer", "https://www.xiezhen.xyz/")
	request.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")

	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("错误: %s \n", err.Error())
		}
	}(response.Body)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("下载文件错误: ", err.Error())
	}
	_ = ioutil.WriteFile(path, body, 0755)
}

func Bootstrap() {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"))
	_ = c.Limit(&colly.LimitRule{DomainGlob: "*.xiezhen.*", Parallelism: 5})

	c.OnHTML(".excerpt h2", func(e *colly.HTMLElement) {
		title := e.Text
		url := e.ChildAttrs("a", "href")[0]
		Downloads(title, url)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Fatal("err: ", err.Error())
	})

	c.OnRequest(func(r *colly.Request) {
		skk.RedOnly(fmt.Sprintf("\n第 %d 页", page))
	})
	for page < 2 {
		url := fmt.Sprintf("https://www.xiezhen.xyz/page/%d", page)
		err := c.Visit(url)
		if err != nil {
			log.Println(err.Error())
			return
		}
		page++
	}
}

func DownloadPool(ch chan Image) {
	for image := range ch {
		Download(image.Src, image.Path)
		wg.Done()
	}
}

func CreateDownloadTask(image Image, ch chan Image) {
	wg.Add(1)
	ch <- image
}

var baseDir string
var page = 1
var regExt = regexp.MustCompile("\\.(png|jpeg|jpg|pjpg)$")
var client = &http.Client{}
var wg = sync.WaitGroup{}
var ch chan Image
var goCnt = 10

type Image struct {
	Src  string
	Path string
}

func main() {
	fmt.Printf(`

	图片下载器 🎉
	@rxxcy

`)

	tempBaseDir := Pwd()
	skk.Magenta("默认文件保存路径", tempBaseDir+"/images/")
	fmt.Printf("自定义保存路径, 回车默认: ")
	_, _ = fmt.Scanf("%s", &baseDir)

	if baseDir == "" {
		baseDir = tempBaseDir + "/images/"
	} else {
		baseDir = baseDir + "/images/"
	}
	skk.Blue("保存路径", baseDir)
	skk.Blue("默认协程数", fmt.Sprintf("%d", goCnt))
	ch = make(chan Image)
	for i := 0; i < goCnt; i++ {
		go DownloadPool(ch)
	}

	//Bootstrap()
	wg.Wait()
}
