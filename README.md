# PRTG/PagerDuty Notification Integration

## Goals

* Create incidents using version 2 of the PagerDuty Events API for triggered PRTG alerts.

* Automatically resolve alerts when status returns to normal or paused in PRTG.

* Pass sensor priority to PagerDuty for event rule processing

## Build & Installation

Build the package

`go get github.com/ccummings-coeur/prtg-pagerduty`

`go build`

Copy the built executable to the PRTG Core Server, and from an Adminstrator powershell session:

`cp pagerduty.exe "C:\Program Files (x86)\PRTG Network Monitor\Notifications\EXE\"`


## Configuring PRTG Notification Template

Create new notification template. Check "EXECUTE PROGRAM" selecting pagerduty.exe from the Program File dropdown.

Populate the parameter field with the following, substituting the service key with your service integration key

`-probe "%probe" -device "%device" -name "%name" -status "%status" -date "%datetime" -linkdevice %linkdevice -message "%message" -priority "%priority"  -custrouting "CDA_Alert_Routing" -servicekey myShineyV2IntegrationKey`

## PRTG Priority to PagerDuty Severity Mapping

PagerDuty severity is mapped from PRTG priority stars as follows: 

| PRTG Priority   | PagerDuty Severity |
|-----------------|--------------------|
| x               | Ignore             |
| xx              | info               |
| xxx             | warning            |
| xxxx            | error              |
| xxxxx / Default | critical           |

PRTG sensors with one star are ignored, and don't get sent to PagerDuty. 