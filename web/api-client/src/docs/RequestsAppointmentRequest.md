
# RequestsAppointmentRequest


## Properties

Name | Type
------------ | -------------
`antiScalpingLevel` | [EntitiesAntiScalpingLevel](EntitiesAntiScalpingLevel.md)
`bookingDuration` | number
`description` | string
`endDate` | string
`endTime` | string
`maxAttendees` | number
`startDate` | string
`startTime` | string
`title` | string
`type` | [EntitiesAppointmentType](EntitiesAppointmentType.md)

## Example

```typescript
import type { RequestsAppointmentRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "antiScalpingLevel": null,
  "bookingDuration": null,
  "description": null,
  "endDate": null,
  "endTime": null,
  "maxAttendees": null,
  "startDate": null,
  "startTime": null,
  "title": null,
  "type": null,
} satisfies RequestsAppointmentRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as RequestsAppointmentRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


