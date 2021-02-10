package m3u8

import (
	"strconv"
)

var (
	_EXT_X_STREAM_INF        = []byte("#EXT-X-STREAM-INF")
	_EXT_X_ALLOW_CACHE       = []byte("#EXT-X-ALLOW-CACHE")
	_EXT_X_PROGRAM_DATE_TIME = []byte("#EXT-X-PROGRAM-DATE-TIME")
	_EXT_X_KEY               = []byte("#EXT-X-KEY")
	_EXT_X_BYTERANGE         = []byte("#EXT-X-BYTERANGE")
)

type Writer struct {
	buf []byte
}

// 写入一个string
func (w *Writer) String(s string) {
	w.buf = append(w.buf, s...)
}

// 写入一个‘\n’
func (w *Writer) Endline() {
	w.buf = append(w.buf, '\n')
}

// 写入一个‘,’
func (w *Writer) Comma() {
	w.buf = append(w.buf, ',')
}

// 表明该文件是一个 m3u8 文件。每个 M3U 文件必须将该标签放置在第一行。
func (w *Writer) EXTM3U() {
	w.buf = w.buf[:0]
	w.buf = append(w.buf, "#EXTM3U\n"...)
}

/*
表示 HLS 的协议版本号，该标签与流媒体的兼容性相关。该标签为全局作用域，使能整个 m3u8 文件；
每个 m3u8 文件内最多只能出现一个该标签定义。如果 m3u8 文件不包含该标签，则默认为协议的第一个版本。
*/
func (w *Writer) EXT_X_VERSION(n int64) {
	w.buf = append(w.buf, "#EXT-X-VERSION:"...)
	w.buf = strconv.AppendInt(w.buf, n, 64)
	w.buf = append(w.buf, '\n')
}

/*
表示每个视频分段最大的时长（单位秒）。
该标签为必选标签。
其格式为：#EXT-X-TARGETDURATION:<s>
其中：参数s表示目标时长（单位：秒）。
*/
func (w *Writer) EXT_X_TARGETDURATION(n float64) {
	w.buf = append(w.buf, "#EXT-X-TARGETDURATION:"...)
	w.buf = strconv.AppendFloat(w.buf, n, 'f', 0, 64)
	w.buf = append(w.buf, '\n')
}

/*
表示播放列表第一个 URL 片段文件的序列号。
每个媒体片段 URL 都拥有一个唯一的整型序列号。
每个媒体片段序列号按出现顺序依次加 1。
如果该标签未指定，则默认序列号从 0 开始。
媒体片段序列号与片段文件名无关。
其格式为：#EXT-X-MEDIA-SEQUENCE:<number>
其中：参数number即为切片序列号。
*/
func (w *Writer) EXT_X_MEDIA_SEQUENCE(n int64) {
	w.buf = append(w.buf, "#EXT-X-MEDIA-SEQUENCE:"...)
	w.buf = strconv.AppendInt(w.buf, n, 64)
	w.buf = append(w.buf, '\n')
}

/*
当以下任一情况变化时，必须使用该标签：
文件格式（file format）
数字（number），类型（type），媒体标识符（identifiers of tracks）
时间戳序列（timestamp sequence）

当以下任一情况变化时，应当使用该标签：
编码参数（encoding parameters）
编码序列（encoding sequence）

注：EXT-X-DISCONTINUITY 的一个经典使用场景就是在视屏流中插入广告，由于视屏流与广告视屏流不是同一份资源，
因此在这两种流切换时使用 EXT-X-DISCONTINUITY 进行指明，客户端看到该标签后，就会处理这种切换中断问题，让体验更佳。
*/
func (w *Writer) EXT_X_DISCONTINUITY() {
	w.buf = append(w.buf, "#EXT-X-DISCONTINUITY\n"...)
}

/*
该标签使能同步相同流的不同 Rendition 和 具备 EXT-X-DISCONTINUITY 标签的不同备份流。
其格式为：

#EXT-X-DISCONTINUITY-SEQUENCE:<number>
其中：参数number为一个十进制整型数值。
如果播放列表未设置 EXT-X-DISCONTINUITY-SEQUENCE 标签，那么对于第一个切片的中断序列号应当为 0。
*/
func (w *Writer) EXT_X_DISCONTINUITY_SEQUENCE(n int64) {
	w.buf = append(w.buf, "#EXT-X-DISCONTINUITY-SEQUENCE:"...)
	w.buf = strconv.AppendInt(w.buf, n, 64)
	w.buf = append(w.buf, '\n')
}

/*
表明 m3u8 文件的结束。
该标签可出现在 m3u8 文件任意位置，一般是结尾。
其格式为：#EXT-X-ENDLIST
*/
func (w *Writer) EXT_X_ENDLIST() {
	w.buf = append(w.buf, "#EXT-X-ENDLIST\n"...)
}

// duration：可以为十进制的整型或者浮点型，其值必须小于或等于 EXT-X-TARGETDURATION 指定的值。
// 如果兼容版本号 EXT-X-VERSION 小于 3，那么必须使用整型。
func (w *Writer) EXT_INF(duration float64, title string) {
	w.buf = append(w.buf, "#EXT-INF:"...)
	// duration
	w.buf = strconv.AppendFloat(w.buf, duration, 'f', 0, 64)
	w.buf = append(w.buf, ","...)
	// title
	if title != "" {
		w.buf = append(w.buf, title...)
		w.buf = append(w.buf, '\n')
	}
}

/*
该标签表示接下来的切片资源是其后 URI 指定的媒体片段资源的局部范围（即截取 URI 媒体资源部分内容作为下一个切片）。
该标签只对其后一个 URI 起作用。其格式为：#EXT-X-BYTERANGE:<n>[@<o>]

n是一个十进制整型，表示截取片段大小（单位：字节）。
可选参数o也是一个十进制整型，指示截取起始位置（以字节表示，在 URI 指定的资源开头移动该字节位置后进行截取）。
如果o未指定，则截取起始位置从上一个该标签截取完成的下一个字节（即上一个n+o+1）开始截取。
如果没有指定该标签，则表明的切分范围为整个 URI 资源片段。
注：使用 EXT-X-BYTERANGE 标签要求兼容版本号 EXT-X-VERSION 大于等于 4。
*/
func (w *Writer) EXT_X_BYTERANGE(n, o int64) {
	w.buf = append(w.buf, "#EXT-X-BYTERANGE:"...)
	w.buf = strconv.AppendInt(w.buf, n, 64)
	if o >= 0 {
		w.buf = append(w.buf, '@')
		w.buf = strconv.AppendInt(w.buf, o, 64)
	}
	w.buf = append(w.buf, '\n')
}

type EXT_X_KEY struct {
	/*
		METHOD：该值是一个可枚举的字符串，指定了加密方法。
		该键是必须参数。其值可为NONE，AES-128，SAMPLE-AES当中的一个。
		其中：
		NONE：表示切片未进行加密（此时其他属性不能出现）；
		AES-128：表示表示使用 AES-128 进行加密。
		SAMPLE-AES：意味着媒体片段当中包含样本媒体，比如音频或视频，它们使用 AES-128 进行加密。这种情况下 IV 属性可以出现也可以不出现。
	*/
	METHOD string
	/*
		URI：指定密钥路径。
		该密钥是一个 16 字节的数据。
		该键是必须参数，除非 METHOD 为NONE。
	*/
	URI string
	/*
		IV：该值是一个 128 位的十六进制数值。
		AES-128 要求使用相同的 16字节 IV 值进行加密和解密。使用不同的 IV 值可以增强密码强度。
		如果属性列表出现 IV，则使用该值；如果未出现，则默认使用媒体片段序列号（即 EXT-X-MEDIA-SEQUENCE）作为其 IV 值，使用大端字节序，往左填充 0 直到序列号满足 16 字节（128 位）。
	*/
	IV string
	/*
	   KEYFORMAT：由双引号包裹的字符串，标识密钥在密钥文件中的存储方式（密钥文件中的 AES-128 密钥是以二进制方式存储的16个字节的密钥）。
	   该属性为可选参数，其默认值为"identity"。
	   使用该属性要求兼容版本号 EXT-X-VERSION 大于等于 5。
	*/
	KEYFORMAT string
	/*
	   KEYFORMATVERSIONS：由一个或多个被/分割的正整型数值构成的带引号的字符串（比如："1"，"1/2"，"1/2/5"）。
	   如果有一个或多特定的 KEYFORMT 版本被定义了，则可使用该属性指示具体版本进行编译。
	   该属性为可选参数，其默认值为"1"。
	   使用该属性要求兼容版本号 EXT-X-VERSION 大于等于 5。
	*/
	KEYFORMATVERSIONS string
}

/*
媒体片段可以进行加密，而该标签可以指定解密方法。
该标签对所有 媒体片段 和 由标签 EXT-X-MAP 声明的围绕其间的所有 媒体初始化块（Meida Initialization Section） 都起作用，
直到遇到下一个 EXT-X-KEY（若 m3u8 文件只有一个 EXT-X-KEY 标签，则其作用于所有媒体片段）。
多个 EXT-X-KEY 标签如果最终生成的是同样的秘钥，则他们都可作用于同一个媒体片段。
该标签使用格式为：#EXT-X-KEY:<attribute-list>
*/
func (w *Writer) EXT_X_KEY(info *EXT_X_KEY) {
	w.buf = append(w.buf, "#EXT-X-KEY:"...)
	w.ext_x_key(info)
}

func (w *Writer) ext_x_key(info *EXT_X_KEY) {
	// METHOD
	w.buf = append(w.buf, "METHOD="...)
	w.buf = append(w.buf, info.METHOD...)
	// URI
	if info.URI != "" {
		w.buf = append(w.buf, ",URI="...)
		w.buf = append(w.buf, info.URI...)
	}
	// IV
	if info.IV != "" {
		w.buf = append(w.buf, ",IV="...)
		w.buf = append(w.buf, info.IV...)
	}
	// KEYFORMAT
	if info.KEYFORMAT != "" {
		w.buf = append(w.buf, ",KEYFORMAT=\""...)
		w.buf = append(w.buf, info.KEYFORMAT...)
		w.buf = append(w.buf, '"')
	}
	// KEYFORMATVERSIONS
	if info.KEYFORMATVERSIONS != "" {
		w.buf = append(w.buf, ",KEYFORMATVERSIONS=\""...)
		w.buf = append(w.buf, info.KEYFORMATVERSIONS...)
		w.buf = append(w.buf, '"')
	}
	w.buf = append(w.buf, '\n')
}

type EXT_X_MAP struct {
	/*
		URI：由引号包裹的字符串，指定了包含媒体初始化块的资源的路径。该属性为必选参数。
	*/
	URI string
	/*
		BYTERANGE：由引号包裹的字符串，指定了媒体初始化块在 URI 指定的资源的位置（片段）。
		该属性指定的范围应当只包含媒体初始化块。
		该属性为可选参数，如果未指定，则表示 URI 指定的资源就是全部的媒体初始化块。
	*/
	BYTERANGE string
}

/*
该标签指明了获取媒体初始化块（Meida Initialization Section）的方法。
该标签对其后所有媒体片段生效，直至遇到另一个 EXT-X-MAP 标签。
其格式为：#EXT-X-MAP:<attribute-list>
*/
func (w *Writer) EXT_X_MAP(info *EXT_X_MAP) {
	w.buf = append(w.buf, "#EXT-X-MAP:"...)
	// URI
	w.buf = append(w.buf, "URI=\""...)
	w.buf = append(w.buf, info.URI...)
	w.buf = append(w.buf, '"')
	// BYTERANGE
	if info.BYTERANGE != "" {
		w.buf = append(w.buf, "BYTERANGE=\""...)
		w.buf = append(w.buf, info.BYTERANGE...)
		w.buf = append(w.buf, '"')
	}
	w.buf = append(w.buf, '\n')
}

/*
该标签使用一个绝对日期/时间表明第一个样本片段的取样时间。
其格式为：#EXT-X-PROGRAM-DATE-TIME:<date-time-msec>
其中，date-time-msec是一个 ISO/IEC 8601:2004 规定的日期格式，形如：YYYY-MM-DDThh:mm:ss.SSSZ。
*/
func (w *Writer) EXT_X_PROGRAM_DATE_TIME(time string) {
	w.buf = append(w.buf, "#EXT-X-PROGRAM-DATE-TIME:"...)
	w.buf = append(w.buf, time...)
	w.buf = append(w.buf, '\n')
}

type EXT_X_DATERANGE struct {
	/*
		ID：双引号包裹的唯一指明日期范围的标识。
		该属性为必选参数。
	*/
	ID string
	/*
		CLASS：双引号包裹的由客户定义的一系列属性与与之对应的语意值。
		所有拥有同一 CLASS 属性的日期范围必须遵守对应的语意。
		该属性为可选参数。
	*/
	CLASS string
	/*
		START-DATE：双引号包裹的日期范围起始值。
		该属性为必选参数。
	*/
	START_DATE string
	/*
		END-DATE：双引号包裹的日期范围结束值。
		该属性值必须大于或等于 START-DATE。
		该属性为可选参数。
	*/
	END_DATE string
	/*
		DURATION：日期范围的持续时间是一个十进制浮点型数值类型（单位：秒）。
		该属性值不能为负数。
		当表达立即时间时，将该属性值设为 0 即可。
		该属性为可选参数。
	*/
	DURATION float64
	/*
		PLANNED-DURATION：该属性为日期范围的期望持续时长。
		其值为一个十进制浮点数值类型（单位：秒）。
		该属性值不能为负数。
		在预先无法得知真实持续时长的情况下，可使用该属性作为日期范围的期望预估时长。
		该属性为可选参数。
	*/
	PLANNED_DURATION float64
}

/*
该标签定义了一系列由属性/值对组成的日期范围。
其格式为：#EXT-X-DATERANGE:<attribute-list>
*/
func (w *Writer) EXT_X_DATERANGE(info *EXT_X_DATERANGE) {
	w.buf = append(w.buf, "#EXT-X-DATERANGE:"...)
	// ID
	w.buf = append(w.buf, "ID=\""...)
	w.buf = append(w.buf, info.ID...)
	w.buf = append(w.buf, '"')
	// CLASS
	if info.CLASS != "" {
		w.buf = append(w.buf, ",CLASS=\""...)
		w.buf = append(w.buf, info.CLASS...)
		w.buf = append(w.buf, '"')
	}
	// START-DATE
	w.buf = append(w.buf, ",START-DATE=\""...)
	w.buf = append(w.buf, info.START_DATE...)
	w.buf = append(w.buf, '"')
	// END-DATE
	if info.END_DATE != "" {
		w.buf = append(w.buf, ",END-DATE=\""...)
		w.buf = append(w.buf, info.END_DATE...)
		w.buf = append(w.buf, '"')
	}
	// DURATION
	if info.DURATION > 0 {
		w.buf = append(w.buf, ",DURATION="...)
		w.buf = strconv.AppendFloat(w.buf, info.DURATION, 'f', 0, 64)
	}
	// PLANNED-DURATION
	if info.PLANNED_DURATION > 0 {
		w.buf = append(w.buf, ",PLANNED-DURATION="...)
		w.buf = strconv.AppendFloat(w.buf, info.PLANNED_DURATION, 'f', 0, 64)
	}
	w.buf = append(w.buf, '\n')
}

/*
该属性值为一个可枚举字符串，其值必须为YES。
该属性表明达到该范围末尾，也即等于后续范围的起始位置 START-DATE。后续范围是指具有相同 CLASS 的，在该标签 START-DATE 之后的具有最早 START-DATE 值的日期范围。
该属性时可选参数。
*/
func (w *Writer) EXT_X_END_ON_NEXT() {
	w.buf = append(w.buf, "#END-ON-NEXT:YES\n"...)
}

/*
表明流媒体类型。全局生效。
该标签为可选标签。
其格式为：#EXT-X-PLAYLIST-TYPE:<type-enum>
其中：type-enum可选值如下：

VOD：即 Video on Demand，表示该视屏流为点播源，因此服务器不能更改该 m3u8 文件；
EVENT：表示该视频流为直播源，因此服务器不能更改或删除该文件任意部分内容（但是可以在文件末尾添加新内容）。

注：VOD 文件通常带有 EXT-X-ENDLIST 标签，因为其为点播源，不会改变；而 EVEVT 文件初始化时一般不会有 EXT-X-ENDLIST 标签，
暗示有新的文件会添加到播放列表末尾，因此也需要客户端定时获取该 m3u8 文件，以获取新的媒体片段资源，直到访问到 EXT-X-ENDLIST 标签才停止）。
*/
func (w *Writer) EXT_X_PLAYLIST_TYPE(_type string) {
	w.buf = append(w.buf, "#EXT-X-PLAYLIST-TYPE:"...)
	w.buf = append(w.buf, _type...)
	w.buf = append(w.buf, '\n')
}

/*
该标签表示每个媒体片段都是一个 I-frame。I-frames 帧视屏编码不依赖于其他帧数，因此可以通过 I-frame 进行快速播放，急速翻转等操作。
该标签全局生效。
其格式为：#EXT-X-I-FRAMES-ONLY

如果播放列表设置了 EXT-X-I-FRAMES-ONLY，那么切片的时长（EXTINF 标签的值）即为当前切片 I-frame 帧开始到下一个 I-frame 帧出现的时长。
媒体资源如果包含 I-frame 切片，那么必须提供媒体初始化块或者通过 EXT-X-MAP 标签提供媒体初始化块的获取途径，这样客户端就能通过这些 I-frame 切片以任意顺序进行加载和解码。
如果 I-frame 切片设置了 EXT-BYTERANGE，那么就绝对不能提供媒体初始化块。
使用 EXT-X-I-FRAMES-ONLY 要求的兼容版本号 EXT-X-VERSION 大于等于 4。
*/
func (w *Writer) EXT_X_I_FRAMES_ONLY() {
	w.buf = append(w.buf, "#EXT-X-I-FRAMES-ONLY\n"...)
}

type EXT_X_MEDIA struct {
	/*
		TYPE：该属性值为一个可枚举字符串。
		其值有如下四种：AUDIO，VIDEO，SUBTITLES，CLOSED-CAPTIONS。
		通常使用的都是CLOSED-CAPTIONS。
		该属性为必选参数。
	*/
	TYPE string
	/*
		URI：双引号包裹的媒体资源播放列表路径。
		如果 TYPE 属性值为 CLOSED-CAPTIONS，那么则不能提供 URI。
		该属性为可选参数。
	*/
	URI string
	/*
		GROUP_ID：双引号包裹的字符串，表示多语言翻译流所属组。
		该属性为必选参数。
	*/
	GROUP_ID string
	/*
		LANGUAGE：双引号包裹的字符串，用于指定流主要使用的语言。
		该属性为可选参数。
	*/
	LANGUAGE string
	/*
		ASSOC_LANGUAGE：双引号包裹的字符串，其内包含一个语言标签，用于提供多语言流的其中一种语言版本。
		该参数为可选参数。
	*/
	ASSOC_LANGUAGE string
	/*
		NAME：双引号包裹的字符串，用于为翻译流提供可读的描述信息。
		如果设置了 LANGUAGE 属性，那么也应当设置 NAME 属性。
		该属性为必选参数。

	*/
	NAME string
	/*
		DEFAULT：该属性值为一个可枚举字符串。
		可选值为YES和NO。
		该属性未指定时默认值为NO。
		如果该属性设为YES，那么客户端在缺乏其他可选信息时应当播放该翻译流。
		该属性为可选参数。
	*/
	DEFAULT string
	/*
		AUTOSELECT：该属性值为一个可枚举字符串。
		其有效值为YES或NO。
		未指定时，默认设为NO。
		如果该属性设置YES，那么客户端在用户没有显示进行设置时，可以选择播放该翻译流，因为其能配置当前播放环境，比如系统语言选择。
		如果设置了该属性，那么当 DEFAULT 设置YES时，该属性也必须设置为YES。
		该属性为可选参数。
	*/
	AUTOSELECT string
	/*
		FORCED：该属性值为一个可枚举字符串。
		其有效值为YES或NO。
		未指定时，默认设为NO。
		只有在设置了 TYPE 为 SUBTITLES 时，才可以设置该属性。
		当该属性设为YES时，则暗示该翻译流包含重要内容。当设置了该属性，客户端应当选择播放匹配当前播放环境最佳的翻译流。
		当该属性设为NO时，则表示该翻译流内容意图用于回复用户显示进行请求。
		该属性为可选参数。
	*/
	FORCED string
	/*
		INSTREAM_ID：由双引号包裹的字符串，用于指示切片的语言（Rendition）版本。
		当 TYPE 设为 CLOSED-CAPTIONS 时，必须设置该属性。
		其可选值为："CC1", "CC2", "CC3", "CC4" 和 "SERVICEn"（n的值为 1~63）。
		对于其他 TYPE 值，该属性绝不能进行设置。
	*/
	INSTREAM_ID string
	/*
		CHARACTERISTICS：由双引号包裹的由一个或多个由逗号分隔的 UTI 构成的字符串。
		每个 UTI 表示一种翻译流的特征。
		该属性可包含私有 UTI。
		该属性为可选参数。
	*/
	CHARACTERISTICS string
	/*
		CHANNELS：由双引号包裹的有序，由反斜杠/分隔的参数列表组成的字符串。
		所有音频 EXT-X-MEDIA 标签应当都设置 CHANNELS 属性。
		如果主播放列表包含两个相同编码但是具有不同数目 channed 的翻译流，则必须设置 CHANNELS 属性；否则，CHANNELS 属性为可选参数。
	*/
	CHANNELS string
}

/*
用于指定相同内容的可替换的多语言翻译播放媒体列表资源。
比如，通过三个 EXT-X-MEIDA 标签，可以提供包含英文，法语和西班牙语版本的相同内容的音频资源，或者通过两个 EXT-X-MEDIA 提供两个不同拍摄角度的视屏资源。
其格式为：#EXT-X-MEDIA:<attribute-list>
*/
func (w *Writer) EXT_X_MEDIA(info *EXT_X_MEDIA) {
	w.buf = append(w.buf, "#EXT-X-MEDIA:"...)
	// TYPE
	w.buf = append(w.buf, "TYPE="...)
	w.buf = append(w.buf, info.TYPE...)
	// URI
	if info.URI != "" {
		w.buf = append(w.buf, ",URI=\""...)
		w.buf = append(w.buf, info.URI...)
		w.buf = append(w.buf, '"')
	}
	// LANGUAGE
	w.buf = append(w.buf, ",LANGUAGE=\""...)
	w.buf = append(w.buf, info.LANGUAGE...)
	w.buf = append(w.buf, '"')
	// ASSOC_LANGUAGE
	if info.ASSOC_LANGUAGE != "" {
		w.buf = append(w.buf, ",ASSOC_LANGUAGE=\""...)
		w.buf = append(w.buf, info.ASSOC_LANGUAGE...)
		w.buf = append(w.buf, '"')
	}
	// NAME
	w.buf = append(w.buf, ",NAME=\""...)
	w.buf = append(w.buf, info.NAME...)
	w.buf = append(w.buf, '"')
	// DEFAULT
	if info.DEFAULT != "" {
		w.buf = append(w.buf, ",DEFAULT=\""...)
		w.buf = append(w.buf, info.DEFAULT...)
		w.buf = append(w.buf, '"')
	}
	// FORCED
	if info.FORCED != "" {
		w.buf = append(w.buf, ",FORCED=\""...)
		w.buf = append(w.buf, info.FORCED...)
		w.buf = append(w.buf, '"')
	}
	// INSTREAM-ID
	if info.INSTREAM_ID != "" {
		w.buf = append(w.buf, ",INSTREAM-ID=\""...)
		w.buf = append(w.buf, info.INSTREAM_ID...)
		w.buf = append(w.buf, '"')
	}
	// CHARACTERISTICS
	if info.CHARACTERISTICS != "" {
		w.buf = append(w.buf, ",CHARACTERISTICS=\""...)
		w.buf = append(w.buf, info.CHARACTERISTICS...)
		w.buf = append(w.buf, '"')
	}
	// CHANNELS
	if info.CHANNELS != "" {
		w.buf = append(w.buf, ",CHANNELS=\""...)
		w.buf = append(w.buf, info.CHANNELS...)
		w.buf = append(w.buf, '"')
	}
	w.buf = append(w.buf, '\n')
}

type EXT_X_I_FRAME_STREAM_INF struct {
	/*
		BANDWIDTH：该属性为每秒传输的比特数，也即带宽。代表该备份流的巅峰速率。
		该属性为必选参数。
	*/
	BANDWIDTH int64
	/*
		AVERAGEB_ANDWIDTH：该属性为备份流的平均切片传输速率。
		该属性为可选参数。
	*/
	AVERAGEB_ANDWIDTH int64
	/*
		CODECS：双引号包裹的包含由逗号分隔的格式列表组成的字符串。
		每个 EXT-X-STREAM-INF 标签都应当携带 CODECS 属性。
	*/
	CODECS string
	/*
		RESOLUTION：该属性描述备份流视屏源的最佳像素方案。
		该属性为可选参数，但对于包含视屏源的备份流建议增加该属性设置。
	*/
	RESOLUTION string
	/*
		FRAME_RATE：该属性用一个十进制浮点型数值作为描述备份流所有视屏最大帧率。
		对于备份流中任意视屏源帧数超过每秒 30 帧的，应当增加该属性设置。
		该属性为可选参数，但对于包含视屏源的备份流建议增加该属性设置。
	*/
	FRAME_RATE float64
	/*
		HDCP_LEVEL：该属性值为一个可枚举字符串。
		其有效值为TYPE-0或NONE。
		值为TYPE-0表示该备份流可能会播放失败，除非输出被高带宽数字内容保护（HDCP）。
		值为NONE表示流内容无需输出拷贝保护。
		使用不同程度的 HDCP 加密备份流应当使用不同的媒体加密密钥。
		该属性为可选参数。在缺乏 HDCP 可能存在播放失败的情况下，应当提供该属性。
	*/
	HDCP_LEVEL string
	/*
		AUDIO：属性值由双引号包裹，其值必须与定义在主播放列表某处的设置了 TYPE 属性值为 AUDIO 的 EXT-X-MEDIA 标签的 GROUP-ID 属性值相匹配。
		该属性为可选参数。
	*/
	AUDIO string
	/*
		VIDEO：属性值由双引号包裹，其值必须与定义在主播放列表某处的设置了 TYPE 属性值为 VIDEO 的 EXT-X-MEDIA 标签的 GROUP-ID 属性值相匹配。
		该属性为可选参数。
	*/
	VIDEO string
	/*
		SUBTITLES：属性值由双引号包裹，其值必须与定义在主播放列表某处的设置了 TYPE 属性值为 SUBTITLES 的 EXT-X-MEDIA 标签的 GROUP-ID 属性值相匹配。
		该属性为可选参数。
	*/
	SUBTITLES string
	/*
		CLOSED_CAPTIONS：该属性值可以是一个双引号包裹的字符串或NONE。
		如果其值为一个字符串，则必须与定义在主播放列表某处的设置了 TYPE 属性值为 CLOSED-CAPTIONS 的 EXT-X-MEDIA 标签的 GROUP-ID 属性值相匹配。
		如果其值为NONE，则所有的 ext-x-stream-info 标签必须同样将该属性设置NONE，表示主播放列表备份流均没有关闭的标题。
		对于某个备份流具备关闭标题，另一个备份流不具备关闭标题可能会触发播放中断。
		该属性为可选参数。
	*/
	CLOSED_CAPTIONS string
}

/*
该标签的属性列表包含了 EXT-X-STREAM-INF 标签同样的属性列表选项，除了 FRAME-RATE，AUDIO，SUBTITLES 和 CLOSED-CAPTIONS
其格式为：

#EXT-X-I-FRAME-STREAM-INF:<attribute-list>
<URI>

URI 指定的媒体播放列表携带了该标签指定的翻译备份源。
URI 为必选参数。
*/
func (w *Writer) EXT_X_I_FRAME_STREAM_INF(info *EXT_X_I_FRAME_STREAM_INF, uri string) {
	w.buf = append(w.buf, "#EXT-X-STREAM-INF:"...)
	// BANDWIDTH
	w.buf = append(w.buf, "BANDWIDTH="...)
	w.buf = strconv.AppendInt(w.buf, info.BANDWIDTH, 64)
	// AVERAGE-BANDWIDTH
	if info.AVERAGEB_ANDWIDTH > 0 {
		w.buf = append(w.buf, ",AVERAGE-BANDWIDTH="...)
		w.buf = strconv.AppendInt(w.buf, info.AVERAGEB_ANDWIDTH, 64)
	}
	// CODECS
	if info.CODECS != "" {
		w.buf = append(w.buf, ",CODECS=\""...)
		w.buf = append(w.buf, info.CODECS...)
		w.buf = append(w.buf, '"')
	}
	// RESOLUTION
	if info.RESOLUTION != "" {
		w.buf = append(w.buf, ",RESOLUTION="...)
		w.buf = append(w.buf, info.RESOLUTION...)
	}
	// FRAME-RATE
	if info.FRAME_RATE > 0 {
		w.buf = append(w.buf, ",FRAME-RATE="...)
		w.buf = strconv.AppendFloat(w.buf, info.FRAME_RATE, 'f', 0, 64)
	}
	// HDCP-LEVEL
	if info.HDCP_LEVEL != "" {
		w.buf = append(w.buf, ",HDCP-LEVEL="...)
		w.buf = append(w.buf, info.HDCP_LEVEL...)
	}
	// AUDIO
	if info.AUDIO != "" {
		w.buf = append(w.buf, ",AUDIO=\""...)
		w.buf = append(w.buf, info.AUDIO...)
		w.buf = append(w.buf, '"')
	}
	// VIDEO
	if info.VIDEO != "" {
		w.buf = append(w.buf, ",VIDEO=\""...)
		w.buf = append(w.buf, info.AUDIO...)
		w.buf = append(w.buf, '"')
	}
	// SUBTITLES
	if info.SUBTITLES != "" {
		w.buf = append(w.buf, ",SUBTITLES="...)
		w.buf = append(w.buf, info.SUBTITLES...)
	}
	// CLOSED-CAPTIONS
	if info.CLOSED_CAPTIONS != "" {
		if info.CLOSED_CAPTIONS == "NONE" {
			w.buf = append(w.buf, ",CLOSED-CAPTIONS="...)
			w.buf = append(w.buf, info.CLOSED_CAPTIONS...)
		} else {
			w.buf = append(w.buf, ",CLOSED-CAPTIONS=\""...)
			w.buf = append(w.buf, info.CLOSED_CAPTIONS...)
			w.buf = append(w.buf, '"')
		}
	}
	// uri
	w.buf = append(w.buf, '\n')
	w.buf = append(w.buf, uri...)
	w.buf = append(w.buf, '\n')
}

type EXT_X_STREAM_INF struct {
	/*
		BANDWIDTH：该属性为每秒传输的比特数，也即带宽。代表该备份流的巅峰速率。
		该属性为必选参数。
	*/
	BANDWIDTH int64
	/*
		AVERAGEB_ANDWIDTH：该属性为备份流的平均切片传输速率。
		该属性为可选参数。
	*/
	AVERAGEB_ANDWIDTH int64
	/*
		CODECS：双引号包裹的包含由逗号分隔的格式列表组成的字符串。
		每个 EXT-X-STREAM-INF 标签都应当携带 CODECS 属性。
	*/
	CODECS string
	/*
		RESOLUTION：该属性描述备份流视屏源的最佳像素方案。
		该属性为可选参数，但对于包含视屏源的备份流建议增加该属性设置。
	*/
	RESOLUTION string
	/*
		HDCP_LEVEL：该属性值为一个可枚举字符串。
		其有效值为TYPE-0或NONE。
		值为TYPE-0表示该备份流可能会播放失败，除非输出被高带宽数字内容保护（HDCP）。
		值为NONE表示流内容无需输出拷贝保护。
		使用不同程度的 HDCP 加密备份流应当使用不同的媒体加密密钥。
		该属性为可选参数。在缺乏 HDCP 可能存在播放失败的情况下，应当提供该属性。
	*/
	HDCP_LEVEL string
	/*
		VIDEO：属性值由双引号包裹，其值必须与定义在主播放列表某处的设置了 TYPE 属性值为 VIDEO 的 EXT-X-MEDIA 标签的 GROUP-ID 属性值相匹配。
		该属性为可选参数。
	*/
	VIDEO string
}

/*
该属性指定了一个备份源。该属性值提供了该备份源的相关信息。
其格式为：

#EXT-X-STREAM-INF:<attribute-list>
<URI>

URI 指定的媒体播放列表携带了该标签指定的翻译备份源。
URI 为必选参数。
*/
func (w *Writer) EXT_X_STREAM_INF(info *EXT_X_I_FRAME_STREAM_INF, uri string) {
	w.buf = append(w.buf, "#EXT-X-I-FRAME-STREAM-INF:"...)
	// BANDWIDTH
	w.buf = append(w.buf, "BANDWIDTH="...)
	w.buf = strconv.AppendInt(w.buf, info.BANDWIDTH, 64)
	// AVERAGE-BANDWIDTH
	if info.AVERAGEB_ANDWIDTH > 0 {
		w.buf = append(w.buf, ",AVERAGE-BANDWIDTH="...)
		w.buf = strconv.AppendInt(w.buf, info.AVERAGEB_ANDWIDTH, 64)
	}
	// CODECS
	if info.CODECS != "" {
		w.buf = append(w.buf, ",CODECS=\""...)
		w.buf = append(w.buf, info.CODECS...)
		w.buf = append(w.buf, '"')
	}
	// RESOLUTION
	if info.RESOLUTION != "" {
		w.buf = append(w.buf, ",RESOLUTION="...)
		w.buf = append(w.buf, info.RESOLUTION...)
	}
	// FRAME-RATE
	if info.FRAME_RATE > 0 {
		w.buf = append(w.buf, ",FRAME-RATE="...)
		w.buf = strconv.AppendFloat(w.buf, info.FRAME_RATE, 'f', 0, 64)
	}
	// HDCP-LEVEL
	if info.HDCP_LEVEL != "" {
		w.buf = append(w.buf, ",HDCP-LEVEL="...)
		w.buf = append(w.buf, info.HDCP_LEVEL...)
	}
	// AUDIO
	if info.AUDIO != "" {
		w.buf = append(w.buf, ",AUDIO=\""...)
		w.buf = append(w.buf, info.AUDIO...)
		w.buf = append(w.buf, '"')
	}
	// VIDEO
	if info.VIDEO != "" {
		w.buf = append(w.buf, ",VIDEO=\""...)
		w.buf = append(w.buf, info.AUDIO...)
		w.buf = append(w.buf, '"')
	}
	// SUBTITLES
	if info.SUBTITLES != "" {
		w.buf = append(w.buf, ",SUBTITLES="...)
		w.buf = append(w.buf, info.SUBTITLES...)
	}
	// CLOSED-CAPTIONS
	if info.CLOSED_CAPTIONS != "" {
		if info.CLOSED_CAPTIONS == "NONE" {
			w.buf = append(w.buf, ",CLOSED-CAPTIONS="...)
			w.buf = append(w.buf, info.CLOSED_CAPTIONS...)
		} else {
			w.buf = append(w.buf, ",CLOSED-CAPTIONS=\""...)
			w.buf = append(w.buf, info.CLOSED_CAPTIONS...)
			w.buf = append(w.buf, '"')
		}
	}
	// uri
	w.buf = append(w.buf, '\n')
	w.buf = append(w.buf, uri...)
	w.buf = append(w.buf, '\n')
}

type EXT_X_SESSION_DATA struct {
	/*
		DATA-ID：由双引号包裹的字符串，代表一个特定的数据值。
		该属性应当使用反向 DNS 进行命名，如"com.example.movie.title"。然而，由于没有中央注册机构，所以可能出现冲突情况。
		该属性为必选参数。
	*/
	DATA_ID string
	/*
		VALUE：该属性值为一个双引号包裹的字符串，其包含 DATA-ID 指定的值。
		如果设置了 LANGUAGE，则 VALUE 应当包含一个用该语言书写的可读字符串。
	*/
	VALUE string
	/*
		URI：由双引号包裹的 URI 字符串。由该 URI 指示的资源必选使用 JSON 格式，否则，客户端可能会解析失败。
	*/
	URI string
	/*
		LANGUAGE：由双引号包裹的，包含一个语言标签的字符串。指示了 VALUE 所使用的语言。
	*/
	LANGUAGE string
}

/*
该标签允许主播放列表携带任意 session 数据。
该标签为可选参数。
其格式为：

#EXT-X-SESSION-DATA:<attribute-list>
*/

func (w *Writer) EXT_X_SESSION_DATA(info *EXT_X_SESSION_DATA) {
	w.buf = append(w.buf, "#EXT-X-SESSION-DATA:"...)
	// DATA_ID
	w.buf = append(w.buf, "DATA-ID=\""...)
	w.buf = append(w.buf, info.DATA_ID...)
	w.buf = append(w.buf, '"')
	// VALUE
	if info.VALUE != "" {
		w.buf = append(w.buf, ",VALUE=\""...)
		w.buf = append(w.buf, info.VALUE...)
		w.buf = append(w.buf, '"')
	}
	// URI
	if info.URI != "" {
		w.buf = append(w.buf, ",URI=\""...)
		w.buf = append(w.buf, info.URI...)
		w.buf = append(w.buf, '"')
	}
	// LANGUAGE
	if info.LANGUAGE != "" {
		w.buf = append(w.buf, ",LANGUAGE=\""...)
		w.buf = append(w.buf, info.LANGUAGE...)
		w.buf = append(w.buf, '"')
	}
	w.buf = append(w.buf, '\n')
}

/*
该标签允许主播放列表（Master Playlist）指定媒体播放列表（Meida Playlist）的加密密钥。这使得客户端可以预先加载这些密钥，而无需从媒体播放列表中获取。
该标签为可选参数。
其格式为：
#EXT-X-SESSION-KEY:<attribute-list>
其属性列表与 EXT-X-KEY 相同，除了 METHOD 属性的值必须不为NONE。
*/
func (w *Writer) EXT_X_SESSION_KEY(info *EXT_X_KEY) {
	w.buf = append(w.buf, "#EXT-X-SESSION-KEY:"...)
	w.ext_x_key(info)
}

/*
该标签表明对于一个媒体片段中的所有媒体样本均可独立进行解码，而无须依赖其他媒体片段信息。
该标签对列表内所有媒体片段均有效。
其格式为：

#EXT-X-INDEPENDENT-SEGMENTS
如果该标签出现在主播放列表中，则其对所有媒体播放列表的所有媒体片段都生效。
*/
func (w *Writer) EXT_X_INDEPENDENT_SEGMENTS() {
	w.buf = append(w.buf, "#EXT-X-INDEPENDENT-SEGMENTS\n"...)
}

type EXT_X_START struct {
	/*
		TIME-OFFSET：该属性值为一个带符号十进制浮点数（单位：秒）。
		一个正数表示以播放列表起始位置开始的时间偏移量。
		一个负数表示播放列表上一个媒体片段最后位置往前的时间偏移量。
		该属性的绝对值应当不超过播放列表的时长。如果超过，则表示到达文件结尾（数值为正数），或者达到文件起始（数值为负数）。
		如果播放列表不包含 EXT-X-ENDLIST 标签，那么 TIME-OFFSET 属性值不应当在播放文件末尾三个切片时长之内。
	*/
	TIME_OFFSET float64
	/*
	   PRECISE：该值为一个可枚举字符串。
	   有效的取值为YES 或 NO。
	   如果值为YES，客户端应当播放包含 TIME-OFFSET 的媒体片段，但不要渲染该块内优先于 TIME-OFFSET 的样本块。
	   如果值为NO，客户端应当尝试渲染在媒体片段内的所有样本块。
	   该属性为可选参数，未指定则认为NO。
	*/
	PRECISE string
}

/*
该标签表示播放列表播放起始位置。
默认情况下，客户端开启一个播放会话时，应当使用该标签指定的位置进行播放。
该标签为可选标签。
其格式为：

#EXT-X-START:<attribute-list>
*/
func (w *Writer) EXT_X_START(info *EXT_X_START) {
	w.buf = append(w.buf, "#EXT-X-START:"...)
	// TIME-OFFSET
	w.buf = append(w.buf, "TIME-OFFSET="...)
	w.buf = strconv.AppendFloat(w.buf, info.TIME_OFFSET, 'f', 0, 64)
	// PRECISE
	if info.PRECISE != "" {
		w.buf = append(w.buf, ",PRECISE="...)
		w.buf = append(w.buf, info.PRECISE...)
	}
	w.buf = append(w.buf, '\n')
}
