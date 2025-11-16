// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CronSchedule represents a parsed cron expression
type CronSchedule struct {
	Minutes  []int
	Hours    []int
	Days     []int
	Months   []int
	Weekdays []int
}

// ParseCronExpression parses a standard cron expression (5 fields)
// Format: minute hour day month weekday
// Examples:
//   "0 0 * * *"     - Every day at midnight
//   "*/5 * * * *"   - Every 5 minutes
//   "0 2 * * 1"     - Every Monday at 2 AM
//   "0 0 1 * *"     - First day of every month at midnight
func ParseCronExpression(expr string) (*CronSchedule, error) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return nil, fmt.Errorf("invalid cron expression: expected 5 fields, got %d", len(fields))
	}

	schedule := &CronSchedule{}
	var err error

	// Parse minute (0-59)
	schedule.Minutes, err = parseField(fields[0], 0, 59)
	if err != nil {
		return nil, fmt.Errorf("invalid minute field: %w", err)
	}

	// Parse hour (0-23)
	schedule.Hours, err = parseField(fields[1], 0, 23)
	if err != nil {
		return nil, fmt.Errorf("invalid hour field: %w", err)
	}

	// Parse day (1-31)
	schedule.Days, err = parseField(fields[2], 1, 31)
	if err != nil {
		return nil, fmt.Errorf("invalid day field: %w", err)
	}

	// Parse month (1-12)
	schedule.Months, err = parseField(fields[3], 1, 12)
	if err != nil {
		return nil, fmt.Errorf("invalid month field: %w", err)
	}

	// Parse weekday (0-7, where 0 and 7 are Sunday)
	schedule.Weekdays, err = parseField(fields[4], 0, 7)
	if err != nil {
		return nil, fmt.Errorf("invalid weekday field: %w", err)
	}

	// Normalize weekday 7 to 0 (Sunday)
	for i, wd := range schedule.Weekdays {
		if wd == 7 {
			schedule.Weekdays[i] = 0
		}
	}

	return schedule, nil
}

// parseField parses a single cron field
func parseField(field string, min, max int) ([]int, error) {
	// Handle asterisk (*)
	if field == "*" {
		return rangeValues(min, max), nil
	}

	// Handle step values (*/n)
	if strings.HasPrefix(field, "*/") {
		step, err := strconv.Atoi(strings.TrimPrefix(field, "*/"))
		if err != nil || step < 1 {
			return nil, fmt.Errorf("invalid step value: %s", field)
		}
		return stepValues(min, max, step), nil
	}

	// Handle ranges (n-m)
	if strings.Contains(field, "-") {
		parts := strings.Split(field, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range: %s", field)
		}
		start, err1 := strconv.Atoi(parts[0])
		end, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil || start < min || end > max || start > end {
			return nil, fmt.Errorf("invalid range: %s", field)
		}
		return rangeValues(start, end), nil
	}

	// Handle comma-separated lists (n,m,o)
	if strings.Contains(field, ",") {
		parts := strings.Split(field, ",")
		var values []int
		for _, part := range parts {
			val, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil || val < min || val > max {
				return nil, fmt.Errorf("invalid value in list: %s", part)
			}
			values = append(values, val)
		}
		return values, nil
	}

	// Handle single value
	val, err := strconv.Atoi(field)
	if err != nil || val < min || val > max {
		return nil, fmt.Errorf("invalid value: %s (must be between %d and %d)", field, min, max)
	}
	return []int{val}, nil
}

// rangeValues returns all values in range [start, end]
func rangeValues(start, end int) []int {
	values := make([]int, end-start+1)
	for i := range values {
		values[i] = start + i
	}
	return values
}

// stepValues returns values in range [min, max] with given step
func stepValues(min, max, step int) []int {
	var values []int
	for i := min; i <= max; i += step {
		values = append(values, i)
	}
	return values
}

// Next returns the next time this cron schedule should run after the given time
func (cs *CronSchedule) Next(after time.Time) time.Time {
	// Start from the next minute
	t := after.Add(time.Minute).Truncate(time.Minute)

	// Try for up to 4 years (to handle edge cases)
	maxIterations := 4 * 365 * 24 * 60
	for i := 0; i < maxIterations; i++ {
		if cs.matches(t) {
			return t
		}
		t = t.Add(time.Minute)
	}

	// Should never reach here with valid cron expression
	return time.Time{}
}

// matches checks if the given time matches the cron schedule
func (cs *CronSchedule) matches(t time.Time) bool {
	// Check minute
	if !contains(cs.Minutes, t.Minute()) {
		return false
	}

	// Check hour
	if !contains(cs.Hours, t.Hour()) {
		return false
	}

	// Check day
	if !contains(cs.Days, t.Day()) {
		return false
	}

	// Check month
	if !contains(cs.Months, int(t.Month())) {
		return false
	}

	// Check weekday
	if !contains(cs.Weekdays, int(t.Weekday())) {
		return false
	}

	return true
}

// contains checks if a slice contains a value
func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// ValidateCronExpression validates a cron expression
func ValidateCronExpression(expr string) error {
	_, err := ParseCronExpression(expr)
	return err
}
