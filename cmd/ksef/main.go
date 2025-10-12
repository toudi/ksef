package main

import (
	"fmt"
	"ksef/cmd/ksef/commands"
	"ksef/internal/logging"
	"os"
)

// var command *commands.Command

// var loggingOutput string = ""
// var configPath string = ""

func main() {
	defer logging.FinishLogging()
	if err := commands.RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// var err error

	// flag.Usage = func() {
	// 	fmt.Printf("Użycie programu: ksef [-c] [-log] [-v] [komenda]\n\n")
	// 	flag.PrintDefaults()
	// 	fmt.Printf("\nDostępne komendy:\n")
	// 	for _, command := range commands.Registry {
	// 		fmt.Printf("%-*s - %s\n", commands.MaxCommandName, command.Name, command.Description)
	// 	}
	// }

	// flag.StringVar(
	// 	&loggingOutput,
	// 	"log",
	// 	loggingOutput,
	// 	"wyjście logowania. Wartość `-` oznacza wyjście standardowe (stdout)",
	// )
	// flag.StringVar(&configPath, "c", configPath, "ścieżka pliku konfiguracyjnego")
	// flag.BoolVar(
	// 	&logging.Verbose,
	// 	"v",
	// 	false,
	// 	"tryb verbose - przełącza wszystkie loggery w tryb debug",
	// )

	// flag.Parse()

	// args := flag.Args()

	// if len(args) < 1 {
	// 	flag.Usage()
	// 	return
	// }

	// if configPath != "" {
	// 	if err = config.ReadConfig(configPath); err != nil {
	// 		fmt.Printf("[ ERROR ] %v\n", err)
	// 		return
	// 	}
	// }

	// if err = logging.InitLogging(loggingOutput); err != nil {
	// 	fmt.Printf("[ ERROR ] Błąd inicjalizacji logowania: %v", err)
	// 	return
	// }

	// defer logging.FinishLogging()

	// logging.SeiLogger.Info("start programu")
	// defer logging.SeiLogger.Info("koniec programu")

	// command = commands.Registry.GetByName(args[0])
	// if command == nil {
	// 	fmt.Printf("[ ERROR ] Nieznana komenda\n")
	// 	flag.Usage()
	// 	return
	// }

	// if err = command.FlagSet.Parse(args[1:]); err != nil {
	// 	return
	// }

	// command.Context = context.Background()

	// err = command.Run(command)
	// if err != nil {
	// 	fmt.Printf("błąd wykonania %s:\n%s\n", args[0], err)
	// 	return
	// }
}
