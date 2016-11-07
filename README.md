# redshirt - A signal handling package for Go

redshirt is a package for handling signals and/or dying gracefully
The package is used by registering a function with one or more signals
When the program receives that signal, the function is run.

Example: If you want to reload your config when a SIGHUP signal
is detected, you can import the redshirt package and do the following:
```
redshirt.Register(myConfigReloadFunction(),redshirt.SIGHUP)
```

redshirt also has anonymous function support by using RegisterFunc:
```
// Register an anonymous function to handle signals
redshirt.RegisterFunc(func(sig os.Signal) error {
	fmt.Println("SIGINT detected.")
	return nil
}, signal.SIGINT)())
```

Multiple functions can be registered to the same signal. The functions
will be executed in the order they were registered.
For example:
```
redshirt.Register(ReloadConfig(),redshirt.SIGHUP)
redshirt.Register(SendEmail(),redshirt.SIGHUP)
redshirt.Register(LogToFile(),redshirt.SIGHUP)
```
