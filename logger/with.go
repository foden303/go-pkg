package logger

func (i *impl) With(tags map[string]interface{}) Logger {
	if len(tags) == 0 {
		return i
	}

	newTags := normalizeTags(tags)
	logger := i.logger
	return &impl{
		logger: logger.With(zapFields(newTags)...),
	}
}
