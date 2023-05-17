package htg

type ITag interface {
	String() string
}

type Tag struct {
	name       string
	attributes []Attribute
	children   []*Tag
	value      string

	start int
	end   int
	line  int
}

func (t *Tag) String() string {
	return t.name
}

type IAttribute interface {
	String() string
}

type Attribute struct {
	name  string
	value string

	start int
	end   int
	line  int
}

func (a *Attribute) String() string {
	return a.name + "=\"" + a.value + "\""
}

type IText interface {
	String() string
}

type Text struct {
	value string
}

func (t *Text) String() string {
	return t.value
}
