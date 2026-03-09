package tools

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/jamespfennell/gtfs"
)

type BARTTripInfo struct {
	RouteID     string
	DirectionID string
	Headsign    string
}

// have to look up trip data separately unfortunately
func getBARTTripLookup() (map[string]BARTTripInfo, error) {
	resp, err := http.Get("https://www.bart.gov/dev/schedules/google_transit.zip")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	staticData, err := gtfs.ParseStatic(b, gtfs.ParseStaticOptions{})
	if err != nil {
		return nil, err
	}

	tripLookup := make(map[string]BARTTripInfo)
	for _, trip := range staticData.Trips {
		tripLookup[trip.ID] = BARTTripInfo{
			RouteID:     trip.Route.Id,
			DirectionID: trip.DirectionId.String(),
			Headsign:    trip.Headsign,
		}
	}

	return tripLookup, nil
}

var BARTStationTool = Tool{
	Name:        "bart_station_list",
	Description: "Returns a map of stations and IDs",
	Exec: func(params map[string]interface{}) (string, error) {
		fmt.Println("bart: station list")
		resp, err := http.Get("https://www.bart.gov/dev/schedules/google_transit.zip")
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		staticData, _ := gtfs.ParseStatic(b, gtfs.ParseStaticOptions{})
		// retrieve ID for each
		stopCount := 0
		var ret strings.Builder
		for _, stop := range staticData.Stops {
			match := regexp.MustCompile(`-\d$|_\d$`).FindString(stop.Id)
			if match != "" {
				continue
			}
			stopCount++

			// if len(stop.Id) != 4 {
			// 	continue
			// }
			ret.WriteString(fmt.Sprintf("Stop `%s` with ID `%s`", stop.Name, stop.Id) + "\n")
		}
		// per stop - CODE: 902101 Name: Lake Meritt ID: A10-1 ZoneID: LAKE
		// fmt.Println(stopCount)
		return ret.String(), nil
	},
}

func getBARTPlatformToStationMap() (map[string]string, error) {
	resp, err := http.Get("https://www.bart.gov/dev/schedules/google_transit.zip")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	staticData, err := gtfs.ParseStatic(b, gtfs.ParseStaticOptions{})
	if err != nil {
		return nil, err
	}

	platformToStation := make(map[string]string)
	for _, stop := range staticData.Stops {
		if len(stop.Id) == 4 {
			continue
		}
		if stop.Parent.ZoneId != "" && len(stop.Parent.ZoneId) == 4 {
			platformToStation[stop.Id] = stop.Parent.ZoneId
		}
	}

	return platformToStation, nil
}

var BARTRealTimeTool = Tool{
	Name:        "bart_realtime",
	Description: "Gets real-time BART info sourced from GTFS",
	Parameters: map[string]any{
		"station_id": map[string]any{
			"type":        "string",
			"description": "The station ID to get info for",
		},
	},
	Exec: func(params map[string]interface{}) (string, error) {
		stationID, ok := params["station_id"].(string)
		if !ok {
			return "", fmt.Errorf("station_id parameter is required")
		}

		fmt.Println("bart: getting realtime data")
		resp, err := http.Get("http://api.bart.gov/gtfsrt/tripupdate.aspx")
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		realtimeData, err := gtfs.ParseRealtime(b, &gtfs.ParseRealtimeOptions{})
		if err != nil {
			return "", err
		}

		platformToStation, err := getBARTPlatformToStationMap()
		if err != nil {
			return "", err
		}

		// fmt.Println(platformToStation)

		var aText strings.Builder
		foundAny := false

		for _, trip := range realtimeData.Trips {

			for _, update := range trip.StopTimeUpdates {
				fmt.Println(trip.ID.ID)
				// fmt.Println(*update.StopID, update.Arrival.Time, update.Departure.Time)
				mstation := platformToStation[*update.StopID]
				// fmt.Println(mstation)
				if mstation == stationID {
					foundAny = true
					if update.Arrival != nil && !update.Arrival.Time.IsZero() {
						arrivalTime := update.Arrival.Time.Format("3:04 PM")

						aText.WriteString(fmt.Sprintf("Train %s arriving at %s\n", trip.ID.ID, arrivalTime))
					}
				}
			}
		}

		if !foundAny {
			return fmt.Sprintf("No trains found for station %s", stationID), nil
		}

		return aText.String(), nil
	},
}
