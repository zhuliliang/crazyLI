package db

import "errors"

var (
    // ErrNotFound signals no rows.
    ErrNotFound = errors.New("db: not found")
)
