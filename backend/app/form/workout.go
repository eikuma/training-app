package form

type (
	// Workout ワークアウトのフォームを表す
	ListWorkout struct {
		ID   int64  `json:"id" form:"id" query:"id" description:"検索したいセッションID"`
		Date string `json:"date" form:"date" query:"date" description:"検索したい日付"`
	}

	CreateWorkoutSession struct {
		Date   string `json:"date" form:"date" query:"date" valid:"required" description:"ワークアウトの日付"`
		UserID int64  `json:"user_id" form:"user_id" query:"user_id" valid:"required" description:"ユーザーID"`
	}

	CreateExercise struct {
		// SessionID    int64  `json:"session_id" form:"session_id" query:"session_id" valid:"required" description:"ワークアウトセッションID"`
		ExerciseName string `json:"exercise_name" form:"exercise_name" query:"exercise_name" valid:"required" description:"エクササイズ名"`
	}

	CreateSet struct {
		// ExerciseID int64   `json:"exercise_id" form:"exercise_id" query:"exercise_id" valid:"required" description:"エクササイズID"`
		SetNumber int64   `json:"set_number" form:"set_number" query:"set_number" valid:"required" description:"セット数"`
		Weight    float64 `json:"weight" form:"weight" query:"weight" valid:"required" description:"重量"`
		Reps      int64   `json:"reps" form:"reps" query:"reps" valid:"required" description:"回数"`
	}
)

func NewListWorkout() *ListWorkout {
	return &ListWorkout{}
}

func NewCreateWorkoutSession() *CreateWorkoutSession {
	return &CreateWorkoutSession{}
}

func NewCreateExercise() *CreateExercise {
	return &CreateExercise{}
}

func NewCreateSet() *CreateSet {
	return &CreateSet{}
}
