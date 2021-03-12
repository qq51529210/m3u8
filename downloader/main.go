package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/qq51529210/m3u8"
)

type downloader struct {
	sync.WaitGroup
	url     string
	dir     string
	urlDir  *url.URL
	routine int
	task    chan string
}

func (d *downloader) Download() error {
	// 创建目录
	err := os.MkdirAll(d.dir, os.ModePerm)
	if err != nil {
		return err
	}
	// url dir
	d.urlDir, err = url.Parse(d.url)
	if err != nil {
		return err
	}
	d.urlDir.Path = path.Dir(d.urlDir.Path)
	// 下载list文件
	ts, err := d.downloadList()
	if err != nil {
		return err
	}
	fmt.Println("download m3u8 list", d.url)
	// 并发
	if d.routine < 1 {
		d.routine = 1
	}
	// 任务
	d.task = make(chan string, d.routine)
	for i := 0; i < d.routine; i++ {
		d.Add(1)
		go d.downloadTSRoutine()
	}
	// 分配
	for _, s := range ts {
		d.task <- s
	}
	close(d.task)
	d.Wait()
	return nil
}

func (d *downloader) downloadList() ([]string, error) {
	// 下载
	rs, err := http.Get(d.url)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()
	if rs.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download list '%s' http status code '%d'", d.url, rs.StatusCode)
	}
	// 保存
	filePath := filepath.Join(d.dir, path.Base(d.url))
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	_, err = io.Copy(file, rs.Body)
	if err != nil {
		return nil, err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	// 读取每一行
	ts := make([]string, 0)
	reader := m3u8.NewReader(file, nil)
	for {
		line, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return ts, nil
			}
			return nil, err
		}
		name := string(line)
		// 只要片段 #EXTINF:10.043356,
		if !strings.HasPrefix(name, m3u8.TagEXTINF) {
			continue
		}
		// :10.043356,
		name = name[len(m3u8.TagEXTINF):]
		// ':'
		if name[0] == ':' {
			name = name[1:]
			// duration
			i := strings.IndexByte(name, ',')
			if i > 0 {
				name = name[i+1:]
				// title 0000.ts
				if len(name) == 0 {
					// 再读一行
					line, err = reader.ReadLine()
					if err != nil {
						return nil, err
					}
					name = string(line)
				}
				ts = append(ts, name)
				continue
			}
		}
		return nil, fmt.Errorf("invalid tag '%s'", name)
	}
}

func (d *downloader) downloadTS(url string) error {
	// 文件是否存在
	tsPath := filepath.Join(d.dir, path.Base(url))
	_, err := os.Stat(tsPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		return nil
	}
	// 下载
	rs, err := http.Get(url)
	if err != nil {
		return err
	}
	defer rs.Body.Close()
	if rs.StatusCode != http.StatusOK {
		return fmt.Errorf("download ts '%s' http status code '%d'", url, rs.StatusCode)
	}
	// 打开文件
	file, err := os.OpenFile(tsPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	// 保存
	_, err = io.Copy(file, rs.Body)
	if err != nil {
		os.Remove(tsPath)
	}
	return err
}

func (d *downloader) downloadTSRoutine() {
	defer d.Done()
	for {
		ts, ok := <-d.task
		if !ok {
			return
		}
		_url, err := url.Parse(ts)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if _url.Scheme != "" {
			_url.Path = path.Join(_url.Path, ts)
			ts = _url.String()
		} else {
			ts = d.urlDir.String() + "/" + ts
		}
		fmt.Println("download ts", ts)
		ok = false
		for {
			err = d.downloadTS(ts)
			if err == nil {
				break
			}
		}
	}
}

func main() {
	// 参数
	d := new(downloader)
	flag.StringVar(&d.url, "url", "", "m3u8 url")
	flag.StringVar(&d.dir, "dir", "", "output dir")
	flag.IntVar(&d.routine, "routine", 5, "concurrent download")
	flag.Parse()
	// 下载
	err := d.Download()
	if err != nil {
		fmt.Println(err)
	}
}
