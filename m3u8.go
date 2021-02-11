// 参考文档：
// https://tools.ietf.org/html/rfc8216
package m3u8

import "io"

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

func (mpl *MediaPlayList) Encode(w io.Writer) (err error) {
	mpl.writer.SetWriter(w)
	_, err = mpl.writer.EXTM3U()
	if err != nil {
		return
	}
	if mpl.EXT_X_VERSION != "" {
		_, err = mpl.writer.EXT_X_VERSION(mpl.EXT_X_VERSION)
		if err != nil {
			return
		}
	}
	if mpl.EXT_X_TARGETDURATION != "" {
		_, err = mpl.writer.EXT_X_TARGETDURATION(mpl.EXT_X_TARGETDURATION)
		if err != nil {
			return
		}
	}
	if mpl.EXT_X_MEDIA_SEQUENCE != "" {
		_, err = mpl.writer.EXT_X_MEDIA_SEQUENCE(mpl.EXT_X_MEDIA_SEQUENCE)
		if err != nil {
			return
		}
	}
	if mpl.EXT_X_DISCONTINUITY_SEQUENCE != "" {
		_, err = mpl.writer.EXT_X_DISCONTINUITY_SEQUENCE(mpl.EXT_X_DISCONTINUITY_SEQUENCE)
		if err != nil {
			return
		}
	}
	if mpl.EXT_X_PLAYLIST_TYPE != "" {
		_, err = mpl.writer.EXT_X_PLAYLIST_TYPE(mpl.EXT_X_PLAYLIST_TYPE)
		if err != nil {
			return
		}
	}
	if mpl.EXT_X_I_FRAMES_ONLY {
		_, err = mpl.writer.EXT_X_I_FRAMES_ONLY()
		if err != nil {
			return
		}
	}
	if mpl.EXT_X_INDEPENDENT_SEGMENTS {
		_, err = mpl.writer.EXT_X_INDEPENDENT_SEGMENTS()
		if err != nil {
			return
		}
	}
	if mpl.EXT_X_START != nil {
		_, err = mpl.writer.EXT_X_START(mpl.EXT_X_START)
		if err != nil {
			return
		}
	}
	for i := 0; i < len(mpl.MediaSegment); i++ {
		_, err = mpl.writer.EXTINF(&mpl.MediaSegment[i].EXTINF)
		if err != nil {
			return
		}
		if mpl.MediaSegment[i].EXT_X_BYTERANGE != nil {
			_, err = mpl.writer.EXT_X_BYTERANGE(mpl.MediaSegment[i].EXT_X_BYTERANGE)
			if err != nil {
				return
			}
		}
		if mpl.MediaSegment[i].EXT_X_DISCONTINUITY {
			_, err = mpl.writer.EXT_X_DISCONTINUITY()
			if err != nil {
				return
			}
		}
		if mpl.MediaSegment[i].EXT_X_KEY != nil {
			_, err = mpl.writer.EXT_X_KEY(mpl.MediaSegment[i].EXT_X_KEY)
			if err != nil {
				return
			}
		}
		if mpl.MediaSegment[i].EXT_X_MAP != nil {
			_, err = mpl.writer.EXT_X_MAP(mpl.MediaSegment[i].EXT_X_MAP)
			if err != nil {
				return
			}
		}
		if mpl.MediaSegment[i].EXT_X_PROGRAM_DATE_TIME != "" {
			_, err = mpl.writer.EXT_X_PROGRAM_DATE_TIME(mpl.MediaSegment[i].EXT_X_PROGRAM_DATE_TIME)
			if err != nil {
				return
			}
		}
		if mpl.MediaSegment[i].EXT_X_DATERANGE != nil {
			_, err = mpl.writer.EXT_X_DATERANGE(mpl.MediaSegment[i].EXT_X_DATERANGE)
			if err != nil {
				return
			}
		}
	}
	if mpl.EXT_X_ENDLIST {
		_, err = mpl.writer.EXT_X_ENDLIST()
		if err != nil {
			return
		}
	}
	return
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

func (mpl *MasterPlayList) Encode(w io.Writer) (err error) {
	mpl.writer.SetWriter(w)
	_, err = mpl.writer.EXTM3U()
	if err != nil {
		return
	}
	if mpl.EXT_X_VERSION != "" {
		_, err = mpl.writer.EXT_X_VERSION(mpl.EXT_X_VERSION)
		if err != nil {
			return
		}
	}
	for i := 0; i < len(mpl.EXT_X_MEDIA); i++ {
		_, err = mpl.writer.EXT_X_MEDIA(&mpl.EXT_X_MEDIA[i])
		if err != nil {
			return
		}
	}
	for i := 0; i < len(mpl.EXT_X_STREAM_INF); i++ {
		_, err = mpl.writer.EXT_X_STREAM_INF(&mpl.EXT_X_STREAM_INF[i])
		if err != nil {
			return
		}
	}
	for i := 0; i < len(mpl.EXT_X_I_FRAME_STREAM_INF); i++ {
		_, err = mpl.writer.EXT_X_I_FRAME_STREAM_INF(&mpl.EXT_X_I_FRAME_STREAM_INF[i])
		if err != nil {
			return
		}
	}
	for i := 0; i < len(mpl.EXT_X_SESSION_DATA); i++ {
		_, err = mpl.writer.EXT_X_SESSION_DATA(&mpl.EXT_X_SESSION_DATA[i])
		if err != nil {
			return
		}
	}
	if mpl.EXT_X_SESSION_KEY != nil {
		_, err = mpl.writer.EXT_X_SESSION_KEY(mpl.EXT_X_SESSION_KEY)
		if err != nil {
			return
		}
	}
	if mpl.EXT_X_INDEPENDENT_SEGMENTS {
		_, err = mpl.writer.EXT_X_INDEPENDENT_SEGMENTS()
		if err != nil {
			return
		}
	}
	if mpl.EXT_X_START != nil {
		_, err = mpl.writer.EXT_X_START(mpl.EXT_X_START)
		if err != nil {
			return
		}
	}
	return
}
