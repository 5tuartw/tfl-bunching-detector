package helpers

func HeadingToDirection(heading int) string {
	switch true {
	case (heading >= 315 || heading <= 45):
		return "northbound"
	case (45 < heading && heading < 135):
		return "eastbound"
	case (135 <= heading && heading <= 225):
		return "southbound"
	case (225 < heading && heading < 315):
		return "westbound"
	default:
		return ""
	}
}
