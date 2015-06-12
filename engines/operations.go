package engines

type Operation struct {
	Name string
}

var Resize = &Operation{
	"resize",
}

var Thumbnail = &Operation{
	"thumbnail",
}

var Rotate = &Operation{
	"rotate",
}

var Flip = &Operation{
	"flip",
}

var Fill = &Operation{
	"fill",
}

var Operations = map[string]*Operation{
	Resize.Name:    Resize,
	Thumbnail.Name: Thumbnail,
	Flip.Name:      Flip,
	Rotate.Name:    Rotate,
	Fill.Name:      Fill,
}
