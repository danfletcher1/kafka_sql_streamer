package main

import (
	"fmt"
	"time"

	KSDriver "./KSDriver"
)

func main() {
	var err error

	err = KSDriver.Connect("root", "password", "mydb", "localhost", "localhost:9092")
	if err != nil {
		fmt.Println(err)
	}

	defer KSDriver.Close()

	err = KSDriver.Update("CREATE TABLE test (id INT AUTO_INCREMENT, name VARCHAR(100), count INT DEFAULT 0, PRIMARY KEY (id))")
	if err != nil {
		fmt.Println(err)
	}
	err = KSDriver.Update("INSERT INTO test (name) VALUES ('dan')")
	if err != nil {
		fmt.Println(err)
	}
	err = KSDriver.Update("UPDATE test SET count=count+1 WHERE name='dan'")
	if err != nil {
		fmt.Println(err)
	}

	mydata, err := KSDriver.Query("SELECT * from test")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Query Immediately")
	// display any returned results, you should be having any
	for i := range mydata {
		for k, v := range mydata[i] {
			fmt.Println("Returned " + k + " : " + v)
		}
	}

	fmt.Println("Wait 3 seconds")
	time.Sleep(3 * time.Second)

	fmt.Println("Query again")

	mydata, err = KSDriver.Query("SELECT * from test")
	if err != nil {
		fmt.Println(err)
	}

	// display any returned results
	for i := range mydata {
		for k, v := range mydata[i] {
			fmt.Println("Returned " + k + " : " + v)
		}
	}
}
