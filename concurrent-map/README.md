
example
	// Create a new map.
	m := cmap.New()

	// Sets item within map, sets "bar" under key "foo"
	m.Set("foo", "bar")

	// Retrieve item from map.
	if tmp, ok := m.Get("foo"); ok {
		bar := tmp.(string)
	}

	// Removes item under key "foo"
	m.Remove("foo")
	