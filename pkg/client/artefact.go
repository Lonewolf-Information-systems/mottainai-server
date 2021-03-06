/*

Copyright (C) 2018  Ettore Di Giacinto <mudler@gentoo.org>
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

package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	storageci "github.com/MottainaiCI/mottainai-server/pkg/storage"
)

func (d *Fetcher) NamespaceFileList(namespace string) []string {
	var fileList []string
	err := d.GetJSONOptions("/api/namespace/"+namespace+"/list", map[string]string{}, &fileList)
	if err != nil {
		return []string{}
	}

	return fileList
}

func (d *Fetcher) StorageFileList(storage string) []string {
	var fileList []string
	err := d.GetJSONOptions("/api/storage/"+storage+"/list", map[string]string{}, &fileList)
	if err != nil {
		return []string{}
	}

	return fileList
}

func (d *Fetcher) TaskFileList(task string) []string {
	var fileList []string
	err := d.GetJSONOptions("/api/tasks/"+task+"/artefacts", map[string]string{}, &fileList)
	if err != nil {
		return []string{}
	}

	return fileList
}

func (d *Fetcher) DownloadArtefactsFromTask(taskid, target string) {
	list := d.TaskFileList(taskid)

	os.MkdirAll(target, os.ModePerm)
	d.AppendTaskOutput("Downloading artefacts from " + taskid)
	for _, file := range list {
		trials := 5
		done := true

		reldir, _ := filepath.Split(file)
		for done {

			if trials < 0 {
				done = false
			}
			os.MkdirAll(filepath.Join(target, reldir), os.ModePerm)
			d.AppendTaskOutput("Downloading  " + d.BaseURL + "/artefact/" + taskid + "/" + file + " to " + filepath.Join(target, file))
			if ok, err := d.Download(d.BaseURL+"/artefact/"+taskid+"/"+file, filepath.Join(target, file)); !ok {
				d.AppendTaskOutput("Downloading failed : " + err.Error())
				trials--
			} else {
				done = false
				d.AppendTaskOutput("Downloading succeeded ")

			}

		}

	}

}

func (d *Fetcher) DownloadArtefactsFromStorage(storage, target string) {
	list := d.StorageFileList(storage)

	var storage_data storageci.Storage
	err := d.GetJSONOptions("/api/storage/"+storage+"/show", map[string]string{}, &storage_data)
	if err != nil {
		d.AppendTaskOutput("Downloading failed during retrieveing storage data : " + err.Error())
		return
	}

	os.MkdirAll(target, os.ModePerm)
	d.AppendTaskOutput("Downloading artefacts from " + storage_data.Name)
	for _, file := range list {
		trials := 5
		done := true

		reldir, _ := filepath.Split(file)
		for done {

			if trials < 0 {
				done = false
			}
			os.MkdirAll(filepath.Join(target, reldir), os.ModePerm)
			d.AppendTaskOutput("Downloading  " + d.BaseURL + "/storage/" + storage_data.Path + "/" + file + " to " + filepath.Join(target, file))
			if ok, err := d.Download(d.BaseURL+"/storage/"+storage_data.Path+"/"+file, filepath.Join(target, file)); !ok {
				d.AppendTaskOutput("Downloading failed : " + err.Error())
				trials--
			} else {
				done = false
				d.AppendTaskOutput("Downloading succeeded ")

			}

		}

	}

}

func (d *Fetcher) DownloadArtefactsFromNamespace(namespace, target string) {
	list := d.NamespaceFileList(namespace)
	os.MkdirAll(target, os.ModePerm)
	d.AppendTaskOutput("Downloading artefacts from " + namespace)
	for _, file := range list {
		trials := 5
		done := true

		reldir, _ := filepath.Split(file)
		for done {

			if trials < 0 {
				done = false
			}
			os.MkdirAll(filepath.Join(target, reldir), os.ModePerm)
			d.AppendTaskOutput("Downloading  " + d.BaseURL + "/namespace/" + namespace + "/" + file + " to " + filepath.Join(target, file))
			if ok, err := d.Download(d.BaseURL+"/namespace/"+namespace+"/"+file, filepath.Join(target, file)); !ok {
				d.AppendTaskOutput("Downloading failed : " + err.Error())
				trials--
			} else {
				done = false
				d.AppendTaskOutput("Downloading succeeded ")

			}

		}

	}

}

func (d *Fetcher) Download(url, where string) (bool, error) {
	fileName := where

	output, err := os.Create(fileName)
	if err != nil {
		return false, err
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (f *Fetcher) UploadStorageFile(storageid, fullpath, relativepath string) error {
	_, file := filepath.Split(fullpath)

	opts := map[string]string{
		"name":      file,
		"path":      relativepath,
		"storageid": storageid,
		//	"namespace": dir,
	}

	request, err := f.Upload("/api/storage/upload", opts, "file", fullpath)

	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	} else {
		var bodyContent []byte
		fmt.Println(strconv.Itoa(resp.StatusCode))
		resp.Body.Read(bodyContent)
		resp.Body.Close()
		fmt.Println(string(bodyContent))
	}
	return nil
}

func (f *Fetcher) UploadArtefact(fullpath, relativepath string) error {
	_, file := filepath.Split(fullpath)

	opts := map[string]string{
		"name":   file,
		"path":   relativepath,
		"taskid": f.docID,
		//	"namespace": dir,
	}

	request, err := f.Upload("/api/tasks/artefact/upload", opts, "file", fullpath)

	if err != nil {
		f.AppendTaskOutput(err.Error())
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		f.AppendTaskOutput(err.Error())
		return err
	} else {
		var bodyContent []byte
		f.AppendTaskOutput(strconv.Itoa(resp.StatusCode))
		resp.Body.Read(bodyContent)
		resp.Body.Close()
		f.AppendTaskOutput(string(bodyContent))
	}
	return nil
}
