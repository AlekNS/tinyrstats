package endpoints

// wrapResponseToData wraps response to sub data object.
// @NOTE: Use only like for json response.
func wrapResponseToData(result interface{}, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"data": result,
	}, nil
}
