```
                                    ,--,    
         ,---._                  ,---.'|    
       .-- -.' \    ,---,   ,---,|   | :    
       |    |   :,`--.' |,`--.' |:   : |    
       :    ;   ||   :  :|   :  :|   ' :    
       :        |:   |  ':   |  ';   ; '    
       |    :   :|   :  ||   :  |'   | |__  
       :         '   '  ;'   '  ;|   | :.'| 
       |    ;   ||   |  ||   |  |'   :    ; 
   ___ |         '   :  ;'   :  ;|   |  ./  
 /    /\    :   :|   |  '|   |  ';   : ;    
/  ../  `..-    ,'   :  |'   :  ||   ,/     
\    \         ; ;   |.' ;   |.' '---'      
 \    \      ,'  '---'   '---'              
  "---....--'                               
```
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=flat-square&logo=go&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?logo=TypeScript&logoColor=FFF&style=flat-square)
![Rust](https://img.shields.io/badge/Rust-%23000000.svg?style=flat-square&logo=rust&logoColor=white)

# Jcloud ☁️

JcloudFile is a client-server application for cloud filePath storage. This project provides a filePath storage system with a backend implemented in Go and a frontend using TypeScript. The backend uses SQLite3 to manage user files and offers a basic API for interacting with stored filePath data.
## 📇 Table of Content

* [📖 Main Functionalities](#-main-functionalities-)
* [🛠️ How it works?](#-how-it-works)
    * [🌟 Overview](#-overview-)
    * [📝 Detailed Steps](#-detailed-steps-)
    * [📄 Code Explanation](#-code-explanation-)


## 📖 Main Functionalities
- [x] Upload, list, delete, edit files
- [x] Unlimited space to keep files
- [x] Opportunity to download files back on local device
- [x] Usage of Sqlite3 database to store files on server
- [x] Rust Desktop App
- [x] TypeScript Web Interface
- [x] Pure Go CLI
- [x] Txt formatted logs with levels


## 🛠️ How it works?

### 🌟 Overview

Jcloud offers multiple ways to interact with its file storage system, providing flexibility to users based on their preferences and needs. You can choose from the following options:

1. **TypeScript Web App** 🌐:
    - **Interactive Interface**: The web application, built with TypeScript, offers a user-friendly graphical interface for managing files.
    - **Features**: Users can upload files, view a list of their uploaded files, and perform actions such as deleting files. The web app communicates with the backend through RESTful API calls, providing real-time updates and a seamless experience.
    - **Technology**: Utilizes modern web technologies to ensure a responsive and engaging user experience.

2. **Rust Application** 🦀:
    - **Desktop or CLI Interface**: The Rust application offers a robust client that can be used as a desktop application or a command-line interface.
    - **Features**: Allows users to perform file operations such as uploading, listing, and deleting files directly from their local environment. The Rust client interacts with the backend API, providing a native experience.
    - **Technology**: Built with Rust for performance and safety, ensuring efficient handling of file operations.

3. **Go CLI** 🛠️:
    - **Command-Line Tool**: The Go CLI provides a powerful command-line interface for interacting with JcloudFile.
    - **Features**: Users can manage files using CLI commands, such as uploading files, listing files, and deleting files. The CLI tool connects to the backend through RESTful API calls, making it a versatile option for advanced users.
    - **Technology**: Implemented in Go, offering a streamlined command-line experience.

### 📝 Detailed Steps
1. **File Upload**:
    - **Process**: When a file is uploaded, regardless of the client interface, an API request is sent to the backend. The backend processes the file and stores it in the SQLite3 database, updating the metadata and file paths accordingly.

2. **Backend Processing**:
    - **Management**: The backend, written in Go, handles all file storage and retrieval operations. It ensures that file data is managed correctly and efficiently, processing requests and interacting with the SQLite3 database.

3. **Database Management**:
    - **Storage**: SQLite3 is utilized to store comprehensive information about each file, including its filename, size, and storage path. This database maintains the integrity and accessibility of file metadata.

4. **API Communication**:
    - **Integration**: RESTful API endpoints are used for communication between the web app, Rust application, Go CLI, and the backend. These endpoints ensure that file operations are consistently handled across all client interfaces.

### 📄 Code Explanation

#### Go CLI written by Cobra-Cli

```go
package main

import (
   "bufio"
   "context"
   "fmt"
   "github.com/JIIL07/jcloud/internal/client/app"
   "github.com/JIIL07/jcloud/internal/client/cmd"
   "github.com/JIIL07/jcloud/internal/client/config"
   "github.com/JIIL07/jcloud/pkg/ctx"
   "log"
   "os"
   "strings"
)

func main() {
   c := config.MustLoad()
   appc, err := app.NewAppContext(c)
   if err != nil {
      log.Fatal(err)
   }
   ctx := jctx.WithContext(context.Background(), "app-context", appc)
   cmd.SetContext(ctx)
   switch {
   case c.Env == "prod":
      cmd.Execute()
   case c.Env == "debug" || c.Env == "local":
      reader := bufio.NewReader(os.Stdin)
      for {
         dir, _ := os.Getwd()
         fmt.Printf("%v>", dir)

         input, _ := reader.ReadString('\n')
         input = strings.TrimSpace(input)

         args := strings.Split(input, " ")
         cmd.RootCmd.SetArgs(args[1:])

         cmd.Execute()
      }
   }
}

```

#### Go Server 

```go
package main

import (
   "context"
   "github.com/JIIL07/jcloud/internal/config"
   "github.com/JIIL07/jcloud/internal/logger"
   "github.com/JIIL07/jcloud/internal/server"
   "github.com/JIIL07/jcloud/internal/storage"
   "github.com/JIIL07/jcloud/pkg/cookies"
   "github.com/JIIL07/jcloud/pkg/env"
   "github.com/JIIL07/jcloud/pkg/log"
   "os"
   "os/signal"
   "syscall"
   "time"
)

func main() {
   jenv.LoadEnv()
   cfg := config.MustLoad()
   log := logger.NewLogger(cfg.Env)
   s, err := storage.InitDatabase(cfg)
   if err != nil {
      log.Error("Failed to initialize database", jlog.Err(err))
      os.Exit(1)
   }
   defer s.CloseDatabase()
   cookies.SetNewCookieStore()
   srv := server.New(cfg.Server, s)
   go func() {
      log.Info("Server starting on port :8080")
      if err := srv.Start(); err != nil {
         log.Error("Server failed to start", jlog.Err(err))
         os.Exit(1)
      }
   }()
   c := make(chan os.Signal, 1)
   signal.Notify(c, os.Interrupt, syscall.SIGTERM)
   <-c
   ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
   defer cancel()
   if err := srv.Stop(ctx); err != nil {
      log.Error("Server shutdown failed", jlog.Err(err))
      os.Exit(1)
   }
   log.Info("Server gracefully stopped")
}

```

