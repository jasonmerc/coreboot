package main

import (
	"fmt"
	"strings"
)

type LPVariant int

const (
	LYNX_POINT_MOBILE LPVariant = iota
	LYNX_POINT_DESKTOP
	LYNX_POINT_SERVER
	LYNX_POINT_ULT
)

type lynxpoint struct {
	variant LPVariant
	node    *DevTreeNode
}

func lp_gpio_to_pirq(g int) int {
	switch g {
	case 8:
		return 0 /* PIRQI */
	case 9:
		return 1 /* PIRQJ */
	case 10:
		return 2 /* PIRQK */
	case 13:
		return 3 /* PIRQL */
	case 14:
		return 4 /* PIRQM */
	case 45:
		return 5 /* PIRQN */
	case 46:
		return 6 /* PIRQO */
	case 47:
		return 7 /* PIRQP */
	case 48:
		return 8 /* PIRQQ */
	case 49:
		return 9 /* PIRQR */
	case 50:
		return 10 /* PIRQS */
	case 51:
		return 11 /* PIRQT */
	case 52:
		return 12 /* PIRQU */
	case 53:
		return 13 /* PIRQV */
	case 54:
		return 14 /* PIRQW */
	case 55:
		return 15 /* PIRQX */
	default:
		return -1
	}
}

func conf0str(conf0 uint32) string {
	if (conf0 & 1) == 0 {
		return "GPIO_MODE_NATIVE"
	} else {
		s := []string{"GPIO_MODE_GPIO"}
		var gpio_output bool
		if ((conf0 >> 2) & 1) == 1 {
			s = append(s, "GPIO_DIR_INPUT")
			gpio_output = false
		} else {
			s = append(s, "GPIO_DIR_OUTPUT")
			gpio_output = true
		}
		if ((conf0 >> 3) & 1) == 1 {
			s = append(s, "GPIO_INVERT")
		}
		if ((conf0 >> 4) & 1) == 1 {
			s = append(s, "GPIO_IRQ_LEVEL")
		}
		if gpio_output {
			if ((conf0 >> 31) & 1) == 1 {
				s = append(s, "GPO_LEVEL_HIGH")
			} else {
				s = append(s, "GPO_LEVEL_LOW")
			}
		}
		return strings.Join(s, " | ")
	}
}

func lpgpio_preset(conf0 uint32, owner uint32, route uint32, irqen uint32, pirq uint32) string {
	if conf0 == 0b1101 { /* MODE_GPIO | INPUT | INVERT */
		if owner == 0 { /* OWNER_ACPI */
			if irqen == 0 && pirq == 0 {
				if route == 0 { /* SCI */
					return "LP_GPIO_ACPI_SCI"
				} else {
					return "LP_GPIO_ACPI_SMI"
				}
			}
			return ""
		} else { /* OWNER_GPIO */
			if route == 0 && irqen == 0 && pirq != 0 {
				return "LP_GPIO_INPUT_INVERT"
			}
			return ""
		}
	}

	if conf0 == 0b101 && owner == 1 { /* MODE_GPIO | INPUT, OWNER_GPIO */
		if route == 0 && irqen == 0 {
			if pirq == 1 {
				return "LP_GPIO_PIRQ"
			} else {
				return "LP_GPIO_INPUT"
			}
		}
		return ""
	}

	if owner == 1 && irqen == 1 {
		if route == 0 && pirq == 0 {
			if conf0 == 0b00101 {
				return "LP_GPIO_IRQ_EDGE"
			}
			if conf0 == 0b10101 {
				return "LP_GPIO_IRQ_LEVEL"
			}
		}
		return ""
	}
	return ""
}

func gpio_str(conf0 uint32, conf1 uint32, owner uint32, route uint32, irqen uint32, reset uint32, blink uint32, pirq uint32) string {
	s := []string{}
	s = append(s, fmt.Sprintf(".conf0 = %s", conf0str(conf0)))
	if conf1 != 0 {
		s = append(s, fmt.Sprintf(".conf1 = 0x%x", conf1))
	}
	if owner != 0 {
		s = append(s, ".owner = GPIO_OWNER_GPIO")
	}
	if route != 0 {
		s = append(s, ".route = GPIO_ROUTE_SMI")
	}
	if irqen != 0 {
		s = append(s, ".irqen = GPIO_IRQ_ENABLE")
	}
	if reset != 0 {
		s = append(s, ".reset = GPIO_RESET_RSMRST")
	}
	if blink != 0 {
		s = append(s, ".blink = GPO_BLINK")
	}
	if pirq != 0 {
		s = append(s, ".pirq = GPIO_PIRQ_APIC_ROUTE")
	}
	return strings.Join(s, ", ")
}

func LPGPIO(ctx Context, inteltool InteltoolData) {
	/* TODO: create it as header files as other boards */
	gpio := Create(ctx, "gpio.c")
	defer gpio.Close()

	AddBootBlockFile("gpio.c", "")
	AddROMStageFile("gpio.c", "")

	Add_gpl(gpio)
	gpio.WriteString(`#include <southbridge/intel/lynxpoint/lp_gpio.h>

const struct pch_lp_gpio_map mainboard_gpio_map[] = {
`)

	for g := 0; g <= 94; g++ {
		fmt.Fprintf(gpio, "\t")
		conf0 := inteltool.GPIO[uint16(0x100+g*8)]
		conf1 := inteltool.GPIO[uint16(0x104+g*8)]
		set := g / 32
		bit := g % 32
		/* owner only effective in GPIO mode */
		owner := (inteltool.GPIO[uint16(0x0+set*4)] >> bit) & 1
		route := (inteltool.GPIO[uint16(0x30+set*4)] >> bit) & 1
		irqen := (inteltool.GPIO[uint16(0x90+set*4)] >> bit) & 1
		reset := (inteltool.GPIO[uint16(0x60+set*4)] >> bit) & 1
		var blink, pirq uint32
		/* blink only effective in GPIO output mode */
		if set == 0 {
			blink = (inteltool.GPIO[uint16(0x18)] >> bit) & 1
		} else {
			blink = 0
		}
		irqset := lp_gpio_to_pirq(g)
		if irqset >= 0 {
			pirq = (inteltool.GPIO[uint16(0x10)] >> irqset) & 1
		} else {
			pirq = 0
		}

		if (conf0 & 1) == 0 {
			fmt.Fprintf(gpio, "LP_GPIO_NATIVE, /* %d: NATIVE */\n", g)
		} else if (conf0 & 4) == 0 {
			/* configured as output */
			var lvl string
			if ((conf0 >> 31) & 1) == 0 {
				lvl = "LOW"
			} else {
				lvl = "HIGH"
			}
			fmt.Fprintf(gpio, "LP_GPIO_OUT_%s, /* %d: OUTPUT %s */\n", lvl, g, lvl)
		} else if (conf1 & 4) != 0 {
			/* configured as input and sensing disabled */
			fmt.Fprintf(gpio, "LP_GPIO_UNUSED, /* %d: UNUSED */\n", g)
		} else {
			is_preset := false
			if conf1 == 0 && reset == 0 && blink == 0 {
				preset := lpgpio_preset(conf0, owner, route, irqen, pirq)
				if preset != "" {
					fmt.Fprintf(gpio, "%s, /* %d: %s */\n", preset, g, preset)
					is_preset = true
				}
			}
			if !is_preset {
				fmt.Fprintf(gpio, "{ %s }, /* %d */\n", gpio_str(conf0, conf1, owner, route, irqen, reset, blink, pirq), g)
			}
		}
	}
	gpio.WriteString(`	LP_GPIO_END
};`)
}

func lpPchGetFlashSize(ctx Context) {
	inteltool := ctx.InfoSource.GetInteltool()
	/* In LP PCH, Boot BIOS Straps field in GCS has only one bit.  */
	switch (inteltool.RCBA[0x3410] >> 10) & 1 {
	case 0:
		ROMProtocol = "SPI"
		highflkb := uint32(0)
		for reg := uint16(0); reg < 5; reg++ {
			fl := (inteltool.RCBA[0x3854+4*reg] >> 16) & 0x1fff
			flkb := (fl + 1) << 2
			if flkb > highflkb {
				highflkb = flkb
			}
		}
		ROMSizeKB = int(highflkb)
		/* Shared with ME. Flashrom is unable to handle it.  */
		FlashROMSupport = "n"
	}
}

func (b lynxpoint) GetGPIOHeader() string {
	return "southbridge/intel/lynxpoint/pch.h"
}

func (b lynxpoint) EnableGPE(in int) {
	if b.variant != LYNX_POINT_ULT {
		b.node.Registers[fmt.Sprintf("gpi%d_routing", in)] = "2"
	}
}

func (b lynxpoint) EncodeGPE(in int) int {
	return in + 0x10
}

func (b lynxpoint) DecodeGPE(in int) int {
	return in - 0x10
}

func (b lynxpoint) NeedRouteGPIOManually() {
	b.node.Comment += ", FIXME: set gpiX_routing for EC support"
}

func (b lynxpoint) Scan(ctx Context, addr PCIDevData) {

	SouthBridge = &b

	inteltool := ctx.InfoSource.GetInteltool()

	if b.variant == LYNX_POINT_ULT {
		LPGPIO(ctx, inteltool)
	} else {
		GPIO(ctx, inteltool)
	}

	KconfigBool["SOUTHBRIDGE_INTEL_LYNXPOINT"] = true
	if b.variant == LYNX_POINT_ULT {
		KconfigBool["INTEL_LYNXPOINT_LP"] = true
	}
	KconfigBool["SERIRQ_CONTINUOUS_MODE"] = true
	KconfigInt["USBDEBUG_HCD_INDEX"] = 2
	KconfigComment["USBDEBUG_HCD_INDEX"] = "FIXME: check this"

	if b.variant == LYNX_POINT_ULT {
		lpPchGetFlashSize(ctx)
	} else {
		ich9GetFlashSize(ctx)
	}

	FADT := ctx.InfoSource.GetACPI()["FACP"]

	sp0dtle_data := (inteltool.IOBP[0xea002550] >> 24) & 0xf
	sp0dtle_edge := (inteltool.IOBP[0xea002554] >> 16) & 0xf
	sp1dtle_data := (inteltool.IOBP[0xea002750] >> 24) & 0xf
	sp1dtle_edge := (inteltool.IOBP[0xea002754] >> 16) & 0xf

	if sp0dtle_data != sp0dtle_edge {
		fmt.Printf("Different SATA Gen3 port0 DTLE data and edge values are used.\n")
	}

	if sp1dtle_data != sp1dtle_edge {
		fmt.Printf("Different SATA Gen3 port1 DTLE data and edge values are used.\n")
	}

	cur := DevTreeNode{
		Chip:    "southbridge/intel/lynxpoint",
		Comment: "Intel Series 8 Lynx Point PCH",

		Registers: map[string]string{
			"gen1_dec":             FormatHexLE32(PCIMap[PCIAddr{Bus: 0, Dev: 0x1f, Func: 0}].ConfigDump[0x84:0x88]),
			"gen2_dec":             FormatHexLE32(PCIMap[PCIAddr{Bus: 0, Dev: 0x1f, Func: 0}].ConfigDump[0x88:0x8c]),
			"gen3_dec":             FormatHexLE32(PCIMap[PCIAddr{Bus: 0, Dev: 0x1f, Func: 0}].ConfigDump[0x8c:0x90]),
			"gen4_dec":             FormatHexLE32(PCIMap[PCIAddr{Bus: 0, Dev: 0x1f, Func: 0}].ConfigDump[0x90:0x94]),
			"sata_ahci":            "1",
			"sata_port_map":        fmt.Sprintf("0x%x", PCIMap[PCIAddr{Bus: 0, Dev: 0x1f, Func: 2}].ConfigDump[0x92]&0x3f),
			"docking_supported":    (FormatBool((FADT[113] & (1 << 1)) != 0)),
			"sata_port0_gen3_dtle": fmt.Sprintf("0x%x", sp0dtle_data),
			"sata_port1_gen3_dtle": fmt.Sprintf("0x%x", sp1dtle_data),
		},
		PCISlots: []PCISlot{
			PCISlot{PCIAddr: PCIAddr{Dev: 0x14, Func: 0}, writeEmpty: false, additionalComment: "xHCI Controller"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x16, Func: 0}, writeEmpty: true, additionalComment: "Management Engine Interface 1"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x16, Func: 1}, writeEmpty: true, additionalComment: "Management Engine Interface 2"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x16, Func: 2}, writeEmpty: true, additionalComment: "Management Engine IDE-R"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x16, Func: 3}, writeEmpty: true, additionalComment: "Management Engine KT"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x19, Func: 0}, writeEmpty: true, additionalComment: "Intel Gigabit Ethernet"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1a, Func: 0}, writeEmpty: true, additionalComment: "USB2 EHCI #2"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1b, Func: 0}, writeEmpty: true, additionalComment: "High Definition Audio"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1c, Func: 0}, writeEmpty: true, additionalComment: "PCIe Port #1"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1c, Func: 1}, writeEmpty: true, additionalComment: "PCIe Port #2"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1c, Func: 2}, writeEmpty: true, additionalComment: "PCIe Port #3"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1c, Func: 3}, writeEmpty: true, additionalComment: "PCIe Port #4"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1c, Func: 4}, writeEmpty: true, additionalComment: "PCIe Port #5"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1c, Func: 5}, writeEmpty: true, additionalComment: "PCIe Port #6"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1c, Func: 6}, writeEmpty: true, additionalComment: "PCIe Port #7"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1c, Func: 7}, writeEmpty: true, additionalComment: "PCIe Port #8"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1d, Func: 0}, writeEmpty: true, additionalComment: "USB2 EHCI #1"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1f, Func: 0}, writeEmpty: true, additionalComment: "LPC bridge"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1f, Func: 2}, writeEmpty: true, additionalComment: "SATA Controller (AHCI)"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1f, Func: 3}, writeEmpty: true, additionalComment: "SMBus"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1f, Func: 5}, writeEmpty: true, additionalComment: "SATA Controller (Legacy)"},
			PCISlot{PCIAddr: PCIAddr{Dev: 0x1f, Func: 6}, writeEmpty: true, additionalComment: "Thermal"},
		},
	}

	b.node = &cur

	PutPCIChip(addr, cur)
	PutPCIDevParent(addr, "", "lpc")

	DSDTIncludes = append(DSDTIncludes, DSDTInclude{
		File: "southbridge/intel/common/acpi/platform.asl",
	})
	DSDTIncludes = append(DSDTIncludes, DSDTInclude{
		File:    "southbridge/intel/lynxpoint/acpi/globalnvs.asl",
		Comment: "global NVS and variables",
	})
	DSDTIncludes = append(DSDTIncludes, DSDTInclude{
		File: "southbridge/intel/common/acpi/sleepstates.asl",
	})
	DSDTPCI0Includes = append(DSDTPCI0Includes, DSDTInclude{
		File: "southbridge/intel/lynxpoint/acpi/pch.asl",
	})

	AddBootBlockFile("bootblock.c", "")
	bb := Create(ctx, "bootblock.c")
	defer bb.Close()
	Add_gpl(bb)
	bb.WriteString(`#include <southbridge/intel/lynxpoint/pch.h>

void mainboard_config_superio(void)
{
}`)

	sb := Create(ctx, "romstage.c")
	defer sb.Close()
	Add_gpl(sb)
	sb.WriteString(`#include <stdint.h>
#include <northbridge/intel/haswell/haswell.h>
#include <northbridge/intel/haswell/raminit.h>
#include <southbridge/intel/lynxpoint/pch.h>

void mainboard_config_rcba(void)
{
`)
	RestoreDIRRoute(sb, "D31IR", uint16(inteltool.RCBA[0x3140]))
	RestoreDIRRoute(sb, "D29IR", uint16(inteltool.RCBA[0x3144]))
	RestoreDIRRoute(sb, "D28IR", uint16(inteltool.RCBA[0x3146]))
	RestoreDIRRoute(sb, "D27IR", uint16(inteltool.RCBA[0x3148]))
	RestoreDIRRoute(sb, "D26IR", uint16(inteltool.RCBA[0x314c]))
	RestoreDIRRoute(sb, "D25IR", uint16(inteltool.RCBA[0x3150]))
	RestoreDIRRoute(sb, "D22IR", uint16(inteltool.RCBA[0x315c]))
	RestoreDIRRoute(sb, "D20IR", uint16(inteltool.RCBA[0x3160]))

	sb.WriteString(`}

/* FIXME: called after romstage_common, remove it if not used */
void mb_late_romstage_setup(void)
{
}

void mb_get_spd_map(uint8_t spd_map[4])
{
	/* FIXME: check this */
	/* note that SPD addresses are left-shifted by 1. */
	spd_map[0] = 0xa0;
	spd_map[1] = 0xa2;
	spd_map[2] = 0xa4;
	spd_map[3] = 0xa6;
}

void mainboard_fill_pei_data(struct pei_data *pei_data)
{
	struct usb2_port_setting usb2_ports[MAX_USB2_PORTS] = {
		/* Length, Enable, OCn#, Location */
`)

	pdo1 := PCIMap[PCIAddr{Bus: 0, Dev: 0x1d, Func: 0}].ConfigDump[0x64]
	ocmap1 := PCIMap[PCIAddr{Bus: 0, Dev: 0x1d, Func: 0}].ConfigDump[0x74:0x78]

	var pdo2 uint8
	var ocmap2 []uint8
	var nPorts int
	if b.variant != LYNX_POINT_ULT {
		pdo2 = PCIMap[PCIAddr{Bus: 0, Dev: 0x1a, Func: 0}].ConfigDump[0x64]
		ocmap2 = PCIMap[PCIAddr{Bus: 0, Dev: 0x1a, Func: 0}].ConfigDump[0x74:0x78]
		nPorts = 14
	} else {
		nPorts = 8
	}

	xusb2pr := GetLE16(PCIMap[PCIAddr{Bus: 0, Dev: 0x14, Func: 0}].ConfigDump[0xd0:0xd4])

	for port := 0; port < nPorts; port++ {
		var port_oc int = -1
		var port_pos string
		var port_disable uint8

		if port < 8 {
			port_disable = ((pdo1 >> port) & (uint8(xusb2pr>>port) ^ 1)) & 1
			for oc := 0; oc < 4; oc++ {
				if (ocmap1[oc] & (1 << port)) != 0 {
					port_oc = oc
					break
				}
			}
		} else {
			port_disable = ((pdo2 >> (port - 8)) & (uint8(xusb2pr>>port) ^ 1)) & 1
			for oc := 0; oc < 4; oc++ {
				if (ocmap2[oc] & (1 << (port - 8))) != 0 {
					port_oc = oc + 4
					break
				}
			}
		}

		/* get USB2 port length and location from IOBP */
		port_iobp := inteltool.IOBP[0xe5004100+uint32(port)*0x100]
		loc_param := (port_iobp >> 8) & 7
		txamp := (port_iobp >> 11) & 7
		var port_length int

		if b.variant == LYNX_POINT_ULT {
			if loc_param == 6 {
				/* back panel or mini pcie, length >= 0x70 */
				port_pos = "USB_PORT_MINI_PCIE"
				if txamp <= 2 {
					port_length = 0x80
				} else {
					port_length = 0x110
				}
			} else if loc_param == 4 {
				port_pos = "USB_PORT_DOCK"
				if txamp <= 1 {
					port_length = 0x40
				} else {
					port_length = 0x80
				}
			} else {
				port_pos = "USB_PORT_BACK_PANEL"
				port_length = 0x40
			}
		} else if b.variant == LYNX_POINT_MOBILE {
			if loc_param == 4 {
				port_pos = "USB_PORT_DOCK"
				if txamp <= 1 {
					port_length = 0x40
				} else {
					port_length = 0x80
				}
			} else if loc_param == 6 {
				/* not internal, not dock, port_length >= 0x70 */
				port_pos = "USB_PORT_BACK_PANEL"
				if txamp <= 2 {
					port_length = 0x80
				} else {
					port_length = 0x110
				}
			} else {
				port_pos = "USB_PORT_BACK_PANEL"
				port_length = 0x40
			}
		} else { /* desktop or server */
			if loc_param == 4 {
				port_pos = "USB_PORT_BACK_PANEL"
				if txamp <= 2 {
					port_length = 0x40
				} else if txamp >= 4 {
					port_length = 0x140
				} else {
					port_length = 0x110
				}
			} else {
				port_pos = "USB_PORT_FLEX"
				port_length = 0x40
			}
		}

		if port_disable == 1 {
			port_pos = "USB_PORT_SKIP"
		}

		if port_oc == -1 {
			fmt.Fprintf(sb, "\t\t{ 0x%04x, %d, USB_OC_PIN_SKIP, %s },\n",
				port_length, (port_disable ^ 1), port_pos)
		} else {
			fmt.Fprintf(sb, "\t\t{ 0x%04x, %d, %d, %s },\n",
				port_length, (port_disable ^ 1), port_oc, port_pos)
		}
	}

	sb.WriteString(`	};
	struct usb3_port_setting usb3_ports[MAX_USB3_PORTS] = {
`)

	xpdo := PCIMap[PCIAddr{Bus: 0, Dev: 0x14, Func: 0}].ConfigDump[0xe8]
	u3ocm := PCIMap[PCIAddr{Bus: 0, Dev: 0x14, Func: 0}].ConfigDump[0xc8:0xd0]

	for port := uint(0); port < 6; port++ {
		var port_oc int = -1
		port_disable := (xpdo >> port) & 1
		for oc := 0; oc < 8; oc++ {
			if (u3ocm[oc] & (1 << port)) != 0 {
				port_oc = oc
				break
			}
		}
		if port_oc == -1 {
			fmt.Fprintf(sb, "\t\t{ %d, USB_OC_PIN_SKIP },\n",
				(port_disable ^ 1))
		} else {
			fmt.Fprintf(sb, "\t\t{ %d, %d },\n",
				(port_disable ^ 1), port_oc)
		}
	}

	sb.WriteString(`	};
	memcpy(pei_data->usb2_ports, usb2_ports, sizeof(usb2_ports));
	memcpy(pei_data->usb3_ports, usb3_ports, sizeof(usb3_ports));
}`)

	gnvs := Create(ctx, "acpi_tables.c")
	defer gnvs.Close()

	Add_gpl(gnvs)
	gnvs.WriteString(`#include <acpi/acpi_gnvs.h>
#include <southbridge/intel/lynxpoint/nvs.h>

/* FIXME: check this function.  */
void acpi_create_gnvs(struct global_nvs *gnvs)
{
	/* The lid is open by default. FIXME: not on desktops? */
	gnvs->lids = 1;

	/* Temperature at which OS will shutdown */
	gnvs->tcrt = 100;
	/* Temperature at which OS will throttle CPU */
	gnvs->tpsv = 90;
}
`)

}

func init() {
	for _, id := range []uint16{
		0x8c41, 0x8c49, 0x8c4b, 0x8c4f,
	} {
		RegisterPCI(0x8086, uint16(id), lynxpoint{variant: LYNX_POINT_MOBILE})
	}

	for _, id := range []uint16{
		0x8c42, 0x8c44, 0x8c46, 0x8c4a,
		0x8c4c, 0x8c4e, 0x8c50, 0x8c5c,
	} {
		RegisterPCI(0x8086, uint16(id), lynxpoint{variant: LYNX_POINT_DESKTOP})
	}

	for _, id := range []uint16{
		0x8c52, 0x8c54, 0x8c56,
	} {
		RegisterPCI(0x8086, uint16(id), lynxpoint{variant: LYNX_POINT_SERVER})
	}

	for _, id := range []uint16{
		0x9c41, 0x9c43, 0x9c45,
	} {
		RegisterPCI(0x8086, uint16(id), lynxpoint{variant: LYNX_POINT_ULT})
	}

	/* PCIe bridge */
	for _, id := range []uint16{
		0x8c10, 0x8c12, 0x8c14, 0x8c16, 0x8c18, 0x8c1a, 0x8c1c, 0x8c1e,
		0x9c10, 0x9c12, 0x9c14, 0x9c16, 0x9c18, 0x9c1a,
	} {
		RegisterPCI(0x8086, id, GenericPCI{})
	}

	/* SMBus controller  */
	RegisterPCI(0x8086, 0x8c22, GenericPCI{MissingParent: "smbus"})
	RegisterPCI(0x8086, 0x9c22, GenericPCI{MissingParent: "smbus"})

	/* SATA */
	for _, id := range []uint16{
		0x8c00, 0x8c02, 0x8c04, 0x8c06, 0x8c08, 0x8c0e,
		0x8c01, 0x8c03, 0x8c05, 0x8c07, 0x8c09, 0x8c0f,
		0x9c03, 0x9c05, 0x9c07, 0x9c0f,
	} {
		RegisterPCI(0x8086, id, GenericPCI{})
	}

	/* EHCI */
	for _, id := range []uint16{
		0x9c26, 0x8c26, 0x8c2d,
	} {
		RegisterPCI(0x8086, id, GenericPCI{})
	}

	/* XHCI */
	RegisterPCI(0x8086, 0x8c31, GenericPCI{})
	RegisterPCI(0x8086, 0x9c31, GenericPCI{})

	/* ME and children */
	for _, id := range []uint16{
		0x8c3a, 0x8c3b, 0x8c3c, 0x8c3d,
		0x9c3a, 0x9c3b, 0x9c3c, 0x9c3d,
	} {
		RegisterPCI(0x8086, id, GenericPCI{})
	}

	/* Ethernet */
	RegisterPCI(0x8086, 0x8c33, GenericPCI{})

	/* Thermal */
	RegisterPCI(0x8086, 0x8c24, GenericPCI{})
	RegisterPCI(0x8086, 0x9c24, GenericPCI{})

	/* LAN Controller on LP PCH (if EEPROM has 0x0000/0xffff in DID) */
	RegisterPCI(0x8086, 0x155a, GenericPCI{})

	/* SDIO */
	RegisterPCI(0x8086, 0x9c35, GenericPCI{})

	/* Smart Sound Technology Controller */
	RegisterPCI(0x8086, 0x9c36, GenericPCI{})

	/* Serial I/O */
	for id := uint16(0x9c60); id <= 0x9c66; id++ {
		RegisterPCI(0x8086, id, GenericPCI{})
	}
}
