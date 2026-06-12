package battery

import (
	"fmt"
	"time"
)

// FormatDuration converts seconds into a human-readable string "XXh XXm XXs"
func FormatDuration(sec float64) string {
	d := time.Duration(sec * float64(time.Second))
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02dh %02dm %02ds", h, m, s)
}
