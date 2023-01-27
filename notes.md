
### How to compile and flash
1. Connect the pico as a storage device by holding down the button on the board, then plug the USB cable in
2. Flash and compile with the tinygo command
```bash
tinygo flash -target=pico
```

### Example of how a double click could be implemented
```go
if input14.Get() {
        // Record the timestamp for double trigger detection
        firstTrigger := time.Now()
        // Set the default control code
        var controlCode uint8 = 0x14

        // Wait for the switch do be released
        timeout := false
        for input14.Get() {
            time.Sleep(time.Millisecond * 10)
            if time.Now().Sub(firstTrigger) > time.Millisecond*500 {
                // if the switch is engaged longer than 500ms, then disregard the double
                timeout = true
                break
            }
        }
        // Check for second trigger if timeout has not occured
        for (time.Now().Sub(firstTrigger)) < (time.Millisecond*500) && !timeout {
            time.Sleep(time.Millisecond * 10)
            if input14.Get() {
                controlCode = 0x15
                break
            }
        }
        fmt.Println("Sending CC:" + strconv.Itoa(int(controlCode)))
        m.SendCC(0, 0, controlCode, 0x7f)
        time.Sleep(time.Millisecond * 500)
}
```

### How to monitor the output on the USB port
Put print statements in the code then use socat to read the output
Printf didn't work for me. Use Println
```go
fmt.println("Something happened")
```
```bash
socat stdio /dev/ttyACM0
```

EEPROM Storage
We need to store the following information

