package runtime

type Colors struct{}

func (g *Colors) Pick(index int) string {
	colorPalette := []string{
		"#E69F00", // Orange
		"#56B4E9", // Sky Blue
		"#009E73", // Bluish Green
		"#F0E442", // Yellow
		"#0072B2", // Blue
		"#D55E00", // Vermillion
		"#CC79A7", // Reddish Purple
		"#8DD3C7", // Light Blue-Green
		"#FDB462", // Soft Orange
		"#B3DE69", // Light Green
		"#FFED6F", // Light Yellow
		"#6A3D9A", // Deep Purple
		"#B15928", // Brownish-Orange
		"#44AA99", // Teal
		"#117733", // Dark Green
		"#999933", // Olive Green
		"#AA4499", // Purple
		"#DDCC77", // Light Tan
		"#882255", // Dark Red
		"#332288", // Dark Blue
	}

	return colorPalette[index%len(colorPalette)]
}
