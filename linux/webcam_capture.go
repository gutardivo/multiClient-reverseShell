// Install: go get -u gocv.io/x/gocv

package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"time"

	"gocv.io/x/gocv"
)

func takeWebcamShot(camera *gocv.VideoCapture) {
	camera.Set(gocv.VideoCaptureFrameWidth, float64(1280))
	camera.Set(gocv.VideoCaptureFrameHeight, float64(720))

	img := gocv.NewMat()
	defer img.Close()

	x := 0
	for x < 3 {
		if camera.IsOpened() {
			for {
				if camera.Read(&img) {
					x++
					if x > 4 {
						imagePath := "webcam_shot.jpg"
						if err := gocv.IMWrite(imagePath, img); err == nil {
							fmt.Printf("Webcam shot saved to %s\n", imagePath)
							break
						}
					}
				} else {
					fmt.Println("Error reading from webcam")
					break
				}
			}
			break
		} else {
			x++
		}
	}
}

func main() {
	camera := gocv.VideoCaptureDevice(0)
	defer camera.Close()

	if !camera.IsOpened() {
		fmt.Println("Error: Camera not found")
		return
	}

	go takeWebcamShot(&camera)

	// Serve the image on a web server
	http.HandleFunc("/image.jpg", func(w http.ResponseWriter, r *http.Request) {
		img := gocv.IMRead("webcam_shot.jpg", gocv.IMReadColor)
		defer img.Close()
		if img.Empty() {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}

		buffer, err := gocv.IMEncode(".jpg", img)
		if err != nil {
			http.Error(w, "Unable to encode image", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", fmt.Sprint(len(buffer)))
		if _, err := w.Write(buffer); err != nil {
			http.Error(w, "Unable to write image", http.StatusInternalServerError)
		}
	})

	go func() {
		// Start the HTTP server to serve the captured image
		http.ListenAndServe(":8080", nil)
	}()

	fmt.Println("Open a web browser and navigate to http://localhost:8080/image.jpg to view the webcam shot.")
	// Give some time to the HTTP server
	time.Sleep(5 * time.Minute)
}

