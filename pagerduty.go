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
	TruncateLength int
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
	var truncatelength = flag.Int("truncatelength", 99, "The length that all inputs should be truncated to to prevent PagerDuty errors.")

	flag.Parse()

	*probe = truncateString(*probe, *truncatelength)
	*device = truncateString(*device, *truncatelength)
	*name = truncateString(*name, *truncatelength)
	*status = truncateString(*status, *truncatelength)
	*date = truncateString(*date, *truncatelength)
	*link = truncateString(*link, *truncatelength)
	*message = truncateString(*message, *truncatelength)
	*serviceKey = truncateString(*serviceKey, *truncatelength)
	*severity = truncateString(*severity, *truncatelength)
	*priority = truncateString(*priority, *truncatelength)
	*custrouting = truncateString(*custrouting, *truncatelength)
	
	pd := &PRTGEvent{
		Probe:       *probe,
		Device:      *device,
		Name:        *name,
		Status:      *status,
		Date:        *date,
		Link:        *link,
		Message:     *message,
		ServiceKey:  *serviceKey,
		IncidentKey: *probe + "-" + *device + "-" + *name,
		Severity:    *severity,
		Priority:    *priority,
		CustRouting: *custrouting,
		TruncateLength: *truncatelength,
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
		RoutingKey: truncateString(prtg.ServiceKey, prtg.TruncateLength),
		Action:     "trigger",
		DedupKey:   truncateString(prtg.IncidentKey, prtg.TruncateLength),
		Client:     "PRTG",
		ClientURL:  truncateString(prtg.Link, prtg.TruncateLength),
		Payload: &event.Payload{
			Summary:   truncateString(prtg.IncidentKey, 254),
			Timestamp: t.Format(layout),
			Source:    truncateString(prtg.Link, prtg.TruncateLength),
			Severity:  translatePriority(prtg.Priority),
			Component: truncateString(prtg.Device, prtg.TruncateLength),
			Group:     truncateString(prtg.Probe, prtg.TruncateLength),
			Class:     truncateString(prtg.Name, prtg.TruncateLength),
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
