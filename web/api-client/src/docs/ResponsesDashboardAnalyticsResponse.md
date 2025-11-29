
# ResponsesDashboardAnalyticsResponse


## Properties

Name | Type
------------ | -------------
`bookingsPerDay` | [Array&lt;ResponsesTimeSeriesPoint&gt;](ResponsesTimeSeriesPoint.md)
`cancellationsPerDay` | [Array&lt;ResponsesTimeSeriesPoint&gt;](ResponsesTimeSeriesPoint.md)
`endDate` | string
`startDate` | string
`totalAppointments` | number
`totalBookings` | number

## Example

```typescript
import type { ResponsesDashboardAnalyticsResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "bookingsPerDay": null,
  "cancellationsPerDay": null,
  "endDate": null,
  "startDate": null,
  "totalAppointments": null,
  "totalBookings": null,
} satisfies ResponsesDashboardAnalyticsResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ResponsesDashboardAnalyticsResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


