// Copyright © 2018 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oracle

import (
	"context"
	"fmt"

	"github.com/banzaicloud/productinfo/logger"

	"github.com/banzaicloud/productinfo/pkg/productinfo"
	"github.com/banzaicloud/productinfo/pkg/productinfo/oracle/client"
)

// Infoer encapsulates the data and operations needed to access external resources
type Infoer struct {
	client           *client.OCI
	shapeSpecs       map[string]ShapeSpecs
	productInfoCache map[string]ITRAProductInfo
}

// ShapeSpecs representation the specs of a certain type of virtual machine
type ShapeSpecs struct {
	PartNumber string
	Cpus       float64 `json:"cpusPerVm"`
	Mem        float64 `json:"memPerVm"`
	NtwPerf    string  `json:"NtwPerf"`
}

const (
	cpu    = "cpu"
	memory = "memory"
)

var regionNames = map[string]string{
	"uk-london-1":    "EU (London)",
	"eu-frankfurt-1": "EU (Frankfurt)",
	"us-ashburn-1":   "US East (Ashburn)",
	"us-phoenix-1":   "US West (Phoenix)",
}

var shapeSpecs = map[string]ShapeSpecs{
	"VM.Standard1.1":  ShapeSpecs{PartNumber: "B88317", Mem: 7, Cpus: 1, NtwPerf: "0.6 Gbps"},
	"VM.Standard2.1":  ShapeSpecs{PartNumber: "B88514", Mem: 15, Cpus: 1, NtwPerf: "1 Gbps"},
	"VM.Standard1.2":  ShapeSpecs{PartNumber: "B88317", Mem: 14, Cpus: 2, NtwPerf: "1.2 Gbps"},
	"VM.Standard2.2":  ShapeSpecs{PartNumber: "B88514", Mem: 30, Cpus: 2, NtwPerf: "2 Gbps"},
	"VM.Standard1.4":  ShapeSpecs{PartNumber: "B88317", Mem: 28, Cpus: 4, NtwPerf: "1.2 Gbps"},
	"VM.Standard2.4":  ShapeSpecs{PartNumber: "B88514", Mem: 60, Cpus: 4, NtwPerf: "4.1 Gbps"},
	"VM.Standard1.8":  ShapeSpecs{PartNumber: "B88317", Mem: 56, Cpus: 8, NtwPerf: "2.4 Gbps"},
	"VM.Standard2.8":  ShapeSpecs{PartNumber: "B88514", Mem: 120, Cpus: 8, NtwPerf: "8.2 Gbps"},
	"VM.Standard1.16": ShapeSpecs{PartNumber: "B88317", Mem: 112, Cpus: 16, NtwPerf: "4.8 Gbps"},
	"VM.Standard2.16": ShapeSpecs{PartNumber: "B88514", Mem: 240, Cpus: 16, NtwPerf: "16.4 Gbps"},
	"VM.Standard2.24": ShapeSpecs{PartNumber: "B88514", Mem: 320, Cpus: 24, NtwPerf: "24.6 Gbps"},
	"VM.DenseIO1.4":   ShapeSpecs{PartNumber: "B88316", Mem: 60, Cpus: 4, NtwPerf: "1.2 Gbps"},
	"VM.DenseIO1.8":   ShapeSpecs{PartNumber: "B88316", Mem: 60, Cpus: 8, NtwPerf: "2.4 Gbps"},
	"VM.DenseIO2.8":   ShapeSpecs{PartNumber: "B88516", Mem: 120, Cpus: 8, NtwPerf: "8.2 Gbps"},
	"VM.DenseIO1.16":  ShapeSpecs{PartNumber: "B88316", Mem: 120, Cpus: 16, NtwPerf: "4.8 Gbps"},
	"VM.DenseIO2.16":  ShapeSpecs{PartNumber: "B88516", Mem: 240, Cpus: 16, NtwPerf: "16.4 Gbps"},
	"VM.DenseIO2.24":  ShapeSpecs{PartNumber: "B88516", Mem: 320, Cpus: 24, NtwPerf: "24.6 Gbps"},
}

// NewInfoer creates a new instance of the infoer
func NewInfoer() (*Infoer, error) {

	oci, err := client.NewOCI()
	if err != nil {
		return nil, err
	}

	return &Infoer{
		client:     oci,
		shapeSpecs: shapeSpecs,
	}, nil
}

// Initialize downloads and parses the SKU list of the Compute Engine service
func (i *Infoer) Initialize(ctx context.Context) (prices map[string]map[string]productinfo.Price, err error) {
	log := logger.Extract(ctx)

	log.Info("initializing price info")

	prices = make(map[string]map[string]productinfo.Price)

	zonesInRegions := make(map[string][]string)
	regions, err := i.GetRegions(ctx)
	if err != nil {
		return nil, err
	}

	for r := range regions {
		zones, err := i.GetZones(ctx, r)
		if err != nil {
			return nil, err
		}
		zonesInRegions[r] = zones
	}

	shapePrices, err := i.GetProductPrices(ctx)
	if err != nil {
		return nil, err
	}

	for region := range regions {
		products, err := i.GetProducts(ctx, region)
		if err != nil {
			return prices, err
		}

		if prices[region] == nil {
			prices[region] = make(map[string]productinfo.Price)
		}
		for _, product := range products {

			shapePrice := shapePrices[product.Type]

			price := prices[region][product.Type]
			price.OnDemandPrice = shapePrice
			prices[region][product.Type] = price
			log.WithField("region", region).Debugf("price info added: [machinetype=%s, price=%v]", product.Type, price)
		}
	}

	log.Debugf("queried zones and regions: %v", zonesInRegions)

	return
}

// GetAttributeValues gets the AttributeValues for the given attribute name
func (i *Infoer) GetAttributeValues(ctx context.Context, attribute string) (values productinfo.AttrValues, err error) {
	log := logger.Extract(ctx)

	log.Debugf("getting %s values", attribute)

	values = make(productinfo.AttrValues, 0)
	uniquemap := make(map[float64]bool)

	shapesInRegions, err := i.client.GetSupportedShapes()
	if err != nil {
		return
	}

	var shapes []string
	for _, shapes = range shapesInRegions {
		for _, shape := range shapes {
			if _, ok := i.shapeSpecs[shape]; !ok {
				continue
			}
			specs := i.shapeSpecs[shape]
			var attr productinfo.AttrValue
			switch attribute {
			case cpu:
				attr = productinfo.AttrValue{
					Value:    specs.Cpus,
					StrValue: fmt.Sprintf("%v", specs.Cpus),
				}
			case memory:
				attr = productinfo.AttrValue{
					Value:    specs.Mem,
					StrValue: fmt.Sprintf("%v", specs.Mem),
				}
			}
			if _, ok := uniquemap[attr.Value]; !ok {
				values = append(values, attr)
				uniquemap[attr.Value] = true
			}
		}
	}

	log.Debugf("found %s values: %v", attribute, values)

	return values, nil
}

// GetCurrentPrices retrieves all the spot prices in a region
func (i *Infoer) GetCurrentPrices(ctx context.Context, region string) (prices map[string]productinfo.Price, err error) {
	return nil, fmt.Errorf("oracle prices cannot be queried on the fly")
}

// GetMemoryAttrName returns the provider representation of the memory attribute
func (i *Infoer) GetMemoryAttrName() string {
	return memory
}

// GetCpuAttrName returns the provider representation of the cpu attribute
func (i *Infoer) GetCpuAttrName() string {
	return cpu
}

// GetNetworkPerformanceMapper returns the network performance mappier implementation for this provider
func (i *Infoer) GetNetworkPerformanceMapper() (mapper productinfo.NetworkPerfMapper, err error) {
	return newNetworkMapper(), nil
}

// GetProductPrices gets prices for available shapes from ITRA
func (i *Infoer) GetProductPrices(ctx context.Context) (prices map[string]float64, err error) {

	prices = make(map[string]float64, 0)
	for shape, specs := range i.shapeSpecs {
		info, _ := i.GetProductInfoFromITRA(ctx, specs.PartNumber)
		prices[shape] = info.GetPrice("PAY_AS_YOU_GO") * specs.Cpus
	}

	return
}

// GetProducts retrieves the available virtual machines types in a region
func (i *Infoer) GetProducts(ctx context.Context, regionId string) (products []productinfo.VmInfo, err error) {

	err = i.client.ChangeRegion(regionId)
	if err != nil {
		return
	}

	shapes, err := i.client.GetSupportedShapesInARegion(regionId)
	if err != nil {
		return
	}

	products = make([]productinfo.VmInfo, 0)
	for _, shape := range shapes {
		s := i.shapeSpecs[shape]
		products = append(products, productinfo.VmInfo{
			Type:    shape,
			NtwPerf: s.NtwPerf,
			Cpus:    s.Cpus,
			Mem:     s.Mem,
		})
	}

	return
}

// GetRegions returns a map with available regions
func (i *Infoer) GetRegions(ctx context.Context) (regions map[string]string, err error) {
	logger.Extract(ctx).Debugf("getting regions")

	c, err := i.client.NewIdentityClient()
	if err != nil {
		return
	}

	_regions, err := c.GetSubscribedRegionNames()
	if err != nil {
		return
	}

	regions = make(map[string]string, 0)
	for _, region := range _regions {
		description := region
		if displayName, ok := regionNames[region]; ok {
			description = displayName
		}
		regions[region] = description
	}

	return
}

// GetZones returns the availability zones in a region
func (i *Infoer) GetZones(ctx context.Context, region string) (zones []string, err error) {
	logger.Extract(ctx).Debug("getting zones")

	err = i.client.ChangeRegion(region)
	if err != nil {
		return
	}

	c, err := i.client.NewIdentityClient()
	if err != nil {
		return
	}

	ads, err := c.GetAvailabilityDomains()
	if err != nil {
		return
	}

	for _, ad := range ads {
		zones = append(zones, *ad.Name)
	}

	return
}

// HasShortLivedPriceInfo - Oracle doesn't have preemptible instances
func (i *Infoer) HasShortLivedPriceInfo() bool {
	return false
}

// GetServices returns the available services on the  given region
func (i *Infoer) GetServices() ([]productinfo.ServiceDescriber, error) {
	services := []productinfo.ServiceDescriber{
		productinfo.NewService("compute"),
		productinfo.NewService("oke")}
	return services, nil
}

// GetService returns the service on the  provider
func (i *Infoer) GetService(ctx context.Context, service string) (productinfo.ServiceDescriber, error) {
	svcs, err := i.GetServices()
	if err != nil {
		return nil, err
	}
	for _, sd := range svcs {
		if service == sd.ServiceName() {
			logger.Extract(ctx).Debugf("found service: %s", service)
			return sd, nil
		}
	}
	return nil, fmt.Errorf("the service [%s] is not supported", service)
}

// GetServiceImages retrieves the images supported by the given service in the given region
func (i *Infoer) GetServiceImages(region, service string) ([]productinfo.ImageDescriber, error) {
	return nil, fmt.Errorf("GetServiceImages - not yet implemented")
}

// GetServiceProducts retrieves the products supported by the given service in the given region
func (i *Infoer) GetServiceProducts(region, service string) ([]productinfo.ProductDetails, error) {
	return nil, fmt.Errorf("GetServiceProducts - not yet implemented")
}

// GetServiceAttributes retrieves the attribute values supported by the given service in the given region for the given attribute
func (i *Infoer) GetServiceAttributes(region, service, attribute string) (productinfo.AttrValues, error) {
	return nil, fmt.Errorf("GetServiceAttributes - not yet implemented")
}