// 参考文档：
// https://tools.ietf.org/html/rfc8216
package m3u8

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sync"
)

const (
	ContentType = "application/vnd.apple.mpegurl"
	// basic tags
	TagEXTM3U        = "#EXTM3U"
	TagEXT_X_VERSION = "#EXT-X-VERSION"
	// media segment tags
	TagEXTINF                  = "#EXTINF"
	TagEXT_X_BYTERANGE         = "#EXT-X-BYTERANGE"
	TagEXT_X_DISCONTINUITY     = "#EXT-X-DISCONTINUITY"
	TagEXT_X_KEY               = "#EXT-X-KEY"
	TagEXT_X_MAP               = "#EXT-X-MAP"
	TagEXT_X_PROGRAM_DATE_TIME = "#EXT-X-PROGRAM-DATE-TIME"
	TagEXT_X_DATERANGE         = "#EXT-X-DATERANGE"
	// media playlist tags
	TagEXT_X_TARGETDURATION         = "#EXT-X-TARGETDURATION"
	TagEXT_X_MEDIA_SEQUENCE         = "#EXT-X-MEDIA-SEQUENCE"
	TagEXT_X_DISCONTINUITY_SEQUENCE = "#EXT-X-DISCONTINUITY-SEQUENCE"
	TagEXT_X_ENDLIST                = "#EXT-X-ENDLIST"
	TagEXT_X_PLAYLIST_TYPE          = "#EXT-X-PLAYLIST-TYPE"
	TagEXT_X_I_FRAMES_ONLY          = "#EXT-X-I-FRAMES-ONLY"
	// master playlist tags
	TagEXT_X_MEDIA                = "#EXT-X-MEDIA"
	TagEXT_X_STREAM_INF           = "#EXT-X-STREAM-INF"
	TagEXT_X_I_FRAME_STREAM_INF   = "#EXT-X-I-FRAME-STREAM-INF"
	TagEXT_X_SESSION_DATA         = "#EXT-X-SESSION-DATA"
	TagEXT_X_SESSION_KEY          = "#EXT-X-SESSION-KEY"
	TagEXT_X_INDEPENDENT_SEGMENTS = "#EXT-X-INDEPENDENT-SEGMENTS"
	TagEXT_X_START                = "#EXT-X-START"
)

var (
	tagEXTINF                  = []byte(TagEXTINF)
	tagEXT_X_BYTERANGE         = []byte(TagEXT_X_BYTERANGE)
	tagEXT_X_DISCONTINUITY     = []byte(TagEXT_X_DISCONTINUITY)
	tagEXT_X_KEY               = []byte(TagEXT_X_KEY)
	tagEXT_X_MAP               = []byte(TagEXT_X_MAP)
	tagEXT_X_PROGRAM_DATE_TIME = []byte(TagEXT_X_PROGRAM_DATE_TIME)
	tagEXT_X_DATERANGE         = []byte(TagEXT_X_DATERANGE)
	// media playlist tags
	tagEXT_X_TARGETDURATION         = []byte(TagEXT_X_TARGETDURATION)
	tagEXT_X_MEDIA_SEQUENCE         = []byte(TagEXT_X_MEDIA_SEQUENCE)
	tagEXT_X_DISCONTINUITY_SEQUENCE = []byte(TagEXT_X_DISCONTINUITY_SEQUENCE)
	tagEXT_X_ENDLIST                = []byte(TagEXT_X_ENDLIST)
	tagEXT_X_PLAYLIST_TYPE          = []byte(TagEXT_X_PLAYLIST_TYPE)
	tagEXT_X_I_FRAMES_ONLY          = []byte(TagEXT_X_I_FRAMES_ONLY)
	// master playlist tags
	tagEXT_X_MEDIA                = []byte(TagEXT_X_MEDIA)
	tagEXT_X_STREAM_INF           = []byte(TagEXT_X_STREAM_INF)
	tagEXT_X_I_FRAME_STREAM_INF   = []byte(TagEXT_X_I_FRAME_STREAM_INF)
	tagEXT_X_SESSION_DATA         = []byte(TagEXT_X_SESSION_DATA)
	tagEXT_X_SESSION_KEY          = []byte(TagEXT_X_SESSION_KEY)
	tagEXT_X_INDEPENDENT_SEGMENTS = []byte(TagEXT_X_INDEPENDENT_SEGMENTS)
	tagEXT_X_START                = []byte(TagEXT_X_START)
	// empty slice
	emptySlice = make([]byte, 0)
)

type EXTINF struct {
	DURATION string
	TITLE    string
}

type EXT_X_BYTERANGE struct {
	N string
	O string
}

type EXT_X_KEY struct {
	METHOD            string
	URI               string
	IV                string
	KEYFORMAT         string
	KEYFORMATVERSIONS string
}

type EXT_X_MAP struct {
	URI       string
	BYTERANGE string
}

type EXT_X_DATERANGE struct {
	ID                 string
	CLASS              string
	START_DATE         string
	END_DATE           string
	DURATION           string
	PLANNED_DURATION   string
	X_CLIENT_ATTRIBUTE map[string]string
	SCTE35_CMD         string
	SCTE35_OUT         string
	SCTE35_IN          string
	END_ON_NEXT        string
}

type EXT_X_MEDIA struct {
	TYPE            string
	URI             string
	GROUP_ID        string
	LANGUAGE        string
	ASSOC_LANGUAGE  string
	NAME            string
	DEFAULT         string
	AUTOSELECT      string
	FORCED          string
	INSTREAM_ID     string
	CHARACTERISTICS string
	CHANNELS        string
}

type EXT_X_STREAM_INF struct {
	BANDWIDTH         string
	AVERAGE_BANDWIDTH string
	CODECS            string
	RESOLUTION        string
	FRAME_RATE        string
	HDCP_LEVEL        string
	AUDIO             string
	VIDEO             string
	SUBTITLES         string
	CLOSED_CAPTIONS   string
	URI               string
}

type EXT_X_I_FRAME_STREAM_INF struct {
	BANDWIDTH         string
	AVERAGE_BANDWIDTH string
	CODECS            string
	RESOLUTION        string
	HDCP_LEVEL        string
	VIDEO             string
	URI               string
}

type EXT_X_SESSION_DATA struct {
	DATA_ID  string
	VALUE    string
	URI      string
	LANGUAGE string
}

type EXT_X_START struct {
	TIME_OFFSET string
	PRECISE     string
}

type MediaSegment struct {
	EXTINF                  EXTINF
	EXT_X_BYTERANGE         *EXT_X_BYTERANGE
	EXT_X_DISCONTINUITY     bool
	EXT_X_KEY               *EXT_X_KEY
	EXT_X_MAP               *EXT_X_MAP
	EXT_X_PROGRAM_DATE_TIME string
	EXT_X_DATERANGE         *EXT_X_DATERANGE
}

type MediaPlayList struct {
	writer                       Writer
	EXT_X_VERSION                string
	EXT_X_TARGETDURATION         string
	EXT_X_MEDIA_SEQUENCE         string
	EXT_X_DISCONTINUITY_SEQUENCE string
	EXT_X_PLAYLIST_TYPE          string
	EXT_X_I_FRAMES_ONLY          bool
	EXT_X_INDEPENDENT_SEGMENTS   bool
	EXT_X_START                  *EXT_X_START
	MediaSegment                 []MediaSegment
	EXT_X_ENDLIST                bool
}

type MasterPlayList struct {
	writer                     Writer
	EXT_X_VERSION              string
	EXT_X_MEDIA                []EXT_X_MEDIA
	EXT_X_STREAM_INF           []EXT_X_STREAM_INF
	EXT_X_I_FRAME_STREAM_INF   []EXT_X_I_FRAME_STREAM_INF
	EXT_X_SESSION_DATA         []EXT_X_SESSION_DATA
	EXT_X_SESSION_KEY          *EXT_X_KEY
	EXT_X_INDEPENDENT_SEGMENTS bool
	EXT_X_START                *EXT_X_START
}

func ParseLine(line []byte) (tag, value []byte) {
	i := bytes.IndexByte(line, ':')
	if i < 0 {
		return line, emptySlice
	}
	return line[:i], line[i+1:]
}

func ParseAttribute(line []byte) (map[string]string, error) {
	m := make(map[string]string)
	i := 0
	for {
		// name=value
		i = bytes.IndexByte(line, '=')
		if i < 0 {
			return nil, fmt.Errorf("incomplete attribute '%s', can't find '='", string(line))
		}
		// name
		name := string(line[:i])
		// value...
		line = line[i+1:]
		if len(line) < 0 {
			// name=
			return nil, fmt.Errorf("incomplete attribute '%s', can't find <value>", string(line))
		}
		if line[0] == '"' {
			i = indexString(line)
			if i < 0 {
				// name="...
				return m, fmt.Errorf("incomplete attribute '%s', can't find end '\"'", string(line))
			} else {
				// name="..."
				m[name] = string(line[:i])
				p := line[i+1:]
				if len(p) > 0 {
					if p[0] != ',' {
						return m, fmt.Errorf("incomplete attribute, can't find ',' after '%s'", line[:i])
					}
					line = p[1:]
				} else {
					return m, nil
				}
			}
		} else {
			i = bytes.IndexByte(line, ',')
			if i < 0 {
				m[name] = string(line)
				return m, nil
			} else {
				m[name] = string(line[:i])
				line = line[i+1:]
				if len(line) <= 0 {
					return m, nil
				}
			}
		}
	}
}

func indexString(s []byte) int {
	for i := 1; i < len(s); i++ {
		if s[i] == '"' && s[i-1] != '\\' {
			return i
		}
	}
	return -1
}

func SimpleDownload(list, dir string, concurrent int) (err error) {
	// 创建目录
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	// 保存m3u8列表文件
	var file *os.File
	file, err = saveM3U8(list, dir)
	defer file.Close()
	urlDir := list[:len(list)-len(path.Base(list))]
	// 并发下载
	exit := make(chan struct{})
	errors := make(chan error, concurrent+1)
	defer close(errors)
	tasks := make(chan string, concurrent+1)
	defer close(tasks)
	var wait sync.WaitGroup
	var exitOnce sync.Once
	for i := 0; i < concurrent; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			defer exitOnce.Do(func() {
				close(exit)
				errors <- nil
			})
			var _url *url.URL
			var err error
			for {
				select {
				case task := <-tasks:
					_url, err = url.Parse(task)
					if err != nil {
						errors <- err
						return
					}
					if task[0] == '/' {
						task = task[1:]
					}
					if _url.Scheme == "" {
						task = urlDir + task
					}
					err = saveTS(task, dir)
					if err != nil {
						errors <- err
						return
					}
				case <-exit:
					return
				}
			}
		}()
	}
	// 读取每一行，创建下载任务
	go func() {
		defer exitOnce.Do(func() {
			close(exit)
			errors <- nil
		})
		var line []byte
		var err error
		reader := NewReader(file, nil)
		for {
			line, err = reader.ReadLine()
			if err != nil {
				errors <- err
				return
			}
			if line == nil {
				return
			}
			if len(line) == 0 {
				continue
			}
			// 只要片段
			if bytes.HasPrefix(line, tagEXTINF) {
				p := line[len(tagEXTINF):]
				// ':'
				if p[0] == ':' {
					p = p[1:]
					// duration
					i := bytes.IndexByte(p, ',')
					if i < 0 {
						errors <- fmt.Errorf("invalid tag '%s', can't find ','", string(line))
						return
					}
					p = p[i+1:]
					// title
					for len(p) == 0 {
						// 再读一行
						line, err = reader.ReadLine()
						if err != nil {
							return
						}
						p = line
					}
					// 添加任务
					select {
					case tasks <- string(p):
					case <-exit:
						return
					}
				} else {
					errors <- fmt.Errorf("invalid tag '%s'", string(line))
					return
				}
			} else if bytes.Compare(line, tagEXT_X_ENDLIST) == 0 {
				return
			}
		}
	}()
	// 等待错误或完成退出
	<-exit
	wait.Wait()
	err = <-errors
	return
}

func saveM3U8(url, dir string) (*os.File, error) {
	fmt.Println("download", url, "to", dir)
	// 下载
	rs, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()
	if rs.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code '%d'", rs.StatusCode)
	}
	// 打开列表文件
	var file *os.File
	file, err = os.OpenFile(filepath.Join(dir, path.Base(url)), os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	// 保存
	_, err = io.Copy(file, rs.Body)
	if err != nil {
		file.Close()
		return nil, err
	}
	file.Seek(0, 0)
	return file, nil
}

func saveTS(url, dir string) error {
	fmt.Println("download", url, "to", dir)
	// 文件是否存在
	tsPath := filepath.Join(dir, path.Base(url))
	_, err := os.Stat(tsPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		fmt.Println(tsPath, "exist")
		return nil
	}
	// 下载
	rs, err := http.Get(url)
	if err != nil {
		return err
	}
	defer rs.Body.Close()
	if rs.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code '%d'", rs.StatusCode)
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
