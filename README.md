# go-midi-controller

# WORK IN PROGRESS

The goal for this project is to create a MIDI footswitch controller using the Raspberry Pi Pico programmed in Go using the TinyGo compiler.  I plan to use this controller to perform tasks in my DAW.  I'm aiming to recreate footswitchable functionality similar to what hardware loopers provide. I'm currently using Bitwig for a DAW because the looping functionality is really good. 

This requires the controller to send out MIDI control codes and the DAW to be able to assign those codes to functions in the DAW.

The controller only needs to send out control codes and the rest is configured in the DAW.

I'm building this fully aware TinyGo is less mature and may be unstable compared to C/C++.  I like Go :)

What I'm thinking I need it to do:
- Play/Stop
- Record/Delete
- Next/Previous Track
- Mute/Solo Track
- Next/Previous Scene
- Tap Tempo

--------------------------
### **Switch Input Modes**

Switches 1-4 are configured as SWITCH_HOLD which will send a different CC when the switch is pressed and when it is held

Switch 5 is configured as SWITCH_BANKSELECT which will select the CC banks listed below

Switch 6 is configured as SWITCH_ONESHOT to be used for tap tempo or any other function that might not play nice when timing is important

--------------------------
### **Control Code Assignment**


|Control Codes  |Bank Selction      |Switch|
|---------------|-------------------|------ |  
|20 (0x14)      |1                  |1      |
|21 (0x15)      |1                  |1 Hold |
|22 (0x16)      |1                  |2      | 
|23 (0x17)      |1                  |2 Hold |
|24 (0x18)      |1                  |3      |
|25 (0x19)      |1                  |3 Hold | 
|26 (0x1a)      |1                  |4      |
|27 (0x1b)      |1                  |4 Hold |
|28 (0x1c)      |1                  |5      |
|29 (0x1d)      |1                  |5 Hold |
|30 (0x1e)      |1                  |6      |
---------------------------------------------
|Control Codes  |Bank Selction      |Switch |
|---------------|-------------------|------ |  
|31 (0x14)      |2                  |1      |
|32 (0x15)      |2                  |1 Hold |
|33 (0x16)      |2                  |2      | 
|34 (0x17)      |2                  |2 Hold |
|35 (0x18)      |2                  |3      |
|36 (0x19)      |2                  |3 Hold | 
|37 (0x1a)      |2                  |4      |
|38 (0x1b)      |2                  |4 Hold |
|39 (0x1c)      |2                  |5      |
|40 (0x1d)      |2                  |5 Hold |
|41 (0x1e)      |2                  |6      |
---------------------------------------------
|Control Codes  |Bank Selction      |Switch |
|---------------|-------------------|------ |  
|42 (0x14)      |3                  |1      |
|43 (0x15)      |3                  |1 Hold |
|44 (0x16)      |3                  |2      | 
|45 (0x17)      |3                  |2 Hold |
|46 (0x18)      |3                  |3      |
|47 (0x19)      |3                  |3 Hold | 
|48 (0x1a)      |3                  |4      |
|49 (0x1b)      |3                  |4 Hold |
|50 (0x1c)      |3                  |5      |
|51 (0x1d)      |3                  |5 Hold |
|52 (0x1e)      |3                  |6      |

