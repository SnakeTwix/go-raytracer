package main

import (
	"fmt"
	"log"
	"os"
	"raytracer/ray"
	"raytracer/vec3"
)

func main() {
	aspectRatio := 16. / 9.
	imageHeight := 1024
	imageWidth := int(aspectRatio * float64(imageHeight))

	// Camera
	focalLength := 1.
	focalVec := vec3.Vec3{
		X: 0,
		Y: 0,
		Z: focalLength,
	}
	viewportHeight := 2.
	viewportWidth := viewportHeight * (float64(imageWidth) / float64(imageHeight))
	cameraCenter := vec3.Point3{}

	viewportU := vec3.Vec3{X: viewportWidth}
	viewportV := vec3.Vec3{Y: -viewportHeight}

	pixelDeltaU := viewportU.DivNew(float64(imageWidth))
	pixelDeltaV := viewportV.DivNew(float64(imageHeight))

	halfViewportU := viewportU.DivNew(2.)
	halfViewportV := viewportV.DivNew(2.)

	viewportUpperLeft := cameraCenter.SubVecNew(&focalVec)
	viewportUpperLeft.SubVec(&halfViewportU)
	viewportUpperLeft.SubVec(&halfViewportV)

	offset := pixelDeltaU.AddVecNew(&pixelDeltaV)
	offset.Mul(0.5)

	startPixel := viewportUpperLeft.AddVecNew(&offset)

	stdout := os.Stdout

	fmt.Printf("P3\n%d %d\n255\n", imageWidth, imageHeight)

	for j := 0; j < imageHeight; j++ {
		log.Println("Scanlines remaining: ", imageHeight-j)
		for i := 0; i < imageWidth; i++ {
			horizontalOffset := pixelDeltaU.MulNew(float64(i))
			verticalOffset := pixelDeltaV.MulNew(float64(j))
			pixelCenter := startPixel.AddVecNew(&horizontalOffset)
			pixelCenter.AddVec(&verticalOffset)

			// Likely don't need to create a new vector
			rayDirection := pixelCenter.SubVecNew(&cameraCenter)
			currentRay := ray.Ray{
				Origin:    &cameraCenter,
				Direction: &rayDirection,
			}

			pixelColor := currentRay.Color()
			pixelColor.Write(stdout)
		}
	}

	log.Println("Done.")
}
