package model

import (
	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/db"
	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
)

type (
	// Exercise ワークアウトのインターフェースを表す
	Exercise interface {
		LoadBySessionID(sessionId int64) (*Exercises, error)
		Load(id int64) (*ExerciseImpl, error)
		Update(attrs map[string]interface{}) (bool, error)
		Create(sessionId int64, exerciseName string) (*ExerciseImpl, error)
	}

	// ExerciseImpl ワークアウトを表す
	ExerciseImpl struct {
		ID           int64  `db:"exercise_id" dbopt:"auto_increment"`
		SessionID    int64  `db:"session_id"`
		ExerciseName string `db:"exercise_name"`
	}

	Exercises []ExerciseImpl
)

func NewExercises() *Exercises {
	return &Exercises{}
}

func NewExercise() Exercise {
	return &ExerciseImpl{}
}

func (r *ExerciseImpl) LoadBySessionID(sessionId int64) (*Exercises, error) {
	return r.LoadBySessionIDTx(db.GetSession("training_db"), sessionId)
	// return nil, nil
}

func (r *ExerciseImpl) LoadBySessionIDTx(tx *dbr.Session, sessionId int64) (*Exercises, error) {
	m := NewExercises()

	builder := tx.Select("*").From("exercises")

	if sessionId != 0 {
		builder = builder.Where("session_id = ?", sessionId)
	}

	if _, err := builder.Load(m); err != nil {
		return nil, errors.Wrapf(err, "couldn't load exercises")
	}
	return m, nil
}

// Load 指定のIDを読み込み
func (m *ExerciseImpl) Load(id int64) (*ExerciseImpl, error) {
	return m.LoadTx(db.GetSession("training_db"), id)
	// return nil, nil
}

// LoadTx トランザクション内で指定のIDを読み込み
func (m *ExerciseImpl) LoadTx(tx dbr.SessionRunner, id int64) (*ExerciseImpl, error) {
	if _, err := tx.Select("*").From("exercises").Where("exercise_id=?", id).Load(m); err != nil {
		return nil, errors.Wrapf(err, "couldn't load exercises")
	}
	return m, nil
}

// Update 更新
func (m *ExerciseImpl) Update(attrs map[string]interface{}) (bool, error) {
	return m.UpdateTx(db.GetSession("training_db"), attrs)
	// return false, nil
}

// UpdateTx トランザクション内で更新
func (m *ExerciseImpl) UpdateTx(tx dbr.SessionRunner, attrs map[string]interface{}) (bool, error) {
	res, err := tx.Update("exercises").SetMap(attrs).Where("exercise_id=?", m.ID).Exec()
	if err != nil {
		return false, errors.Wrapf(err, "couldn't update exercises")
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, errors.Wrapf(err, "couldn't fetch result")
	}

	return rows == 1, nil
}

// Create 作成
func (r *ExerciseImpl) Create(sessionId int64, exerciseName string) (*ExerciseImpl, error) {
	return r.CreateTx(db.GetSession("training_db"), sessionId, exerciseName)
	// return nil, nil
}

// CreateTx トランザクション内で作成
func (r *ExerciseImpl) CreateTx(tx dbr.SessionRunner, sessionId int64, exerciseName string) (*ExerciseImpl, error) {
	m := &ExerciseImpl{
		SessionID:    sessionId,
		ExerciseName: exerciseName,
	}

	res, err := tx.InsertInto("exercises").
		Columns("session_id", "exercise_name").
		Record(m).
		Exec()

	if err != nil {
		return nil, errors.Wrapf(err, "couldn't create exercises")
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get last insert id for exercises")
	}
	m.ID = lastID
	return m, nil
}
