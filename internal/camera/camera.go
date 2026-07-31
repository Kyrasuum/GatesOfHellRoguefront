package camera

import (
	"roguefront/pkg/lmath"

	pub_object "roguefront/pkg/object"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var ()

type Cam struct {
	parent pub_object.Object
	camera rl.Camera

	dist  float32
	yaw   float32
	pitch float32
	roll  float32
}

func NewCam() (s *Cam, err error) {
	s = &Cam{}
	err = s.Init()

	return s, err
}

func (s *Cam) Init() error {
	if s == nil {
		return nil
	}

	s.camera = rl.Camera3D{}
	s.camera.Position = rl.NewVector3(0.0, 0.0, 0.0)
	s.camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	s.camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	s.camera.Fovy = 45.0
	s.camera.Projection = rl.CameraPerspective

	s.dist = 4
	s.roll = 0
	s.pitch = -30
	s.yaw = 45
	s.UpdateCam()

	return nil
}

func (s *Cam) OnResize(w int32, h int32) {
	if s == nil {
		return
	}
}

func (s *Cam) Prerender() []func() {
	if s == nil {
		return []func(){}
	}

	return []func(){}
}

func (s *Cam) Render() []func() {
	if s == nil {
		return []func(){}
	}

	rlx.BeginMode3D(s.camera)
	return []func(){rlx.EndMode3D}
}

func (s *Cam) Postrender() []func() {
	if s == nil {
		return []func(){}
	}

	return []func(){}
}

func (s *Cam) Update(dt float32) {
	if s == nil {
		return
	}
}

func (s *Cam) GetPos() rl.Vector3 {
	if s == nil {
		return rl.Vector3{}
	}

	return s.camera.Position
}

func (s *Cam) SetPos(pos rl.Vector3) {
	if s == nil {
		return
	}

	s.camera.Position = pos
}

func (s *Cam) GetTar() rl.Vector3 {
	if s == nil {
		return rl.Vector3{}
	}

	return s.camera.Target
}

func (s *Cam) SetTar(tar rl.Vector3) {
	if s == nil {
		return
	}

	s.camera.Target = tar
}

func (s *Cam) GetModelMatrix() rl.Matrix {
	if s == nil {
		return rl.Matrix{}
	}

	ql := lmath.Quat{}
	ql = *ql.FromEuler(0-float64(s.pitch), float64(s.yaw), float64(s.roll))

	view := ql.RotateVec3(lmath.Vec3{0, 0, float64(s.dist)})
	view.X += float64(s.camera.Target.X)
	view.Y += float64(s.camera.Target.Y)
	view.Z += float64(s.camera.Target.Z)

	camMat := rl.QuaternionToMatrix(rl.NewQuaternion(float32(ql.X), float32(ql.Y), float32(ql.Z), float32(ql.W)))
	camMat = rl.MatrixMultiply(camMat, rl.MatrixTranslate(float32(view.X), float32(view.Y), float32(view.Z)))

	return camMat
}

func (s *Cam) GetCameraMatrix() rl.Matrix {
	if s == nil {
		return rl.Matrix{}
	}

	return rlx.GetCameraMatrix(s.camera)
}

func (s *Cam) GetWorldToScreen(pos rl.Vector3) rl.Vector2 {
	if s == nil {
		return rl.Vector2{}
	}

	return rlx.GetWorldToScreen(pos, s.camera)
}

func (s *Cam) GetProjMatrix(width int32, height int32) rl.Matrix {
	if s == nil {
		return rl.Matrix{}
	}

	return rlx.GetCameraProjectionMatrix(&s.camera, float32(width)/float32(height))
}

func (s *Cam) GetViewProjMatrix() rl.Matrix {
	if s == nil {
		return rl.Matrix{}
	}

	aspect := float32(rlx.GetScreenWidth()) / float32(rlx.GetScreenHeight())
	view := rlx.GetCameraViewMatrix(&s.camera)
	proj := rlx.GetCameraProjectionMatrix(&s.camera, aspect)
	return rl.MatrixMultiply(view, proj)
}

func (s *Cam) GetDist() float32 {
	if s == nil {
		return 0
	}

	return s.dist
}

func (s *Cam) SetDist(d float32) {
	if s == nil {
		return
	}

	s.dist = d
}

func (s *Cam) GetPitch() float32 {
	if s == nil {
		return 0
	}

	return s.pitch
}

func (s *Cam) SetPitch(p float32) {
	if s == nil {
		return
	}

	s.pitch = p
}

func (s *Cam) GetYaw() float32 {
	if s == nil {
		return 0
	}

	return s.yaw
}

func (s *Cam) SetYaw(y float32) {
	if s == nil {
		return
	}

	s.yaw = y
}

func (s *Cam) GetRoll() float32 {
	if s == nil {
		return 0
	}

	return s.roll
}

func (s *Cam) SetRoll(r float32) {
	if s == nil {
		return
	}

	s.roll = r
}

func (s *Cam) ZoomCam(zoom float32) {
	if s == nil {
		return
	}

	s.dist -= zoom
}

func (s *Cam) RotateCam(pitch float32, yaw float32) {
	if s == nil {
		return
	}

	s.pitch += pitch
	s.yaw += yaw
}

func (s *Cam) MoveCam(move lmath.Vec3) {
	if s == nil {
		return
	}

	s.camera.Target.X += float32(move.X)
	s.camera.Target.Y += float32(move.Y)
	s.camera.Target.Z += float32(move.Z)
}

func (s *Cam) UpdateCam() {
	if s == nil {
		return
	}

	ql := lmath.Quat{}
	ql = *ql.FromEuler(rl.Deg2rad*float64(s.pitch), rl.Deg2rad*float64(s.yaw), float64(s.roll))

	view := ql.RotateVec3(lmath.Vec3{0, 0, float64(s.dist)})

	s.camera.Position.X = float32(view.X) + s.camera.Target.X
	s.camera.Position.Y = float32(view.Y) + s.camera.Target.Y
	s.camera.Position.Z = float32(view.Z) + s.camera.Target.Z
}

func (s *Cam) OnAdd(obj pub_object.Object) {
	if s == nil {
		return
	}
	s.parent = obj
}

func (s *Cam) OnRemove() {
	if s == nil {
		return
	}
	s.parent = nil
}

func (c *Cam) GetParent() pub_object.Object {
	if c == nil {
		return nil
	}
	return c.parent
}
