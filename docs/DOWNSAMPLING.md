# Downsampling

Scrutiny collects alot of data, that can cause the database to grow unbounded. 

- Smart data
- Smart test data
- Temperature data
- Disk metrics (capacity/usage)
- etc

This data must be accurate in the short term, and is useful for doing trend analysis in the long term.
However, for trend analysis we only need aggregate data, individual data points are not as useful.

Scrutiny will automatically downsample data on a schedule to ensure that the database size stays reasonable, while still
ensuring historical data is present for comparisons.


| Bucket Name | Retention Period | Downsampling Range | Downsampling Aggregation Window | Downsampling Cron | Comments |
| --- | --- | --- | --- | --- | --- |
| `metrics` | 15 days | `-2w -1w` | `1w` | weekly on Sunday at 1:00am |
| `metrics_weekly` | 9 weeks | `-2mo -1mo` | `1mo` | monthly on first day of the month at 1:30am
| `metrics_monthly` | 25 months | `-2y -1y` | `1y` | yearly on the first day of the year at 2:00am
| `metrics_yearly` | forever | - | - | - | |


After 5 months, here's how may data points should exist in each bucket for one disk

| Bucket Name | Datapoints | Comments |
| --- | --- | --- |
| `metrics` | 15 | 7 daily datapoints , up to 7 pending data, 1 buffer data point |
| `metrics_weekly` | 9 | 4 aggregated weekly data points, 4 pending datapoints, 1 buffer data point |
| `metrics_monthly` | 3 | 3 aggregated monthly data points | 
| `metrics_yearly` | 0 | |

After 5 years, here's how may data points should exist in each bucket for one disk

| Bucket Name | Datapoints | Comments |
| --- | --- | --- |
| `metrics` | - | - |
| `metrics_weekly` | - | 
| `metrics_monthly` | - |
| `metrics_yearly` | - | 
