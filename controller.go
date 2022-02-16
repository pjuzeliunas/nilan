package nilan

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/goburrow/modbus"
)

// Controller is used for communicating with Nilan CTS700 heatpump over
// Modbus TCP.
type Controller struct {
	Config Config
}

func (c *Controller) getHandler(slaveID byte) *modbus.TCPClientHandler {
	// Modbus TCP
	handler := modbus.NewTCPClientHandler(c.Config.NilanAddress)
	handler.Timeout = 10 * time.Second
	handler.SlaveId = slaveID
	err := handler.Connect()

	if err != nil {
		panic(err)
	}

	return handler
}

// FetchValue from register
func (c *Controller) FetchValue(slaveID byte, register Register) (uint16, error) {
	handler := c.getHandler(slaveID)
	defer handler.Close()
	client := modbus.NewClient(handler)
	resultBytes, error := client.ReadHoldingRegisters(uint16(register), 1)
	if error != nil {
		return 0, error
	}
	if len(resultBytes) == 2 {
		return binary.BigEndian.Uint16(resultBytes), nil
	} else {
		return 0, errors.New("cannot read register value")
	}
}

// FetchRegisterValues from slave
func (c *Controller) FetchRegisterValues(slaveID byte, registers []Register) (map[Register]uint16, error) {
	m := make(map[Register]uint16)

	handler := c.getHandler(slaveID)
	defer handler.Close()
	client := modbus.NewClient(handler)

	for _, register := range registers {
		resultBytes, err := client.ReadHoldingRegisters(uint16(register), 1)
		if err != nil {
			return m, err
		}
		if len(resultBytes) == 2 {
			resultWord := binary.BigEndian.Uint16(resultBytes)
			m[register] = resultWord
		} else {
			return m, errors.New("no result bytes")
		}
	}

	return m, nil
}

// SetRegisterValues on slave
func (c *Controller) SetRegisterValues(slaveID byte, values map[Register]uint16) error {
	handler := c.getHandler(slaveID)
	defer handler.Close()
	client := modbus.NewClient(handler)

	for register, value := range values {
		_, error := client.WriteSingleRegister(uint16(register), value)
		if error != nil {
			return error
		}
	}
	return nil
}

// Register is address of register on client
type Register uint16

const (
	// FanSpeedRegister is ID of register holding desired FanSpeed value
	FanSpeedRegister Register = 20148
	// DesiredRoomTemperatureRegister is ID of register holding desired room temperature in C times 10.
	// Example: 23.5 C is stored as 235.
	DesiredRoomTemperatureRegister Register = 20260
	// MasterTemperatureSensorSettingRegister is ID of register holding either 0 (read temperature from T3)
	// or 1 (read temperature from Text)
	MasterTemperatureSensorSettingRegister Register = 20263
	// T3ExtractAirTemperatureRegister is ID of register holding room temperature value when
	// MasterTemperatureSensorSettingRegister is 0
	T3ExtractAirTemperatureRegister Register = 20286
	// TextRoomTemperatureRegister is ID of register holding room temperature value when
	// MasterTemperatureSensorSettingRegister is 1
	TextRoomTemperatureRegister Register = 20280
	// OutdoorTemperatureRegister is ID of register outdoor temperature
	OutdoorTemperatureRegister Register = 20282
	// AverageHumidityRegister is ID of register holding average humidity value
	AverageHumidityRegister Register = 20164
	// ActualHumidityRegister is ID of register holding actual humidity value
	ActualHumidityRegister Register = 21776
	// DHWTopTankTemperatureRegister is ID of register holding T11 top DHW tank temperature
	DHWTopTankTemperatureRegister Register = 20520
	// DHWBottomTankTemperatureRegister is ID of register holding T11 bottom DHW tank temperature
	DHWBottomTankTemperatureRegister Register = 20522
	// DHWSetPointRegister is ID of register holding desired DHW temperature
	DHWSetPointRegister Register = 20460
	// DHWPauseRegister is ID of register holding DHW pause flag
	DHWPauseRegister Register = 20440
	// DHWPauseDurationRegister is ID of register holding DHW pause duration value
	DHWPauseDurationRegister Register = 20441
	// CentralHeatingPauseRegister is ID of register holding central heating pause flag
	CentralHeatingPauseRegister Register = 20600
	// CentralHeatingPauseDurationRegister is ID of register holding central heating pause duration value
	CentralHeatingPauseDurationRegister Register = 20601
	// CentralHeatingPowerRegister is ID of register holding On/Off value of central heating
	CentralHeatingPowerRegister Register = 20602
	// VentilationModeRegister is ID of register holding ventilation mode value (0, 1 or 2).
	VentilationModeRegister Register = 20120
	// VentilationPauseRegister is ID of register holding ventilation pause flag
	VentilationPauseRegister Register = 20100
	// SetpointSupplyTemperatureRegisterAIR9 is ID of register holding setpoint supply temperature
	// on AIR9 models
	SetpointSupplyTemperatureRegisterAIR9 Register = 20680
	// SetpointSupplyTemperatureRegisterGEO is ID of register holding setpoint supply temperature
	// on GEO models
	SetpointSupplyTemperatureRegisterGEO Register = 20640
	// DeviceTypeGEOReigister is ID of register that holds number 8 on GEO models
	DeviceTypeGEOReigister Register = 21839
	// DeviceTypeAIR9Register is ID of register that holds number 9 on AIR9 models
	DeviceTypeAIR9Register Register = 21899
	// T18ReadingRegisterGEO is ID of register holding T18 supply flow temperature reading
	// on GEO models
	T18ReadingRegisterGEO Register = 20653
	// T18ReadingRegisterAIR9 is ID of register holding T18 supply flow temperature reading
	// on AIR9 models
	T18ReadingRegisterAIR9 Register = 20686
	// EventOutdoorFilterWarningRegister is ID of register holding outdoor filter warning presence value
	EventOutdoorFilterWarningRegister Register = 22507
	// EventExtractFilterWarningRegister is ID of register holding extract filter warning presence value
	EventExtractFilterWarningRegister Register = 22508
	// EventHeaterOverHeatAlarmRegister is ID of register holding overheat alarm presence value
	EventHeaterOverHeatAlarmRegister Register = 22512
	// EventHeaterFrostWarningRegister is ID of register holding frost warning presence value
	EventHeaterFrostWarningRegister Register = 22514
	// EventHeaterFrostLongAlarmRegister is ID of register holding long frost alarm presence value
	EventHeaterFrostLongAlarmRegister Register = 22515
	// EventHeaterFrostAlarmRegister is ID of register holding frost alarm presence value
	EventHeaterFrostAlarmRegister Register = 22516
	// EventFireThermAlarmRegister is ID of register holding brandindgang activation status value
	EventFireThermAlarmRegister Register = 22521
	// EventKlixonWarningRegister is ID of register holding klixon warning presence value
	EventKlixonWarningRegister Register = 22578
	// EventCompressHighPressWarning is ID of register holding high presure warning presence value
	EventCompressHighPressWarning Register = 22579
)

type DeviceType int

const (
	DeviceTypeAir9 DeviceType = 0
	DeviceTypeGeo  DeviceType = 1
)

func (c *Controller) GetDeviceType() (DeviceType, error) {
	geoRegValue, e1 := c.FetchValue(4, DeviceTypeGEOReigister)
	if e1 != nil {
		return 0, e1
	}

	air9RegValue, e2 := c.FetchValue(4, DeviceTypeAIR9Register)
	if e2 != nil {
		return 0, e2
	}

	switch {
	case geoRegValue == 8:
		return DeviceTypeGeo, nil
	case air9RegValue == 9:
		return DeviceTypeAir9, nil
	default:
		return 0, errors.New("cannot determine device type")
	}
}

func (c *Controller) supplyFlowSetpointTemperatureRegister() (Register, error) {
	deviceType, err := c.GetDeviceType()
	if err != nil {
		return 0, err
	}

	switch deviceType {
	case DeviceTypeGeo:
		return SetpointSupplyTemperatureRegisterGEO, nil
	case DeviceTypeAir9:
		return SetpointSupplyTemperatureRegisterAIR9, nil
	default:
		return 0, errors.New("cannot determine supply flow setpoint register")
	}
}

// FetchSettings of Nilan
func (c *Controller) FetchSettings() (*Settings, error) {
	supplyTemperatureRegister, e1 := c.supplyFlowSetpointTemperatureRegister()
	if e1 != nil {
		return nil, e1
	}

	client1Registers := []Register{
		FanSpeedRegister,
		DesiredRoomTemperatureRegister,
		DHWSetPointRegister,
		DHWPauseRegister,
		DHWPauseDurationRegister,
		VentilationModeRegister,
		VentilationPauseRegister}
	client4Registers := []Register{
		CentralHeatingPauseRegister,
		CentralHeatingPauseDurationRegister,
		CentralHeatingPowerRegister,
		supplyTemperatureRegister}

	client1RegisterValues, e1 := c.FetchRegisterValues(1, client1Registers)
	if e1 != nil {
		return nil, e1
	}
	client4RegisterValues, e2 := c.FetchRegisterValues(4, client4Registers)
	if e2 != nil {
		return nil, e2
	}

	fanSpeed := new(FanSpeed)
	*fanSpeed = FanSpeed(client1RegisterValues[FanSpeedRegister])

	desiredRoomTemperature := new(int)
	*desiredRoomTemperature = int(client1RegisterValues[DesiredRoomTemperatureRegister])

	desiredDHWTemperature := new(int)
	*desiredDHWTemperature = int(client1RegisterValues[DHWSetPointRegister])

	dhwPaused := new(bool)
	*dhwPaused = client1RegisterValues[DHWPauseRegister] == 1

	dhwPauseDuration := new(int)
	*dhwPauseDuration = int(client1RegisterValues[DHWPauseDurationRegister])

	centralHeatingPaused := new(bool)
	*centralHeatingPaused = client4RegisterValues[CentralHeatingPauseRegister] == 1

	centralHeatingPauseDuration := new(int)
	*centralHeatingPauseDuration = int(client4RegisterValues[CentralHeatingPauseDurationRegister])

	centralHeatingOn := new(bool)
	*centralHeatingOn = client4RegisterValues[CentralHeatingPowerRegister] == 1

	ventilationMode := new(int)
	*ventilationMode = int(client1RegisterValues[VentilationModeRegister])

	ventilationPause := new(bool)
	*ventilationPause = client1RegisterValues[VentilationPauseRegister] == 1

	setpointTemperature := new(int)
	*setpointTemperature = int(client4RegisterValues[supplyTemperatureRegister])

	settings := &Settings{FanSpeed: fanSpeed,
		DesiredRoomTemperature:      desiredRoomTemperature,
		DesiredDHWTemperature:       desiredDHWTemperature,
		DHWProductionPaused:         dhwPaused,
		DHWProductionPauseDuration:  dhwPauseDuration,
		CentralHeatingPaused:        centralHeatingPaused,
		CentralHeatingPauseDuration: centralHeatingPauseDuration,
		CentralHeatingIsOn:          centralHeatingOn,
		VentilationMode:             ventilationMode,
		VentilationOnPause:          ventilationPause,
		SetpointSupplyTemperature:   setpointTemperature}

	return settings, nil
}

// SendSettings of Nilan
func (c *Controller) SendSettings(settings Settings) error {
	settingsStr := spew.Sprintf("%+v", settings)
	log.Printf("Sending new settings to Nialn (<nil> values will be ignored): %+v\n", settingsStr)
	client1RegisterValues := make(map[Register]uint16)
	client4RegisterValues := make(map[Register]uint16)

	if settings.FanSpeed != nil {
		fanSpeed := uint16(*settings.FanSpeed)
		client1RegisterValues[FanSpeedRegister] = fanSpeed
	}

	if settings.DesiredRoomTemperature != nil {
		desiredRoomTemperature := uint16(*settings.DesiredRoomTemperature)
		client1RegisterValues[DesiredRoomTemperatureRegister] = desiredRoomTemperature
	}

	if settings.DesiredDHWTemperature != nil {
		desiredDHWTemperature := uint16(*settings.DesiredDHWTemperature)
		client1RegisterValues[DHWSetPointRegister] = desiredDHWTemperature
	}

	if settings.DHWProductionPaused != nil {
		if *settings.DHWProductionPaused {
			client1RegisterValues[DHWPauseRegister] = uint16(1)
		} else {
			client1RegisterValues[DHWPauseRegister] = uint16(0)
		}
	}

	if settings.DHWProductionPauseDuration != nil {
		pauseDuration := uint16(*settings.DHWProductionPauseDuration)
		client1RegisterValues[DHWPauseDurationRegister] = pauseDuration
	}

	if settings.CentralHeatingPaused != nil {
		if *settings.CentralHeatingPaused {
			client4RegisterValues[CentralHeatingPauseRegister] = uint16(1)
		} else {
			client4RegisterValues[CentralHeatingPauseRegister] = uint16(0)
		}
	}

	if settings.CentralHeatingPauseDuration != nil {
		pauseDuration := uint16(*settings.CentralHeatingPauseDuration)
		client4RegisterValues[CentralHeatingPauseDurationRegister] = pauseDuration
	}

	if settings.VentilationMode != nil {
		ventilationMode := *settings.VentilationMode
		if ventilationMode != 0 && ventilationMode != 1 && ventilationMode != 2 {
			return errors.New("unsupported VentilationMode value")
		}
		ventilationModeVal := uint16(ventilationMode)
		client1RegisterValues[VentilationModeRegister] = ventilationModeVal
	}

	if settings.VentilationOnPause != nil {
		if *settings.VentilationOnPause {
			client1RegisterValues[VentilationPauseRegister] = uint16(1)
		} else {
			client1RegisterValues[VentilationPauseRegister] = uint16(0)
		}
	}

	if settings.SetpointSupplyTemperature != nil {
		setpointTempeature := uint16(*settings.SetpointSupplyTemperature)
		client4RegisterValues[SetpointSupplyTemperatureRegisterAIR9] = setpointTempeature
		client4RegisterValues[SetpointSupplyTemperatureRegisterGEO] = setpointTempeature
	}

	e1 := c.SetRegisterValues(1, client1RegisterValues)
	if e1 != nil {
		return e1
	}

	e2 := c.SetRegisterValues(4, client4RegisterValues)
	if e2 != nil {
		return e2
	}

	return nil
}

func (c *Controller) roomTemperatureRegister() (Register, error) {
	masterSensorSetting, error := c.FetchValue(1, MasterTemperatureSensorSettingRegister)
	if error != nil {
		return 0, error
	}
	if masterSensorSetting == 0 {
		return T3ExtractAirTemperatureRegister, nil
	} else {
		return TextRoomTemperatureRegister, nil
	}
}

func (c *Controller) t18ReadingRegister() (Register, error) {
	deviceType, err := c.GetDeviceType()
	if err != nil {
		return 0, err
	}

	switch deviceType {
	case DeviceTypeGeo:
		return T18ReadingRegisterGEO, nil
	case DeviceTypeAir9:
		return T18ReadingRegisterAIR9, nil
	default:
		return 0, errors.New("cannot determine T18 reading register")
	}
}

// FetchReadings of Nilan sensors
func (c *Controller) FetchReadings() (*Readings, error) {
	roomTemperatureRegister, e1 := c.roomTemperatureRegister()
	if e1 != nil {
		return nil, e1
	}

	t18Register, e2 := c.t18ReadingRegister()
	if e2 != nil {
		return nil, e2
	}

	client1Registers := []Register{roomTemperatureRegister,
		OutdoorTemperatureRegister,
		AverageHumidityRegister,
		ActualHumidityRegister,
		DHWTopTankTemperatureRegister,
		DHWBottomTankTemperatureRegister}

	client4Registers := []Register{t18Register}

	client1ReadingsRaw, e1 := c.FetchRegisterValues(1, client1Registers)
	if e1 != nil {
		return nil, e1
	}
	client4ReadingsRaw, e2 := c.FetchRegisterValues(4, client4Registers)
	if e2 != nil {
		return nil, e2
	}

	roomTemperature := int(client1ReadingsRaw[roomTemperatureRegister])
	outdoorTemperature := int(client1ReadingsRaw[OutdoorTemperatureRegister])
	averageHumidity := int(client1ReadingsRaw[AverageHumidityRegister])
	actualHumidity := int(client1ReadingsRaw[ActualHumidityRegister])
	dhwTopTemperature := int(client1ReadingsRaw[DHWTopTankTemperatureRegister])
	dhwBottomTemperature := int(client1ReadingsRaw[DHWBottomTankTemperatureRegister])
	supplyFlowTemperature := int(client4ReadingsRaw[t18Register])

	readings := &Readings{
		RoomTemperature:          roomTemperature,
		OutdoorTemperature:       outdoorTemperature,
		AverageHumidity:          averageHumidity,
		ActualHumidity:           actualHumidity,
		DHWTankTopTemperature:    dhwTopTemperature,
		DHWTankBottomTemperature: dhwBottomTemperature,
		SupplyFlowTemperature:    supplyFlowTemperature}

	if readings.AverageHumidity == 0 {
		fmt.Println("what?")
	}

	return readings, nil
}

func (c *Controller) FetchErrors() (*Errors, error) {
	registers := []Register{
		EventOutdoorFilterWarningRegister,
		EventExtractFilterWarningRegister,
		EventHeaterOverHeatAlarmRegister,
		EventHeaterFrostWarningRegister,
		EventHeaterFrostLongAlarmRegister,
		EventHeaterFrostAlarmRegister,
		EventFireThermAlarmRegister,
		EventKlixonWarningRegister,
		EventCompressHighPressWarning,
	}
	readings, err := c.FetchRegisterValues(1, registers)
	if err != nil {
		return nil, err
	}

	outdoorFilterWarn := int(readings[EventOutdoorFilterWarningRegister]) == 1
	extractFilterWarn := int(readings[EventExtractFilterWarningRegister]) == 1
	heaterOverheatAlarm := int(readings[EventHeaterOverHeatAlarmRegister]) == 1
	heaterFrostWarning := int(readings[EventHeaterFrostWarningRegister]) == 1
	heaterLongFrostAlarm := int(readings[EventHeaterFrostLongAlarmRegister]) == 1
	heaterFrostAlarm := int(readings[EventHeaterFrostAlarmRegister]) == 1
	fireThermAlarm := int(readings[EventFireThermAlarmRegister]) == 1
	klixonWarn := int(readings[EventKlixonWarningRegister]) == 1
	highPressureWarn := int(readings[EventCompressHighPressWarning]) == 1

	oldFilterWarning := outdoorFilterWarn || extractFilterWarn
	otherErrors := heaterOverheatAlarm || heaterFrostWarning || heaterLongFrostAlarm || heaterFrostAlarm || fireThermAlarm || klixonWarn || highPressureWarn

	errors := Errors{
		OldFilterWarning: oldFilterWarning,
		OtherErrors:      otherErrors,
	}

	return &errors, nil
}
