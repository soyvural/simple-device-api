package service

import (
	"fmt"
	"net/http"
	"regexp"
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
	re *regexp.Regexp
}

func newDeviceService(DB db) *device {
	return &device{db: DB, re: regexp.MustCompile(`^\w+$`)}
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
	if err := d.validateInput(&dev); err != nil {
		c.JSON(http.StatusBadRequest, err)
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
	id := c.Param("id")
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
	id := c.Param("id")
	if device := d.db.Get(id); device != nil {
		c.IndentedJSON(http.StatusOK, device)
		return
	}
	c.Writer.WriteHeader(http.StatusNotFound)
}

func (d *device) validateInput(dev *types.Device) error {
	var errs []error
	if !d.re.MatchString(dev.Name) {
		errs = append(errs, fmt.Errorf("invalid device name %q.", dev.Name))
	}
	if !d.re.MatchString(dev.Brand) {
		errs = append(errs, fmt.Errorf("invalid device brand %q.", dev.Brand))
	}
	if len(errs) == 0 {
		return nil
	}

	msgs := make([]string, 0, len(errs))
	for _, e := range errs {
		msgs = append(msgs, fmt.Sprintf("error: %v", e))
	}
	return fmt.Errorf(strings.Join(msgs, "\n"))
}
