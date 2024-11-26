package main

import (
	"gonum.org/v1/gonum/mat"
	"log"
	"os"
	"path/filepath"
	"raytracer/camera"
	"raytracer/objects"
	"raytracer/ray"
)

func main() {

	// World
	world := ray.NewHittableList()

	// Sphere 1
	sphere1 := objects.NewSphere(mat.NewVecDense(3, []float64{0, 0, -1}), 0.5)
	world.Add(&sphere1)

	// Sphere 2
	sphere2 := objects.NewSphere(mat.NewVecDense(3, []float64{0, -100.5, -1}), 100)
	world.Add(&sphere2)

	sphere3 := objects.NewSphere(mat.NewVecDense(3, []float64{5, 0, -5}), 0.5)
	world.Add(&sphere3)

	sphere4 := objects.NewSphere(mat.NewVecDense(3, []float64{-5, 0, -5}), 1)
	world.Add(&sphere4)

	file := getFile()

	// Camera
	c := camera.NewDefaultCamera(file)
	c.Render(&world)

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
