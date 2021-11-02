package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/soyvural/simple-device-api/types"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
)

type fakeDB struct {
	mockVal *types.Device
}

func (f *fakeDB) Insert(d types.Device) *types.Device {
	return f.mockVal
}
func (f *fakeDB) Get(id string) *types.Device {
	return f.mockVal
}
func (f *fakeDB) Delete(id string) *types.Device {
	return f.mockVal
}

func TestCreateDevice(t *testing.T) {
	tests := []struct {
		desc             string
		device           *types.Device
		deviceJSON       string
		returningVal     *types.Device
		wantedStatusCode int
	}{
		{
			desc: "success",
			device: &types.Device{
				Name:  "Phone",
				Brand: "Apple",
				Model: "13 Pro Max",
			},
			returningVal: &types.Device{
				ID:    "1234",
				Name:  "Phone",
				Brand: "Apple",
				Model: "13 Pro Max",
			},
			wantedStatusCode: http.StatusCreated,
		},
		{
			desc: "invalid JSON",
			device: &types.Device{
				Name:  "!name$",
				Brand: "Apple",
			},
			wantedStatusCode: http.StatusBadRequest,
		},
		{
			desc:             "invalid device",
			deviceJSON:       `{"name": "12334"}`,
			wantedStatusCode: http.StatusBadRequest,
		},
		{
			desc: "service unavailable",
			device: &types.Device{
				Name:  "Phone",
				Brand: "Apple",
				Model: "13 Pro Max",
			},
			returningVal:     nil,
			wantedStatusCode: http.StatusServiceUnavailable,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			router := gin.Default()
			svc := &Service{router: router, deviceSvc: newDeviceService(&fakeDB{mockVal: tc.returningVal})}
			svc.SetRoute_v1()

			var r io.Reader = strings.NewReader(tc.deviceJSON)
			if tc.device != nil {
				out, err := json.Marshal(tc.device)
				if err != nil {
					t.Fatalf("err: %v", err)
				}
				r = bytes.NewReader(out)
			}

			req := httptest.NewRequest("POST", "/api/v1/device", r)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if diff := cmp.Diff(tc.wantedStatusCode, rr.Code); diff != "" {
				t.Fatalf("Status code mismatch (-want +got): %s\n", diff)
			}
		})
	}
}

func TestGetDevice(t *testing.T) {
	tests := []struct {
		desc             string
		id               string
		returningVal     *types.Device
		wantedStatusCode int
	}{
		{
			desc: "success",
			id:   "1234",
			returningVal: &types.Device{
				ID:    "1234",
				Name:  "Phone",
				Brand: "Apple",
				Model: "13 Pro Max",
			},
			wantedStatusCode: http.StatusOK,
		},
		{
			desc:             "not found",
			id:               "1234",
			returningVal:     nil,
			wantedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			router := gin.Default()
			svc := &Service{router: router, deviceSvc: newDeviceService(&fakeDB{mockVal: tc.returningVal})}
			svc.SetRoute_v1()

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/device/%s", tc.id), nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if diff := cmp.Diff(tc.wantedStatusCode, rr.Code); diff != "" {
				t.Fatalf("Status code mismatch (-want +got): %s\n", diff)
			}

			if tc.wantedStatusCode != 200 {
				return
			}

			d := types.Device{}
			if err := json.Unmarshal(rr.Body.Bytes(), &d); err != nil {
				t.Fatalf("Failed to read response body, err: %v", err)
			}

			if diff := cmp.Diff(*tc.returningVal, d); diff != "" {
				t.Fatalf("Returning device mismatch (-want +got): %s\n", diff)
			}
		})
	}
}

func TestDeleteDevice(t *testing.T) {
	tests := []struct {
		desc             string
		id               string
		returningVal     *types.Device
		wantedStatusCode int
	}{
		{
			desc: "success",
			id:   "1234",
			returningVal: &types.Device{
				ID:    "1234",
				Name:  "Phone",
				Brand: "Apple",
				Model: "13 Pro Max",
			},
			wantedStatusCode: http.StatusOK,
		},
		{
			desc:             "not found",
			id:               "1234",
			returningVal:     nil,
			wantedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			router := gin.Default()
			svc := &Service{router: router, deviceSvc: newDeviceService(&fakeDB{mockVal: tc.returningVal})}
			svc.SetRoute_v1()

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/device/%s", tc.id), nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if diff := cmp.Diff(tc.wantedStatusCode, rr.Code); diff != "" {
				t.Fatalf("Status code mismatch (-want +got): %s\n", diff)
			}

			if tc.wantedStatusCode != 200 {
				return
			}

			d := types.Device{}
			if err := json.Unmarshal(rr.Body.Bytes(), &d); err != nil {
				t.Fatalf("Failed to read response body, err: %v", err)
			}

			if diff := cmp.Diff(*tc.returningVal, d); diff != "" {
				t.Fatalf("Returning device mismatch (-want +got): %s\n", diff)
			}
		})
	}
}
