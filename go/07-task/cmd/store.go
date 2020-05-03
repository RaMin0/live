package cmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

const (
	dbFileName = "tasks.db"
)

var (
	tasksBucketName = []byte("tasks")
)

type Task struct {
	ID        int    `json:"id"`
	Details   string `json:"details"`
	Completed bool   `json:"completed"`
}

func CreateTask(task *Task) error {
	return withDB(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(tasksBucketName)

			id, err := bucket.NextSequence()
			if err != nil {
				return err
			}
			task.ID = int(id)

			b, err := json.Marshal(&task)
			if err != nil {
				return err
			}

			return bucket.Put(itob(task.ID), b)
		})
	})
}

func ListTasks(completed bool) ([]*Task, error) {
	var tasks []*Task
	return tasks, withDB(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(tasksBucketName)

			return bucket.ForEach(func(_, b []byte) error {
				var task Task
				if err := json.Unmarshal(b, &task); err != nil {
					return err
				}
				if task.Completed != completed {
					return nil
				}
				tasks = append(tasks, &task)
				return nil
			})
		})
	})
}

func MarkTaskAsCompleted(task *Task) error {
	return withDB(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(tasksBucketName)

			b := bucket.Get(itob(task.ID))
			if b == nil {
				return fmt.Errorf("task not found with ID=%d", task.ID)
			}
			if err := json.Unmarshal(b, task); err != nil {
				return err
			}
			task.Completed = true
			b, err := json.Marshal(&task)
			if err != nil {
				return err
			}
			return bucket.Put(itob(task.ID), b)
		})
	})
}

func DeleteTask(task *Task) error {
	return withDB(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(tasksBucketName)

			b := bucket.Get(itob(task.ID))
			if b == nil {
				return fmt.Errorf("task not found with ID=%d", task.ID)
			}

			if err := json.Unmarshal(b, task); err != nil {
				return err
			}
			return bucket.Delete(itob(task.ID))
		})
	})
}

func withDB(fn func(*bolt.DB) error) error {
	db, err := bolt.Open(dbFileName, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(tasksBucketName)
		return err
	}); err != nil {
		return err
	}

	return fn(db)
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
