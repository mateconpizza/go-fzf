package menu

// Border draw border around the finder.
type Border string

const (
	BorderRounded    Border = "rounded"
	BorderSharp      Border = "sharp"
	BorderBold       Border = "bold"
	BorderDouble     Border = "double"
	BorderBlock      Border = "block"
	BorderThinblock  Border = "thinblock"
	BorderHorizontal Border = "horizontal"
	BorderVertical   Border = "vertical"
	BorderLine       Border = "line"
	BorderTop        Border = "top"
	BorderBottom     Border = "bottom"
	BorderLeft       Border = "left"
	BorderRight      Border = "right"
	BorderDashed     Border = "dashed"
	BorderNone       Border = "none"
)

func (b Border) Arg(s string) string {
	return s + "=" + string(b)
}

type Layout string

const (
	LayoutDefault     Layout = "default"      // LayoutDefault display from the bottom of the screen.
	LayoutReverse     Layout = "reverse"      // LayoutReverse display from the top of the screen.
	LayoutReverseList Layout = "reverse-list" // LayoutReverseList display from the top of the screen, prompt at the bottom.
)

// InfoStyle determines the display style of the finder info. (e.g. match
// counter, loading indicator, etc.)
type InfoStyle string

const (
	InfoStyleDefault     InfoStyle = "default"      // InfoStyleDefault on the left end of the horizontal separator.
	InfoStyleRight       InfoStyle = "right"        // InfoStyleRight on the right end of the horizontal separator.
	InfoStyleHidden      InfoStyle = "hidden"       // InfoStyleHidden do not display finder info.
	InfoStyleInline      InfoStyle = "inline"       // InfoStyleInline after the prompt with the default prefix ' < '.
	InfoStyleInlineRight InfoStyle = "inline-right" // InfoStyleInlineRight on the right end of the prompt line.
)

// ColorValue represents an ANSI color understood by fzf.
type ColorValue string

const (
	// ColorDefault default or original terminal color.
	ColorDefault ColorValue = "-1"

	ColorBlack         ColorValue = "black"
	ColorRed           ColorValue = "red"
	ColorGreen         ColorValue = "green"
	ColorYellow        ColorValue = "yellow"
	ColorBlue          ColorValue = "blue"
	ColorMagenta       ColorValue = "magenta"
	ColorCyan          ColorValue = "cyan"
	ColorWhite         ColorValue = "white"
	ColorBrightBlack   ColorValue = "bright-black"
	ColorBrightRed     ColorValue = "bright-red"
	ColorBrightGreen   ColorValue = "bright-green"
	ColorBrightYellow  ColorValue = "bright-yellow"
	ColorBrightBlue    ColorValue = "bright-blue"
	ColorBrightMagenta ColorValue = "bright-magenta"
	ColorBrightCyan    ColorValue = "bright-cyan"
	ColorBrightWhite   ColorValue = "bright-white"

	AttributeRegular         ColorValue = "regular"
	AttributeStrip           ColorValue = "strip"
	AttributeBold            ColorValue = "bold"
	AttributeUnderline       ColorValue = "underline"
	AttributeUnderlineDouble ColorValue = "underline-double"
	AttributeUnderlineCurly  ColorValue = "underline-curly"
	AttributeUnderlineDotted ColorValue = "underline-dotted"
	AttributeUnderlineDashed ColorValue = "underline-dashed"
	AttributeReverse         ColorValue = "reverse"
	AttributeDim             ColorValue = "dim"
	AttributeItalic          ColorValue = "italic"
	AttributeStrikethrough   ColorValue = "strikethrough"
)
