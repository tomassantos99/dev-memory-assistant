package pkg

func CropString(s string, max int) string {
    if len(s) <= max {
        return s
    }
    if max <= 3 {
        return s[:max]
    }
    return s[:max-3] + "..."
}
