package utils

func IsContained(userID int64, userIDs []int64) bool {
	for _, id := range userIDs {
		if id == userID {
			return true
		}
	}
	return false
}
