# Dell Latitude E7240

This page is about the notebook [Dell Latitude E7240].

## Release status

Dell Latitude E7240 was released in 2013 and is now end of life.
It can be bought from a secondhand market like Taobao or eBay.

## Required proprietary blobs

The following blobs are required to operate the hardware:
1. mrc.bin
2. Intel ME firmware

Memory reference code in mrc.bin is used to initialize the Haswell platform.
You need this blob to build a working coreboot image. Please read
[mrc.bin](../../northbridge/intel/haswell/mrc.bin) for instructions on
retrieving and using it.

Intel ME firmware is in the flash chip. It is not needed when building coreboot.
It can be extracted from the OEM firmware. You can also flash only the BIOS
region to leave Intel ME firmware untouched.

## Programming

To flash with a programmer, you need to remove the battery and the base cover.

For more details have a look at the general [flashing tutorial].

It is also possible to flash internally with flashrom when the IFD is unlocked.

## Debugging

The board can be debugged with EHCI debug. The EHCI debug port is next to the miniDP port.

There's a serial port on dock, but it's not supported yet in coreboot.

Schematic of this laptop can be found on [Lab One]. The board name is Compal LA-9431P.

## Test status

### Not working

- EC ACPI
- Touchpad and the buttons on it
- SD/MMC card reader (kernel reports "Timeout waiting for hardware cmd interrupt.")
- No internal display before booting to OS when connected with a dock

### Working

- Integrated graphics init with libgfxinit
- mSATA
- WLAN
- USB
- Keyboard
- Dock: all USB ports, DisplayPort, eSATA
- Internal flashing when IFD is unlocked


## Technology

```eval_rst
+------------------+-----------------------------+
| CPU              | Intel Haswell-ULT           |
+------------------+-----------------------------+
| PCH              | Intel Lynx Point Low Power  |
+------------------+-----------------------------+
| EC               | SMSC MEC5075                |
+------------------+-----------------------------+
| Super I/O        | SMSC ECE5048                |
+------------------+-----------------------------+
| Coprocessor      | Intel Management Engine     |
+------------------+-----------------------------+
```

[Dell Latitude E7240]: https://www.dell.com/support/home/en-us/product-support/product/latitude-e7240-ultrabook/docs
[flashing tutorial]: ../../flash_tutorial/ext_power.md
[Lab One]: https://www.laboneinside.com/dell-latitude-e7240-schematic-diagram/
