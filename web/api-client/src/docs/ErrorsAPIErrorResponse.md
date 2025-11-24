
# ErrorsAPIErrorResponse


## Properties

Name | Type
------------ | -------------
`code` | string
`fields` | [Array&lt;ErrorsFieldError&gt;](ErrorsFieldError.md)
`message` | string
`meta` | object
`requestId` | string
`status` | number

## Example

```typescript
import type { ErrorsAPIErrorResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "code": null,
  "fields": null,
  "message": null,
  "meta": null,
  "requestId": null,
  "status": null,
} satisfies ErrorsAPIErrorResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ErrorsAPIErrorResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


