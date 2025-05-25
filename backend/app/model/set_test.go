package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestSetLoadByExerciseID(t *testing.T) {

// 	type args struct {
// 		exerciseId int64
// 	}
// 	tests := []struct {
// 		testCase string
// 		args     args
// 		want     *Sets
// 	}{
// 		{
// 			testCase: "全件取得",
// 			args: args{
// 				exerciseId: 0,
// 			},
// 			want: &Sets{
// 				{},
// 			},
// 		},
// 		{
// 			testCase: "idで検索",
// 			args: args{
// 				exerciseId: int64(5),
// 			},
// 			want: &Sets{
// 				{},
// 			},
// 		},
// 	}

// 	// 比較時に順序は保証されないのでソートオプションを利用
// 	opts := cmp.Options{
// 		cmpopts.SortSlices(func(a, b SetImpl) bool {
// 			return a.ID < b.ID
// 		}),
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.testCase, func(t *testing.T) {
// 			got, err := NewSet().LoadByExerciseID(tt.args.exerciseId)
// 			assert.NoError(t, err)
// 			if diff := cmp.Diff(*tt.want, *got, opts); diff != "" {
// 				t.Errorf("LoadByIDAndDate mismatch (-want +got):\n%s", diff)
// 			}
// 		})
// 	}
// }

func TestSetLoad(t *testing.T) {
	s, err := NewSet().Create(int64(5), int64(1), float64(35.0), int64(10))
	assert.NoError(t, err)

	m, err := new(SetImpl).Load(s.ID)

	if assert.NoError(t, err) {
		assert.Equal(t, s.ID, m.ID)
		assert.Equal(t, s.ExerciseID, m.ExerciseID)
		assert.Equal(t, s.SetNumber, m.SetNumber)
		assert.Equal(t, s.Weight, m.Weight)
		assert.Equal(t, s.Reps, m.Reps)
	}
}

func TestSetUpdate(t *testing.T) {
	s, err := NewSet().Create(int64(4), int64(1), float64(35.0), int64(10))
	assert.NoError(t, err)

	updated, err := s.Update(map[string]interface{}{"set_number": int64(2), "weight": float64(40.0), "reps": int64(12)})

	if assert.NoError(t, err) {
		assert.True(t, updated)
	}
}

func TestSetCreate(t *testing.T) {
	s, err := NewSet().Create(int64(4), int64(1), float64(35.0), int64(10))

	if assert.NoError(t, err) {
		assert.Equal(t, int64(4), s.ExerciseID)
		assert.Equal(t, int64(1), s.SetNumber)
		assert.Equal(t, float64(35.0), s.Weight)
		assert.Equal(t, int64(10), s.Reps)
	}
}
