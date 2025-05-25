package response

import (
	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/model"
)

type (
	WorkoutSession struct {
		ID     int64  `json:"id"`
		Date   string `json:"date"`
		UserID int64  `json:"user_id"`
	}

	WorkoutSessions []WorkoutSession

	Exercise struct {
		ID           int64  `json:"exercise_id"`
		SessionID    int64  `json:"session_id"`
		ExerciseName string `json:"exercise_name"`
		Sets         Sets   `json:"sets"`
	}

	Exercises []Exercise

	Set struct {
		ID         int64   `json:"set_id"`
		ExerciseID int64   `json:"exercise_id"`
		SetNumber  int64   `json:"set_number"`
		Weight     float64 `json:"weight"`
		Reps       int64   `json:"reps"`
	}

	Sets []Set

	GetWorkoutSession struct {
		ID        int64     `json:"id"`
		Date      string    `json:"date"`
		UserID    int64     `json:"user_id"`
		Exercises Exercises `json:"exercises"`
	}
)

func NewWorkoutSession() *WorkoutSession {
	return &WorkoutSession{}
}

func NewGetWorkoutSession() *GetWorkoutSession {
	return &GetWorkoutSession{}
}

func NewExercise() *Exercise {
	return &Exercise{}
}

func (r *WorkoutSession) WorkoutSessionFromModel(m *model.WorkoutSessionImpl) *WorkoutSession {
	r.ID = m.ID
	r.Date = m.Date.Format("2006-01-02")
	r.UserID = m.UserID
	return r
}

func (r *GetWorkoutSession) GetWorkoutSessionFromModel(workoutSession *model.WorkoutSessionImpl, exercises Exercises) *GetWorkoutSession {
	r.ID = workoutSession.ID
	r.Date = workoutSession.Date.Format("2006-01-02")
	r.UserID = workoutSession.UserID
	r.Exercises = exercises
	return r
}

func (r *Exercise) ExerciseFromModel(exercise *model.ExerciseImpl, sets *model.Sets) *Exercise {
	r.ID = exercise.ID
	r.SessionID = exercise.SessionID
	r.ExerciseName = exercise.ExerciseName
	r.Sets = *r.SetFromModel(sets)
	return r
}

func (r *Exercise) SetFromModel(sets *model.Sets) *Sets {
	var responseSets Sets
	if sets == nil {
		return &responseSets
	}
	for _, set := range *sets {
		responseSets = append(responseSets, Set{
			ID:         set.ID,
			ExerciseID: set.ExerciseID,
			SetNumber:  set.SetNumber,
			Weight:     set.Weight,
			Reps:       set.Reps,
		})
	}
	return &responseSets
}
