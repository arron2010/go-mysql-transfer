package glogger







func Debugf(template string, args ...interface{}) {

	logger.Debugf(template, args...)
}


func Infof(template string, args ...interface{}) {

	logger.Infof(template, args...)
}


func Warnf(template string, args ...interface{}) {

	logger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {

	logger.Errorf(template, args...)
}
