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
# File Storage System in Go using SQLite3

This project is a file storage system that uses SQLite3 as a database to manage user files. Implemented in the Go programming language, this system provides a basic API for interacting with file data stored in the database.

## Features

- **Table Creation and Initialization**: The program automatically creates tables in the SQLite3 database for storing files.
- **File Adding**: Users can add files to the database through the API.
- **File Retrieval**: Users can retrieve and download their files on demand.
- **File Deletion**: Functionality is provided for deleting files from the database.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

You will need the following installed on your system:

- Go (version 1.14 or higher)
- SQLite3

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/JIIL07/sql.git
   cd sql

2. **Build**
    ```bash
    go build -o file.exe #You can use any other name for .exe file
    