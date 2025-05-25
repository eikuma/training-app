package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestExerciseLoadBySessionID(t *testing.T) {

// 	type args struct {
// 		sessionId int64
// 	}
// 	tests := []struct {
// 		testCase string
// 		args     args
// 		want     *Exercises
// 	}{
// 		{
// 			testCase: "全件取得",
// 			args: args{
// 				sessionId: 0,
// 			},
// 			want: &Exercises{
// 				{},
// 			},
// 		},
// 		{
// 			testCase: "idで検索",
// 			args: args{
// 				sessionId: int64(23),
// 			},
// 			want: &Exercises{
// 				{},
// 			},
// 		},
// 	}

// 	// 比較時に順序は保証されないのでソートオプションを利用
// 	opts := cmp.Options{
// 		cmpopts.SortSlices(func(a, b ExerciseImpl) bool {
// 			return a.ID < b.ID
// 		}),
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.testCase, func(t *testing.T) {
// 			got, err := NewExercise().LoadBySessionID(tt.args.sessionId)
// 			assert.NoError(t, err)
// 			if diff := cmp.Diff(*tt.want, *got, opts); diff != "" {
// 				t.Errorf("LoadByIDAndDate mismatch (-want +got):\n%s", diff)
// 			}
// 		})
// 	}
// }

func TestExerciseLoad(t *testing.T) {
	e, err := NewExercise().Create(int64(23), "チェストプレス")
	assert.NoError(t, err)

	m, err := new(ExerciseImpl).Load(e.ID)

	if assert.NoError(t, err) {
		assert.Equal(t, e.ID, m.ID)
		assert.Equal(t, e.SessionID, m.SessionID)
		assert.Equal(t, e.ExerciseName, m.ExerciseName)
	}
}

func TestExerciseUpdate(t *testing.T) {
	e, err := NewExercise().Create(int64(23), "チェストプレス")
	assert.NoError(t, err)

	updated, err := e.Update(map[string]interface{}{"exercise_name": "ベンチプレス"})

	if assert.NoError(t, err) {
		assert.True(t, updated)
	}
}

func TestExerciseCreate(t *testing.T) {
	e, err := NewExercise().Create(int64(22), "チェストプレス")

	if assert.NoError(t, err) {
		assert.Equal(t, int64(22), e.SessionID)
		assert.Equal(t, "チェストプレス", e.ExerciseName)
	}
}
