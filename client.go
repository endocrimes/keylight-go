package keylight

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (k *KeyLight) httpGet(ctx context.Context, path string, target interface{}) error {
	url := fmt.Sprintf("http://%s:%d/%s", k.DNSAddr, k.Port, path)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func (k *KeyLight) httpPut(ctx context.Context, path string, body interface{}, target interface{}) error {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:%d/%s", k.DNSAddr, k.Port, path)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, buf)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

// FetchSettings allows you to retrieve general device settings.
func (k *KeyLight) FetchSettings(ctx context.Context) (*DeviceSettings, error) {
	s := &DeviceSettings{}
	err := k.httpGet(ctx, "elgato/lights/settings", s)
	return s, err
}

// FetchDeviceInfo returns metadata for the accessory.
func (k *KeyLight) FetchDeviceInfo(ctx context.Context) (*DeviceInfo, error) {
	i := &DeviceInfo{}
	err := k.httpGet(ctx, "elgato/accessory-info", i)
	return i, err
}

// FetchLightGroup returns all of the individual lights that are owned by an
// accessory. This in conjunction with UpdateLightGroup will allow you to
// control your lights.
func (k *KeyLight) FetchLightGroup(ctx context.Context) (*LightGroup, error) {
	o := &LightGroup{Lights: make([]*Light, 0)}
	err := k.httpGet(ctx, "elgato/lights", o)
	return o, err
}

// UpdateLightGroup allows you to update the settings for individual lights
// in an accessory. It returns the updated options.
func (k *KeyLight) UpdateLightGroup(ctx context.Context, newOptions *LightGroup) (*LightGroup, error) {
	o := &LightGroup{Lights: make([]*Light, 0)}
	err := k.httpPut(ctx, "elgato/lights", newOptions, o)
	return o, err
}
