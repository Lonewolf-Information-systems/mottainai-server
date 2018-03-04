/*

Copyright (C) 2017-2018  Ettore Di Giacinto <mudler@gentoo.org>
Credits goes also to Gogs authors, some code portions and re-implemented design
are also coming from the Gogs project, which is using the go-macaron framework
and was really source of ispiration. Kudos to them!

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/

package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/MottainaiCI/mottainai-server/pkg/artefact"
	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
	"github.com/MottainaiCI/mottainai-server/pkg/tasks"
)

var TaskColl = "Tasks"

func (d *Database) CreateTask(t map[string]interface{}) (int, error) {

	return d.InsertDoc(TaskColl, t)
}

func (d *Database) DeleteTask(docID int) error {

	artefacts, err := d.GetTaskArtefacts(docID)
	if err != nil {
		return err
	}
	for _, artefact := range artefacts {
		artefact.CleanFromTask()
		d.DeleteArtefact(artefact.ID)
	}
	os.RemoveAll(filepath.Join(setting.Configuration.ArtefactPath, strconv.Itoa(docID)))

	return d.DeleteDoc(TaskColl, docID)
}

func (d *Database) UpdateTask(docID int, t map[string]interface{}) error {
	return d.UpdateDoc(TaskColl, docID, t)
}

func (d *Database) GetTask(docID int) (agenttasks.Task, error) {
	doc, err := d.GetDoc(TaskColl, docID)
	if err != nil {
		return agenttasks.Task{}, err
	}
	t := agenttasks.NewTaskFromMap(doc)
	t.ID = docID
	return t, err
}

func (d *Database) GetTaskArtefacts(id int) ([]artefact.Artefact, error) {
	queryResult, err := d.FindDoc(ArtefactColl, `[{"eq": `+strconv.Itoa(id)+`, "in": ["task"]}]`)
	var res []artefact.Artefact
	if err != nil {
		panic(err)
		//return res, err
	}

	fmt.Println("Getting artefacts")
	// Query result are document IDs
	for docid := range queryResult {
		fmt.Println("Got", docid)

		// Read document
		art, err := d.GetArtefact(docid)
		if err != nil {
			panic(err)
			//return res, err
		}

		res = append(res, art)
	}
	return res, nil
}

func (d *Database) ListTasks() []DocItem {
	return d.ListDocs(TaskColl)
}

func (d *Database) AllTasks() []agenttasks.Task {
	tasks := d.DB().Use(TaskColl)
	tasks_id := make([]agenttasks.Task, 0)

	tasks.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		t := agenttasks.NewTaskFromJson(docContent)
		t.ID = id
		tasks_id = append(tasks_id, t)
		return true
	})
	return tasks_id
}