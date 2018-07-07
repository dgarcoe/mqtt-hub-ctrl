#include <errno.h>
#include <usb.h>
#include <stdio.h>
#include <string.h>

#define USB_RT_HUB			(USB_TYPE_CLASS | USB_RECIP_DEVICE)
#define USB_RT_PORT			(USB_TYPE_CLASS | USB_RECIP_OTHER)
#define USB_PORT_FEAT_POWER		8
#define USB_PORT_FEAT_INDICATOR         22
#define USB_DIR_IN			0x80		/* to host */

#define COMMAND_SET_NONE  0
#define COMMAND_SET_LED   1
#define COMMAND_SET_POWER 2
#define HUB_LED_GREEN 2

#define HUB_CHAR_LPSM		0x0003
#define HUB_CHAR_PORTIND        0x0080

struct usb_hub_descriptor {
  unsigned char bDescLength;
  unsigned char bDescriptorType;
  unsigned char bNbrPorts;
  unsigned char wHubCharacteristics[2];
  unsigned char bPwrOn2PwrGood;
  unsigned char bHubContrCurrent;
  unsigned char data[0];
};

#define CTRL_TIMEOUT 1000
#define USB_STATUS_SIZE 4

#define MAX_HUBS 128
struct hub_info {
  int busnum, devnum;
  struct usb_device *dev;
  int nport;
  int indicator_support;
};

void hub_port_status (usb_dev_handle *uh, int nport);

int usb_find_hubs (struct hub_info *hubs,int listing, int verbose, int busnum, int devnum, int hub);

int get_hub (struct hub_info *hubs,int busnum, int devnum);

int send_command(int hub, int request, int feature, int index);
