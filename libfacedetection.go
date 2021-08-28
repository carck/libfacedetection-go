// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a Apache-style
// license that can be found in the LICENSE file.

package libfacedetection

/*
#cgo LDFLAGS: -lfacedetection
#include "libfacedetection_capi.h"
*/
import "C"
import (
	"image"
	"unsafe"
)

type Face struct {
	Confidence int
	X         int
	Y         int
	W         int
	H         int
	Landmarks []int
}

func DetectFaceRGBA(m *image.RGBA) []Face {
	if m == nil {
		return nil
	}

	rv := C.libfacedetection_capi_facedetect_rgba(
		(*C.uint8_t)(unsafe.Pointer(&m.Pix[0])),
		C.int(m.Rect.Dx()),
		C.int(m.Rect.Dy()),
		C.int(m.Stride),
	)
	defer C.libfacedetection_capi_result_free(rv)

	n := int(C.libfacedetection_capi_result_len(rv))
	if n <= 0 {
		return nil
	}

	face := make([]Face, n)
	for i := 0; i < n; i++ {
		var t C.libfacedetection_capi_face_t
		C.libfacedetection_capi_result_get(rv, C.int(i), &t)
		
		face[i].Confidence = int(t.confidence)
		face[i].X = int(t.x)
		face[i].Y = int(t.y)
		face[i].W = int(t.w)
		face[i].H = int(t.h)
		face[i].Landmarks = make([]int, 10)
		for j := 0; j < 10; j++ {
			face[i].Landmarks[j] = int(t.landmarks[j])
		}
	}

	return face
}

func DetectFaceRGB(rgb []byte, w, h, stride int) []Face {
	if len(rgb) == 0 {
		return nil
	}

	rv := C.libfacedetection_capi_facedetect_rgb(
		(*C.uint8_t)(unsafe.Pointer(&rgb[0])),
		C.int(w),
		C.int(h),
		C.int(stride),
	)
	defer C.libfacedetection_capi_result_free(rv)

	n := int(C.libfacedetection_capi_result_len(rv))
	if n <= 0 {
		return nil
	}

	face := make([]Face, n)
	for i := 0; i < n; i++ {
		var t C.libfacedetection_capi_face_t
		C.libfacedetection_capi_result_get(rv, C.int(i), &t)

		face[i].Confidence = int(t.confidence)
		face[i].X = int(t.x)
		face[i].Y = int(t.y)
		face[i].W = int(t.w)
		face[i].H = int(t.h)
		face[i].Landmarks = make([]int, 10)
                for j := 0; j < 10; j++ {
                        face[i].Landmarks[j] = int(t.landmarks[j])
                }
	}

	return face
}

func NewRGBImageFrom(m image.Image) (rgb []byte, w, h int) {
	b := m.Bounds()

	w = b.Dx()
	h = b.Dy()
	rgb = make([]byte, w*h*3)

	off := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			pr, pg, pb, _ := m.At(x, y).RGBA()
			rgb[off+0] = uint8(pb >> 8)
			rgb[off+1] = uint8(pg >> 8)
			rgb[off+2] = uint8(pr>> 8)
			off += 3
		}
	}

	return
}
