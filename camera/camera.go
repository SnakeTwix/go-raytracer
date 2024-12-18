package camera

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"io"
	"log"
	"os"
	"raytracer/ray"
	"raytracer/util"
	"strings"
	"sync"
	"sync/atomic"
)

type Camera struct {
	AspectRatio float64
	ImageWidth  int
	ImageHeight int

	center              *mat.VecDense
	startPixel          *mat.VecDense
	pixelDeltaU         *mat.VecDense
	pixelDeltaV         *mat.VecDense
	pixelSamplesScale   float64
	pixelSamplePerPixel int
	output              io.Writer
	maxDepth            int
}

type renderLine struct {
	scanlineNumber int
	color          []*mat.VecDense
}

func NewDefaultCamera(fileOutput *os.File) Camera {
	// Some default options, honestly don't feel like documenting all of this yet
	aspectRatio := 16. / 9.
	imageHeight := 1080
	pixelSamplePerPixel := 10

	imageWidth := int(aspectRatio * float64(imageHeight))

	focalLength := 1.
	focalVec := mat.NewVecDense(3, []float64{0, 0, focalLength})

	viewportHeight := 2.
	viewportWidth := viewportHeight * (float64(imageWidth) / float64(imageHeight))
	cameraCenter := util.NewZeroVector()

	viewportU := mat.NewVecDense(3, []float64{viewportWidth, 0, 0})
	viewportV := mat.NewVecDense(3, []float64{0, -viewportHeight, 0})

	pixelDeltaU := util.NewZeroVector()
	pixelDeltaU.ScaleVec(1/float64(imageWidth), viewportU)

	pixelDeltaV := util.NewZeroVector()
	pixelDeltaV.ScaleVec(1/float64(imageHeight), viewportV)

	halfViewportU := util.NewZeroVector()
	halfViewportU.ScaleVec(1/2., viewportU)

	halfViewportV := util.NewZeroVector()
	halfViewportV.ScaleVec(1/2., viewportV)

	viewportUpperLeft := util.NewZeroVector()
	viewportUpperLeft.SubVec(cameraCenter, focalVec)
	viewportUpperLeft.SubVec(viewportUpperLeft, halfViewportU)
	viewportUpperLeft.SubVec(viewportUpperLeft, halfViewportV)

	offset := util.NewZeroVector()
	offset.AddVec(pixelDeltaU, pixelDeltaV)
	offset.ScaleVec(0.5, offset)

	startPixel := util.NewZeroVector()
	startPixel.AddVec(viewportUpperLeft, offset)

	return Camera{
		AspectRatio:         aspectRatio,
		ImageHeight:         imageHeight,
		ImageWidth:          imageWidth,
		output:              fileOutput,
		center:              cameraCenter,
		startPixel:          startPixel,
		pixelDeltaU:         pixelDeltaU,
		pixelDeltaV:         pixelDeltaV,
		pixelSamplePerPixel: pixelSamplePerPixel,
		pixelSamplesScale:   1.0 / float64(pixelSamplePerPixel),
		maxDepth:            10,
	}
}

// WriteColor writes a single pixel to the file
func WriteColor(color *mat.VecDense, buffer io.Writer) {
	r := color.AtVec(0)
	g := color.AtVec(1)
	b := color.AtVec(2)

	intensity := util.NewInterval(0, 0.999)

	ir := int(255.999 * intensity.Clamp(r))
	ig := int(255.999 * intensity.Clamp(g))
	ib := int(255.999 * intensity.Clamp(b))

	str := fmt.Sprintf("%d %d %d\n", ir, ig, ib)

	_, err := buffer.Write([]byte(str))
	if err != nil {
		log.Fatal("Failed to write to buffer: ", err)
	}
}

// WriteLineColor writes an entire line to the file, as opposed to one pixel
func WriteLineColor(line []*mat.VecDense, buffer io.Writer) {
	var builder strings.Builder

	for _, color := range line {
		r := util.LinearToGamma(color.AtVec(0))
		g := util.LinearToGamma(color.AtVec(1))
		b := util.LinearToGamma(color.AtVec(2))

		intensity := util.NewInterval(0, 0.999)

		ir := int(255.999 * intensity.Clamp(r))
		ig := int(255.999 * intensity.Clamp(g))
		ib := int(255.999 * intensity.Clamp(b))

		builder.WriteString(fmt.Sprintf("%d %d %d\n", ir, ig, ib))
	}

	_, err := buffer.Write([]byte(builder.String()))
	if err != nil {
		log.Fatal("Failed to write to buffer: ", err)
	}
}

func (c *Camera) Render(world ray.Hittable) {
	// Okay, this one's a big one

	// Channel for sending a processed line
	renderLineCh := make(chan renderLine, 10)
	// Receives a value when all of the lines have been processed (i.e. written to the file)
	renderFinishCh := make(chan struct{})
	go c.processRenderedLines(renderLineCh, renderFinishCh)
	lineCounter := atomic.Uint64{}

	log.Println("Rendering started")

	// Just write out a ppm header, should probably refactor and put into somewhere else
	ppmHeader := fmt.Sprintf("P3\n%d %d\n255\n", c.ImageWidth, c.ImageHeight)
	_, err := c.output.Write([]byte(ppmHeader))
	if err != nil {
		log.Fatal("Failed to write the header of ppm", err)
	}

	var wg sync.WaitGroup

	// The main processor of every line
	lineWorker := func(lineNumberCh <-chan int) {

		// Take each pixel of an assigned line and compute the color of it
		// With regard to pixel sampling (antialiasing)
		for j := range lineNumberCh {
			scanline := make([]*mat.VecDense, 0, c.ImageWidth)

			for i := 0; i < c.ImageWidth; i++ {
				pixelColor := util.NewZeroVector()

				for _ = range c.pixelSamplePerPixel {
					currentRay := c.getRay(i, j)
					pixelColor.AddVec(pixelColor, currentRay.Color(world, c.maxDepth))
				}

				pixelColor.ScaleVec(c.pixelSamplesScale, pixelColor)
				scanline = append(scanline, pixelColor)
			}

			renderLineCh <- renderLine{scanlineNumber: j, color: scanline}

			log.Printf("Scanline %d done with id %d\n", lineCounter.Load(), j)
			lineCounter.Add(1)
		}

		wg.Done()
	}

	lineNumberCh := make(chan int, 10)
	for _ = range 16 {
		wg.Add(1)
		go lineWorker(lineNumberCh)
	}

	for j := 0; j < c.ImageHeight; j++ {
		lineNumberCh <- j
	}
	close(lineNumberCh)
	wg.Wait()
	// Need to ensure we've processed all the lines before closing the channel
	close(renderLineCh)
	<-renderFinishCh
}

func (c *Camera) processRenderedLines(renderChannel <-chan renderLine, finish chan<- struct{}) {
	lines := make(map[int][]*mat.VecDense)
	lastLineOutput := -1

	for renderLine := range renderChannel {
		lines[renderLine.scanlineNumber] = renderLine.color

		for {
			nextLine, exists := lines[lastLineOutput+1]
			if !exists {
				break
			}

			WriteLineColor(nextLine, c.output)
			lastLineOutput++
			delete(lines, lastLineOutput)
		}
	}

	finish <- struct{}{}
}

func (c *Camera) getRay(x, y int) ray.Ray {
	offset := c.sampleSquare()

	tempVec := util.NewZeroVector()
	rayDirection := util.NewZeroVector()
	rayDirection.AddVec(rayDirection, c.startPixel)

	tempVec.CopyVec(c.pixelDeltaU)
	tempVec.ScaleVec(float64(x)+offset.AtVec(0), tempVec)
	rayDirection.AddVec(rayDirection, tempVec)

	tempVec.Zero()
	tempVec.CopyVec(c.pixelDeltaV)
	tempVec.ScaleVec(float64(y)+offset.AtVec(1), tempVec)
	rayDirection.AddVec(rayDirection, tempVec)

	rayOrigin := util.NewZeroVector()
	rayOrigin.CopyVec(c.center)

	rayDirection.SubVec(rayDirection, rayOrigin)

	return ray.Ray{
		Origin:    rayOrigin,
		Direction: rayDirection,
	}
}

func (c *Camera) sampleSquare() *mat.VecDense {
	return mat.NewVecDense(3, []float64{util.RandomF64() - 0.5, util.RandomF64() - 0.5, 0})
}
