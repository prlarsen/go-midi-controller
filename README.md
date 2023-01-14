# go-midi-controller

# WORK IN PROGRESS

The goal for this project is to create a MIDI footswitch controller using the Raspberry Pi Pico programmed in Go using the TinyGo compiler.  I plan to use this controller to perform tasks in my DAW.  I'm aiming to recreate footswitchable functionality similar to what hardware loopers provide. I'm currently using Bitwig for a DAW because the looping functionality is really good. 

This requires the controller to send out MIDI control codes and the DAW to be able to assign those codes to functions in the DAW.

What I'm thinking I need it to do:
- Play/Stop
- Record/Delete
- Next/Previous Track
- Mute/Solo Track
- Next/Previous Scene
- Tap Tempo

To keep the switch count and physical footprint down, I will likely make the switches dual function, either by double pressing the switch or short press/long press

The controller only needs to send out control codes and the rest is configured in the DAW.

I'm building this fully aware TinyGo is less mature and may be unstable compared to C/Python.

|Control Codes  |Function           |Switch|
|---------------|-------------------|------ |  
|20 (0x14)      |Next Track         |1      |
|21 (0x15)      |Previous Track     |1 Hold |
|22 (0x16)      |Mute Selected Track|2      | 
|23 (0x17)      |Solo Selected Track|2 Hold |
|24 (0x18)      |Tap Tempo          |3      |
|25 (0x19)      |Unassigned         |       | 
|26 (0x1a)      |Play               |4      |
|27 (0x1b)      |Stop               |4 Hold |
|28 (0x1c)      |Record             |5      |
|29 (0x1d)      |Delete Recording   |5 Hold |
|30 (0x1e)      |Arm Record         |6      |

