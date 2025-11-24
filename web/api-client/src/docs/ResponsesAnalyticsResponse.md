
# ResponsesAnalyticsResponse


## Properties

Name | Type
------------ | -------------
`appointmentsByType` | { [key: string]: number; }
`avgAttendeesPerBooking` | number
`avgLeadTimeHours` | number
`bookingsByStatus` | { [key: string]: number; }
`bookingsPerDay` | [Array&lt;ResponsesTimeSeriesPoint&gt;](ResponsesTimeSeriesPoint.md)
`cancellationsPerDay` | [Array&lt;ResponsesTimeSeriesPoint&gt;](ResponsesTimeSeriesPoint.md)
`distinctCustomers` | number
`endDate` | string
`guestVsRegistered` | { [key: string]: number; }
`medianLeadTimeHours` | number
`partyCapacity` | [ResponsesAnalyticsResponsePartyCapacity](ResponsesAnalyticsResponsePartyCapacity.md)
`peakDays` | [Array&lt;ResponsesBucketCount&gt;](ResponsesBucketCount.md)
`peakHours` | [Array&lt;ResponsesBucketCount&gt;](ResponsesBucketCount.md)
`rejectionsPerDay` | [Array&lt;ResponsesTimeSeriesPoint&gt;](ResponsesTimeSeriesPoint.md)
`repeatCustomers` | number
`slotUtilizationPercent` | number
`startDate` | string
`topAppointments` | [Array&lt;ResponsesTopAppointment&gt;](ResponsesTopAppointment.md)
`totalAppointments` | number
`totalBookings` | number

## Example

```typescript
import type { ResponsesAnalyticsResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "appointmentsByType": null,
  "avgAttendeesPerBooking": null,
  "avgLeadTimeHours": null,
  "bookingsByStatus": null,
  "bookingsPerDay": null,
  "cancellationsPerDay": null,
  "distinctCustomers": null,
  "endDate": null,
  "guestVsRegistered": null,
  "medianLeadTimeHours": null,
  "partyCapacity": null,
  "peakDays": null,
  "peakHours": null,
  "rejectionsPerDay": null,
  "repeatCustomers": null,
  "slotUtilizationPercent": null,
  "startDate": null,
  "topAppointments": null,
  "totalAppointments": null,
  "totalBookings": null,
} satisfies ResponsesAnalyticsResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ResponsesAnalyticsResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


