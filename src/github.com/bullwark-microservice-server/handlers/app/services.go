package app

import (
    "github.com/bullwark-microservice-server/common/cli"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "path/filepath"

    "github.com/bullwark-microservice-server/common"
    "github.com/bullwark-microservice-server/serviceproviders"
)

type UploadResponse struct {
    Filepath string `json:"filepath"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}

const (
    FileNameLength = 64
    MaxUploadSize  = 1024 * 1024 * 1024 * 10 // 10 GB
    BinaryDirPerm  = 0755
)

var AllowedFileTypes = [...]string{
    "application/java-archive",
    "application/octet-stream",
}

var serviceProvider serviceproviders.ServiceProvider

func CreateService(w http.ResponseWriter, r *http.Request) {

    var registrationRequest serviceproviders.Request
    var serviceId = common.RandString(FileNameLength)
    var encoder = json.NewEncoder(w)

    // TODO: Replace with DI
    //
    if serviceProvider == nil {
        serviceProvider = serviceproviders.NewDockerServiceProvider()
    }

    // Parse body
    //
    err := json.NewDecoder(r.Body).Decode(&registrationRequest)

    if err != nil {
        fmt.Println(err)
        _ = encoder.Encode(serviceproviders.Response{
            Status: "error",
            Reason: fmt.Sprintf("%s", err),
        })
        return
    }

    // Register with service registration impl
    //
    err = serviceProvider.Create(registrationRequest)

    if err != nil {
        fmt.Println(err)
        _ = encoder.Encode(serviceproviders.Response{
            Status: "error",
            Reason: fmt.Sprintf("%s", err),
        })
        return
    }

    // Success
    //
    _ = encoder.Encode(serviceproviders.Response{
        ServiceId: serviceId,
        Status:    "success",
    })
    return
}

func UploadService(w http.ResponseWriter, r *http.Request) {

    var encoder = json.NewEncoder(w)
    var coreConfig = common.Configs[cli.CoreConfigType].(cli.CoreConfig)

    r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

    // Check file size against max upload size
    //
    if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
        fmt.Println(err)
        writeError(w, "File too large.", http.StatusBadRequest)
        return
    }

    // Get file type
    //
    file, header, err := r.FormFile("file")
    if err != nil {
        fmt.Println(err)
        writeError(w, "Invalid file.", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Read all file bytes
    //
    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Println(err)
        writeError(w, "Invalid file.", http.StatusBadRequest)
        return
    }

    // Validate file type
    //
    filetype := http.DetectContentType(fileBytes)
    for _, b := range AllowedFileTypes {
        if b == filetype {
            writeError(w, "Invalid file type.", http.StatusBadRequest)
            return
        }
    }

    // Generate new name
    //
    fileName := common.RandString(FileNameLength)
    fileExt := filepath.Ext(header.Filename)

    newPath := filepath.Join(*coreConfig.BinSavePath, fileExt, fileName + fileExt)
    fmt.Printf("FileType: %s, File: %s\n", fileExt, newPath)

    // Create save directory if it doesn't already exist
    //
    if _, err := os.Stat(filepath.Join(*coreConfig.BinSavePath, fileExt)); os.IsNotExist(err) {
        _ = os.MkdirAll(filepath.Join(*coreConfig.BinSavePath, fileExt), BinaryDirPerm)
    }

    // Write to system
    //
    newFile, err := os.Create(newPath)
    if err != nil {
        fmt.Println(err)
        writeError(w, "Can't write file.", http.StatusInternalServerError)
        return
    }
    defer newFile.Close()
    if _, err := newFile.Write(fileBytes); err != nil {
        fmt.Println(err)
        writeError(w, "Can't write file.", http.StatusInternalServerError)
        return
    }

    // Success
    //
    _ = encoder.Encode(UploadResponse{
        Filepath: fileName + fileExt,
    })
    return
}

func writeError(w http.ResponseWriter, error string, code int) {
    errorMessage, _ := json.Marshal(ErrorResponse{
        Error: error,
    })
    http.Error(w, string(errorMessage), code)
}
