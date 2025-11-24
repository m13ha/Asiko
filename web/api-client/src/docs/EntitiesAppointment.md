
# EntitiesAppointment


## Properties

Name | Type
------------ | -------------
`antiScalpingLevel` | [EntitiesAntiScalpingLevel](EntitiesAntiScalpingLevel.md)
`appCode` | string
`attendeesBooked` | number
`bookingDuration` | number
`bookings` | [Array&lt;EntitiesBooking&gt;](EntitiesBooking.md)
`createdAt` | string
`deletedAt` | Date
`description` | string
`endDate` | string
`endTime` | string
`id` | string
`maxAttendees` | number
`ownerId` | string
`startDate` | string
`startTime` | string
`status` | [EntitiesAppointmentStatus](EntitiesAppointmentStatus.md)
`title` | string
`type` | [EntitiesAppointmentType](EntitiesAppointmentType.md)
`updatedAt` | string

## Example

```typescript
import type { EntitiesAppointment } from ''

// TODO: Update the object below with actual values
const example = {
  "antiScalpingLevel": null,
  "appCode": null,
  "attendeesBooked": null,
  "bookingDuration": null,
  "bookings": null,
  "createdAt": null,
  "deletedAt": null,
  "description": null,
  "endDate": null,
  "endTime": null,
  "id": null,
  "maxAttendees": null,
  "ownerId": null,
  "startDate": null,
  "startTime": null,
  "status": null,
  "title": null,
  "type": null,
  "updatedAt": null,
} satisfies EntitiesAppointment

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as EntitiesAppointment
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


