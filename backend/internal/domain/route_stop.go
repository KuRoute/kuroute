package domain

import "github.com/google/uuid"

// Model

// RouteStop is one node in a computed route an immutable snapshot
// of where the courier must go and in what order.
// lat/lng are copied from Package at computation time to keep
// the route immutable even if the package address is later corrected.
type RouteStop struct {
	ID        uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RouteID   uuid.UUID     `gorm:"type:uuid;not null;index:idx_route_stop_route"  json:"routeId"`
	PackageID uuid.UUID     `gorm:"type:uuid;not null"                             json:"packageId"`
	StopOrder int16         `gorm:"not null;index:idx_route_stop_route"            json:"stopOrder"`
	Lat       float64       `gorm:"type:numeric(9,6);not null"                     json:"lat"`
	Lng       float64       `gorm:"type:numeric(9,6);not null"                     json:"lng"`
	Status    PackageStatus `gorm:"type:package_status;not null;default:'assigned'" json:"status"`

	Route           Route            `gorm:"foreignKey:RouteID"    json:"-"`
	Package         Package          `gorm:"foreignKey:PackageID"  json:"-"`
	DeliveryReports []DeliveryReport `gorm:"foreignKey:RouteStopID" json:"-"`
}

func (RouteStop) TableName() string { return "route_stop" }

// Response

// RouteStop has no public Create/Update request — it is created
// exclusively by the route solver via ComputeRouteRequest.
// Status updates are side effects of DeliveryReport creation.
type RouteStopResponse struct {
	ID        uuid.UUID     `json:"id"`
	RouteID   uuid.UUID     `json:"routeId"`
	PackageID uuid.UUID     `json:"packageId"`
	StopOrder int16         `json:"stopOrder"`
	Lat       float64       `json:"lat"`
	Lng       float64       `json:"lng"`
	Status    PackageStatus `json:"status"`

	RecipientName string `json:"recipientName,omitempty"`
	AddressText   string `json:"addressText,omitempty"`
}

type RouteStopDetailResponse struct {
	ID        uuid.UUID     `json:"id"`
	RouteID   uuid.UUID     `json:"routeId"`
	PackageID uuid.UUID     `json:"packageId"`
	StopOrder int16         `json:"stopOrder"`
	Lat       float64       `json:"lat"`
	Lng       float64       `json:"lng"`
	Status    PackageStatus `json:"status"`

	RecipientName  string                  `json:"recipientName,omitempty"`
	AddressText    string                  `json:"addressText,omitempty"`
	DeliveryReport *DeliveryReportResponse `json:"deliveryReport,omitempty"`
}

func NewRouteStopResponse(rs *RouteStop) RouteStopResponse {
	resp := RouteStopResponse{
		ID:        rs.ID,
		RouteID:   rs.RouteID,
		PackageID: rs.PackageID,
		StopOrder: rs.StopOrder,
		Lat:       rs.Lat,
		Lng:       rs.Lng,
		Status:    rs.Status,
	}
	if rs.Package.ID != (uuid.UUID{}) {
		resp.RecipientName = rs.Package.RecipientName
		resp.AddressText = rs.Package.AddressText
	}
	return resp
}

func NewRouteStopDetailResponse(rs *RouteStop) RouteStopDetailResponse {
	resp := RouteStopDetailResponse{
		ID:        rs.ID,
		RouteID:   rs.RouteID,
		PackageID: rs.PackageID,
		StopOrder: rs.StopOrder,
		Lat:       rs.Lat,
		Lng:       rs.Lng,
		Status:    rs.Status,
	}
	if rs.Package.ID != (uuid.UUID{}) {
		resp.RecipientName = rs.Package.RecipientName
		resp.AddressText = rs.Package.AddressText
	}
	if len(rs.DeliveryReports) > 0 {
		latest := rs.DeliveryReports[0]
		for _, report := range rs.DeliveryReports[1:] {
			if report.ReportedAt.After(latest.ReportedAt) {
				latest = report
			}
		}
		deliveryReport := NewDeliveryReportResponse(&latest)
		resp.DeliveryReport = &deliveryReport
	}
	return resp
}
