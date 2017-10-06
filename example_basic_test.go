package logrus_test

import (
	"github.com/xitonix/logrus"
	"os"
)

func Example_basic() {
	var log = logrus.New(logrus.DebugLevel)
	log.SetFormatter(new(logrus.JSONFormatter))
	log.SetFormatter(&logrus.TextFormatter{
		DisableSorting:   false,
		DisableTimestamp: true,
	})
	log.SetOutput(os.Stdout)

	// file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	// if err == nil {
	// 	log.SetOutput(file)
	// } else {
	// 	log.Info("Failed to log to file, using default stderr")
	// }

	defer func() {
		err := recover()
		if err != nil {
			entry := *err.(**logrus.Entry)
			log.AsError().WithFields(logrus.Fields{
				"omg":         true,
				"err_animal":  entry.Data["animal"],
				"err_size":    entry.Data["size"],
				"err_level":   entry.Level,
				"err_message": entry.Message,
				"number":      100,
			}).Write("The ice breaks!") // or use Fatal() to force the process to exit with a nonzero code
		}
	}()

	log.AsDebug().WithFields(logrus.Fields{
		"animal": "walrus",
		"number": 8,
	}).Write("Started observing beach")

	log.AsInfo().WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Write("A group of walrus emerges from the ocean")

	log.AsWarning().WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	}).Write("The group's number increased tremendously!")

	log.AsDebug().WithFields(logrus.Fields{
		"temperature": -4,
	}).Write("Temperature changes")

	log.AsPanic().WithFields(logrus.Fields{
		"animal": "orca",
		"size":   9009,
	}).Write("It's over 9000!")

	// Output:
	// level=debug msg="Started observing beach" animal=walrus number=8
	// level=info msg="A group of walrus emerges from the ocean" animal=walrus size=10
	// level=warning msg="The group's number increased tremendously!" number=122 omg=true
	// level=debug msg="Temperature changes" temperature=-4
	// level=panic msg="It's over 9000!" animal=orca size=9009
	// level=error msg="The ice breaks!" err_animal=orca err_level=panic err_message="It's over 9000!" err_size=9009 number=100 omg=true
}
