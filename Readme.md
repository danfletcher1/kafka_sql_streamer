	Name: Dan Fletcher
	Date: 26/06/18
	Title: Kafka_sql_streamer
	Codebase: https://github.com/danfletcher1/kafka_sql_streamer
	Description:

	Currently this is DEVEOPMENT status
    
    I learned an interesting use of Kafka, used in some of the big companies, is as a DB transactional buffer. They use kafka to absorb the DB updates, and the DB to apply updates in a controlled way, which may not be immediately.
    Where you are able to live with eventual consistance, and you are able to write your DB updates in transational way, that does not depend on any immediate return like an insert record number. This style is really the norm where possible, and you should consider carefully how you write your SQL with this method. 

    Depending on your chosen DB, MySQL supports multimaster to a point, and unlimited replicas.  Postgres has a single master, the use of multiple replicas are common, but does have the advantage of delayed replica, and point in time recovery. 
    In most systems we make many more DB queries than updates, it makes allot of sense to send your DB update statements through a different work queue to your queries, and send the work of queries to the almost unlimited resources of site local replicas. Increasingly important in widely distributed networks.

    By using Kafka to absorb the updates, and then spooling to the DB masters at a rate they can ingest you have many advantages:-
	(I'm showing whats possible it doesn't do this yet)

	You can support bursty data
	You can buffer updates, and read from replicas whilst you resolve an issue, your don't need to be 'down'
	You can freeze, delay or replay your transaction stream, delayed and replay streams are extremely useful should some unforcean incident occure.
	You can delete an update from a delayed stream should you freeze in time.
	You can alter transactions midflow, say you want to drop risky transactions.
	You can log/alert/report on certain transactions or transacton volumes/aggregations, without loading your db more.
	You can do point in time recovery from a backup then stream update until the point an issue occures, and continue omiting by the problem update.
	You can have multiple masters now, each applying the same updates in the same order, but they don't need to be on the same update at any given time, eventual consistancy.
	You can freeze a master, do an version change, and continue applying updates from where you left off.
	You ccan apply the same updates to mysql/postgres/oracle for example
	You can have a customer version and a development version of your database, both receiving and applying the same live data.
	Your code does not need to know anything about the underline DB structure. 



	Structure
	---------
	The first driver will send all updates (any SQL statement) to a Kafka queue, a 2nd 'service' will read the Kafka queue and apply the sql statements to the DB Master. The driver will do queries direct from an SQL slave.
	The driver will be incorporated into your code. Servie will run as a program on your SQL master or if you need the resources a machine next to your SQL master, or multiple masters, thats the beauty.   

	Reading :- 
	This will return any SQL query you want from the SQL slave. Whilst you can do a statement that updates something, don't, if you are using read only slaves, at best it will be ignored, worst it will error. Its may be easier to use stored procedures for everything.

	Returned is always a map[int]map[string]string each returned row having an int, and each containing key/value pairs for every column

	Writing:-
	This will return nothing, you are advised to do updates in a single transactional way or a stored procedure is best. By transational I mean don't update table 1 with a name and in another statement update the address, make it a single statement that either completes in full or fails completely. You really don't want something half updated. Using stored procedures is the best, its easier and more reliable than complex statements.  



	USAGE:
	// Always connect first
	sql.Connect(dbuser, dbpass, db, dbhost, kafkahost)
	defer sql.Close()


	// Currently Old style Query
	// This risks SQL injection attacks but allows you to write any SQL statement
	// All items are returned as strings (thats the driver not me)

	err = KSDriver.Update("UPDATE test SET count=count+1 WHERE name='dan'")
	if err != nil {
		fmt.Println(err)
	}
	

	mydata, err := KSDriver.Query("SELECT * from test")
	if err != nil {
		fmt.Println(err)
	}

	// display any returned results
	for i := range mydata {
		for k, v := range mydata[i] {
			fmt.Println("Returned " + k + " : " + v)
		}
	}