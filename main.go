package main

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"nmbr.one/sql-backups/discord"
	"nmbr.one/sql-backups/files"
)

func prepareBackupFile(backupFileName string) (*os.File, error) {
	var file, openError = os.Create(backupFileName)

	if openError != nil {
		return nil, openError
	}

	return file, nil
}

func startBackup(sqlHost string, sqlDB string, sqlUser string, sqlPassword string, file *os.File) error {
	var command = exec.Command("mysqldump", "-h", sqlHost, "-u", sqlUser, "--password="+sqlPassword, sqlDB)

	command.Stdout = file

	var runError = command.Run()

	if runError != nil {
		return runError
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
func main() {
	var config = files.NewPropertiesFile("config.properties")

	if config.Exists() {
		config.Load()
	} else {
		println("No config.properties file found.")
		os.Exit(105)
	}

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
				var backupFile, fileError = prepareBackupFile(sqlDB + ".sql")

				if fileError != nil {
					println("Can't open file " + sqlDB + ".sql.")
					os.Exit(101)
				}

				var backupError = startBackup(sqlHost, sqlDB, sqlUser, sqlPassword, backupFile)

				backupFile.Close()

				if backupError != nil {
					println("Can't find mysqldump executable. Please review your mysql installation.")
					os.Exit(102)
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
