package service

import (
	"fmt"
	"time"

	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/model"
	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/response"
)

type (
	// Workout ワークアウトのサービスを表す
	Workout interface {
		List(requestingUserID int64, sessionID int64, date time.Time) (response.WorkoutSessions, error) // Modified signature
		Get(id int64) (*response.GetWorkoutSession, error)
		CreateWorkoutSession(date time.Time, userId int64) (*response.WorkoutSession, error)
		CreateExercise(sessionId int64, exerciseName string) (*response.Exercise, error)
		CreateSet(exerciseID int64, setNumber int64, weight float64, reps int64) (*response.Sets, error)
	}

	// WorkoutImpl ワークアウトのサービスを表す
	WorkoutImpl struct {
		WorkoutSession model.WorkoutSession
		Exercise       model.Exercise
		Set            model.Set
	}
)

func NewWorkout() Workout {
	return &WorkoutImpl{
		WorkoutSession: model.NewWorkoutSession(),
		Exercise:       model.NewExercise(),
		Set:            model.NewSet(),
	}
}

// List ワークアウトの一覧を取得
func (s *WorkoutImpl) List(requestingUserID int64, sessionID int64, date time.Time) (response.WorkoutSessions, error) { // Modified signature
	workoutSessionsImpl, err := s.WorkoutSession.LoadByIDAndDate(requestingUserID, sessionID, date) // Pass requestingUserID, use sessionID
	if err != nil {
		return nil, err
	}

	var responseWorkoutSessions response.WorkoutSessions
	// Ensure correct processing of workoutSessionsImpl (pointer to slice)
	if workoutSessionsImpl != nil {
		for _, workoutSession := range *workoutSessionsImpl {
			responseWorkoutSessions = append(responseWorkoutSessions, *response.NewWorkoutSession().WorkoutSessionFromModel(&workoutSession))
		}
	}

	return responseWorkoutSessions, nil
}

// Get ワークアウトの詳細を取得
func (s *WorkoutImpl) Get(id int64) (*response.GetWorkoutSession, error) {
	workoutSession, err := s.WorkoutSession.Load(id)
	if err != nil {
		return nil, err
	}
	if workoutSession.ID != id || workoutSession.ID == 0 {
		return nil, fmt.Errorf("workout session not found. id %d", id)
	}

	exercises, err := s.Exercise.LoadBySessionID(workoutSession.ID)
	if err != nil {
		return nil, err
	}

	var responseExercises response.Exercises
	for _, exercise := range *exercises {
		sets, err := s.Set.LoadByExerciseID(exercise.ID)
		if err != nil {
			return nil, err
		}
		responseExercises = append(responseExercises, *response.NewExercise().ExerciseFromModel(&exercise, sets))
	}

	return response.NewGetWorkoutSession().GetWorkoutSessionFromModel(workoutSession, responseExercises), nil
}

func (s *WorkoutImpl) CreateWorkoutSession(date time.Time, userId int64) (*response.WorkoutSession, error) {
	workoutSession, err := s.WorkoutSession.Create(date, userId)
	if err != nil {
		return nil, err
	}

	workoutSession, err = s.WorkoutSession.Load(workoutSession.ID)
	if err != nil {
		return nil, err
	}

	return response.NewWorkoutSession().WorkoutSessionFromModel(workoutSession), nil
}

func (s *WorkoutImpl) CreateExercise(sessionId int64, exerciseName string) (*response.Exercise, error) {
	exercise, err := s.Exercise.Create(sessionId, exerciseName)
	if err != nil {
		return nil, err
	}

	exercise, err = s.Exercise.Load(exercise.ID)
	if err != nil {
		return nil, err
	}

	return response.NewExercise().ExerciseFromModel(exercise, nil), nil
}

func (s *WorkoutImpl) CreateSet(exerciseID int64, setNumber int64, weight float64, reps int64) (*response.Sets, error) {
	set, err := s.Set.Create(exerciseID, setNumber, weight, reps)
	if err != nil {
		return nil, err
	}

	sets, err := s.Set.LoadByExerciseID(set.ExerciseID)
	if err != nil {
		return nil, err
	}

	return response.NewExercise().SetFromModel(sets), nil
}
