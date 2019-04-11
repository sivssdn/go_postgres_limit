package main

import(
  "fmt"
  "os"
  "sync"
  "strconv"

  "database/sql"
  _ "github.com/lib/pq"
  "github.com/joho/godotenv"
)

func main(){
  err := godotenv.Load()
  checkErr(err)

  fmt.Println("Connection established")

  var wg sync.WaitGroup
  maxConnections, err := strconv.Atoi(os.Getenv("MAX_CONNECTIONS"))
  checkErr(err)
  wg.Add(maxConnections)
  maxConnectionNumber := 0

  for i:=0; i<maxConnections; i++{
    go func(connection_number int){

        dbCreds := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DB_NAME"))
        db, err := sql.Open("postgres",dbCreds)
        checkErr(err)

        defer func(db *sql.DB) {
                defer wg.Done()
                //defer db.Close()

                if maxConnectionNumber < connection_number{
                  maxConnectionNumber = connection_number
                }

                if err := recover(); err != nil {
                    fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
                    fmt.Println(connection_number)

                    fmt.Println("Total Connections : ", i)
                    os.Exit(1)
                }
        }(db)

    }(i)
  }

  wg.Wait()
  fmt.Println("\n Program Execution Finished with Max Connection Number : ", maxConnectionNumber)
}

func checkErr(err error){
  if err != nil{
    panic(err)
  }
}
