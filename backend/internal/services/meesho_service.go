package services

import (
	"meesho-clone/internal/models"
	"net/http"
	"strings"
	"time"
)

// MeeshoService handles integration with Meesho API
type MeeshoService struct {
	client  *http.Client
	baseURL string
}

// NewMeeshoService creates a new Meesho service
func NewMeeshoService() *MeeshoService {
	return &MeeshoService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://app.stg.meeshogcp.in/api/1.0",
	}
}

// FetchHomescreenData fetches homescreen data from Meesho API
func (s *MeeshoService) FetchHomescreenData(userID string) (*models.MeeshoAPIResponse, error) {
	// Return mock data directly embedded in the code
	response := &models.MeeshoAPIResponse{
		TopNavBar: models.TopNavBar{
			ID:    1,
			Title: "New Navigation Bar India",
			Tiles: []models.NavTile{
				{
					ID:    1,
					Title: "Categories",
					Image: "https://images-dev.meesho.com/images/marketing/1594489152649.jpg",
					DestinationData: map[string]interface{}{
						"screen":              "categories",
						"has_location_filter": false,
					},
				},
				{
					ID:    228,
					Title: "Kurti & Dress Materials",
					Image: "https://images-dev.meesho.com/images/marketing/1649688502928.png",
					DestinationData: map[string]interface{}{
						"screen":              "single_collection",
						"collection_id":       36176,
						"collection_name":     "Kurti & Sets",
						"has_location_filter": false,
					},
				},
				{
					ID:    234,
					Title: "Kids & Toys ",
					Image: "https://images-dev.meesho.com/images/marketing/1649689217815.png",
					DestinationData: map[string]interface{}{
						"screen":              "single_collection",
						"collection_id":       32448,
						"collection_name":     "Kids",
						"has_location_filter": false,
					},
				},
				{
					ID:    230,
					Title: "Westernwear",
					Image: "https://images-dev.meesho.com/images/marketing/1649690440106.jpeg",
					DestinationData: map[string]interface{}{
						"screen":              "single_collection",
						"collection_id":       30994,
						"collection_name":     "Women Fashion",
						"has_location_filter": false,
					},
				},
				{
					ID:    236,
					Title: "Home",
					Image: "https://images-dev.meesho.com/images/marketing/1687854960403.jpg",
					DestinationData: map[string]interface{}{
						"screen":                    "catalog_listing_page",
						"catalog_listing_page_id":   172088,
						"payload":                   "",
						"catalog_listing_page_name": "Home Decor & Furnishings",
						"has_location_filter":       false,
					},
				},
				{
					ID:    231,
					Title: "Men Clothing",
					Image: "https://images-dev.meesho.com/images/marketing/1689675132726.png",
					DestinationData: map[string]interface{}{
						"screen":              "single_collection",
						"collection_id":       48370,
						"collection_name":     "Men Fashion",
						"has_location_filter": false,
					},
				},
				{
					ID:    227,
					Title: "Saree",
					Image: "https://images-dev.meesho.com/images/marketing/1628672353857.jpg",
					DestinationData: map[string]interface{}{
						"screen":              "single_collection",
						"collection_id":       30978,
						"collection_name":     "Sarees",
						"has_location_filter": false,
					},
				},
				{
					ID:    232,
					Title: "Suits",
					Image: "https://images-dev.meesho.com/images/marketing/1628672428427.jpg",
					DestinationData: map[string]interface{}{
						"screen":              "single_collection",
						"collection_id":       30990,
						"collection_name":     "Suits",
						"has_location_filter": false,
					},
				},
				{
					ID:    8,
					Title: "Beauty",
					Image: "https://images-dev.meesho.com/images/marketing/1651505214223.png",
					DestinationData: map[string]interface{}{
						"screen":                    "catalog_listing_page",
						"catalog_listing_page_id":   2,
						"payload":                   "",
						"catalog_listing_page_name": "Solid",
						"has_location_filter":       false,
					},
				},
				{
					ID:    229,
					Title: "Jewellery",
					Image: "https://images-dev.meesho.com/images/marketing/1649689138272.jpeg",
					DestinationData: map[string]interface{}{
						"screen":              "single_collection",
						"collection_id":       30984,
						"collection_name":     "Jewellery",
						"has_location_filter": false,
					},
				},
				{
					ID:    235,
					Title: "Home Textiles",
					Image: "https://images-dev.meesho.com/images/marketing/1687854996278.jpg",
					DestinationData: map[string]interface{}{
						"screen":                    "catalog_listing_page",
						"catalog_listing_page_id":   68288,
						"payload":                   "",
						"catalog_listing_page_name": " Home Textiles",
						"has_location_filter":       false,
					},
				},
				{
					ID:    238,
					Title: "Kitchen",
					Image: "https://images-dev.meesho.com/images/marketing/1687854973673.jpg",
					DestinationData: map[string]interface{}{
						"screen":                    "catalog_listing_page",
						"catalog_listing_page_id":   86271,
						"payload":                   "",
						"catalog_listing_page_name": "Kitchen Utility",
						"has_location_filter":       false,
					},
				},
			},
		},
		WidgetGroups: []models.WidgetGroup{
			{
				ID:              42109,
				MongoWGID:       "6712220c9433e105c4bb6d1b",
				Position:        0,
				Title:           "Travel Bags Upto 70% Offf",
				Tag:             "",
				Type:            5,
				BackgroundColor: "#FFFFFF",
				Dynamic:         false,
				DSEnabled:       false,
				AdsEnabled:      false,
				Priority:        10000000,
				WidgetGroupType: "FIXED",
				Widgets: []models.Widget{
					{
						ID:               382108,
						Title:            "Laptop Bags",
						Image:            "https://images-dev.meesho.com/images/widgets/UHQP2/sq1an.jpg",
						ImageAspectRatio: 1.0,
						Screen:           "catalog_listing_page",
						Type:             23,
						DestinationID:    180679,
						Fixed:            false,
						Priority:         6,
						Data: map[string]interface{}{
							"catalog_listing_page_id":   180679,
							"payload":                   "",
							"catalog_listing_page_name": "Laptop Bags New Mall",
							"with_oauth":                false,
							"has_location_filter":       false,
							"clp_gst_type":              "gst",
							"is_dynamic":                false,
							"screen":                    "catalog_listing_page",
							"view_template_id":          1,
						},
					},
					{
						ID:               382109,
						Title:            "Duffel Bags",
						Image:            "https://images-dev.meesho.com/images/widgets/A3U85/bq29c.jpg",
						ImageAspectRatio: 1.0,
						Screen:           "catalog_listing_page",
						Type:             23,
						DestinationID:    186321,
						Fixed:            false,
						Priority:         5,
						Data: map[string]interface{}{
							"catalog_listing_page_id":   186321,
							"payload":                   "",
							"catalog_listing_page_name": "Duffel Bags",
							"with_oauth":                false,
							"has_location_filter":       false,
							"clp_gst_type":              "gst",
							"is_dynamic":                false,
							"screen":                    "catalog_listing_page",
							"view_template_id":          1,
						},
					},
					{
						ID:               382110,
						Title:            "Small Travel Bags",
						Image:            "https://images-dev.meesho.com/images/widgets/R2XUE/xgdrn.jpg",
						ImageAspectRatio: 1.0,
						Screen:           "catalog_listing_page",
						Type:             23,
						DestinationID:    180448,
						Fixed:            false,
						Priority:         4,
						Data: map[string]interface{}{
							"catalog_listing_page_id":   180448,
							"payload":                   "",
							"catalog_listing_page_name": "Small Travel Bags Mall",
							"with_oauth":                false,
							"has_location_filter":       false,
							"clp_gst_type":              "gst",
							"is_dynamic":                false,
							"screen":                    "catalog_listing_page",
							"view_template_id":          1,
						},
					},
					{
						ID:               382111,
						Title:            "Rucksacks",
						Image:            "https://images-dev.meesho.com/images/widgets/Z6MX8/zzxp9.jpg",
						ImageAspectRatio: 1.0,
						Screen:           "catalog_listing_page",
						Type:             23,
						DestinationID:    179100,
						Fixed:            false,
						Priority:         3,
						Data: map[string]interface{}{
							"catalog_listing_page_id":   179100,
							"payload":                   "",
							"catalog_listing_page_name": "Rucksacks",
							"with_oauth":                false,
							"has_location_filter":       false,
							"clp_gst_type":              "gst",
							"is_dynamic":                false,
							"screen":                    "catalog_listing_page",
							"view_template_id":          1,
						},
					},
					{
						ID:               382112,
						Title:            "Pouches",
						Image:            "https://images-dev.meesho.com/images/widgets/DEHQA/a63ts.jpg",
						ImageAspectRatio: 1.0,
						Screen:           "catalog_listing_page",
						Type:             23,
						DestinationID:    180438,
						Fixed:            false,
						Priority:         2,
						Data: map[string]interface{}{
							"catalog_listing_page_id":   180438,
							"payload":                   "",
							"catalog_listing_page_name": "Pouches Mall",
							"with_oauth":                false,
							"has_location_filter":       false,
							"clp_gst_type":              "gst",
							"is_dynamic":                false,
							"screen":                    "catalog_listing_page",
							"view_template_id":          1,
						},
					},
					{
						ID:               382113,
						Title:            "Travel Accessories",
						Image:            "https://images-dev.meesho.com/images/widgets/9LYKG/lvjqr.jpg",
						ImageAspectRatio: 1.0,
						Screen:           "catalog_listing_page",
						Type:             23,
						DestinationID:    186322,
						Fixed:            false,
						Priority:         1,
						Data: map[string]interface{}{
							"catalog_listing_page_id":   186322,
							"payload":                   "",
							"catalog_listing_page_name": "Travel Accessories & Others",
							"with_oauth":                false,
							"has_location_filter":       false,
							"clp_gst_type":              "gst",
							"is_dynamic":                false,
							"screen":                    "catalog_listing_page",
							"view_template_id":          1,
						},
					},
				},
			},
		},
	}

	return response, nil
}

// SearchProducts searches for products (placeholder for future implementation)
func (s *MeeshoService) SearchProducts(query string, userID string) (interface{}, error) {
	// This would be implemented when we add search functionality
	// For now, return empty result
	return map[string]interface{}{
		"message": "Search functionality coming soon",
		"query":   query,
		"user_id": userID,
	}, nil
}

// GetProductDetails gets details for a specific product (placeholder)
func (s *MeeshoService) GetProductDetails(productID string, userID string) (interface{}, error) {
	// This would be implemented when we add product details functionality
	return map[string]interface{}{
		"message":    "Product details functionality coming soon",
		"product_id": productID,
		"user_id":    userID,
	}, nil
}

// FormatHomescreenResponse formats the response for frontend consumption
func (s *MeeshoService) FormatHomescreenResponse(apiResponse *models.MeeshoAPIResponse, userID string) map[string]interface{} {
	// Add user context to the response
	formattedResponse := map[string]interface{}{
		"user_id":       userID,
		"top_nav_bar":   apiResponse.TopNavBar,
		"widget_groups": apiResponse.WidgetGroups,
		"timestamp":     time.Now().Unix(),
		"success":       true,
	}

	// Process categories for easier frontend consumption
	if len(apiResponse.TopNavBar.Tiles) > 0 {
		categories := make([]map[string]interface{}, 0)
		for _, tile := range apiResponse.TopNavBar.Tiles {
			category := map[string]interface{}{
				"id":          tile.ID,
				"title":       tile.Title,
				"image":       tile.Image,
				"destination": tile.DestinationData,
			}
			categories = append(categories, category)
		}
		formattedResponse["categories"] = categories
	}

	// Process products/widgets for easier frontend consumption
	if len(apiResponse.WidgetGroups) > 0 {
		products := make([]map[string]interface{}, 0)
		for _, group := range apiResponse.WidgetGroups {
			for _, widget := range group.Widgets {
				product := map[string]interface{}{
					"id":          widget.ID,
					"title":       widget.Title,
					"image":       widget.Image,
					"screen":      widget.Screen,
					"group_title": group.Title,
					"group_id":    group.ID,
					"destination": widget.Data,
				}
				products = append(products, product)
			}
		}
		formattedResponse["products"] = products
	}

	return formattedResponse
}

// ExtractImageURL safely extracts image URL from various sources
func (s *MeeshoService) ExtractImageURL(sources ...string) string {
	for _, source := range sources {
		if source != "" && strings.HasPrefix(source, "http") {
			return source
		}
	}
	return ""
}
