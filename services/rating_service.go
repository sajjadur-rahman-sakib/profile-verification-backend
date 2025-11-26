package services

import (
	"errors"

	"verify/config"
	"verify/models"

	"gorm.io/gorm"
)

type RatingService struct{}

func NewRatingService() *RatingService {
	return &RatingService{}
}

func (s *RatingService) GiveRating(raterEmail, ratedEmail string, rating int) error {
	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	if raterEmail == ratedEmail {
		return errors.New("you cannot rate yourself")
	}

	var ratedUser models.User
	if err := config.DB.Where("email = ? AND is_verified = ?", ratedEmail, true).First(&ratedUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("rated user not found or not verified")
		}
		return err
	}

	var rater models.User
	if err := config.DB.Where("email = ?", raterEmail).First(&rater).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("rater not found")
		}
		return err
	}

	var existingRating models.Rating
	result := config.DB.Where("rater_email = ? AND rated_email = ?", raterEmail, ratedEmail).First(&existingRating)

	if result.Error == nil {
		existingRating.Rating = rating
		if err := config.DB.Save(&existingRating).Error; err != nil {
			return err
		}
	} else {
		newRating := models.Rating{
			RaterEmail: raterEmail,
			RatedEmail: ratedEmail,
			Rating:     rating,
		}
		if err := config.DB.Create(&newRating).Error; err != nil {
			return err
		}
	}

	return s.updateUserAverageRating(ratedEmail)
}

func (s *RatingService) updateUserAverageRating(email string) error {
	var ratings []models.Rating
	if err := config.DB.Where("rated_email = ?", email).Find(&ratings).Error; err != nil {
		return err
	}

	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	if len(ratings) == 0 {
		user.AverageRating = 0
	} else {
		totalRating := 0
		for _, rating := range ratings {
			totalRating += rating.Rating
		}
		user.AverageRating = float64(totalRating) / float64(len(ratings))
	}

	return config.DB.Save(&user).Error
}

func (s *RatingService) GetAverageRating(email string) (float64, int, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, 0, nil
		}
		return 0, 0, err
	}

	var totalRatings int64
	config.DB.Model(&models.Rating{}).Where("rated_email = ?", email).Count(&totalRatings)

	return user.AverageRating, int(totalRatings), nil
}

func (s *RatingService) GetUserRatings(email string) ([]models.Rating, error) {
	var ratings []models.Rating
	if err := config.DB.Where("rated_email = ?", email).Find(&ratings).Error; err != nil {
		return nil, err
	}
	return ratings, nil
}
