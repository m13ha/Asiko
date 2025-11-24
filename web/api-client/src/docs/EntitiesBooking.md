
# EntitiesBooking


## Properties

Name | Type
------------ | -------------
`appCode` | string
`appointmentId` | string
`attendeeCount` | number
`available` | boolean
`bookingCode` | string
`capacity` | number
`createdAt` | string
`date` | string
`deletedAt` | Date
`description` | string
`email` | string
`endTime` | string
`id` | string
`isSlot` | boolean
`name` | string
`notificationChannel` | string
`notificationStatus` | string
`phone` | string
`seatsBooked` | number
`startTime` | string
`status` | string
`updatedAt` | string
`userId` | string

## Example

```typescript
import type { EntitiesBooking } from ''

// TODO: Update the object below with actual values
const example = {
  "appCode": null,
  "appointmentId": null,
  "attendeeCount": null,
  "available": null,
  "bookingCode": null,
  "capacity": null,
  "createdAt": null,
  "date": null,
  "deletedAt": null,
  "description": null,
  "email": null,
  "endTime": null,
  "id": null,
  "isSlot": null,
  "name": null,
  "notificationChannel": null,
  "notificationStatus": null,
  "phone": null,
  "seatsBooked": null,
  "startTime": null,
  "status": null,
  "updatedAt": null,
  "userId": null,
} satisfies EntitiesBooking

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as EntitiesBooking
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


