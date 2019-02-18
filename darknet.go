package godarknet

/*
#include "core.h"
#cgo LDFLAGS: -ldl
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

func Open(libpath string) (dnet *Darknet, err error) {
	dnet = new(Darknet)
	clibpath := C.CString(libpath)
	mu.Lock()
	defer func() {
		C.free(unsafe.Pointer(clibpath))
		mu.Unlock()
	}()
	if handle := C.dlopen(clibpath, C.RTLD_NOW); handle != nil {
		dnet.handle = handle
	} else {
		err = dnet.catchErr()
	}
	return
}

func (this *Darknet) Init(config, weights string, gpu int) (err error) {
	cfg := C.CString(config)
	wgt := C.CString(weights)
	cfunc := C.CString("init")
	mu.Lock()
	defer func() {
		C.free(unsafe.Pointer(cfg))
		C.free(unsafe.Pointer(wgt))
		C.free(unsafe.Pointer(cfunc))
		mu.Unlock()
	}()
	if fn := C.dlsym(this.handle, cfunc); fn != nil {
		if int(C.call_init(fn, cfg, wgt, C.int(gpu))) != 1 {
			err = errors.New("init error")
		}
	} else {
		err = this.catchErr()
	}
	return
}

func (this *Darknet) Detect(imagePath string) (BboxList, error) {
	cimg := C.CString(imagePath)
	cfunc := C.CString("detect_image")
	mu.Lock()
	defer func() {
		C.free(unsafe.Pointer(cimg))
		C.free(unsafe.Pointer(cfunc))
		mu.Unlock()
	}()
	var err error
	bboxes := BboxList{{}}
	if fn := C.dlsym(this.handle, cfunc); fn != nil {
		var cbboxes C.struct_bbox_t_container
		if ret := int(C.call_detect_image(fn, cimg, &cbboxes)); ret == 1 {
			bboxes = (*CbboxList)(&cbboxes).ToGo()
		}
	} else {
		err = this.catchErr()
	}
	return bboxes, err
}

func (this *Darknet) Close() error {
	var err error
	if this.handle != nil {
		mu.Lock()
		defer mu.Unlock()
		err = this.dispose()
		if err == nil {
			if C.dlclose(this.handle) == 0 {
				this.handle = nil
			} else {
				err = this.catchErr()
			}
		}
	}
	return err
}

func (this *Darknet) dispose() error {
	cfunc := C.CString("dispose")
	this.mu.Lock()
	defer func() {
		C.free(unsafe.Pointer(cfunc))
		this.mu.Unlock()
	}()
	var err error
	if fn := C.dlsym(this.handle, cfunc); fn != nil {
		if int(C.call_dispose(fn, this.handle)) != 1 {
			err = errors.New("dispose error")
		}
	} else {
		err = this.catchErr()
	}
	return err
}

func (this *Darknet) catchErr() error {
	return errors.New(C.GoString(C.dlerror()))
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
