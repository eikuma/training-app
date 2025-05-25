package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/form"
	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/service"
)

type (
	// Recommendation トレーニングメニュー提案のハンドラを表す
	Recommendation interface {
		ProposeTrainingMenu(c echo.Context) error
	}

	// RecommendationImpl トレーニングメニュー提案のハンドラ実装
	RecommendationImpl struct {
		RecommendationService service.Recommendation
	}
)

// コンストラクタ
func NewRecommendation() Recommendation {
	return &RecommendationImpl{
		RecommendationService: service.NewRecommendation(),
	}
}

// ユーザが選択した条件に応じてトレーニングメニューを提案
func (h *RecommendationImpl) ProposeTrainingMenu(c echo.Context) error {
	f := form.NewProposeTrainingMenu()
	fmt.Println(f)
	if err := c.Bind(f); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid form: "+err.Error())
	}

	// OpenAI API等を利用して提案を行う
	result, err := h.RecommendationService.ProposeTrainingMenu(
		f.TrainingGoal,
		f.TargetParts,
		f.ExperienceLevel,
		f.AvailableTime,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"recommendation": result,
	})
}
