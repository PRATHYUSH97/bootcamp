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

func CreateComplete(task string) (int, error) {
	var id int
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(completebucket)
		//id64, _ := b.NextSequence()
		//id = int(id64)
		//key := itob(id)
		time := time.Now().String()
		return b.Put([]byte(time), []byte(task))
	})
	if err != nil {
		return -1, err
	}
	return id, nil

}

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

func DeleteTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key))
	})
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
