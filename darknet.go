package godarknet

/*
#cgo CFLAGS: -I./
#cgo LDFLAGS: -L./ -ldarknet
#include "core.h"
*/
import "C"
import (
	"errors"
	"sync"
	"unsafe"
)

var (
	mu sync.Mutex
)

func Init(conf, weights string, gpu int) int {
	mu.Lock()
	C.dispose()
	cconf := C.CString(conf)
	cweights := C.CString(weights)
	defer func() {
		C.free(unsafe.Pointer(cconf))
		C.free(unsafe.Pointer(cweights))
		mu.Unlock()
	}()
	return int(C.init(cconf, cweights, C.int(gpu)))
}

func DetectImage(path string) (BboxList, error) {
	bl := BboxList{}
	if !fileExists(path) {
		return bl, errors.New("no such file or dictionary")
	}
	mu.Lock()
	cpath := C.CString(path)
	defer func() {
		C.free(unsafe.Pointer(cpath))
		mu.Unlock()
	}()
	var cbboxes C.struct_bbox_t_container
	if int(C.detect_image(cpath, &cbboxes)) > 0 {
		bl = (*CbboxList)(&cbboxes).ToGo()
	}
	return bl, nil
}

type CbboxList C.struct_bbox_t_container

func (this *CbboxList) ToGo() BboxList {
	bl := BboxList{}
	for _, cbbox := range this.candidates {
		if uint(cbbox.w) == 0 && uint(cbbox.h) == 0 {
			continue
		}
		b := Bbox{
			X:             uint(cbbox.x),
			Y:             uint(cbbox.y),
			Width:         uint(cbbox.w),
			Height:        uint(cbbox.h),
			Probability:   float32(cbbox.prob),
			ObjectId:      uint(cbbox.obj_id),
			TrackId:       uint(cbbox.track_id),
			FramesCounter: uint(cbbox.frames_counter),
		}
		bl = append(bl, b)
	}
	return bl
}
