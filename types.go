package godarknet

type Bbox struct {
	X             uint     `json:"x"`
	Y             uint     `json:"y"`
	Width         uint     `json:"width"`
	Height        uint     `json:"height"`
	Probability   float32 `json:"prob"`
	ObjectId      uint    `json:"obj_id"`
	TrackId       uint    `json:"track_id"`
	FramesCounter uint    `json:"frames_counter"`
}

type BboxList []Bbox
