# AnalyticsApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**getDashboardAnalytics**](AnalyticsApi.md#getdashboardanalytics) | **GET** /analytics/dashboard | Get user dashboard analytics |
| [**getUserAnalytics**](AnalyticsApi.md#getuseranalytics) | **GET** /analytics | Get user analytics |



## getDashboardAnalytics

> ResponsesDashboardAnalyticsResponse getDashboardAnalytics(startDate, endDate)

Get user dashboard analytics

Get minimal analytics for dashboard display. Includes totals and daily series only.

### Example

```ts
import {
  Configuration,
  AnalyticsApi,
} from '';
import type { GetDashboardAnalyticsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new AnalyticsApi(config);

  const body = {
    // string | Start date (YYYY-MM-DD)
    startDate: startDate_example,
    // string | End date (YYYY-MM-DD)
    endDate: endDate_example,
  } satisfies GetDashboardAnalyticsRequest;

  try {
    const data = await api.getDashboardAnalytics(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **startDate** | `string` | Start date (YYYY-MM-DD) | [Defaults to `undefined`] |
| **endDate** | `string` | End date (YYYY-MM-DD) | [Defaults to `undefined`] |

### Return type

[**ResponsesDashboardAnalyticsResponse**](ResponsesDashboardAnalyticsResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Invalid date format or missing parameters |  -  |
| **401** | Unauthorized |  -  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getUserAnalytics

> ResponsesAnalyticsResponse getUserAnalytics(startDate, endDate)

Get user analytics

Get analytics for the authenticated user over a date window. Includes totals, breakdowns (by type/status, guest vs registered), utilization, lead-time stats, daily series, peak hours/days, and top appointments.

### Example

```ts
import {
  Configuration,
  AnalyticsApi,
} from '';
import type { GetUserAnalyticsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new AnalyticsApi(config);

  const body = {
    // string | Start date (YYYY-MM-DD)
    startDate: startDate_example,
    // string | End date (YYYY-MM-DD)
    endDate: endDate_example,
  } satisfies GetUserAnalyticsRequest;

  try {
    const data = await api.getUserAnalytics(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **startDate** | `string` | Start date (YYYY-MM-DD) | [Defaults to `undefined`] |
| **endDate** | `string` | End date (YYYY-MM-DD) | [Defaults to `undefined`] |

### Return type

[**ResponsesAnalyticsResponse**](ResponsesAnalyticsResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | Invalid date format or missing parameters |  -  |
| **401** | Unauthorized |  -  |
| **500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

