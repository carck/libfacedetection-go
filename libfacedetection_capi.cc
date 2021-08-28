// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a Apache-style
// license that can be found in the LICENSE file.

extern "C" {
	#include "libfacedetection_capi.h"
}

#include <facedetection/facedetectcnn.h>

#include <string>

// define the buffer size. Do not change the size!
#define DETECT_BUFFER_SIZE 0x20000

struct libfacedetection_capi_result_t {
	std::string *sBuffer;
	int * result;

	libfacedetection_capi_result_t(std::string* s, int* p): sBuffer(s), result(p) {
		//
	}
	~libfacedetection_capi_result_t() {
		delete this->sBuffer;
	}
};

void libfacedetection_capi_result_free(
	libfacedetection_capi_result_t* p
) {
	delete p;
}

libfacedetection_capi_result_t* libfacedetection_capi_facedetect_rgb(
	uint8_t * rgb, int width, int height, int step
) {
	std::string* sBuffer = new std::string();
	sBuffer->resize(DETECT_BUFFER_SIZE);

	unsigned char* pBuffer = (unsigned char *)sBuffer->data();
	int* pResults = facedetect_cnn(pBuffer, rgb, width, height, step);
	return new libfacedetection_capi_result_t(sBuffer, pResults);
}

libfacedetection_capi_result_t* libfacedetection_capi_facedetect_rgba(
	uint8_t* rgba, int width, int height, int step
) {
	std::string rgbBuffer;
	rgbBuffer.resize(width*height*3);

	uint8_t* rgb = (uint8_t*)rgbBuffer.data();
	for(int y = 0; y < height; y++) {
		for(int x = 0; x < width; x++) {
			*(rgb + 2) = *rgba++; // r
			*(rgb + 1) = *rgba++; // g
			*rgb       = *rgba++; // b
			rgba++;               // a
			rgb += 3;
		}
	}

	std::string* sBuffer = new std::string();
	sBuffer->resize(DETECT_BUFFER_SIZE);

	unsigned char* pBuffer = (unsigned char *)sBuffer->data();
	int* pResults = facedetect_cnn(pBuffer, rgb, width, height, step);
	return new libfacedetection_capi_result_t(sBuffer, pResults);
}

int libfacedetection_capi_result_len(libfacedetection_capi_result_t* self) {
	int* pResults = self->result;
	return pResults? *pResults: 0;
}

libfacedetection_capi_bool_t libfacedetection_capi_result_get(
	libfacedetection_capi_result_t* self, int i,
	libfacedetection_capi_face_t* face
) {
	int* pResults = self->result;
	int n = pResults? *pResults: 0;

	if(i < 0 || i >= n) return 0;

	short * p = ((short*)(pResults+1))+142*i;
	face->confidence = p[0];
	face->x = p[1];
	face->y = p[2];
	face->w = p[3];
	face->h = p[4];
	for(int j = 0; j < 10; j++) {
		face->landmarks[j] = p[5+j];
	}

	return 1;
}
