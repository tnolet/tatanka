package control

import (
	"fmt"
)

func FatalErrorMail(error string) string {
	return fmt.Sprintf(`Hi,

It seems I encountered a fatal error. That sucks. The error goes like:

%v

Please send help,

Tatanka
`, error)

}


func InitMail(instanceId string, region string, lifeTime int, bidRegion string, reqId string) string {
	return fmt.Sprintf(`Hi,

I was just spawned on instance %v in region %v. 
My lifetime will hopefully be %v seconds. I already created spot request %v in %v.

Bye for now,

Tatanka
`, instanceId, region, lifeTime, reqId, bidRegion)

}

func EvacuationMail(instanceId string, region string, bidRegion string) string {
return fmt.Sprintf(`Hi,

I'm evacuating instance %v in region %v. I created a new spot request in %v.
Hope I see you again soon.

Bye for now,

Tatanka
`, instanceId, region, bidRegion)

}

func CasualUpdateMail(upTime int,region string) string {
	return fmt.Sprintf(`Hi,

Just letting you know I'm doing fine here, running for %v seconds in %v.

Bye for now,

Tatanka
`, upTime, region)

}

