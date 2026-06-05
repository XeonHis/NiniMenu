package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	xdraw "golang.org/x/image/draw"
)

var sizes = []int{32, 64, 128, 180, 256, 512}

func main() {
	threshold := flag.Int("t", 240, "white threshold (0-255), higher = stricter")
	outDir := flag.String("o", "", "output directory (default: same as input)")
	noCrop := flag.Bool("no-crop", false, "skip cropping to content bounds after removing white edges")
	noTrim := flag.Bool("no-trim", false, "skip white edge removal and cropping, only generate resized icons")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: imgtool [options] <input_image>")
		fmt.Println()
		fmt.Println("Removes white edges from image (makes transparent),")
		fmt.Println("crops to content, and generates resized icons.")
		fmt.Println()
		fmt.Println("Sizes: 32, 64, 128, 180, 256, 512")
		fmt.Println()
		flag.PrintDefaults()
		os.Exit(1)
	}

	inputPath := flag.Arg(0)

	f, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding: %v\n", err)
		os.Exit(1)
	}

	rgba := toRGBA(img)

	if !*noTrim {
		removeWhiteEdges(rgba, *threshold)
		fmt.Printf("White edges removed (threshold=%d)\n", *threshold)

		if !*noCrop {
			cropped := cropToContent(rgba)
			if cropped != nil {
				rgba = cropped
				fmt.Printf("Cropped to content: %dx%d\n", rgba.Bounds().Dx(), rgba.Bounds().Dy())
			}
		}
	} else {
		fmt.Println("Skipping white edge removal and cropping (-no-trim)")
	}

	dir := *outDir
	if dir == "" {
		dir = filepath.Dir(inputPath)
	}
	os.MkdirAll(dir, 0755)

	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))

	if !*noTrim {
		outPath := filepath.Join(dir, base+"_trimmed.png")
		savePNG(rgba, outPath)
		fmt.Printf("Saved: %s (%dx%d)\n", outPath, rgba.Bounds().Dx(), rgba.Bounds().Dy())
	}

	for _, sz := range sizes {
		resized := resizeImage(rgba, sz, sz)
		outPath := filepath.Join(dir, fmt.Sprintf("%d.png", sz))
		savePNG(resized, outPath)
		fmt.Printf("Saved: %s\n", outPath)
	}

	fmt.Println("Done!")
}

func toRGBA(img image.Image) *image.RGBA {
	b := img.Bounds()
	rgba := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}
	return rgba
}

func isWhite(c color.Color, threshold int) bool {
	r, g, b, _ := c.RGBA()
	return int(r>>8) >= threshold && int(g>>8) >= threshold && int(b>>8) >= threshold
}

func removeWhiteEdges(img *image.RGBA, threshold int) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	visited := make([][]bool, h)
	for i := range visited {
		visited[i] = make([]bool, w)
	}

	type point struct{ x, y int }
	queue := []point{}

	for x := 0; x < w; x++ {
		if isWhite(img.At(x+b.Min.X, b.Min.Y), threshold) && !visited[0][x] {
			visited[0][x] = true
			queue = append(queue, point{x, 0})
		}
		if isWhite(img.At(x+b.Min.X, b.Min.Y+h-1), threshold) && !visited[h-1][x] {
			visited[h-1][x] = true
			queue = append(queue, point{x, h - 1})
		}
	}
	for y := 0; y < h; y++ {
		if isWhite(img.At(b.Min.X, y+b.Min.Y), threshold) && !visited[y][0] {
			visited[y][0] = true
			queue = append(queue, point{0, y})
		}
		if isWhite(img.At(b.Min.X+w-1, y+b.Min.Y), threshold) && !visited[y][w-1] {
			visited[y][w-1] = true
			queue = append(queue, point{w - 1, y})
		}
	}

	dx := [4]int{-1, 1, 0, 0}
	dy := [4]int{0, 0, -1, 1}

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		img.SetRGBA(p.x+b.Min.X, p.y+b.Min.Y, color.RGBA{0, 0, 0, 0})

		for i := 0; i < 4; i++ {
			nx, ny := p.x+dx[i], p.y+dy[i]
			if nx >= 0 && nx < w && ny >= 0 && ny < h && !visited[ny][nx] {
				if isWhite(img.At(nx+b.Min.X, ny+b.Min.Y), threshold) {
					visited[ny][nx] = true
					queue = append(queue, point{nx, ny})
				}
			}
		}
	}
}

func cropToContent(img *image.RGBA) *image.RGBA {
	b := img.Bounds()
	minX, minY := b.Max.X, b.Max.Y
	maxX, maxY := b.Min.X, b.Min.Y

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	if minX > maxX || minY > maxY {
		return nil
	}

	cropRect := image.Rect(minX, minY, maxX+1, maxY+1)
	cropped := image.NewRGBA(cropRect)
	for y := cropRect.Min.Y; y < cropRect.Max.Y; y++ {
		for x := cropRect.Min.X; x < cropRect.Max.X; x++ {
			cropped.Set(x, y, img.At(x, y))
		}
	}
	return cropped
}

func resizeImage(src *image.RGBA, w, h int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), xdraw.Over, nil)
	return dst
}

func savePNG(img *image.RGBA, path string) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	defer f.Close()

	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(f, img); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding PNG: %v\n", err)
	}
}

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}
