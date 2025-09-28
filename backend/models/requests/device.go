package requests

type DeviceTokenRequest struct {
	DeviceID string `json:"device_id" validate:"required"`
}
