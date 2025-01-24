package images

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// ImageListHandler lists all images
// @Summary List of Docker Image
// @Description Get a list of docker image
// @Tags Docker Image Management
// @Accept json
// @Produce json
// @Success 200 {array} Image
// @Failure 400 {object} map[string]string
// @Router /api/images [get]
func (h *ImageHandler) ImageListHandler(c echo.Context) error {

	ctx := c.Request().Context()

	output, err := h.ImageService.ListImages(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, output)
}
