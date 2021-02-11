package m3u8

import "io"

type Writer struct {
	buff   []byte
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	p := new(Writer)
	p.writer = writer
	return p
}

func (w *Writer) flushBuffer() (int, error) {
	n, err := w.writer.Write(w.buff)
	w.buff = w.buff[:0]
	return n, err
}

func (w *Writer) writeNoEmptyString(s string) {
	if s != "" {
		w.buff = append(w.buff, s...)
	}
}

func (w *Writer) writeAttributeNoEmptyString(a, s string) {
	if s != "" {
		w.buff = append(w.buff, a...)
		w.buff = append(w.buff, s...)
	}
}

func (w *Writer) TagEXTM3U() (int, error) {
	w.buff = w.buff[:0]
	w.buff = append(w.buff, "#EXTM3U\n"...)
	// 输出
	return w.flushBuffer()
}

func (w *Writer) TagEXT_X_VERSION(tag string) (int, error) {
	w.buff = append(w.buff, "#EXT-X-VERSION:"...)
	w.buff = append(w.buff, tag...)
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

func (w *Writer) TagEXTINF(tag *EXTINF) (int, error) {
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

func (w *Writer) TagEXT_X_BYTERANGE(tag *EXT_X_BYTERANGE) (int, error) {
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

func (w *Writer) TagEXT_X_DISCONTINUITY() (int, error) {
	w.buff = append(w.buff, "#EXT-X-DISCONTINUITY\n"...)
	// 输出
	return w.flushBuffer()
}

func (w *Writer) TagEXT_X_KEY(tag *EXT_X_KEY) (int, error) {
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

func (w *Writer) TagEXT_X_MAP(tag *EXT_X_MAP) (int, error) {
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

func (w *Writer) TagEXT_X_PROGRAM_DATE_TIME(dateTime string) (int, error) {
	w.buff = append(w.buff, "#EXT-X-PROGRAM-DATE-TIME:"...)
	w.buff = append(w.buff, dateTime...)
	// 换行
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

func (w *Writer) TagEXT_X_DATERANGE(tag *EXT_X_DATERANGE) (int, error) {
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

func (w *Writer) TagEXT_X_TARGETDURATION(tag string) (int, error) {
	w.buff = append(w.buff, "#EXT-X-TARGETDURATION:"...)
	w.buff = append(w.buff, tag...)
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

func (w *Writer) TagEXT_X_MEDIA_SEQUENCE(tag string) (int, error) {
	w.buff = append(w.buff, "#EXT-X-MEDIA-SEQUENCE:"...)
	w.buff = append(w.buff, tag...)
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

func (w *Writer) TagEXT_X_DISCONTINUITY_SEQUENCE(tag string) (int, error) {
	w.buff = append(w.buff, "#EXT-X-DISCONTINUITY-SEQUENCE:"...)
	w.buff = append(w.buff, tag...)
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

func (w *Writer) TagEXT_X_ENDLIST(tag string) (int, error) {
	w.buff = append(w.buff, "#EXT-X-ENDLIST\n"...)
	// 输出
	return w.flushBuffer()
}

func (w *Writer) TagEXT_X_PLAYLIST_TYPE(tag string) (int, error) {
	w.buff = append(w.buff, "#EXT-X-PLAYLIST-TYPE:"...)
	w.buff = append(w.buff, tag...)
	w.buff = append(w.buff, '\n')
	// 输出
	return w.flushBuffer()
}

func (w *Writer) TagEXT_X_I_FRAMES_ONLY(tag string) (int, error) {
	w.buff = append(w.buff, "#EXT-X-I-FRAMES-ONLY\n"...)
	// 输出
	return w.flushBuffer()
}

func (w *Writer) TagEXT_X_MEDIA(tag *EXT_X_MEDIA) (int, error) {
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

func (w *Writer) TagEXT_X_STREAM_INF(tag *EXT_X_STREAM_INF) (int, error) {
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

func (w *Writer) TagEXT_X_I_FRAME_STREAM_INF(tag *EXT_X_I_FRAME_STREAM_INF) (int, error) {
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

func (w *Writer) TagEXT_X_SESSION_DATA(tag *EXT_X_SESSION_DATA) (int, error) {
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

func (w *Writer) TagEXT_X_SESSION_KEY(tag *EXT_X_KEY) (int, error) {
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

func (w *Writer) TagEXT_X_I_INDEPENDENT_SEGMENTS(tag string) (int, error) {
	w.buff = append(w.buff, "#EXT-X-I-INDEPENDENT-SEGMENTS\n"...)
	// 输出
	return w.flushBuffer()
}

func (w *Writer) TagEXT_X_START(tag *EXT_X_START) (int, error) {
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
