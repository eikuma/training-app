package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// func TestWorkoutSessionLoadByIDAndDate(t *testing.T) {
// 	today := time.Date(2022, 10, 1, 0, 0, 0, 0, time.UTC)
// 	tomorrow := today.Add(24 * time.Hour)

// 	// Create でレコードを挿入
// 	ws1, err := NewWorkoutSession().Create(today, int64(11))
// 	assert.NoError(t, err)
// 	ws2, err := NewWorkoutSession().Create(today, int64(22))
// 	assert.NoError(t, err)
// 	ws3, err := NewWorkoutSession().Create(tomorrow, int64(33))
// 	assert.NoError(t, err)

// 	type args struct {
// 		id   int64
// 		date time.Time
// 	}
// 	tests := []struct {
// 		testCase string
// 		args     args
// 		want     *WorkoutSessions
// 	}{
// 		{
// 			testCase: "全件取得",
// 			args: args{
// 				id:   0,
// 				date: time.Time{}, // zero 値ならフィルタ条件に掛からない
// 			},
// 			want: &WorkoutSessions{
// 				*ws1, *ws2, *ws3,
// 			},
// 		},
// 		{
// 			testCase: "idで検索",
// 			args: args{
// 				id:   ws1.ID,
// 				date: time.Time{},
// 			},
// 			want: &WorkoutSessions{
// 				*ws1,
// 			},
// 		},
// 		{
// 			testCase: "dateで検索",
// 			args: args{
// 				id:   0,
// 				date: today,
// 			},
// 			want: &WorkoutSessions{
// 				*ws1, *ws2,
// 			},
// 		},
// 	}

// 	// 比較時に順序は保証されないのでソートオプションを利用
// 	opts := cmp.Options{
// 		cmpopts.SortSlices(func(a, b WorkoutSessionImpl) bool {
// 			return a.ID < b.ID
// 		}),
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.testCase, func(t *testing.T) {
// 			got, err := NewWorkoutSession().LoadByIDAndDate(tt.args.id, tt.args.date)
// 			assert.NoError(t, err)
// 			if diff := cmp.Diff(*tt.want, *got, opts); diff != "" {
// 				t.Errorf("LoadByIDAndDate mismatch (-want +got):\n%s", diff)
// 			}
// 		})
// 	}
// }

func TestWorkoutSessionLoad(t *testing.T) {
	date := time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)
	ws, err := NewWorkoutSession().Create(date, 1)
	assert.NoError(t, err)

	m, err := new(WorkoutSessionImpl).Load(ws.ID)

	if assert.NoError(t, err) {
		assert.Equal(t, ws.ID, m.ID)
		assert.Equal(t, ws.Date, m.Date)
		assert.Equal(t, ws.UserID, m.UserID)
	}
}

func TestWorkoutSessionUpdate(t *testing.T) {
	date := time.Now().Truncate(24 * time.Hour)
	ws, err := NewWorkoutSession().Create(date, 1)
	assert.NoError(t, err)

	// user_id を更新
	newUserID := int64(99)
	updated, err := ws.Update(map[string]interface{}{"user_id": newUserID})

	if assert.NoError(t, err) {
		assert.True(t, updated)
	}
}

func TestWorkoutSessionCreate(t *testing.T) {
	date := time.Now().Truncate(24 * time.Hour)
	m, err := NewWorkoutSession().Create(date, int64(42))

	if assert.NoError(t, err) {
		assert.Equal(t, date, m.Date)
		assert.Equal(t, int64(42), m.UserID)
	}
}
