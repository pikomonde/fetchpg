package tools

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
)

func Screenshot(page *rod.Page, e *rod.Element, path string) {
	a, _ := e.Eval(`
	()=>{
		x=this.getBoundingClientRect()
		console.log(this,  x)
		return {bottom: x.bottom, height: x.height, left:x.left, right:x.right, top:x.top, width:x.width, x:x.x, y:x.y}
	}
	`)

	opts := &proto.PageCaptureScreenshot{
		Format: "png",
		Clip: &proto.PageViewport{
			X:      a.Value.Get("x").Num(),
			Y:      a.Value.Get("y").Num(),
			Width:  a.Value.Get("width").Num(),
			Height: a.Value.Get("height").Num(),
			Scale:  1,
		},
		FromSurface: true,
	}
	bin, _ := page.Screenshot(true, opts)
	utils.OutputFile(path, bin)
}
