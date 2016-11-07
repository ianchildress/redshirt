# redshirt - A signal handling package for Go

*This is a quick and dirty readme and thus can't guarantee 100% accuracy.  I'll polish it up shortly.*

redshirt is a package for handling signals and/or dying gracefully.
The package is used by registering a function with one or more signals
When the program receives that signal, the function is run.

Example: If you want to reload your config when a SIGHUP signal
is detected, you can import the redshirt package and do the following:
```
type Config struct {
  file string
} 

// The Signal method fulfills the interface requirement. Add a method
// to your type with the following signature.
func Signal(sig os.Signal) error {
	// This is my callback function. In here I can do
	// things to reload my config file.
	return nil
}

// Pass in the interface along with the signals to the Register function.
redshirt.Register(conf,redshirt.SIGHUP) 

// When a SIGHUP is sent to the program, the conf.Signal function will be run.
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
redshirt.Register(foo,redshirt.SIGHUP)
redshirt.Register(bar,redshirt.SIGHUP)
redshirt.Register(boo,redshirt.SIGHUP)
```
