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

package storagesapi

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/MottainaiCI/mottainai-server/pkg/context"
	"github.com/MottainaiCI/mottainai-server/pkg/db"
	"github.com/MottainaiCI/mottainai-server/pkg/settings"
	"github.com/MottainaiCI/mottainai-server/pkg/utils"
)

func StorageCreate(ctx *context.Context, db *database.Database) (string, error) {
	name := ctx.Params(":name")
	name, _ = utils.Strip(name)

	docID, err := db.CreateStorage(map[string]interface{}{
		"name": name,
		"path": name,
	})
	//
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(filepath.Join(setting.Configuration.StoragePath, name), os.ModePerm)
	if err != nil {
		return ":(", err
	}

	return strconv.Itoa(docID), nil
}

type StorageForm struct {
	ID         int                   `form:"storageid" binding:"Required"`
	Name       string                `form:"name"`
	Path       string                `form:"path"`
	FileUpload *multipart.FileHeader `form:"file"`
}

func StorageUpload(uf StorageForm, ctx *context.Context, db *database.Database) string {

	file, err := uf.FileUpload.Open()

	storage, err := db.GetStorage(uf.ID)
	defer file.Close()
	if err != nil {
		return err.Error()
	}

	os.MkdirAll(filepath.Join(setting.Configuration.StoragePath, storage.Path, uf.Path), os.ModePerm)
	f, err := os.OpenFile(filepath.Join(setting.Configuration.StoragePath, storage.Path, uf.Path, uf.Name), os.O_WRONLY|os.O_CREATE, os.ModePerm)

	defer f.Close()
	io.Copy(f, file)

	return "OK"
}
