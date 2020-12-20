package main

import (
	"fmt"
	iot "github.com/ctlove0523/huaweicloud-iot-device-sdk-go"
	"time"
)

func main() {
	// 创建设备并初始化
	device := iot.CreateIotDevice("5fdb75cccbfe2f02ce81d4bf_go-mqtt", "123456789", "tcp://iot-mqtts.cn-north-4.myhuaweicloud.com:1883")
	device.Init()
	fmt.Printf("device connected: %v\n", device.IsConnected())

	// 注册平台设置属性callback,当应用通过API设置设备属性时，会调用此callback，支持注册多个callback
	device.AddPropertiesSetHandler(func(propertiesSetRequest iot.DevicePropertyDownRequest) bool {
		fmt.Println("I get property set command")
		fmt.Printf("request is %s", iot.Interface2JsonString(propertiesSetRequest))
		return true
	})

	// 注册平台查询设备属性callback，当平台查询设备属性时此callback被调用，仅支持设置一个callback
	device.SetPropertyQueryHandler(func(query iot.DevicePropertyQueryRequest) iot.ServicePropertyEntry {
		return iot.ServicePropertyEntry{
			ServiceId: "value",
			Properties: DemoProperties{
				Value:   "QUERY RESPONSE",
				MsgType: "query property",
			},
			EventTime: "2020-12-19 02:23:24",
		}
	})

	// 设备上报属性
	props := iot.ServicePropertyEntry{
		ServiceId: "value",
		EventTime: iot.DataCollectionTime(),
		Properties: DemoProperties{
			Value:   "chen tong",
			MsgType: "23",
		},
	}

	var content []iot.ServicePropertyEntry
	content = append(content, props)
	services := iot.ServiceProperty{
		Services: content,
	}
	device.ReportProperties(services)

	// 设备查询设备影子数据
	device.QueryDeviceShadow(iot.DevicePropertyQueryRequest{
		ServiceId: "value",
	}, func(response iot.DevicePropertyQueryResponse) {
		fmt.Printf("query device shadow success.\n,device shadow data is %s\n", iot.Interface2JsonString(response))
	})

	// 批量上报子设备属性
	subDevice1 := iot.DeviceService{
		DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-1",
		Services: content,
	}
	subDevice2 := iot.DeviceService{
		DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-2",
		Services: content,
	}

	subDevice3 := iot.DeviceService{
		DeviceId: "5fdb75cccbfe2f02ce81d4bf_sub-device-3",
		Services: content,
	}

	var devices []iot.DeviceService
	devices = append(devices, subDevice1, subDevice2, subDevice3)

	device.BatchReportSubDevicesProperties(iot.DevicesService{
		Devices: devices,
	})
	time.Sleep(1 * time.Minute)
}

type DemoProperties struct {
	Value   string `json:"value"`
	MsgType string `json:"msgType"`
}