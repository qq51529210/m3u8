package m3u8

import "io"

type Writer struct {
	buff   []byte    // 缓存
	writer io.Writer // 输出目标
}

// writer是接收输出的数据
func NewWriter(writer io.Writer) *Writer {
	p := new(Writer)
	p.writer = writer
	return p
}

// 设置writer
func (w *Writer) SetWriter(writer io.Writer) {
	w.writer = writer
}

// 刷新数据，buff写到writer中
func (w *Writer) flushBuffer() (int, error) {
	n, err := w.writer.Write(w.buff)
	w.buff = w.buff[:0]
	return n, err
}

// s有数据才添加到缓存buff
func (w *Writer) writeNoEmptyString(s string) {
	if s != "" {
		w.buff = append(w.buff, s...)
	}
}

// s有数据才添加到缓存buff，先添加a，在添加s
func (w *Writer) writeAttributeNoEmptyString(a, s string) {
	if s != "" {
		w.buff = append(w.buff, a...)
		w.buff = append(w.buff, s...)
	}
}

// #EXTM3U
func (w *Writer) EXTM3U() (int, error) {
	w.buff = append(w.buff, "#EXTM3U\n"...)
	// 输出
	return w.flushBuffer()
}

// #EXT-X-VERSION:<n>
func (w *Writer) EXT_X_VERSION(tag string) (int, error) {
	w.buff = append(w.buff, "#EXT-X-VERSION:"...)
	w.buff = append(w.buff, tag...)
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXTINF:<duration>,[<title>]
func (w *Writer) EXTINF(tag *EXTINF) (int, error) {
	w.buff = append(w.buff, "#EXTINF:"...)
	// duration
	w.buff = append(w.buff, tag.DURATION...)
	w.buff = append(w.buff, ',')
	// title
	w.writeNoEmptyString(tag.TITLE)
	// 换行
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXT-X-BYTERANGE:<n>[@<o>]
func (w *Writer) EXT_X_BYTERANGE(tag *EXT_X_BYTERANGE) (int, error) {
	w.buff = append(w.buff, "#EXT-X-BYTERANGE:"...)
	// n
	w.buff = append(w.buff, tag.N...)
	// o
	if tag.O != "" {
		w.buff = append(w.buff, '@')
		w.buff = append(w.buff, tag.O...)
	}
	// 换行
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXT-X-DISCONTINUITY
func (w *Writer) EXT_X_DISCONTINUITY() (int, error) {
	w.buff = append(w.buff, "#EXT-X-DISCONTINUITY\n"...)
	// 输出
	return w.flushBuffer()
}

// #EXT-X-KEY:<attribute-list>
func (w *Writer) EXT_X_KEY(tag *EXT_X_KEY) (int, error) {
	w.buff = append(w.buff, "#EXT-X-KEY:"...)
	// METHOD
	w.writeAttributeNoEmptyString("METHOD=", tag.METHOD)
	// URI
	w.writeAttributeNoEmptyString(",URI=", tag.URI)
	// IV
	w.writeAttributeNoEmptyString(",IV=", tag.IV)
	// KEYFORMAT
	w.writeAttributeNoEmptyString(",KEYFORMAT=", tag.KEYFORMAT)
	// KEYFORMATVERSIONS
	w.writeAttributeNoEmptyString(",KEYFORMATVERSIONS=", tag.KEYFORMATVERSIONS)
	// 换行
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXT-X-MAP:<attribute-list>
func (w *Writer) EXT_X_MAP(tag *EXT_X_MAP) (int, error) {
	w.buff = append(w.buff, "#EXT-X-MAP:"...)
	// URI
	w.writeAttributeNoEmptyString("URI=", tag.URI)
	// BYTERANGE
	w.writeAttributeNoEmptyString(",BYTERANGE=", tag.BYTERANGE)
	// 换行
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXT-X-PROGRAM-DATE-TIME:<date-time-msec>
func (w *Writer) EXT_X_PROGRAM_DATE_TIME(dateTime string) (int, error) {
	w.buff = append(w.buff, "#EXT-X-PROGRAM-DATE-TIME:"...)
	w.buff = append(w.buff, dateTime...)
	// 换行
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXT-X-DATERANGE:<attribute-list>
func (w *Writer) EXT_X_DATERANGE(tag *EXT_X_DATERANGE) (int, error) {
	w.buff = append(w.buff, "#EXT-X-DATERANGE:"...)
	// ID
	w.writeAttributeNoEmptyString("ID=", tag.ID)
	// CLASS
	w.writeAttributeNoEmptyString(",CLASS=", tag.CLASS)
	// START_DATE
	w.writeAttributeNoEmptyString(",START_DATE=", tag.START_DATE)
	// END_DATE
	w.writeAttributeNoEmptyString(",END_DATE=", tag.END_DATE)
	// DURATION
	w.writeAttributeNoEmptyString(",DURATION=", tag.DURATION)
	// PLANNED_DURATION
	w.writeAttributeNoEmptyString(",PLANNED_DURATION=", tag.PLANNED_DURATION)
	for k, v := range tag.X_CLIENT_ATTRIBUTE {
		w.buff = append(w.buff, ',')
		w.buff = append(w.buff, k...)
		w.buff = append(w.buff, '=')
		w.buff = append(w.buff, v...)
	}
	// SCTE35_CMD
	w.writeAttributeNoEmptyString(",SCTE35_CMD=", tag.SCTE35_CMD)
	// SCTE35_OUT
	w.writeAttributeNoEmptyString(",SCTE35_OUT=", tag.SCTE35_OUT)
	// SCTE35_IN
	w.writeAttributeNoEmptyString(",SCTE35_IN=", tag.SCTE35_IN)
	// END_ON_NEXT
	w.writeAttributeNoEmptyString(",END_ON_NEXT=", tag.END_ON_NEXT)
	// 换行
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXT-X-TARGETDURATION:<s>
func (w *Writer) EXT_X_TARGETDURATION(tag string) (int, error) {
	w.buff = append(w.buff, "#EXT-X-TARGETDURATION:"...)
	w.buff = append(w.buff, tag...)
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXT-X-MEDIA-SEQUENCE:<s>
func (w *Writer) EXT_X_MEDIA_SEQUENCE(tag string) (int, error) {
	w.buff = append(w.buff, "#EXT-X-MEDIA-SEQUENCE:"...)
	w.buff = append(w.buff, tag...)
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXT-X-DISCONTINUITY-SEQUENCE:<number>
func (w *Writer) EXT_X_DISCONTINUITY_SEQUENCE(tag string) (int, error) {
	w.buff = append(w.buff, "#EXT-X-DISCONTINUITY-SEQUENCE:"...)
	w.buff = append(w.buff, tag...)
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXT-X-ENDLIST
func (w *Writer) EXT_X_ENDLIST() (int, error) {
	w.buff = append(w.buff, "#EXT-X-ENDLIST\n"...)
	// 输出
	return w.flushBuffer()
}

// #EXT-X-PLAYLIST-TYPE:<type-enum>
func (w *Writer) EXT_X_PLAYLIST_TYPE(tag string) (int, error) {
	w.buff = append(w.buff, "#EXT-X-PLAYLIST-TYPE:"...)
	w.buff = append(w.buff, tag...)
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXT-X-I-FRAMES-ONLY
func (w *Writer) EXT_X_I_FRAMES_ONLY() (int, error) {
	w.buff = append(w.buff, "#EXT-X-I-FRAMES-ONLY\n"...)
	// 输出
	return w.flushBuffer()
}

// #EXT-X-MEDIA:<attribute-list>
func (w *Writer) EXT_X_MEDIA(tag *EXT_X_MEDIA) (int, error) {
	w.buff = append(w.buff, "#EXT-X-MEDIA:"...)
	// TYPE
	w.writeAttributeNoEmptyString("TYPE=", tag.TYPE)
	// URI
	w.writeAttributeNoEmptyString(",URI=", tag.URI)
	// GROUP_ID
	w.writeAttributeNoEmptyString(",GROUP_ID=", tag.GROUP_ID)
	// LANGUAGE
	w.writeAttributeNoEmptyString(",LANGUAGE=", tag.LANGUAGE)
	// ASSOC_LANGUAGE
	w.writeAttributeNoEmptyString(",ASSOC_LANGUAGE=", tag.ASSOC_LANGUAGE)
	// NAME
	w.writeAttributeNoEmptyString(",NAME=", tag.NAME)
	// DEFAULT
	w.writeAttributeNoEmptyString(",DEFAULT=", tag.DEFAULT)
	// AUTOSELECT
	w.writeAttributeNoEmptyString(",AUTOSELECT=", tag.AUTOSELECT)
	// FORCED
	w.writeAttributeNoEmptyString(",FORCED=", tag.FORCED)
	// INSTREAM_ID
	w.writeAttributeNoEmptyString(",INSTREAM_ID=", tag.INSTREAM_ID)
	// CHARACTERISTICS
	w.writeAttributeNoEmptyString(",CHARACTERISTICS=", tag.CHARACTERISTICS)
	// CHANNELS
	w.writeAttributeNoEmptyString(",CHANNELS=", tag.CHANNELS)
	// 换行
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXT-X-STREAM-INF:<attribute-list>
// <URI>
func (w *Writer) EXT_X_STREAM_INF(tag *EXT_X_STREAM_INF) (int, error) {
	w.buff = append(w.buff, "#EXT-X-STREAM-INF:"...)
	// BANDWIDTH
	w.writeAttributeNoEmptyString("BANDWIDTH=", tag.BANDWIDTH)
	// AVERAGE_BANDWIDTH
	w.writeAttributeNoEmptyString(",AVERAGE_BANDWIDTH=", tag.AVERAGE_BANDWIDTH)
	// CODECS
	w.writeAttributeNoEmptyString(",CODECS=", tag.CODECS)
	// RESOLUTION
	w.writeAttributeNoEmptyString(",RESOLUTION=", tag.RESOLUTION)
	// FRAME_RATE
	w.writeAttributeNoEmptyString(",FRAME_RATE=", tag.FRAME_RATE)
	// HDCP_LEVEL
	w.writeAttributeNoEmptyString(",HDCP_LEVEL=", tag.HDCP_LEVEL)
	// AUDIO
	w.writeAttributeNoEmptyString(",AUDIO=", tag.AUDIO)
	// VIDEO
	w.writeAttributeNoEmptyString(",VIDEO=", tag.VIDEO)
	// SUBTITLES
	w.writeAttributeNoEmptyString(",SUBTITLES=", tag.SUBTITLES)
	// CLOSED_CAPTIONS
	w.writeAttributeNoEmptyString(",CLOSED_CAPTIONS=", tag.CLOSED_CAPTIONS)
	// 换行
	w.buff = append(w.buff, '\n')
	// URI
	w.buff = append(w.buff, tag.URI...)
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXT-X-I-FRAME-STREAM-INF:<attribute-list>
// <URI>
func (w *Writer) EXT_X_I_FRAME_STREAM_INF(tag *EXT_X_I_FRAME_STREAM_INF) (int, error) {
	w.buff = append(w.buff, "#EXT-X-I-FRAME-STREAM-INF:"...)
	// BANDWIDTH
	w.writeAttributeNoEmptyString("BANDWIDTH=", tag.BANDWIDTH)
	// AVERAGE_BANDWIDTH
	w.writeAttributeNoEmptyString(",AVERAGE_BANDWIDTH=", tag.AVERAGE_BANDWIDTH)
	// CODECS
	w.writeAttributeNoEmptyString(",CODECS=", tag.CODECS)
	// RESOLUTION
	w.writeAttributeNoEmptyString(",RESOLUTION=", tag.RESOLUTION)
	// HDCP_LEVEL
	w.writeAttributeNoEmptyString(",HDCP_LEVEL=", tag.HDCP_LEVEL)
	// VIDEO
	w.writeAttributeNoEmptyString(",VIDEO=", tag.VIDEO)
	// 换行
	w.buff = append(w.buff, '\n')
	// URI
	w.buff = append(w.buff, tag.URI...)
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXT-X-SESSION-DATA:<attribute-list>
func (w *Writer) EXT_X_SESSION_DATA(tag *EXT_X_SESSION_DATA) (int, error) {
	w.buff = append(w.buff, "#EXT-X-SESSION-DATA:"...)
	// DATA_ID
	w.writeAttributeNoEmptyString("DATA_ID=", tag.DATA_ID)
	// VALUE
	w.writeAttributeNoEmptyString(",VALUE=", tag.VALUE)
	// URI
	w.writeAttributeNoEmptyString(",URI=", tag.URI)
	// LANGUAGE
	w.writeAttributeNoEmptyString(",LANGUAGE=", tag.LANGUAGE)
	// 换行
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXT-X-SESSION-KEY:<attribute-list>
func (w *Writer) EXT_X_SESSION_KEY(tag *EXT_X_KEY) (int, error) {
	w.buff = append(w.buff, "#EXT-X-SESSION-KEY:"...)
	// METHOD
	w.writeAttributeNoEmptyString("METHOD=", tag.METHOD)
	// URI
	w.writeAttributeNoEmptyString(",URI=", tag.URI)
	// IV
	w.writeAttributeNoEmptyString(",IV=", tag.IV)
	// KEYFORMAT
	w.writeAttributeNoEmptyString(",KEYFORMAT=", tag.KEYFORMAT)
	// KEYFORMATVERSIONS
	w.writeAttributeNoEmptyString(",KEYFORMATVERSIONS=", tag.KEYFORMATVERSIONS)
	// 换行
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

// #EXT-X-INDEPENDENT-SEGMENTS
func (w *Writer) EXT_X_INDEPENDENT_SEGMENTS() (int, error) {
	w.buff = append(w.buff, "#EXT-X-INDEPENDENT-SEGMENTS\n"...)
	// 输出
	return w.flushBuffer()
}

// #EXT-X-START:<attribute-list>
func (w *Writer) EXT_X_START(tag *EXT_X_START) (int, error) {
	w.buff = append(w.buff, "#EXT-X-START:"...)
	// TIME_OFFSET
	w.writeAttributeNoEmptyString("TIME_OFFSET=", tag.TIME_OFFSET)
	// PRECISE
	w.writeAttributeNoEmptyString(",PRECISE=", tag.PRECISE)
	// 换行
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}
