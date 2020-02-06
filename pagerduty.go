package main

import (
	"flag"
	"github.com/coeurmining/prtg-pagerduty/event"
	"log"
	"strings"
	"time"
	"os"
)

// PRTGEvent represents the data passed by prtg via flags
type PRTGEvent struct {
	Probe       string
	Device      string
	Name        string
	Status      string
	Date        string
	Link        string
	Message     string
	ServiceKey  string
	IncidentKey string
	Severity    string
	Priority    string
	CustRouting string
}

func main() {
	var probe = flag.String("probe", "local", "The PRTG probe name")
	var device = flag.String("device", "device", "The PRTG device name")
	var name = flag.String("name", "name", "The PRTG sensor name for the device")
	var status = flag.String("status", "status", "The current status for the event")
	var date = flag.String("date", "date", "The date time for the triggered event")
	var link = flag.String("linkdevice", "http://localhost", "The link to the triggering sensor")
	var message = flag.String("message", "message", "The PRTG message for the alert")
	var serviceKey = flag.String("servicekey", "myServiceKey", "The PagerDuty v2 service integration key")
	var severity = flag.String("severity", "error", "The severity level of the incident (critical, error, warning, or info)")
	var priority = flag.String("priority", "priority", "The Priority of the Sensor in PRTG")
	var custrouting = flag.String("custrouting", "custrouting", "The custom routing identifier for PD Event Rules")
	
	flag.Parse()

	*probe = truncateString(*probe, 100)
	*device = truncateString(*device, 100)
	*name = truncateString(*name, 100)
	*status = truncateString(*status, 100)
	*date = truncateString(*date, 100)
	*link = truncateString(*link, 100)
	*message = truncateString(*message, 100)
	*serviceKey = truncateString(*serviceKey, 100)
	*severity = truncateString(*severity, 100)
	*priority = truncateString(*priority, 100)
	*custrouting = truncateString(*custrouting, 100)
	
	pd := &PRTGEvent{
		Probe:       *probe,
		Device:      *device,
		Name:        *name,
		Status:      *status,
		Date:        *date,
		Link:        *link,
		Message:     *message,
		ServiceKey:  *serviceKey,
		IncidentKey: truncateString(*probe + "-" + *device + "-" + *name, 50),
		Severity:    *severity,
		Priority:    *priority,
		CustRouting: *custrouting,
	}

	if strings.Contains(pd.Status, "Up") || strings.Contains(pd.Status, "ended") {
		resolveEvent(pd)
	} else {
		event, err := triggerEvent(pd)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(event)
	}
}

func truncateString(str string, num int) string {
	stringtotruncate := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		stringtotruncate = str[0:num] + "..."
	}
	return stringtotruncate
}

func triggerEvent(prtg *PRTGEvent) (*event.Response, error) {
	const layout = "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, prtg.Date)
	if err != nil {
		t = time.Now()
	}
	newEvent := &event.Event{
		RoutingKey: truncateString(prtg.ServiceKey, 100),
		Action:     "trigger",
		DedupKey:   truncateString(prtg.IncidentKey, 100),
		Client:     "PRTG",
		ClientURL:  truncateString(prtg.Link, 100),
		Payload: &event.Payload{
			Summary:   truncateString(prtg.IncidentKey, 255),
			Timestamp: t.Format(layout),
			Source:    truncateString(prtg.Link, 100),
			Severity:  translatePriority(prtg.Priority),
			Component: truncateString(prtg.Device, 100),
			Group:     truncateString(prtg.Probe, 100),
			Class:     truncateString(prtg.Name, 100),
			Details: "Link: " + prtg.Link +
				"\nIncidentKey: " + prtg.IncidentKey +
				"\nStatus: " + prtg.Status +
				"\nDate: " + prtg.Date +
				"\nMessage: " + prtg.Message +
				"\nCustom Routing: " + prtg.CustRouting,
		},
	}
	res, err := event.ManageEvent(*newEvent)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func resolveEvent(prtg *PRTGEvent) (*event.Response, error) {
	triggeredEvent := &event.Event{
		RoutingKey: prtg.ServiceKey,
		Action:     "resolve",
		DedupKey:   prtg.IncidentKey,
	}
	res, err := event.ManageEvent(*triggeredEvent)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func translatePriority(priority string) (string){
    switch priority {
    case "*":
        os.Exit(0)
        return "info" //function requires a return
    case "**":
        return "info"
    case "***":
        return "warning"
    case "****":
        return "error"
    default:
        return "critical"
    }
}
