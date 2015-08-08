package control

import (
	"fmt"
)

func InitMail(instanceId string, region string, lifeTime int, bidRegion string) string {
	return fmt.Sprintf(`Hi,

I was just spawned on instance %v in region %v. 
My lifetime will hopefully be %v seconds. I already create a spot request in %v.

Bye for now,

Tatanka
`, instanceId, region, lifeTime, bidRegion)

}

func EvacuationMail(instanceId string, region string, bidRegion string) string {
	return fmt.Sprintf(`Hi,

I'm evacuating instance %v in region %v. I created a new spot request in %v. 
Hope I see you again soon.

Bye for now,

Tatanka
`, instanceId, region, bidRegion)

}
