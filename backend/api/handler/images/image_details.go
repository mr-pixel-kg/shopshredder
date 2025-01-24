package images

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// ImageDetailsHandler returns information about a image
// @Summary Details about a Docker Image
// @Description Get details about a docker image
// @Tags Docker Image Management
// @Accept json
// @Produce json
// @Param id path string true "Image ID" example(a407dee395ed97ead1e40c7537395d6271c07cc89c317f8eda1c19f6fc783695)
// @Success 200 {object} Image
// @Failure 400 {object} map[string]string
// @Router /api/images/{id} [get]
func (h *ImageHandler) ImageDetailsHandler(c echo.Context) error {

	ctx := c.Request().Context()
	imageId := c.Param("id")

	output, err := h.ImageService.GetImage(ctx, imageId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, output)
}
