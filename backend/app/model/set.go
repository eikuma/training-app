package model

import (
	"github.com/everytv/pre-employment-training-2024/final/ikuma.esaki/backend/db"
	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
)

type (
	// Set ワークアウトのインターフェースを表す
	Set interface {
		LoadByExerciseID(exerciseId int64) (*Sets, error)
		Load(id int64) (*SetImpl, error)
		Update(attrs map[string]interface{}) (bool, error)
		Create(exerciseID int64, setNumber int64, weight float64, reps int64) (*SetImpl, error)
	}

	// SetImpl ワークアウトを表す
	SetImpl struct {
		ID         int64   `db:"set_id" dbopt:"auto_increment"`
		ExerciseID int64   `db:"exercise_id"`
		SetNumber  int64   `db:"set_number"`
		Weight     float64 `db:"weight"`
		Reps       int64   `db:"reps"`
	}

	Sets []SetImpl
)

func NewSets() *Sets {
	return &Sets{}
}

func NewSet() Set {
	return &SetImpl{}
}

func (r *SetImpl) LoadByExerciseID(exerciseId int64) (*Sets, error) {
	return r.LoadByExerciseIDTx(db.GetSession("training_db"), exerciseId)
	// return nil, nil
}

func (r *SetImpl) LoadByExerciseIDTx(tx *dbr.Session, exerciseId int64) (*Sets, error) {
	m := NewSets()

	builder := tx.Select("*").From("sets")

	if exerciseId != 0 {
		builder = builder.Where("exercise_id = ?", exerciseId)
	}

	if _, err := builder.Load(m); err != nil {
		return nil, errors.Wrapf(err, "couldn't load sets")
	}
	return m, nil
}

// Load 指定のIDを読み込み
func (m *SetImpl) Load(id int64) (*SetImpl, error) {
	return m.LoadTx(db.GetSession("training_db"), id)
	// return nil, nil
}

// LoadTx トランザクション内で指定のIDを読み込み
func (m *SetImpl) LoadTx(tx dbr.SessionRunner, id int64) (*SetImpl, error) {
	if _, err := tx.Select("*").From("sets").Where("set_id=?", id).Load(m); err != nil {
		return nil, errors.Wrapf(err, "couldn't load sets")
	}
	return m, nil
}

// Update 更新
func (m *SetImpl) Update(attrs map[string]interface{}) (bool, error) {
	return m.UpdateTx(db.GetSession("training_db"), attrs)
	// return false, nil
}

// UpdateTx トランザクション内で更新
func (m *SetImpl) UpdateTx(tx dbr.SessionRunner, attrs map[string]interface{}) (bool, error) {
	res, err := tx.Update("sets").SetMap(attrs).Where("set_id=?", m.ID).Exec()
	if err != nil {
		return false, errors.Wrapf(err, "couldn't update sets")
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, errors.Wrapf(err, "couldn't fetch result")
	}

	return rows == 1, nil
}

// Create 作成
func (r *SetImpl) Create(exerciseID int64, setNumber int64, weight float64, reps int64) (*SetImpl, error) {
	return r.CreateTx(db.GetSession("training_db"), exerciseID, setNumber, weight, reps)
	// return nil, nil
}

// CreateTx トランザクション内で作成
func (r *SetImpl) CreateTx(tx dbr.SessionRunner, exerciseID int64, setNumber int64, weight float64, reps int64) (*SetImpl, error) {
	m := &SetImpl{
		ExerciseID: exerciseID,
		SetNumber:  setNumber,
		Weight:     weight,
		Reps:       reps,
	}

	res, err := tx.InsertInto("sets").
		Columns("exercise_id", "set_number", "weight", "reps").
		Record(m).
		Exec()

	if err != nil {
		return nil, errors.Wrapf(err, "couldn't create sets")
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get last insert id for exercises")
	}
	m.ID = lastID
	return m, nil
}
