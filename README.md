# PRTG/PagerDuty Notification Integration

## Goals

* Create incidents using version 2 of the PagerDuty Events API for triggered PRTG alerts.

* Automatically resolve alerts when status returns to normal or paused in PRTG.


## Build & Installation

Build the package

`go get github.com/ccummings-coeur/prtg-pagerduty`

`go build`

From an Adminstrator powershell session:

`cp pagerduty.exe "C:\Program Files (x86)\PRTG Network Monitor\Notifications\EXE\"`


## Configuring PRTG Notification Template

Create new notification template. Check "EXECUTE PROGRAM" selecting pagerduty.exe from the Program File dropdown.

Populate the parameter field with the following, substituting the service key with your service integration key

`-probe "%probe" -device "%device" -name "%name" -status "%status" -date "%datetime" -linkdevice %linkdevice -message "%message" -priority "%priority"  -custrouting "CDA_Alert_Routing" -servicekey myShineyV2IntegrationKey`

If you want the severity set on the alert, add `-severity "critical"` the value can be replaced with the alert severity that you want set, in this case, critical is used.
