package webp

//#include "webp_util.h"
import "C"
import (
	"io"
	"sync"
	"unsafe"
)

//WebPPicture pointer -> Writer for the WebPPicture
var wh *writeHandler = &writeHandler{
	wm: make(map[unsafe.Pointer]io.Writer),
}

type writeHandler struct {
	m  sync.Mutex
	wm map[unsafe.Pointer]io.Writer
}

func (h *writeHandler) add(p unsafe.Pointer, w io.Writer) {
	h.m.Lock()
	defer h.m.Unlock()
	h.wm[p] = w
}

func (h *writeHandler) get(p unsafe.Pointer) io.Writer {
	h.m.Lock()
	defer h.m.Unlock()
	return h.wm[p]
}

func (h *writeHandler) del(p unsafe.Pointer) {
	h.m.Lock()
	defer h.m.Unlock()
	delete(h.wm, p)
}

//writeTo is a Go implementation of the C function we defined in our webp_util header file
//the purpose of this is to be able to use this callback function and write to a native Go writer.
//We map a writer to the WebPPicture's pointer so we can grab it when this function is called

//export writeTo
func writeTo(data *C.uint8_t, data_size C.size_t, picture *C.struct_WebPPicture) C.int {
	w := wh.get(unsafe.Pointer(picture))

	if w != nil {
		read, err := w.Write(C.GoBytes(unsafe.Pointer(data), C.int(data_size)))

		if err != nil {
			return 0
		}

		return C.int(read)
	}

	return 0
}
