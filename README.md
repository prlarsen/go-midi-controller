# go-midi-controller

The goal for this project is to create a MIDI footswitch controller using the Raspberry Pi Pico programmed in Go using the TinyGo compiler.
I built a 6 switch controller, but it can easily be expanded or reduced.  It uses an LCD to display labels for each switch. 
Switch functions are organized into switch banks.  The 6th switch is reserved for selecting the bank.  5 functions per bank.  10 functions if using hold switches. 
I'm using it with Helix midi codes to be compatible with my OpenRig system.  See [OpenRig](openrig.eu) for more details.

Note: eeprom.go is not used.  It's there for managing I2C EEPROM devices, which will be useful for building an app to change configurations over Bluetooth. 

### Switch Types
- Oneshot: Sends a single control code
- Hold: A hold switch can send 2 different control codes based on whether the switch is pressed quickly or held down
- Toggle: Toggles between on and off states for a given control code
- Program Change: Sends a program change sequence for selecting presets. 
--------------------------
### Parts/Materials
- Raspberry Pi Pico
- GeekPi I2C 1602 LCD Display
- Momentary FootSwitches
- Aluminum Chassis
- Dupont Cable Jumpers

### How to compile and flash
1. Connect the pico as a storage device by holding down the button on the board, then plug the USB cable in
2. Flash and compile with the tinygo command
```bash
tinygo flash -target=pico
```



### How to monitor the output on the USB port
Directly monitor control codes sent over USB
```bash
# Find port that pico is using
 aseqdump -l

Port    Client name                      Port name
 0:0    System                           Timer
 0:1    System                           Announce
14:0    Midi Through                     Midi Through Port-0
16:0    Ensoniq AudioPCI                 ES1371
20:0    Pico                             Pico MIDI 1
```
```bash
# Monitor events on the port
aseqdump -p 20:0

Waiting for data. Press Ctrl+C to end.
Source  Event                  Ch  Data
 20:0   Control change          0, controller 32, value 0
 20:0   Program change          0, program 0
 20:0   Control change          0, controller 32, value 0
```
#### Another way to do it
Put print statements in the code then use socat to read the output
Printf didn't work for me. Use Println
```go
fmt.println("Something happened")
```
```bash
socat stdio /dev/ttyACM0
```

