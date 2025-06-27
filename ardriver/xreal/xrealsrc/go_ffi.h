#include "device_imu.h"
extern void goIMUEventHandler(uint64_t, device_imu_event_type, device_imu_ahrs_type*);
void imuEventHandler(uint64_t timestamp, device_imu_event_type event, const device_imu_ahrs_type* ahrs);
