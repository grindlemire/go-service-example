# log
a simple wrapper for [uber-go/zap](https://github.com/uber-go/zap) that simplifies configuration and handles log rotation


# Why?
Because I often need a logger that is extremely simple to hook up and just use. However I also want callers, logfile configuration, and structured logs. Zap is a good choice except I also want to have log rotation which zap doesn't support natively.

This package is intended to be a very lightweight wrapper on top of zap while greatly simplifying the configuration, providing capbaility to easily produce colored logs or json logs, and seamlessly integrate with rotating log files.

# Usage:
```Golang
func main() {
	// configuration is compatible with env loading or json if you want to store it in a file
	// I also provide log.Default for convenience

	// Simply call Init with the configuration and all the files you want to log to (if any)
	err := log.Init(log.Default, "./test.log")
	if err != nil {
		log.Fatalf("unable to initialize logger: %v", err)
	}

	log.Debug("no formatting")
	log.Debugf("formatting [%s]", "here")
	log.Debugw("extra fields", log.Fields{"like": "this", "or_numbers": 1})

	log.Info("no formatting")
	log.Infof("formatting [%s]", "here")
	log.Infow("extra fields", log.Fields{"like": "this", "or_numbers": 1})

	log.Warn("no formatting")
	log.Warnf("formatting [%s]", "here")
	log.Warnw("extra fields", log.Fields{"like": "this", "or_numbers": 1})

	log.Error("no formatting")
	log.Errorf("formatting [%s]", "here")
	log.Errorw("extra fields", log.Fields{"like": "this", "or_numbers": 1})

	log.Fatal("All done")
}
```
