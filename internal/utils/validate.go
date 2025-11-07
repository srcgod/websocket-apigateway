package utils

/*
func ValidateGetPostsParams(c *gin.Context) (int32, int32, int32 error) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit, err := ConvertToInt32(limitStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid limit: %w", err)
	}
	if limit == 0 {
		limit = 100 // default
	}
ё
	offset, err := ConvertToInt32(offsetStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid offset: %w", err)
	}

	return limit, offset, nil
}*/
