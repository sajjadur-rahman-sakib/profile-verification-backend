package handlers

import (
"net/http"
"strconv"

"github.com/labstack/echo/v4"
"main.go/services"
)

type RatingHandler struct {
ratingService *services.RatingService
}

func NewRatingHandler(ratingService *services.RatingService) *RatingHandler {
return &RatingHandler{ratingService}
}

func (h *RatingHandler) GiveRating(c echo.Context) error {
raterEmail := c.FormValue("rater_email")
ratedEmail := c.FormValue("rated_email")
ratingStr := c.FormValue("rating")

if raterEmail == "" || ratedEmail == "" || ratingStr == "" {
return c.JSON(http.StatusBadRequest, map[string]string{
"error": "rater_email, rated_email, and rating are required",
})
}

rating, err := strconv.Atoi(ratingStr)
if err != nil {
return c.JSON(http.StatusBadRequest, map[string]string{
"error": "rating must be a valid number",
})
}

if err := h.ratingService.GiveRating(raterEmail, ratedEmail, rating); err != nil {
return c.JSON(http.StatusBadRequest, map[string]string{
"error": err.Error(),
})
}

return c.JSON(http.StatusOK, map[string]string{
"message": "Rating submitted successfully",
})
}

func (h *RatingHandler) GetUserRatings(c echo.Context) error {
email := c.QueryParam("email")

if email == "" {
email = c.FormValue("email")
}

if email == "" {
return c.JSON(http.StatusBadRequest, map[string]string{
"error": "email is required",
})
}

ratings, err := h.ratingService.GetUserRatings(email)
if err != nil {
return c.JSON(http.StatusInternalServerError, map[string]string{
"error": "failed to fetch ratings",
})
}

averageRating, totalRatings, err := h.ratingService.GetAverageRating(email)
if err != nil {
return c.JSON(http.StatusInternalServerError, map[string]string{
"error": "failed to calculate average rating",
})
}

return c.JSON(http.StatusOK, map[string]interface{}{
"email":          email,
"average_rating": averageRating,
"total_ratings":  totalRatings,
"ratings":        ratings,
})
}
