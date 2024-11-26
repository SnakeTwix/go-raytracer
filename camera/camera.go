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
}

type renderLine struct {
	scanlineNumber int
	color          []*mat.VecDense
}

func NewDefaultCamera(fileOutput *os.File) Camera {
	aspectRatio := 16. / 9.
	imageHeight := 1080
	pixelSamplePerPixel := 10

	imageWidth := int(aspectRatio * float64(imageHeight))

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
	}
}

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

func WriteLineColor(line []*mat.VecDense, buffer io.Writer) {
	var builder strings.Builder

	for _, color := range line {
		r := color.AtVec(0)
		g := color.AtVec(1)
		b := color.AtVec(2)

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
	sendCh := make(chan renderLine, 10)
	renderFinishCh := make(chan struct{})
	go c.processRenderedLines(sendCh, renderFinishCh)

	lineCounter := atomic.Uint64{}

	log.Println("Rendering started")

	ppmHeader := fmt.Sprintf("P3\n%d %d\n255\n", c.ImageWidth, c.ImageHeight)
	_, err := c.output.Write([]byte(ppmHeader))
	if err != nil {
		log.Fatal("Failed to write the header of ppm", err)
	}

	var wg sync.WaitGroup
	wg.Add(c.ImageHeight)

	processLine := func(lineNumberCh <-chan int) {
		for j := range lineNumberCh {
			scanline := make([]*mat.VecDense, 0, c.ImageWidth)

			for i := 0; i < c.ImageWidth; i++ {
				pixelColor := mat.NewVecDense(3, nil)

				for _ = range c.pixelSamplePerPixel {
					currentRay := c.getRay(i, j)
					pixelColor.AddVec(pixelColor, currentRay.Color(world))
				}

				pixelColor.ScaleVec(c.pixelSamplesScale, pixelColor)
				scanline = append(scanline, pixelColor)
			}

			sendCh <- renderLine{scanlineNumber: j, color: scanline}
			wg.Done()

			log.Printf("Scanline %d done with id %d\n", lineCounter.Load(), j)
			lineCounter.Add(1)
		}
	}

	lineNumberCh := make(chan int, 10)
	for _ = range 8 {
		go processLine(lineNumberCh)
	}

	for j := 0; j < c.ImageHeight; j++ {
		//log.Println("Scanlines remaining: ", c.ImageHeight-j)
		lineNumberCh <- j
	}

	close(lineNumberCh)
	wg.Wait()
	close(sendCh)

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

	tempVec := mat.NewVecDense(3, nil)
	rayDirection := mat.NewVecDense(3, nil)
	rayDirection.AddVec(rayDirection, c.startPixel)

	tempVec.CopyVec(c.pixelDeltaU)
	tempVec.ScaleVec(float64(x)+offset.AtVec(0), tempVec)
	rayDirection.AddVec(rayDirection, tempVec)

	tempVec.Zero()
	tempVec.CopyVec(c.pixelDeltaV)
	tempVec.ScaleVec(float64(y)+offset.AtVec(1), tempVec)
	rayDirection.AddVec(rayDirection, tempVec)

	rayOrigin := mat.NewVecDense(3, nil)
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
