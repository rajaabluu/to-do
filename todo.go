package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
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

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "status"},
			{Align: simpletable.AlignRight, Text: "Created At"},
			{Align: simpletable.AlignRight, Text: "Completed At"},
		},
	}

	var cells [][]*simpletable.Cell

	for i, item := range *t {

		var status string
		completedAt := item.CompletedAt.Format(time.RFC822)

		if item.CompletedAt.IsZero() {
			completedAt = "Not Completed"
		}

		task := blue(item.Task)

		if item.Done {
			task = green(item.Task)
			status = "Done"
		} else {
			status = "Pending"
		}

		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", i+1)},
			{Text: task},
			{Text: status},
			{Text: item.CreatedAt.Format(time.RFC822)},
			{Text: completedAt},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: red(fmt.Sprintf("You have %d pending todos", t.CountPending()))},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func (t *Todos) CountPending() int {
	total := 0

	for _, item := range *t {
		if !item.Done {
			total++
		}
	}

	return total
}
