
# EntitiesNotification


## Properties

Name | Type
------------ | -------------
`createdAt` | string
`eventType` | string
`id` | string
`isRead` | boolean
`message` | string
`resourceId` | string
`userId` | string

## Example

```typescript
import type { EntitiesNotification } from ''

// TODO: Update the object below with actual values
const example = {
  "createdAt": null,
  "eventType": null,
  "id": null,
  "isRead": null,
  "message": null,
  "resourceId": null,
  "userId": null,
} satisfies EntitiesNotification

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as EntitiesNotification
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


