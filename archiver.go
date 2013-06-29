// archiver.go
package neat

import (
    "encoding/xml"
    "errors"
	"io/ioutil"
    "os"
    "path"
    "strconv"
    "strings"
)

type Archiver interface {
    Archive(experiment *Experiment) (err error)
    Restore() (experiment *Experiment, err error)
}

func NewArchiver(path, prefix string) (archiver Archiver, err error) {

    // Verify the path exists and is a directory
    var info os.FileInfo
    info, err = os.Stat(path)
    if os.IsNotExist(err) {
        return
    }

    if !info.IsDir() {
        err = errors.New("FileArchive path must be a directory")
    }

    // Return the new Archive
    archiver = &xmlArchiver{path, prefix}
    return

}

// Default archive type. Stores as an XML file
type xmlArchiver struct {
    path   string
    prefix string
}

// Archives the experiment to an XML file
func (archiver xmlArchiver) Archive(experiment *Experiment) (err error) {

    // Create a new encoder
    var file *os.File
    file, err = os.Create(path.Join(archiver.path,
        strings.Join([]string{archiver.prefix, "-",
            strconv.FormatUint(uint64(experiment.Current.Generation), 10),
            ".xml"}, "")))
    if err != nil {
        return
    }
    encoder := xml.NewEncoder(file)

    // Archive the Experiment
    encoder.Encode(experiment)

    // Close the archive
    encoder.Flush()
    file.Close()
    return
}

// Restores the latest version of an experiment from the archive
func (archiver xmlArchiver) Restore() (experiment *Experiment, err error) {

    // Load the most recent generation from the archive path. ReadDir returns
    // the files sorted by name.
    var files []os.FileInfo
    files, err = ioutil.ReadDir(archiver.path)
    if err != nil {
        return
    }

    if len(files) == 0 {
        err = errors.New("There are no archive files to load.")
        return
    }

    current := files[len(files)-1] // This should be the latest file

    // Load the experiement
    r, err := os.Open(current.Name())
    if err != nil {
        return
    }
    defer r.Close()

    xml := xml.NewDecoder(r)

    experiment = &Experiment{}
    err = xml.Decode(experiment)

    // Return the restored experiment
    return

}
