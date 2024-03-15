package alfajor

func (alfajor Alfajor) SimpleLog(typeOfLog string, errString string) {
	alfajor.Log(typeOfLog, empty, errString)
}

func (alfajor Alfajor) Log(typeOfLog string, where string, errString string) {

	var (
		logString  string
		fileName   string
		exist      error
		folderName string
	)

	//Get the currente date and time
	dateTime := getDateTime()
	time := getHourAndMinutes()

	if alfajor.AlfajorParams.SeparateAlfajor {

		folderDateTime := alfajor.AlfajorParams.RootDir + slash + dateTime
		exist = checIfExist(folderDateTime)
		if exist != nil {
			alfajor.mkdir(folderDateTime)
		}

		folderName = alfajor.AlfajorParams.RootDir + slash + dateTime + slash + typeOfLog
		exist = checIfExist(folderName)
		if exist != nil {
			alfajor.mkdir(folderName)
		}

		//File name
		fileName = folderName + slash + dateTime + extension
		exist = checIfExist(fileName)
		if exist != nil {
			alfajor.createfile(fileName)
		}

	} else {
		//File name
		fileName = alfajor.AlfajorParams.RootDir + slash + dateTime + extension
		exist = checIfExist(fileName)
		if exist != nil {
			alfajor.createfile(fileName)
		}
	}

	//Define text line
	if alfajor.AlfajorParams.PrettyAlfajor {
		logString = alfajor.defineTextPretty(typeOfLog, time, where, errString)
	} else {
		logString = alfajor.defineText(typeOfLog, time, where, errString)
	}

	//Append log
	alfajor.appendText(logString, fileName)

}

func (alfajor Alfajor) defineText(typeOfLog string, time string, where string, errString string) string {

	var logText string

	if where == empty {
		logText = time + space + or + typeOfLog + or + errString
	} else {
		logText = time + space + or + typeOfLog + or + where + or + errString
	}

	return logText

}

func (alfajor Alfajor) defineTextPretty(typeOfLog string, time string, where string, errString string) string {

	var logText string

	if where == empty {
		logText = time + space + or + space + typeOfLog + space + or + space + errString
	} else {
		logText = time + space + or + space + typeOfLog + space + or + space + where + colon + space + errString
	}

	return logText

}
