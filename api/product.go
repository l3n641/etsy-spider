package api

type ProductReq struct {
	Url         string   `json:"url" binding:"required"`
	Title       string   `json:"title" binding:"required" `
	HtmlContent string   `json:"htmlContent" binding:"required" `
	Images      []string `json:"images" binding:"required"`
}
