package m3u8

const (
	ContentType = "application/vnd.apple.mpegurl"
)

type MediaSegment struct {
	EXTINF                  *EXTINF
	EXT_X_BYTERANGE         *EXT_X_BYTERANGE
	EXT_X_DISCONTINUITY     bool
	EXT_X_KEY               *EXT_X_KEY
	EXT_X_MAP               *EXT_X_MAP
	EXT_X_PROGRAM_DATE_TIME string
	EXT_X_DATERANGE         *EXT_X_DATERANGE
}

type MediaPlayList struct {
	write                        Writer
	EXT_X_VERSION                string
	EXT_X_TARGETDURATION         string
	EXT_X_MEDIA_SEQUENCE         string
	EXT_X_DISCONTINUITY_SEQUENCE string
	EXT_X_ENDLIST                bool
	EXT_X_PLAYLIST_TYPE          string
	EXT_X_I_FRAMES_ONLY          bool
	EXT_X_INDEPENDENT_SEGMENTS   bool
	EXT_X_START                  *EXT_X_START
	MediaSegment                 []MediaSegment
}

type MasterPlayList struct {
	write                      Writer
	EXT_X_VERSION              string
	EXT_X_MEDIA                []EXT_X_MEDIA
	EXT_X_STREAM_INF           []EXT_X_STREAM_INF
	EXT_X_I_FRAME_STREAM_INF   []EXT_X_I_FRAME_STREAM_INF
	EXT_X_SESSION_DATA         []EXT_X_SESSION_DATA
	EXT_X_SESSION_KEY          *EXT_X_KEY
	EXT_X_INDEPENDENT_SEGMENTS bool
	EXT_X_START                *EXT_X_START
}
