package model

import (
	"time"

	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/db"
	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
)

type (
	// WorkoutSession ワークアウトのインターフェースを表す
	WorkoutSession interface {
		LoadByIDAndDate(id int64, date time.Time) (*WorkoutSessions, error)
		Load(id int64) (*WorkoutSessionImpl, error)
		Update(attrs map[string]interface{}) (bool, error)
		Create(date time.Time, userId int64) (*WorkoutSessionImpl, error)
	}

	// WorkoutSessionImpl ワークアウトを表す
	WorkoutSessionImpl struct {
		ID     int64     `db:"session_id" dbopt:"auto_increment"`
		Date   time.Time `db:"training_date"`
		UserID int64     `db:"user_id"`
	}

	WorkoutSessions []WorkoutSessionImpl
)

func NewWorkoutSessions() *WorkoutSessions {
	return &WorkoutSessions{}
}

func NewWorkoutSession() WorkoutSession {
	return &WorkoutSessionImpl{}
}

func (r *WorkoutSessionImpl) LoadByIDAndDate(id int64, date time.Time) (*WorkoutSessions, error) {
	return r.LoadByIDAndDateTx(db.GetSession("training_db"), id, date)
	// return nil, nil
}

func (r *WorkoutSessionImpl) LoadByIDAndDateTx(tx *dbr.Session, id int64, date time.Time) (*WorkoutSessions, error) {
	m := NewWorkoutSessions()

	builder := tx.Select("*").From("workout_sessions")

	if id != 0 {
		builder = builder.Where("session_id = ?", id)
	}
	if !date.IsZero() {
		builder = builder.Where("training_date = ?", date)
	}

	if _, err := builder.Load(m); err != nil {
		return nil, errors.Wrapf(err, "couldn't load workout_sessions")
	}
	return m, nil
}

// Load 指定のIDを読み込み
func (m *WorkoutSessionImpl) Load(id int64) (*WorkoutSessionImpl, error) {
	return m.LoadTx(db.GetSession("training_db"), id)
	// return nil, nil
}

// LoadTx トランザクション内で指定のIDを読み込み
func (m *WorkoutSessionImpl) LoadTx(tx dbr.SessionRunner, id int64) (*WorkoutSessionImpl, error) {
	if _, err := tx.Select("*").From("workout_sessions").Where("session_id=?", id).Load(m); err != nil {
		return nil, errors.Wrapf(err, "couldn't load workout_sessions")
	}
	return m, nil
}

// Update 更新
func (m *WorkoutSessionImpl) Update(attrs map[string]interface{}) (bool, error) {
	return m.UpdateTx(db.GetSession("training_db"), attrs)
	// return false, nil
}

// UpdateTx トランザクション内で更新
func (m *WorkoutSessionImpl) UpdateTx(tx dbr.SessionRunner, attrs map[string]interface{}) (bool, error) {
	res, err := tx.Update("workout_sessions").SetMap(attrs).Where("session_id=?", m.ID).Exec()
	if err != nil {
		return false, errors.Wrapf(err, "couldn't update workout_sessions")
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, errors.Wrapf(err, "couldn't fetch result")
	}

	return rows == 1, nil
}

// Create 作成
func (r *WorkoutSessionImpl) Create(date time.Time, userId int64) (*WorkoutSessionImpl, error) {
	return r.CreateTx(db.GetSession("training_db"), date, userId)
	// return nil, nil
}

// CreateTx トランザクション内で作成
func (r *WorkoutSessionImpl) CreateTx(tx dbr.SessionRunner, date time.Time, userId int64) (*WorkoutSessionImpl, error) {
	m := &WorkoutSessionImpl{
		Date:   date,
		UserID: userId,
	}

	res, err := tx.InsertInto("workout_sessions").
		Columns("training_date", "user_id").
		Record(m).
		Exec()
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't create workout_sessions")
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get last insert id for workout_sessions")
	}
	m.ID = lastID
	return m, nil
}
