// +build !linux
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"github.com/StackExchange/wmi"
)

const wmqlProcessor = "SELECT Manufacturer, Name, NumberOfLogicalProcessors, NumberOfCores FROM Win32_Processor"

type win32Processor struct {
	Manufacturer              *string
	Name                      *string
	NumberOfLogicalProcessors uint32
	NumberOfCores             uint32
}

func (ctx *context) cpuFillInfo(info *CPUInfo) error {
	// Getting info from WMI
	var win32descriptions []win32Processor
	if err := wmi.Query(wmqlProcessor, &win32descriptions); err != nil {
		return err
	}
	// Converting into standard structures
	info.Processors = ctx.processorsGet(win32descriptions)
	var totCores uint32
	var totThreads uint32
	for _, p := range info.Processors {
		totCores += p.NumCores
		totThreads += p.NumThreads
	}
	info.TotalCores = totCores
	info.TotalThreads = totThreads
	return nil
}

func (ctx *context) processorsGet(win32descriptions []win32Processor) []*Processor {
	var procs []*Processor
	// Converting into standard structures
	for index, description := range win32descriptions {
		p := &Processor{
			Id:         index, // TODO: how to get a decent "Physical ID" to use ?
			Model:      *description.Name,
			Vendor:     *description.Manufacturer,
			NumCores:   description.NumberOfCores,
			NumThreads: description.NumberOfLogicalProcessors,
		}
		procs = append(procs, p)
	}
	return procs
}
