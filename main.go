package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"

	"nmbr.one/sql-backups/discord"
	"nmbr.one/sql-backups/files"
)

func prepareBackupFile(backupFileName string) *os.File {
	var file, openError = os.Create(backupFileName)

	if openError != nil {
		println("Can't open file '" + backupFileName + "'. You should check your filesystem permissions.")
		os.Exit(101)
	}

	return file
}

func closeBackupFile(backupFile *os.File) {
	var closeError = backupFile.Close()

	if closeError != nil {
		println("Can't close file '" + backupFile.Name() + "'. You should check your filesystem permissions and disk usage.")
		os.Exit(106)
	}
}

func testMysqldump(mysqldumpPath string) {
	var command = exec.Command(mysqldumpPath, "--version")

	var runError = command.Run()

	if runError != nil {
		println("Can't find mysqldump at path '" + mysqldumpPath + "'. You should check the mysqldump-path option in config.")
		os.Exit(102)
	}
}

func startBackup(mysqldumpPath string, sqlHost string, sqlDB string, sqlUser string, sqlPassword string, file *os.File) error {
	var command *exec.Cmd

	if sqlPassword == "" {
		command = exec.Command(mysqldumpPath, "-h", sqlHost, "-u", sqlUser, sqlDB)
	} else {
		command = exec.Command(mysqldumpPath, "-h", sqlHost, "-u", sqlUser, "--password="+sqlPassword, sqlDB)
	}

	var errorMessageBuffer = &bytes.Buffer{}

	command.Stdout = file
	command.Stderr = errorMessageBuffer

	var runError = command.Run()

	if runError != nil {
		var errorMessage = errorMessageBuffer.String()

		return errors.New(errorMessage)
	}

	return nil
}

func parseDatabases(configEntry string) []string {
	return strings.Split(configEntry, ",")
}

// EXIT STATUS DESCRIPTION
// 100 = Invalid backup-interval
// 101 = Can not open sql file
// 102 = Can not find mysqldump executable
// 104 = Invalid initial-backup-time
// 105 = No config.properties file found
// 106 = Can not close sql file
func main() {
	var config = files.NewPropertiesFile("config.properties")

	if config.Exists() {
		config.Load()
	} else {
		println("No config.properties file found.")
		os.Exit(105)
	}

	var mysqldumpPath = config.Get("mysqldump-path")

	testMysqldump(mysqldumpPath)

	var logger = discord.Logger{
		WebhookURL: config.Get("webhook-url"),
	}

	var backupInterval, parseDurationError = time.ParseDuration(config.Get("backup-interval"))

	if parseDurationError != nil {
		println("Could not parse duration of backup-interval from config. Example duration: 24h")
		os.Exit(100)
	}

	var initialBackupTime, parseTimeError = time.Parse(time.TimeOnly, config.Get("initial-backup-time"))

	if parseTimeError != nil {
		println("Could not parse time of initial-backup-time from config. Example time: 06:00:00")
		os.Exit(104)
	}

	var now = time.Now()
	var nextBackup = time.Date(now.Year(), now.Month(), now.Day(),
		initialBackupTime.Hour(), initialBackupTime.Minute(), initialBackupTime.Second(),
		0, now.Location())

	for now.After(nextBackup) {
		nextBackup = nextBackup.Add(backupInterval)
	}

	println("First backup will run at", nextBackup.Format(time.RFC1123))

	var sqlHost = config.Get("host")
	var sqlUser = config.Get("user")
	var sqlPassword = config.Get("password")
	var sqlDatabases = parseDatabases(config.Get("databases"))

	for {
		if time.Now().After(nextBackup) {
			println("Backup Time")

			for _, sqlDB := range sqlDatabases {
				var backupFile = prepareBackupFile(sqlDB + ".sql")
				var backupError = startBackup(mysqldumpPath, sqlHost, sqlDB, sqlUser, sqlPassword, backupFile)

				closeBackupFile(backupFile)

				if backupError != nil {
					println("Backup for database", sqlDB, "failed due to wrong credentials provided by the host, user and password options from config.")
					println("Error: " + backupError.Error())

					logger.Log(discord.Embed{
						Title: "SQL-Backup wurde fehlerhaft abgeschlossen",
						Color: discord.ColorRed,
						Fields: []discord.Field{
							{
								Name:   "Datenbank",
								Value:  sqlDB,
								Inline: false,
							},
							{
								Name:   "Fehlerbeschreibung",
								Value:  backupError.Error(),
								Inline: false,
							},
						},
					})
				} else {
					logger.LogFile(discord.Embed{
						Title: "SQL-Backup erfolgreich abgeschlossen",
						Color: discord.ColorGreen,
						Fields: []discord.Field{
							{
								Name:   "Datenbank",
								Value:  sqlDB,
								Inline: true,
							},
							{
								Name:   "Erstellungszeitpunkt",
								Value:  time.Now().Format(time.RFC1123),
								Inline: true,
							},
						},
					}, sqlDB+".sql")
				}

				os.Remove(sqlDB + ".sql")
			}

			nextBackup = nextBackup.Add(backupInterval)

			println("Next backup will run at", nextBackup.Format(time.RFC1123))
		}

		time.Sleep(time.Millisecond * 10)
	}
}
