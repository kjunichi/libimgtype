package main

import (
	"image"
	"image/color"
	"regexp"
	"unsafe"

	"github.com/nfnt/resize"
	termbox "github.com/nsf/termbox-go"
)

/*
#include <stdlib.h>

extern void termPutImage(unsigned char*,int,int);
extern void termInit();
extern void termFlush();
extern void termClose();
extern void termWait();

inline static void foo() {
	unsigned char *data = (unsigned char*)malloc(sizeof(char)*3*128*128);
	for(int i=0; i< 3*128*128; i++) {
		//if(i%3==0) {
		data[i]=(unsigned char)i%255;
		//}
	}
	termInit();
	termPutImage(data,128,128);
	termFlush();
	termWait();
	termClose();
}
*/
import "C"

func isPng(name string) bool {
	r := regexp.MustCompile(`.png$`)
	if r.MatchString(name) {
		return true
	}
	return false
}
func getTypeFromPath(path string) string {
	if isPng(path) {
		return "png"
	}
	return "jpeg"
}

func rgb2hsv(r, g, b uint32) (h, s, v int32) {
	max := r
	min := r
	if max < g {
		max = g
	}
	if max < b {
		max = b
	}
	if min > g {
		min = g
	}
	if min > b {
		min = b
	}
	v = int32(max)
	if max == 0 {
		s = 0
	} else {
		s = int32((3 * (max - min)) / max)
	}
	if max == min {
		h = -1
		//fmt.Println("rgb2hsv = :", h, s, v)
		return
	}
	if max == r {
		h = int32((60.0 * (int32(g) - int32(b))) / int32(max-min))
	} else if max == g {
		h = int32((60.0 * (int32(b) - int32(r)) / int32(max-min)) + 120.0)
	} else {
		h = int32((60.0 * (int32(r) - int32(g)) / int32(max-min)) + 240.0)
	}

	if h < 0 {
		h += 360
	}

	return
}

func getCellInfo(r, g, b uint32) (chi, cidx, cbidx int32) {
	//fmt.Printf("in : r,g,b = %d,%d,%d\n",r,g,b)
	chi = 10
	cbidx = 0
	h, s, v := rgb2hsv(r, g, b)
	//fmt.Printf("hsv : %d,%d,%d\n",h,s,v)
	if h >= 330 || (h < 30 && h >= 0) {
		// red
		cidx = 2
		if s > 2 {
			if h > 0 {
				chi = 11
				cbidx = 4
			}
			if h >= 330 {
				chi = 11
				cbidx = 6
			}
		} else {
			chi = 11
			if s == 2 {
				cbidx = 8
			} else {
				chi = 1
				cbidx = 8
				cidx = 2
			}
		}
	}
	if h >= 30 && h < 90 {
		// yellow
		cidx = 4
		if s > 2 {
			if h < 60 {
				chi = 11
				cbidx = 2
			}
			if h > 60 {
				chi = 11
				cbidx = 3
			}
		} else {
			chi = 11
			if s == 2 {
				cbidx = 8
			} else {
				chi = 1
				cidx = 4
				cbidx = 8
			}
		}
	}
	if h >= 90 && h < 150 {
		// green
		cidx = 3
		if s > 2 {
			if h < 120 {
				chi = 11
				cbidx = 4
			}
			if h > 120 {
				chi = 11
				cbidx = 7
			}
		} else {
			chi = 11
			if s == 2 {
				cbidx = 8
			} else {
				chi = 1
				cbidx = 8
				cidx = 3
			}
		}
	}
	if h >= 150 && h <= 210 {
		// cyan
		cidx = 7
		if s > 2 {
			if h < 180 {
				chi = 11
				cbidx = 3
			}
			if h > 180 {
				chi = 11
				cbidx = 5
			}
		} else {
			chi = 11
			if s == 2 {
				cbidx = 8
			} else {
				chi = 1
				cbidx = 8
				cidx = 7
			}
		}
	}
	if h > 210 && h < 270 {
		// blue
		cidx = 5
		if s > 2 {
			if h < 240 {
				chi = 11
				cbidx = 7
			}
			if h > 240 {
				chi = 11
				cbidx = 6
			}
		} else {
			chi = 11
			if s == 2 {
				cbidx = 8
			} else {
				chi = 1
				cbidx = 8
				cidx = 5
			}
		}
	}
	if h >= 270 && h < 330 {
		// magenta
		cidx = 6
		if s > 2 {
			if h < 300 {
				chi = 11
				cbidx = 5
			}
			if h > 300 {
				chi = 11
				cbidx = 2
			}
		} else {
			chi = 11
			if s == 2 {
				cbidx = 8
			} else {
				chi = 1
				cbidx = 8
				cidx = 6
			}
		}
	}
	if h < 0 {
		chi = 10
		if v == 0 {
			cidx = 1
		}
		if v == 1 {
			chi = 0
			cidx = 8
			cbidx = 8
		}
		if v == 2 {
			chi = 1
			cbidx = 8 + 512
			cidx = 8
		}
		if v == 3 {
			chi = 0
			cidx = 0
			cbidx = 8 + 512
		}

		return
	}
	if v > 2 {
		cidx += 512
		if cbidx > 0 {
			cbidx += 512
		}
	}
	if v == 2 && s == 3 && chi == 10 {
		chi = 12
		cbidx = cidx + 512
	}
	//fmt.Printf("out : chi, cidx, cbidx = %d,%d,%d\n",chi,cidx,cbidx)
	return
}

//export termInit
func termInit() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
}

//export termFlush
func termFlush() {
	termbox.Flush()
}

//export termClose
func termClose() {
	termbox.Close()
}

//export termWait
func termWait() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				//Escキー押した時の処理
				return
			case termbox.KeyArrowUp:
				//上キー押した時の処理
			case termbox.KeyArrowDown:
				//下キー押した時の処理
			default:
			}
		default:
		}
	}
}

//export termPutImage
func termPutImage(imageData *C.uchar, width, height C.int) {
	len := 3 * width * height
	slice := (*[1 << 30]C.uchar)(unsafe.Pointer(imageData))[:len:len]

	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	for h := 0; h < int(height); h++ {
		for w := 0; w < int(width); w++ {
			r := uint8(slice[3*(h*int(width)+w)])
			g := uint8(slice[3*(h*int(width)+w)+1])
			b := uint8(slice[3*(h*int(width)+w)+2])
			c := color.RGBA{r, g, b, 0}
			img.SetRGBA(w, h, c)
		}
	}
	w, h := termbox.Size()
	imgText := resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	termbox.SetOutputMode(termbox.Output256)
	const mojiretsu = "█▓23456789 ░▒#" //{' ', '░', '▒', '▓', '█'}
	nihongoRune := []rune(mojiretsu)
	termbox.SetCell(0, 0, nihongoRune[1], termbox.ColorRed, termbox.ColorBlue)
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			r, g, b, _ := imgText.At(j, i).RGBA()
			if r > 0 {
				r = uint32((r - 1) * 4 / 65535)
			} else {
				r = 0
			}
			if g > 0 {
				g = uint32((g - 1) * 4 / 65535)
			} else {
				g = 0
			}
			if b > 0 {
				b = uint32((b - 1) * 4 / 65535)
			} else {
				b = 0
			}

			chi, cidx, cbidx := getCellInfo(r, g, b)
			//fmt.Println(r, g, b, chi, cidx, cbidx)
			termbox.SetCell(j, i, nihongoRune[chi], termbox.Attribute(cbidx), termbox.Attribute(cidx))
		}
	}
}

func main() {
	C.foo()
}
