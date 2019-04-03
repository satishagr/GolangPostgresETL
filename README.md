# GolangPostgresETL

This simple application can be used to process large files (over 500 Mb) and load the data into a Postgres database using a Golang API. Entire ETL process has been implemented with additional featured such as pagination, search filter and average calculation on subsets. 

![Go-Postgres](https://cdn-images-1.medium.com/max/870/1*pboqJIblWykgo2F3nnfS1g.png)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

#### 1) Setup Golang on your system.
Please follow the below link to setup Golang on your OS and make sure that all the enviroment paths(GO_PATH, GO_ROOT) have been properly setup.

https://golang.org/doc/install


#### 2) Download Postgres and setup a Postgres server.

https://www.postgresql.org/download/


#### 3) Download and extract the "GolangPostgresETL-master.zip" inside the $GO_PATH/src/ and rename the folder to something easier such as GoETL(will be referred later as well).

https://github.com/satishagr/GolangPostgresETL/archive/master.zip


#### 4) Download the 'Building Footprints Dataset' as 'building.csv' from the below link and move this file to $GO_PATH/src/GoETL/

https://data.cityofnewyork.us/api/views/mtik-6c5q/rows.csv?accessType=DOWNLOAD



### Installing

#### 1) Open PgAdmin app and setup a Postgres server with a database named 'buildings'(for easier implementation since same has been used in the code). 
Note : Keep the Postgres username and password handy as this info will have to be updated in the main.go file under GoETL directory.

#### 2) Copy and Run all the commands mentioned in the db_commands.txt file inside GoETL directory on the Postgres buildings database, which will create a table named 'building'.

#### 3) Open a cmd/bash terminal and move to the src/GoETL dir. Run the below command to start the application.

C:\Users\User\Go\src\GoETL> go run main.go

#### 4) Open a browser and open the 'http://localhost:8080/' page. This should render a page like given below.

![Sample_Ui](https://github.com/satishagr/GolangPostgresETL/blob/master/misc/Sample_Ui.png)

The API should also load the records in the database like shown below.

![Sample_Db](https://github.com/satishagr/GolangPostgresETL/blob/master/misc/Sample_Db.png)

Note : The port (currently '8080') can re-configured using the main.go file.

Hence, a ETL Process has been successfully implemented using Golang and Postgres. The average of each numeric field is displayed at the last row of the table. Search box can be used to filter records based on keywords from any of the fields provided. Pagination has been implemented for the ease of the user.

## Happy Coding :)

## Authors

* **Satish Agrawal** - *Initial work* - [satishagr](https://github.com/satishagr)


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
