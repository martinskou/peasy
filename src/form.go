package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type FieldType int64

const (
	Text FieldType = iota
	Password
	Number
	Boolean
	File
	Files
	Select
	Hidden
	Space
)

type Option struct {
	Text  string
	Value string
}
type FieldFile struct {
	Filename    string
	TmpFilename string
}
type Field struct {
	Title       string
	Placeholder string
	Name        string
	Type        FieldType
	Value       string
	Options     []Option
	// File upload
	Files []FieldFile
	//	Filename     string
	//	Filenames    []string
	//	TmpFilename  string
	//	TmpFilenames []string
}

func (f Field) HasFiles() bool {
	return len(f.Files) > 0
}

func (f Field) Render() string {
	if f.Type == Space {
		return fmt.Sprintf("<br>")
	}
	if f.Type == Text {
		return fmt.Sprintf("<div class='mb-3'><label class='form-label'>%s<br> <input class='form-control' type='text' name='%s'  placeholder='%s' value='%s'></label></div>", f.Title, f.Name, f.Placeholder, f.Value)
	}
	if f.Type == Password {
		return fmt.Sprintf("<div class='mb-3'><label class='form-label'>%s<br> <input class='form-control' type='password' name='%s' value='%s'></label></div>", f.Title, f.Name, f.Value)
	}
	if f.Type == Number {
		return fmt.Sprintf("<div class='mb-3'><label class='form-label'>%s<br> <input class='form-control' type='number' name='%s' value='%s'></label></div>", f.Title, f.Name, f.Value)
	}
	if f.Type == Boolean {
		// if f.Value : checked
		return fmt.Sprintf("<div class='mb-3'><label class='form-label'>%s<br> <input class='form-control' type='checkbox' name='%s' value='1'></label></div>", f.Title, f.Name)
	}
	if f.Type == File {
		if f.Value != "" {
			return fmt.Sprintf("<div class='mb-3'><label class='form-label'>%s (Current: %s)<br> <input class='form-control' type='file' name='%s'></label></div>", f.Title, f.Value, f.Name)
		} else {
			return fmt.Sprintf("<div class='mb-3'><label class='form-label'>%s<br> <input class='form-control' type='file' name='%s'></label></div>", f.Title, f.Name)
		}
	}
	if f.Type == Files {
		if f.Value != "" {
			return fmt.Sprintf("<div class='mb-3'><label class='form-label'>%s (Current: %s)<br> <input class='form-control' type='file' name='%s' multiple='multiple'></label></div>", f.Title, f.Value, f.Name)
		} else {
			return fmt.Sprintf("<div class='mb-3'><label class='form-label'>%s<br> <input class='form-control' type='file' name='%s' multiple='multiple'></label></div>", f.Title, f.Name)
		}
	}
	if f.Type == Select {
		ol := ""
		for _, o := range f.Options {
			ol += fmt.Sprintf("<option value='%s'>%s</option>", o.Value, o.Text)
		}
		return fmt.Sprintf("<div class='mb-3'><label class='form-label'>%s</label><br> <select class='form-control' name='%s'>%s</select></div>", f.Title, f.Name, ol)
	}
	if f.Type == Hidden {
		return fmt.Sprintf("<input type='hidden' name='%s' value='%s'>", f.Name, f.Value)
	}
	return "Unknown field type"
}

func (f *Field) Parse(r *http.Request) {
	if f.Type == Files {
		fhs := r.MultipartForm.File[f.Name]
		for _, fh := range fhs {
			nf, err := fh.Open()
			if err != nil {
				fmt.Println(err)
				return
			}
			defer nf.Close()
			//	fmt.Println(fh.Filename)

			ff := FieldFile{}

			ff.Filename = fh.Filename

			tempFile, err := ioutil.TempFile("temp", "upload-*.png")
			if err != nil {
				fmt.Println(err)
				//	f.Value = ""
				return
			}
			defer tempFile.Close()

			// read all of the contents of our uploaded file into a
			// byte array
			fileBytes, err := ioutil.ReadAll(nf)
			if err != nil {
				fmt.Printf("Field.Parse, unable: %s\n", err)
				//	f.Value = ""
				return
			}
			// write this byte array to our temporary file
			tempFile.Write(fileBytes)

			ff.TmpFilename = tempFile.Name()

			f.Files = append(f.Files, ff)
		}

	}

	if f.Type == File {
		file, handler, err := r.FormFile(f.Name)
		if err != nil {
			//			fmt.Println("Error Retrieving the File")
			//			fmt.Println(err)
			return
		}
		defer file.Close()

		ff := FieldFile{}

		ff.Filename = handler.Filename

		//		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		//		fmt.Printf("File Size: %+v\n", handler.Size)
		//		fmt.Printf("MIME Header: %+v\n", handler.Header)

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern
		tempFile, err := ioutil.TempFile("temp", "upload-*.png")
		if err != nil {
			fmt.Println(err)
			//	f.Value = ""
			return
		}
		defer tempFile.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Printf("Field.Parse, unable: %s\n", err)
			//	f.Value = ""
			return
		}
		// write this byte array to our temporary file
		tempFile.Write(fileBytes)

		ff.TmpFilename = tempFile.Name()
		f.Files = append(f.Files, ff)

	} else {
		//		fmt.Printf(" %s %s \n", f.Name, r.FormValue(f.Name))

		f.Value = strings.TrimSpace(r.FormValue(f.Name))
	}
}

func NewField(title string, placeholder string, name string, tpe FieldType, value string) Field {
	f := Field{}
	f.Title = title
	f.Placeholder = placeholder
	f.Name = name
	f.Type = tpe
	f.Value = value
	return f
}

type Form struct {
	Name   string
	Fields []Field
	Submit string
	Action string
}

func (f Form) Render() string {
	var sb strings.Builder
	has_file := false
	for _, fld := range f.Fields {
		if fld.Type == File || fld.Type == Files {
			has_file = true
		}
	}
	if has_file {
		sb.WriteString(fmt.Sprintf("<form action='%s' class='%s' method='post' enctype='multipart/form-data'>", f.Action, strings.ToLower(f.Name)))
	} else {
		sb.WriteString(fmt.Sprintf("<form action='%s' class='%s' method='post'>", f.Action, strings.ToLower(f.Name)))
	}
	sb.WriteString(fmt.Sprintf("<h2>%s</h2>", f.Name))
	//sb.WriteString("<fieldset>")
	for _, fld := range f.Fields {
		sb.WriteString(fld.Render())
	}
	//sb.WriteString("</fieldset>")
	sb.WriteString("<div class='col-auto'><button type='submit' class='btn btn-primary mb-3'>" + f.Submit + "</button></div>")
	sb.WriteString("</form>")
	return sb.String()
}

func (f *Form) AddField(fld Field) {
	f.Fields = append(f.Fields, fld)
}

func (f *Form) GetField(name string) *Field {
	for _, fld := range f.Fields {
		if fld.Name == name {
			return &fld
		}
	}
	return nil
}

func (f *Form) Parse(r *http.Request) {
	err := r.ParseMultipartForm(1024 * 1024 * 50)
	if err == http.ErrNotMultipart {
		r.ParseForm()
	}

	for i := 0; i < len(f.Fields); i++ {
		f.Fields[i].Parse(r)
	}
}

func NewForm(name, submit, action string) Form {
	f := Form{}
	f.Name = name
	f.Submit = submit
	f.Action = action
	return f
}
