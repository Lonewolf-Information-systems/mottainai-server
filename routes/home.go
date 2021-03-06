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

package routes

import (
	"github.com/MottainaiCI/mottainai-server/pkg/context"
	"github.com/MottainaiCI/mottainai-server/pkg/db"
	"github.com/MottainaiCI/mottainai-server/pkg/mottainai"
	"github.com/MottainaiCI/mottainai-server/pkg/template"

	"github.com/MottainaiCI/mottainai-server/routes/api"

	"github.com/MottainaiCI/mottainai-server/routes/nodes"
	"github.com/MottainaiCI/mottainai-server/routes/tasks"
)

func NotFound(c *context.Context) {
	c.Data["Title"] = "Page Not Found"
	c.NotFound()
}

func Setup(m *mottainai.Mottainai) {

	m.NotFound(NotFound)

	// setup templates
	// m.Use(macaron.Renderer())

	m.Get("/", func(ctx *context.Context, db *database.Database) {
		//ctx.Data["Name"] = "jeremy"
		rtasks, _ := db.FindDoc("Tasks", `[{"eq": "running", "in": ["status"]}]`)
		running_tasks := len(rtasks)
		wtasks, _ := db.FindDoc("Tasks", `[{"eq": "waiting", "in": ["status"]}]`)
		waiting_tasks := len(wtasks)
		etasks, _ := db.FindDoc("Tasks", `[{"eq": "error", "in": ["result"]}]`)
		error_tasks := len(etasks)

		ctx.Data["TotalTasks"] = db.DB().Use("Tasks").ApproxDocCount()
		if ctx.Data["TotalTasks"] == 0 {
			ctx.Data["TotalTasks"] = len(db.ListDocs("Tasks"))
		}
		ctx.Data["RunningTasks"] = running_tasks
		ctx.Data["WaitingTasks"] = waiting_tasks
		ctx.Data["ErroredTasks"] = error_tasks

		template.TemplatePreview(ctx, "index")
	})

	tasks.Setup(m)
	nodesroute.Setup(m)
	api.Setup(m)
}
