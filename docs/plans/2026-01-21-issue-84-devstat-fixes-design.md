# Design: Fix Device Statistics False Failures (Issue #84)

## Problem

ATA Device Statistics attributes added in v1.8.0 are causing false failure alerts:
- Some drives report impossibly high values (e.g., 420 billion mechanical failures)
- Device status is not updated when devstat attributes fail (missing propagation)
- Notification logic cannot handle string-based attribute IDs when filtering by critical
- No way for users to ignore problematic attributes

## Solution Overview

Four changes to fix the issue:

1. **Add device status propagation** to `ProcessAtaDeviceStatistics`
2. **Add value sanity checks** - mark values > 1 million as invalid
3. **Fix notify.go** to handle string-based devstat attribute IDs
4. **Add config-based attribute ignore lists** (UI tracked in #97)

## Detailed Design

### 1. Device Status Propagation

**File:** `webapp/backend/pkg/models/measurements/smart.go`

Add status propagation to `ProcessAtaDeviceStatistics` matching `ProcessAtaSmartInfo`:

```go
func (sm *Smart) ProcessAtaDeviceStatistics(cfg config.Interface, deviceStatistics collector.SmartInfo) {
    for _, page := range deviceStatistics.AtaDeviceStatistics.Pages {
        for _, stat := range page.Table {
            if !stat.Flags.Valid {
                continue
            }

            attrId := fmt.Sprintf("devstat_%d_%d", page.Number, stat.Offset)
            attrModel := SmartAtaDeviceStatAttribute{
                AttributeId: attrId,
                Value:       stat.Value,
            }

            attrModel.PopulateAttributeStatus()
            sm.Attributes[attrId] = &attrModel

            // NEW: Check ignore list
            var ignored bool
            if cfg != nil {
                ignoredAttrs := cfg.GetStringSlice("failures.ignored.devstat")
                for _, ignoredId := range ignoredAttrs {
                    if attrId == ignoredId {
                        ignored = true
                        break
                    }
                }
            }

            // NEW: Propagate failure status to device (matching ProcessAtaSmartInfo)
            if pkg.AttributeStatusHas(attrModel.Status, pkg.AttributeStatusFailedScrutiny) && !ignored {
                sm.Status = pkg.DeviceStatusSet(sm.Status, pkg.DeviceStatusFailedScrutiny)
            }
        }
    }
}
```

**Caller update:** Update call site in `FromCollectorSmartInfo` to pass `cfg`.

### 2. Value Sanity Checks

**File:** `webapp/backend/pkg/models/measurements/smart_ata_devstat_attribute.go`

Add a new status for invalid values and sanity checking:

```go
const (
    // Maximum reasonable value for failure count attributes
    // Any value above this is considered invalid/corrupted data
    MaxReasonableFailureCount = 1_000_000
)

func (sa *SmartAtaDeviceStatAttribute) PopulateAttributeStatus() *SmartAtaDeviceStatAttribute {
    if metadata, ok := thresholds.AtaDeviceStatsMetadata[sa.AttributeId]; ok {
        // Sanity check: reject impossibly high values for failure count attributes
        if metadata.Ideal == thresholds.ObservedThresholdIdealLow &&
           sa.Value > MaxReasonableFailureCount {
            sa.Status = pkg.AttributeStatusSet(sa.Status, pkg.AttributeStatusInvalidValue)
            sa.StatusReason = "Value exceeds reasonable maximum, likely corrupted data"
            sa.TransformedValue = sa.Value
            return sa
        }

        if metadata.Critical {
            threshold := int64(100)
            if sa.Threshold > 0 {
                threshold = sa.Threshold
            }

            if metadata.Ideal == thresholds.ObservedThresholdIdealLow && sa.Value >= threshold {
                sa.Status = pkg.AttributeStatusSet(sa.Status, pkg.AttributeStatusFailedScrutiny)
                sa.StatusReason = "Device statistic exceeds recommended threshold"
            }
        }
    }

    sa.TransformedValue = sa.Value
    return sa
}
```

**File:** `webapp/backend/pkg/constants.go`

Add new status constant:

```go
const (
    AttributeStatusPassed          AttributeStatus = 0
    AttributeStatusFailedSmart     AttributeStatus = 1
    AttributeStatusWarningScrutiny AttributeStatus = 2
    AttributeStatusFailedScrutiny  AttributeStatus = 4
    AttributeStatusInvalidValue    AttributeStatus = 8  // NEW
)
```

### 3. Fix notify.go for String-Based Attribute IDs

**File:** `webapp/backend/pkg/notify/notify.go`

Update the critical attribute check to handle devstat IDs:

```go
if statusFilterAttributes == pkg.MetricsStatusFilterAttributesCritical {
    critical := false
    if device.IsScsi() {
        critical = thresholds.ScsiMetadata[attrId].Critical
    } else if device.IsNvme() {
        critical = thresholds.NmveMetadata[attrId].Critical
    } else {
        // ATA: handle both numeric IDs and string-based devstat IDs
        if strings.HasPrefix(attrId, "devstat_") {
            if metadata, ok := thresholds.AtaDeviceStatsMetadata[attrId]; ok {
                critical = metadata.Critical
            }
        } else {
            attrIdInt, err := strconv.Atoi(attrId)
            if err != nil {
                continue
            }
            critical = thresholds.AtaMetadata[attrIdInt].Critical
        }
    }
    if !critical {
        continue
    }
}
```

### 4. Config-Based Attribute Ignore Lists

**File:** `example.scrutiny.yaml`

Add new config section:

```yaml
failures:
  # Existing transient config (unchanged)
  transient:
    ata: []

  # NEW: Permanently ignored attributes (not counted toward failure status)
  ignored:
    ata: []              # numeric SMART attribute IDs, e.g., [1, 9, 199]
    devstat: []          # string-based device statistics, e.g., ["devstat_1_40"]
    nvme: []             # NVME attribute IDs
    scsi: []             # SCSI attribute IDs
```

**Files to update:**
- `ProcessAtaSmartInfo` - check `failures.ignored.ata`
- `ProcessAtaDeviceStatistics` - check `failures.ignored.devstat`
- `ProcessNvmeSmartInfo` - check `failures.ignored.nvme`
- `ProcessScsiSmartInfo` - check `failures.ignored.scsi`

## Testing Strategy

1. **Unit tests** for sanity check logic with boundary values
2. **Unit tests** for notify.go devstat handling
3. **Integration test** with mock smartctl data containing high values
4. **Manual test** with config ignore list

## Files Changed

| File | Change |
|------|--------|
| `webapp/backend/pkg/constants.go` | Add `AttributeStatusInvalidValue` |
| `webapp/backend/pkg/models/measurements/smart.go` | Add status propagation, ignore list support |
| `webapp/backend/pkg/models/measurements/smart_ata_devstat_attribute.go` | Add sanity checks |
| `webapp/backend/pkg/notify/notify.go` | Handle string-based attribute IDs |
| `example.scrutiny.yaml` | Document ignore list config |

## Follow-up

- Issue #97: Add UI for configuring ignored attributes
