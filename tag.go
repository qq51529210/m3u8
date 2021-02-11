// 参考文档：
// https://tools.ietf.org/html/rfc8216
package m3u8

const (
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
