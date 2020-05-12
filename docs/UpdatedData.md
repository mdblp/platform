Here we document the data types that have to be created or udpated with new fields: 

- meal/wizard
- bolus
- physical activity
- deviceEvent - Alarm 
- food
- deviceEvent - Zen mode
- deviceEvent - Private Mode

## wizard 

The wizard object comes with an optional `recommended` structure that can be leveraged for our purpose. This structure is composed of 3 optional floating point value fields:
- carb: amount of insulin to cover the the total grams of carbohydrate input (`carbInput`)
- correction: amount of insulin recommended by the device to bring the PWD to their target blood glucose.
- net: total amount of recommended insulin

Here is an example of what can be sent with the related meaning:
- `recommended.net` is the system recommendation
- `bolus.normal` is the value delivered by the insulin pump.
- `bolus.expectedNormal` is the original value that has been requested to the insulin pump.
- `bolus.prescriptor` is a new field that describes the origin of the bolus. Details are defined in the below food section. 

And the additional field we would need:
- `entryTime` is a UTC string timestamp that defines at what time the patient has entered the meal. This field is optional. It takes the same format as `time` field.

```json
{
  "time": "2020-05-12T08:50:08.000Z",
  "timezoneOffset": 120,
  "deviceTime": "2020-05-12T08:50:08",
  "entryTime": "2020-05-12T08:45:08.000Z",
  "deviceId": "IdOfTheDevice",
  "type": "wizard",
  "carbInput": 50,
  "insulinOnBoard": 5.0,
  "recommended": {
    "net": 5
  },
  "units": "mg/dL",
  "bolus": {
    "time": "2020-05-12T08:50:08.000Z",
    "timezoneOffset": 120,
    "deviceTime": "2020-05-12T08:50:08",
    "deviceId": "IdOfTheDevice",
    "type": "bolus",
    "subType": "normal",
    "normal": 3.5,
    "expectedNormal": 4.0, 
    "prescriptor": "auto-altered"
  }
}
```

## food 

As of now we don't have the information of the origin of the rescueCarbs value, is it a patient decision, is it a suystem recommendation, and in that case what was the recommendation vs the actual value.

Here we are introducing 2 new fields in the food object:
- `prescribedNutrition`: same structure as nutrition. It's an optional field. It gives the value that has been recommended by the system. 
- `prescriptor`: is the origin of the `rescuecarbs` object. This field is optional in most of the cases. 
    - range of values: `auto | manual | auto-altered`
    - `auto`: nutrition and prescribedNutrition are equal
    - `manual`: prescribedNutrition is ignored
    - `auto-altered`: nutrition and prescribedNutrition are __not equal__, `prescribedNutrition` is mandatory in that case. 

```json
{
  "type": "food",
  "meal": "rescuecarbs",
  "nutrition": {
    "carbohydrate" : {
      "net": 20,
      "units": "grams"
    }
  },
  "prescribedNutrition": {
    "carbohydrate" : {
      "net": 30,
      "units": "grams"
    }
  },
  "prescriptor": "auto-altered",
  "meal": "rescuecarbs",
  "deviceId": "IdOfTheDevice",
  "deviceTime": "2020-05-12T06:50:08",
  "time": "2020-05-12T06:50:08.000Z",
  "timezoneOffset": 120
}

```

## bolus

3 types of bolus events are available as of now in the system:
- normal
- square
- dual/square

Here we are introducing 2 new fields in the bolus objects:
- `prescriptor`: same as above in `food`. This field is optional. 
- `insulinOnBoard`: amount of active insulin estimated by the system. This field will be accepted when `prescriptor` is either `auto` or `auto-altered`. It will be ignored for `manual` entries.

```json
{
  "time": "2020-05-12T08:50:08.000Z",
  "timezoneOffset": 120,
  "deviceTime": "2020-05-12T08:50:08",
  "deviceId": "IdOfTheDevice",
  "type": "bolus",
  "subType": "normal",
  "normal": 3.5,
  "expectedNormal": 4.0, 
  "prescriptor": "system-altered"
}
```

## physical activity

We need additional fields to get the time at which the physical activity is created, and the last time it was updated by the patient:
- `entryTime` is a UTC string timestamp that defines at what time the patient has entered the physical activity. This field is optional. It takes the same format as `time` field.
- `lastUpdatedTime` is a UTC string timestamp that gives the last time the patient has updated the physical activity. This field is optional. It takes the same format as `time` field.
  - `lastUpdatedTime` >= `entryTime`

```json
{
    "type": "physicalActivity",
    "reportedIntensity": "medium",
    "duration": { 
    	"value": 60,
    	"units": "minutes"
    },
    "clockDriftOffset": 0,
    "conversionOffset": 0,
    "deviceId": "DexG5MobRec_DX72101079",
    "deviceTime": "2016-07-12T23:52:47",
    "entryTime": "2020-05-12T08:00:08.000Z",
    "lastUpdatedTime": "2020-05-12T08:30:08.000Z",
    "time": "2020-05-12T08:50:08.000Z",
    "timezoneOffset": 60
}
```

## Alarm events
Leveraging the `deviceEvent` type with the already defined `alarm` subType. We add couple of fields to get more details on alarms and acknowledgement. 

- `alarmLevel`: `alarm | alert` 
- `alarmCode`: code of the alarm. This field is optional. 
- `alarmLabel`: label or description of the alarm. This field is optional. 
- `eventId`: unique Id of the event generated by the client system. This ID will be used to reconciliate data for the same event. 
- `eventType`: `start | acknowledge` is the type of event for the given alarm.
  - `start`: alarm generated by the system
  - `acknowledge`: the system has received the patient acknowledge. 

For a given alarm that has been acknowledged by the patient, we will receive 2 deviceEvents:
- the first one that gives the creation time on system, `eventType`: `start`
- the second one that gives the patient acknowledge, `eventType`: `acknowledge`

```json
{
  "type": "deviceEvent",
  "subType": "alarm",
  "alarmType": "handset",
  "alarmLevel": "alarm", 
  "alarmCode": "123456",
  "alarmLabel": "Label of the alarm",
  "eventId": "123456789",
  "eventType": "alarm",
  "deviceId": "Id12345",
  "deviceTime": "2018-02-01T00:00:00",
  "guid": "e3c82e3d-23ba-4048-9056-7b1b3c5aa4cc",
  "id": "d5ed640dd8f74e6cb1a6bff796de3ba2",
  "time": "2020-05-12T08:50:08.000Z",
  "timezoneOffset": 60
}
```


## Zen mode
Leveraging the `deviceEvent` type and creating a new `zen` subType.

- `subType`: `zen`
- `endTime`: is a UTC string timestamp that indicates at what time the event is terminated. 

```json
{
  "type": "deviceEvent",
  "subType": "zen",
  "endTime": "2020-05-12T09:50:08.000Z",
  "deviceId": "Id12345",
  "deviceTime": "2018-02-01T00:00:00",
  "guid": "e3c82e3d-23ba-4048-9056-7b1b3c5aa4cc",
  "id": "d5ed640dd8f74e6cb1a6bff796de3ba2",
  "time": "2020-05-12T08:50:08.000Z",
  "timezoneOffset": 60
}
```

## Confidential mode
Leveraging the `deviceEvent` type and creating a new `confidential` subType.

- `subType`: `confidential`
- `endTime`: is a UTC string timestamp that indicates at what time the event is terminated. 

```json
{
  "type": "deviceEvent",
  "subType": "confidential",
  "endTime": "2020-05-12T10:50:08.000Z",
  "deviceId": "Id12345",
  "deviceTime": "2018-02-01T00:00:00",
  "guid": "e3c82e3d-23ba-4048-9056-7b1b3c5aa4cc",
  "id": "d5ed640dd8f74e6cb1a6bff796de3ba2",
  "time": "2020-05-12T08:50:08.000Z",
  "timezoneOffset": 60
}
```
