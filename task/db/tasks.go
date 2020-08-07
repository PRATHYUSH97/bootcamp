package db

import (
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("tasks")
var completebucket = []byte("complete")
var db *bolt.DB

type Task struct {
	Key   int
	Value string
}

type Complete struct {
	Time  time.Time
	Value string
}

//connects to boltdb and creates taskBucket and completebucket
func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})

	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(completebucket)
		return err
	})

	if err != nil {
		return err
	}
	return nil
}

//inserts into taskBucket and returns the id with which task was inserted
func CreateTask(task string) (int, error) {
	var id int
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)
		return b.Put(key, []byte(task))
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

//inserts completed task into completebucket
func CreateComplete(task string) error {

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(completebucket)
		time := time.Now().String()
		return b.Put([]byte(time), []byte(task))
	})
	if err != nil {
		return err
	}
	return nil

}

//returns all the pending tasks
func AllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

//returns all the tasks completed within 14 hours
func Completedtoday() ([]string, error) {
	var tasks []string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(completebucket)
		c := b.Cursor()
		timestart := time.Now()
		layout := "Mon Jan 02 2006 15:04:05 GMT-0700"
		for k, v := c.First(); k != nil; k, v = c.Next() {

			key := string(k)
			timenow, _ := time.Parse(layout, key)
			hours := timenow.Sub(timestart).Hours()
			if hours < 14.00 {
				tasks = append(tasks, string(v))
			}

		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil

}

// deletes a task from taskbucket
func DeleteTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key))
	})
}

//converts integer to byte slice
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

//converts byte slice to integer
func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
