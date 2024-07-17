package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todos []item

func (t *Todos) Add(task string) {

	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo)
}

func (t *Todos) Complete(i int) error {

	list := *t
	if i <= 0 || i > len(list) {
		return errors.New("invalid task")
	}

	list[i-1].CompletedAt = time.Now()
	list[i-1].Done = true

	return nil

}

func (t *Todos) Delete(i int) error {

	list := *t
	if i <= 0 || i > len(list) {
		return errors.New("invalid task")
	}

	*t = append(list[:i-1], list[i:]...)

	return nil

}

func (t *Todos) Load(filename string) error {

	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return err
	}

	result := json.Unmarshal(file, t)
	if result != nil {
		return result
	}

	return nil

}

func (t *Todos) Store(filename string) error {

	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func (t *Todos) Print() {

	if len(*t) > 0 {
		for i, item := range *t {
			fmt.Printf("%d - %s\n", i+1, item.Task)
		}
	} else {
		fmt.Println("No tasks added")
	}
}
