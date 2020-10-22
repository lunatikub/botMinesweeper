package bot

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/kbinani/screenshot"
)

var blockColor = []color.RGBA{
	{222, 222, 220, 255}, // 0
	{221, 250, 195, 255}, // 1
	{236, 237, 191, 255}, // 2
	{237, 218, 180, 255}, // 3
	{237, 195, 138, 255}, // 4
	{247, 161, 162, 255}, // 5
	{254, 167, 133, 255}, // 6
	{255, 125, 96, 255},  // 7
	{255, 50, 60, 255},   // 8
	{186, 189, 182, 255}, // unreveal
	{46, 52, 54, 255},    // mine
}

// Recognition a grid configuration for the gnome-mines game
type recognition struct {
	x, y    int             // top left cell
	sz      int             // cell width/height
	spacing int             // space between 2 cells
	h, w    int             // width/height (number of cell))
	bounds  image.Rectangle // bounds of the screen
}

func (r *recognition) getBounds(screenID int) {
	screenshot.NumActiveDisplays()
	r.bounds = screenshot.GetDisplayBounds(screenID)
	log.Printf("[recognition] screenID:%d, height:%d, widht:%d",
		screenID, r.bounds.Dy(), r.bounds.Dx())
}

// get the top left cell
func (r *recognition) getTopLeftCell(img *image.RGBA) error {
	for y := 0; y < r.bounds.Dy(); y++ {
		for x := 0; x < r.bounds.Dx(); x++ {
			if blockColor[covered] == img.At(x, y) {
				r.x = x
				r.y = y
				log.Printf("[recognition] top left corner: {y:%d,x:%d}", r.y, r.x)
				return nil
			}
		}
	}
	return fmt.Errorf("cannot get top left cell")
}

// get the size (height/widht) of a cell
func (r *recognition) getCellSize(img *image.RGBA) error {
	x := r.x
	for {
		if x == r.bounds.Dx() {
			return fmt.Errorf("cannot get size of a cell")
		}
		if blockColor[covered] != img.At(x, r.y) {
			break
		}
		x++
	}
	r.sz = x - r.x
	log.Printf("[recognition] cell size:%d", r.sz)
	return nil
}

// get the space between 2 cells
func (r *recognition) getSpacing(img *image.RGBA) error {
	x := r.x + r.sz
	for {
		if x == r.bounds.Dx() {
			return fmt.Errorf("cannot get spacing between 2 cells")
		}
		if blockColor[covered] == img.At(x, r.y) {
			break
		}
		x++
	}
	r.spacing = x - (r.x + r.sz)
	log.Printf("[recognition] spacing:%d", r.spacing)
	return nil
}

// get the width number of cells
func (r *recognition) getHorizontalCell(img *image.RGBA) error {
	x := r.x
	shift := r.sz + r.spacing
	for {
		if blockColor[covered] != img.At(x, r.y) {
			break
		}
		r.w++
		x += shift
	}
	log.Printf("[recognition] width number of cells:%d", r.w)
	return nil
}

// get the height number of cells
func (r *recognition) getVerticalCell(img *image.RGBA) error {
	y := r.y
	shift := r.sz + r.spacing
	for {
		if blockColor[covered] != img.At(r.x, y) {
			break
		}
		r.h++
		y += shift
	}
	log.Printf("[recognition] height number of cells:%d", r.h)
	return nil
}

func (r *recognition) getDims() {
	var img *image.RGBA
	var err error
	if img, err = screenshot.CaptureRect(r.bounds); err != nil {
		panic(err)
	}
	if err = r.getTopLeftCell(img); err != nil {
		panic(err)
	}
	if err = r.getCellSize(img); err != nil {
		panic(err)
	}
	if err = r.getSpacing(img); err != nil {
		panic(err)
	}
	r.getHorizontalCell(img)
	r.getVerticalCell(img)
}

func (r *recognition) get(img *image.RGBA, y, x int) int {
	x = r.x + x*(r.sz+r.spacing) + 6
	y = r.y + y*(r.sz+r.spacing) + 6
	for i, v := range blockColor {
		if v == img.At(x, y) {
			return i
		}
	}
	return covered
}

// GetConfiguration get the current configuration of the grid

func (r *recognition) getConfiguration(g *grid) {
	img, _ := screenshot.CaptureRect(r.bounds)
	for y := 0; y < r.h; y++ {
		for x := 0; x < r.w; x++ {
			if g.cells[y][x] == covered {
				g.set(y, x, r.get(img, y, x))
			}
		}
	}
}

// Create a new instance of the recognition
func newRecognition(screenID int) *recognition {
	r := new(recognition)
	r.getBounds(screenID)
	r.getDims()
	return r
}
