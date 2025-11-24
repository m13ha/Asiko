
# ResponsesAppointmentResponse


## Properties

Name | Type
------------ | -------------
`appCode` | string
`bookingDuration` | number
`createdAt` | string
`description` | string
`endDate` | string
`endTime` | string
`id` | string
`maxAttendees` | number
`startDate` | string
`startTime` | string
`status` | [EntitiesAppointmentStatus](EntitiesAppointmentStatus.md)
`title` | string
`type` | [EntitiesAppointmentType](EntitiesAppointmentType.md)
`updatedAt` | string

## Example

```typescript
import type { ResponsesAppointmentResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "appCode": null,
  "bookingDuration": null,
  "createdAt": null,
  "description": null,
  "endDate": null,
  "endTime": null,
  "id": null,
  "maxAttendees": null,
  "startDate": null,
  "startTime": null,
  "status": null,
  "title": null,
  "type": null,
  "updatedAt": null,
} satisfies ResponsesAppointmentResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ResponsesAppointmentResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


