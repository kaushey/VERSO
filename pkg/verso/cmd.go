package verso

var Commands = map[string]func([]string) error{
	"init":     InitHandler,
	"commit":   CommitHandler,
	"add":      AddHandler,
	"cat-file": CatFileHandler,
	"log":      LogHandler,
	"status":   StatusHandler,
}

func Usage() string {
	s := "Usage: verso [command] [options]\nAvailable commands:\n"
	for k := range Commands {
		s += " - " + k + "\n"
	}
	return s
}
