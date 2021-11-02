package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/soyvural/simple-device-api/types"

	"github.com/gin-gonic/gin"
)

type db interface {
	Get(id string) *types.Device
	Insert(d types.Device) *types.Device
	Delete(id string) *types.Device
}

type device struct {
	db db
}

func newDeviceService(DB db) *device {
	return &device{db: DB}
}

// CreateDevice .
// @BasePath /api/v1
// @Description Creates a device and returns device object recently created in store.
// @Summary Show a list of students
// @Accept  json
// @Produce  json
// @Param device body types.Device true "Device definition"
// @Success 201 {object} types.Device
// @Router /api/v1/device [post]
func (d *device) CreateDevice(c *gin.Context) {
	dev := types.Device{}
	if err := c.BindJSON(&dev); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	dev.ID = uuid.NewString()
	newDev := d.db.Insert(dev)
	if newDev == nil {
		c.Writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	c.IndentedJSON(http.StatusCreated, newDev)
}

// DeleteDevice .
// @BasePath /api/v1
// @Description Deletes the device by given id.
// @Summary Show a list of students
// @Accept  json
// @Produce  json
// @Param id path string true "Device ID"
// @Success 200 {object} types.Device
// @Router /api/v1/device [delete]
func (d *device) DeleteDevice(c *gin.Context) {
	id, ok := d.mustHaveParam("id", c)
	if !ok {
		return
	}
	if device := d.db.Delete(id); device != nil {
		c.IndentedJSON(http.StatusOK, device)
		return
	}
	c.Writer.WriteHeader(http.StatusNotFound)
}

// GetDevice .
// @BasePath /api/v1
// @Description retrieves the Device and returns it by given ID.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Device ID"
// @Success 200
// @Router /api/v1/device/{id} [get]
func (d *device) GetDevice(c *gin.Context) {
	id, ok := d.mustHaveParam("id", c)
	if !ok {
		return
	}
	if device := d.db.Get(id); device != nil {
		c.IndentedJSON(http.StatusOK, device)
		return
	}
	c.Writer.WriteHeader(http.StatusNotFound)
}

func (d *device) mustHaveParam(name string, c *gin.Context) (string, bool) {
	val := c.Param(name)
	if strings.Trim(val, " ") == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("%s param not found", name))
		return "", false
	}
	return val, true
}
