package main

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"io"
	"log"
	"os"
	"path/filepath"
	"raytracer/objects"
	"raytracer/ray"
)

func WriteColor(color *mat.VecDense, buffer io.Writer) {
	r := color.AtVec(0)
	g := color.AtVec(1)
	b := color.AtVec(2)

	ir := int(255.999 * r)
	ig := int(255.999 * g)
	ib := int(255.999 * b)

	str := fmt.Sprintf("%d %d %d\n", ir, ig, ib)

	_, err := buffer.Write([]byte(str))

	if err != nil {
		log.Fatal("Failed to write to buffer: ", err)
	}

}

func main() {
	aspectRatio := 16. / 9.
	imageHeight := 2160
	imageWidth := int(aspectRatio * float64(imageHeight))

	// World
	world := ray.NewHittableList()
	sphere1 := objects.NewSphere(mat.NewVecDense(3, []float64{0, 0, -1}), 0.5)
	world.Add(&sphere1)

	sphere2 := objects.NewSphere(mat.NewVecDense(3, []float64{0, -100.5, -1}), 100)
	world.Add(&sphere2)

	// Camera
	focalLength := 1.
	focalVec := mat.NewVecDense(3, []float64{0, 0, focalLength})

	viewportHeight := 2.
	viewportWidth := viewportHeight * (float64(imageWidth) / float64(imageHeight))
	cameraCenter := mat.NewVecDense(3, nil)

	viewportU := mat.NewVecDense(3, []float64{viewportWidth, 0, 0})
	viewportV := mat.NewVecDense(3, []float64{0, -viewportHeight, 0})

	pixelDeltaU := mat.NewVecDense(3, nil)
	pixelDeltaU.ScaleVec(1/float64(imageWidth), viewportU)

	pixelDeltaV := mat.NewVecDense(3, nil)
	pixelDeltaV.ScaleVec(1/float64(imageHeight), viewportV)

	halfViewportU := mat.NewVecDense(3, nil)
	halfViewportU.ScaleVec(1/2., viewportU)

	halfViewportV := mat.NewVecDense(3, nil)
	halfViewportV.ScaleVec(1/2., viewportV)

	viewportUpperLeft := mat.NewVecDense(3, nil)
	viewportUpperLeft.SubVec(cameraCenter, focalVec)
	viewportUpperLeft.SubVec(viewportUpperLeft, halfViewportU)
	viewportUpperLeft.SubVec(viewportUpperLeft, halfViewportV)

	offset := mat.NewVecDense(3, nil)
	offset.AddVec(pixelDeltaU, pixelDeltaV)
	offset.ScaleVec(0.5, offset)

	startPixel := mat.NewVecDense(3, nil)
	startPixel.AddVec(viewportUpperLeft, offset)

	file := getFile()

	_, err := file.WriteString(fmt.Sprintf("P3\n%d %d\n255\n", imageWidth, imageHeight))
	if err != nil {
		log.Fatal("Failed to write the header of ppm", err)
	}

	for j := 0; j < imageHeight; j++ {
		log.Println("Scanlines remaining: ", imageHeight-j)
		for i := 0; i < imageWidth; i++ {
			horizontalOffset := mat.NewVecDense(3, nil)
			horizontalOffset.ScaleVec(float64(i), pixelDeltaU)

			verticalOffset := mat.NewVecDense(3, nil)
			verticalOffset.ScaleVec(float64(j), pixelDeltaV)

			pixelCenter := mat.NewVecDense(3, nil)
			pixelCenter.AddVec(startPixel, horizontalOffset)
			pixelCenter.AddVec(pixelCenter, verticalOffset)

			rayDirection := mat.NewVecDense(3, nil)
			rayDirection.SubVec(pixelCenter, cameraCenter)

			currentRay := ray.Ray{
				Origin:    cameraCenter,
				Direction: rayDirection,
			}

			pixelColor := currentRay.Color(&world)
			WriteColor(pixelColor, file)
		}
	}

	log.Println("Done.")
}

func getFile() *os.File {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal("Couldn't get the current working path", err)
	}

	filePath := filepath.Join(currentPath, "image.ppm")
	if len(os.Args) == 2 {
		filePath = filepath.Join(currentPath, os.Args[1])
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return file
}
