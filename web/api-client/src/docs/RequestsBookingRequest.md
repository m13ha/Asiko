
# RequestsBookingRequest


## Properties

Name | Type
------------ | -------------
`appCode` | string
`attendeeCount` | number
`date` | string
`description` | string
`deviceToken` | string
`email` | string
`endTime` | string
`name` | string
`phone` | string
`startTime` | string

## Example

```typescript
import type { RequestsBookingRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "appCode": null,
  "attendeeCount": null,
  "date": null,
  "description": null,
  "deviceToken": null,
  "email": null,
  "endTime": null,
  "name": null,
  "phone": null,
  "startTime": null,
} satisfies RequestsBookingRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as RequestsBookingRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


