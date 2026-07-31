package goray

/*
#include "raylib.h"
#include "rlgl.h"
#include "raymath.h"

// Update buffers
void UpdateModelUVs(Model* mdl) {
	UpdateMeshBuffer(mdl->meshes[0], 1, &mdl->meshes->texcoords[0], mdl->meshes->vertexCount*2*sizeof(float), 0);
}
*/
import "C"
import (
	"unsafe"

	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-gl/gl/v2.1/gl"
	// "github.com/go-gl/gl/v3.1/gles2"
	// "github.com/go-gl/gl/v4.6-core/gl"
)

var (
	NEVER   = uint32(gl.NEVER)
	ALWAYS  = uint32(gl.ALWAYS)
	EQUAL   = uint32(gl.EQUAL)
	NEQUAL  = uint32(gl.NOTEQUAL)
	LEQUAL  = uint32(gl.LEQUAL)
	GEQUAL  = uint32(gl.GEQUAL)
	LESS    = uint32(gl.LESS)
	GREATER = uint32(gl.GREATER)

	ZERO    = uint32(gl.ZERO)
	KEEP    = uint32(gl.KEEP)
	INCR    = uint32(gl.INCR)
	DECR    = uint32(gl.DECR)
	INVERT  = uint32(gl.INVERT)
	REPLACE = uint32(gl.REPLACE)
	INCR_W  = uint32(gl.INCR_WRAP)
	DECR_W  = uint32(gl.DECR_WRAP)
)

func UpdateModelUVs(mdl *raylib.Model) {
	C.UpdateModelUVs((*C.Model)(unsafe.Pointer(mdl)))
}

func DrawRenderBatchActive() {
	C.rlDrawRenderBatchActive()
}

func EnableStencilTest(enable bool) {
	if enable {
		gl.Enable(gl.STENCIL_TEST)
	} else {
		gl.Disable(gl.STENCIL_TEST)
	}
}

func EnableDepthTest(enable bool) {
	if enable {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
}

func EnableDepthMask(enable bool) {
	gl.DepthMask(enable)
}

func SetStencilFunc(test uint32, ref int32, mask uint32) {
	gl.StencilFunc(test, ref, mask)
}

func SetStencilOp(sfail, dpfail, dppass uint32) {
	gl.StencilOp(sfail, dpfail, dppass)
}

func SetColorMask(r, g, b, a bool) {
	gl.ColorMask(r, g, b, a)
}

func ClearStencilBuffer() {
	gl.Clear(gl.STENCIL_BUFFER_BIT)
}

func ClearDepthBuffer() {
	gl.Clear(gl.DEPTH_BUFFER_BIT)
}

func SetStencilMask(mask uint32) {
	gl.StencilMask(mask)
}

func SetStencilClearValue(value int32) {
	gl.ClearStencil(value)
}

func SetStencilFuncSeparate(face uint32, fun uint32, ref int32, mask uint32) {
	gl.StencilFuncSeparate(face, fun, ref, mask)
}
