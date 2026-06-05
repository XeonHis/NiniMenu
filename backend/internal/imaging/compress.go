package imaging

import (
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"

	imgproc "github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
)

func ProcessUpload(srcPath, dstDir, baseName string, maxDim, quality int) (string, error) {
	img, err := imgproc.Open(srcPath, imgproc.AutoOrientation(true))
	if err != nil {
		return "", err
	}

	b := img.Bounds()
	if b.Dx() > maxDim || b.Dy() > maxDim {
		img = imgproc.Fit(img, maxDim, maxDim, imgproc.Lanczos)
		b = img.Bounds()
	}

	rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(rgba, rgba.Bounds(), image.White, image.Point{}, draw.Src)
	draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Over)

	dstPath := filepath.Join(dstDir, baseName+".jpg")
	f, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := jpeg.Encode(f, rgba, &jpeg.Options{Quality: quality}); err != nil {
		os.Remove(dstPath)
		return "", err
	}

	return dstPath, nil
}
