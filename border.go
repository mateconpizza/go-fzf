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
	BorderNone       Border = "none"
)

func (b Border) Arg(s string) string {
	return s + "=" + string(b)
}
