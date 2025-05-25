package service

import (
	"errors"
	"testing"
	"time"

	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/model"
	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/model/mock_model"
	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/app/response"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWorkoutList(t *testing.T) {
	t.Parallel()
	type fields struct {
		WorkoutSession model.WorkoutSession
	}
	type args struct {
		id   int64
		date time.Time
	}
	tests := []struct {
		testCase  string
		args      args
		fields    func(ctrl *gomock.Controller) fields
		assertion func(r response.WorkoutSessions, err error)
	}{
		{
			testCase: "正常系",
			args: args{
				id:   int64(0),
				date: time.Time{},
			},
			fields: func(ctrl *gomock.Controller) fields {
				WorkoutSession := mock_model.NewMockWorkoutSession(ctrl)
				WorkoutSession.EXPECT().LoadByIDAndDate(int64(0), time.Time{}).Return(&model.WorkoutSessions{
					{ID: int64(1), Date: time.Now(), UserID: int64(1)},
					{ID: int64(2), Date: time.Now(), UserID: int64(1)},
					{ID: int64(3), Date: time.Now(), UserID: int64(1)},
				}, nil)
				return fields{
					WorkoutSession: WorkoutSession,
				}
			},
			assertion: func(r response.WorkoutSessions, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, r)
				assert.Len(t, r, 3)
			},
		},
		{
			testCase: "正常系(1件取得id指定)",
			args: args{
				id:   int64(1),
				date: time.Time{},
			},
			fields: func(ctrl *gomock.Controller) fields {
				WorkoutSession := mock_model.NewMockWorkoutSession(ctrl)
				WorkoutSession.EXPECT().LoadByIDAndDate(int64(1), time.Time{}).Return(&model.WorkoutSessions{
					{ID: int64(1), Date: time.Now(), UserID: int64(1)},
				}, nil)
				return fields{
					WorkoutSession: WorkoutSession,
				}
			},
			assertion: func(r response.WorkoutSessions, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, r)
				assert.Len(t, r, 1)
			},
		},
		{
			testCase: "正常系(1件取得日付指定)",
			args: args{
				id:   int64(0),
				date: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
			},
			fields: func(ctrl *gomock.Controller) fields {
				WorkoutSession := mock_model.NewMockWorkoutSession(ctrl)
				WorkoutSession.EXPECT().LoadByIDAndDate(int64(0), time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC)).Return(&model.WorkoutSessions{
					{ID: int64(1), Date: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC), UserID: int64(1)},
				}, nil)
				return fields{
					WorkoutSession: WorkoutSession,
				}
			},
			assertion: func(r response.WorkoutSessions, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, r)
				assert.Len(t, r, 1)
			},
		},
		{
			testCase: "正常系(空)",
			args: args{
				id:   int64(100),
				date: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
			},
			fields: func(ctrl *gomock.Controller) fields {
				WorkoutSession := mock_model.NewMockWorkoutSession(ctrl)
				WorkoutSession.EXPECT().LoadByIDAndDate(int64(100), time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC)).Return(&model.WorkoutSessions{}, nil)
				return fields{
					WorkoutSession: WorkoutSession,
				}
			},
			assertion: func(r response.WorkoutSessions, err error) {
				assert.NoError(t, err)
				assert.Nil(t, r)
				assert.Len(t, r, 0)
			},
		},
		{
			testCase: "エラー",
			args: args{
				id:   int64(100),
				date: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
			},
			fields: func(ctrl *gomock.Controller) fields {
				WorkoutSession := mock_model.NewMockWorkoutSession(ctrl)
				WorkoutSession.EXPECT().LoadByIDAndDate(int64(100), time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC)).Return(nil, errors.New("couldn't load workout_sessions"))
				return fields{
					WorkoutSession: WorkoutSession,
				}
			},
			assertion: func(r response.WorkoutSessions, err error) {
				assert.Error(t, err)
				assert.Nil(t, r)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.testCase, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			fields := tt.fields(ctrl)
			w := &WorkoutImpl{
				WorkoutSession: fields.WorkoutSession,
			}
			tt.assertion(w.List(tt.args.id, tt.args.date))
		})
	}
}

func TestWorkoutGet(t *testing.T) {
	t.Parallel()
	type fields struct {
		WorkoutSession model.WorkoutSession
		Exercise       model.Exercise
		Set            model.Set
	}
	type args struct {
		id int64
	}
	tests := []struct {
		testCase  string
		args      args
		fields    func(ctrl *gomock.Controller) fields
		assertion func(r *response.GetWorkoutSession, err error)
	}{
		{
			testCase: "正常系",
			args: args{
				id: int64(1),
			},
			fields: func(ctrl *gomock.Controller) fields {
				WorkoutSession := mock_model.NewMockWorkoutSession(ctrl)
				WorkoutSession.EXPECT().Load(int64(1)).Return(&model.WorkoutSessionImpl{ID: int64(1), Date: time.Now(), UserID: int64(1)}, nil)
				Exercise := mock_model.NewMockExercise(ctrl)
				Exercise.EXPECT().LoadBySessionID(int64(1)).Return(&model.Exercises{
					{ID: int64(1), SessionID: int64(1), ExerciseName: "test"},
					{ID: int64(2), SessionID: int64(1), ExerciseName: "test"},
				}, nil)
				Set := mock_model.NewMockSet(ctrl)
				Set.EXPECT().LoadByExerciseID(int64(1)).Return(&model.Sets{
					{ID: int64(1), ExerciseID: int64(1), SetNumber: int64(1), Weight: float64(10), Reps: int64(10)},
					{ID: int64(2), ExerciseID: int64(1), SetNumber: int64(2), Weight: float64(10), Reps: int64(10)},
				}, nil)
				Set.EXPECT().LoadByExerciseID(int64(2)).Return(&model.Sets{
					{ID: int64(1), ExerciseID: int64(2), SetNumber: int64(1), Weight: float64(10), Reps: int64(10)},
					{ID: int64(2), ExerciseID: int64(2), SetNumber: int64(2), Weight: float64(10), Reps: int64(10)},
				}, nil)
				return fields{
					WorkoutSession: WorkoutSession,
					Exercise:       Exercise,
					Set:            Set,
				}
			},
			assertion: func(r *response.GetWorkoutSession, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, r)
				assert.Equal(t, int64(1), r.ID)
				assert.Equal(t, int64(1), r.UserID)
				assert.Len(t, r.Exercises, 2)
				assert.Len(t, r.Exercises[0].Sets, 2)
				assert.Len(t, r.Exercises[1].Sets, 2)
			},
		},
		{
			testCase: "エラー",
			args: args{
				id: int64(0),
			},
			fields: func(ctrl *gomock.Controller) fields {
				WorkoutSession := mock_model.NewMockWorkoutSession(ctrl)
				WorkoutSession.EXPECT().Load(int64(0)).Return(nil, errors.New("couldn't load workout_session"))
				return fields{
					WorkoutSession: WorkoutSession,
				}
			},
			assertion: func(r *response.GetWorkoutSession, err error) {
				assert.Error(t, err)
				assert.Nil(t, r)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.testCase, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			fields := tt.fields(ctrl)
			w := &WorkoutImpl{
				WorkoutSession: fields.WorkoutSession,
				Exercise:       fields.Exercise,
				Set:            fields.Set,
			}
			tt.assertion(w.Get(tt.args.id))
		})
	}
}

func TestWorkoutCreateWorkoutSession(t *testing.T) {
	t.Parallel()
	type fields struct {
		WorkoutSession model.WorkoutSession
	}
	type args struct {
		date   time.Time
		userId int64
	}
	tests := []struct {
		testCase  string
		args      args
		fields    func(ctrl *gomock.Controller) fields
		assertion func(r *response.WorkoutSession, err error)
	}{
		{
			testCase: "正常系",
			args: args{
				date:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
				userId: int64(1),
			},
			fields: func(ctrl *gomock.Controller) fields {
				WorkoutSession := mock_model.NewMockWorkoutSession(ctrl)
				WorkoutSession.EXPECT().Create(time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC), int64(1)).Return(&model.WorkoutSessionImpl{ID: int64(1), Date: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC), UserID: int64(1)}, nil)
				WorkoutSession.EXPECT().Load(int64(1)).Return(&model.WorkoutSessionImpl{ID: int64(1), Date: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC), UserID: int64(1)}, nil)
				return fields{
					WorkoutSession: WorkoutSession,
				}
			},
			assertion: func(r *response.WorkoutSession, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, r)
				assert.Equal(t, int64(1), r.ID)
				assert.Equal(t, time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC).Format("2006-01-02"), r.Date)
				assert.Equal(t, int64(1), r.UserID)
			},
		},
		{
			testCase: "エラー",
			args: args{
				date:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
				userId: int64(0),
			},
			fields: func(ctrl *gomock.Controller) fields {
				WorkoutSession := mock_model.NewMockWorkoutSession(ctrl)
				WorkoutSession.EXPECT().Create(time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC), int64(0)).Return(nil, errors.New("couldn't create workout_session"))
				return fields{
					WorkoutSession: WorkoutSession,
				}
			},
			assertion: func(r *response.WorkoutSession, err error) {
				assert.Error(t, err)
				assert.Nil(t, r)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.testCase, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			fields := tt.fields(ctrl)
			w := &WorkoutImpl{
				WorkoutSession: fields.WorkoutSession,
			}
			tt.assertion(w.CreateWorkoutSession(tt.args.date, tt.args.userId))
		})
	}
}

func TestWorkoutCreateExercise(t *testing.T) {
	t.Parallel()
	type fields struct {
		Exercise model.Exercise
	}
	type args struct {
		sessionId    int64
		exerciseName string
	}
	tests := []struct {
		testCase  string
		args      args
		fields    func(ctrl *gomock.Controller) fields
		assertion func(r *response.Exercise, err error)
	}{
		{
			testCase: "正常系",
			args: args{
				sessionId:    int64(1),
				exerciseName: "test",
			},
			fields: func(ctrl *gomock.Controller) fields {
				Exercise := mock_model.NewMockExercise(ctrl)
				Exercise.EXPECT().Create(int64(1), "test").Return(&model.ExerciseImpl{ID: int64(1), SessionID: int64(1), ExerciseName: "test"}, nil)
				Exercise.EXPECT().Load(int64(1)).Return(&model.ExerciseImpl{ID: int64(1), SessionID: int64(1), ExerciseName: "test"}, nil)
				return fields{
					Exercise: Exercise,
				}
			},
			assertion: func(r *response.Exercise, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, r)
				assert.Equal(t, int64(1), r.ID)
				assert.Equal(t, int64(1), r.SessionID)
				assert.Equal(t, "test", r.ExerciseName)
			},
		},
		{
			testCase: "エラー",
			args: args{
				sessionId:    int64(0),
				exerciseName: "",
			},
			fields: func(ctrl *gomock.Controller) fields {
				Exercise := mock_model.NewMockExercise(ctrl)
				Exercise.EXPECT().Create(int64(0), "").Return(nil, errors.New("couldn't create exercise"))
				return fields{
					Exercise: Exercise,
				}
			},
			assertion: func(r *response.Exercise, err error) {
				assert.Error(t, err)
				assert.Nil(t, r)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.testCase, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			fields := tt.fields(ctrl)
			w := &WorkoutImpl{
				Exercise: fields.Exercise,
			}
			tt.assertion(w.CreateExercise(tt.args.sessionId, tt.args.exerciseName))
		})
	}
}

func TestWorkoutCreateSet(t *testing.T) {
	t.Parallel()
	type fields struct {
		Set model.Set
	}
	type args struct {
		exerciseID int64
		setNumber  int64
		weight     float64
		reps       int64
	}
	tests := []struct {
		testCase  string
		args      args
		fields    func(ctrl *gomock.Controller) fields
		assertion func(r *response.Sets, err error)
	}{
		{
			testCase: "正常系",
			args: args{
				exerciseID: int64(1),
				setNumber:  int64(1),
				weight:     float64(10),
				reps:       int64(10),
			},
			fields: func(ctrl *gomock.Controller) fields {
				Set := mock_model.NewMockSet(ctrl)
				Set.EXPECT().Create(int64(1), int64(1), float64(10), int64(10)).Return(&model.SetImpl{
					ID: int64(1), ExerciseID: int64(1), SetNumber: int64(1), Weight: float64(10), Reps: int64(10),
				}, nil)
				Set.EXPECT().LoadByExerciseID(int64(1)).Return(&model.Sets{
					{ID: int64(1), ExerciseID: int64(1), SetNumber: int64(1), Weight: float64(10), Reps: int64(10)},
					{ID: int64(2), ExerciseID: int64(1), SetNumber: int64(2), Weight: float64(10), Reps: int64(10)},
				}, nil)
				return fields{
					Set: Set,
				}
			},
			assertion: func(r *response.Sets, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, r)
			},
		},
		{
			testCase: "エラー",
			args: args{
				exerciseID: int64(1),
				setNumber:  int64(1),
				weight:     float64(10),
				reps:       int64(10),
			},
			fields: func(ctrl *gomock.Controller) fields {
				Set := mock_model.NewMockSet(ctrl)
				Set.EXPECT().Create(int64(1), int64(1), float64(10), int64(10)).Return(nil, errors.New("couldn't create set"))
				return fields{
					Set: Set,
				}
			},
			assertion: func(r *response.Sets, err error) {
				assert.Error(t, err)
				assert.Nil(t, r)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.testCase, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			fields := tt.fields(ctrl)
			w := &WorkoutImpl{
				Set: fields.Set,
			}
			tt.assertion(w.CreateSet(tt.args.exerciseID, tt.args.setNumber, tt.args.weight, tt.args.reps))
		})
	}
}
