#logrotate file

`logrotated` can be configured to send a `SIGHUP` signal to a process after rotating it's logs.  This library reopens the underlying `os.File` when a `SIGHUP` is received by the app.  

###Example
This is will enable all log calls to output to the log file without interruption when `logrotated` rotates the file.

	logfile, err := logrotate.NewFile("/log/path/here")
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logfile)


ref: [http://linux.die.net/man/8/logrotate](http://linux.die.net/man/8/logrotate)