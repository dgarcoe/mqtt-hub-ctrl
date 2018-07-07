/*
 * Copyright (C) 2006 Free Software Initiative of Japan
 *
 * Author: NIIBE Yutaka  <gniibe at fsij.org>
 *
 * This file can be distributed under the terms and conditions of the
 * GNU General Public License version 2 (or later).
 *
 */

#include "hub-ctrl.h"

static struct hub_info hubs[MAX_HUBS];
static int number_of_hubs_with_feature;

void
hub_port_status (usb_dev_handle *uh, int nport)
{
  int i;

  printf(" Hub Port Status:\n");
  for (i = 0; i < nport; i++)
    {
      char buf[USB_STATUS_SIZE];
      int ret;

      ret = usb_control_msg (uh,
			     USB_ENDPOINT_IN | USB_TYPE_CLASS | USB_RECIP_OTHER,
			     USB_REQ_GET_STATUS,
			     0, i + 1,
			     buf, USB_STATUS_SIZE,
			     CTRL_TIMEOUT);
      if (ret < 0)
	{
	  fprintf (stderr,
		   "cannot read port %d status, %s (%d)\n",
		   i + 1, strerror(errno), errno);
	  break;
	}

      printf("   Port %d: %02x%02x.%02x%02x", i + 1,
	     buf[3], buf [2],
	     buf[1], buf [0]);

      printf("%s%s%s%s%s",
	     (buf[2] & 0x10) ? " C_RESET" : "",
	     (buf[2] & 0x08) ? " C_OC" : "",
	     (buf[2] & 0x04) ? " C_SUSPEND" : "",
	     (buf[2] & 0x02) ? " C_ENABLE" : "",
	     (buf[2] & 0x01) ? " C_CONNECT" : "");

      printf("%s%s%s%s%s%s%s%s%s%s\n",
	     (buf[1] & 0x10) ? " indicator" : "",
	     (buf[1] & 0x08) ? " test" : "",
	     (buf[1] & 0x04) ? " highspeed" : "",
	     (buf[1] & 0x02) ? " lowspeed" : "",
	     (buf[1] & 0x01) ? " power" : "",
	     (buf[0] & 0x10) ? " RESET" : "",
	     (buf[0] & 0x08) ? " oc" : "",
	     (buf[0] & 0x04) ? " suspend" : "",
	     (buf[0] & 0x02) ? " enable" : "",
	     (buf[0] & 0x01) ? " connect" : "");
    }
}

int
usb_find_hubs (int listing, int verbose, int busnum, int devnum, int hub)
{
  struct usb_bus *busses;
  struct usb_bus *bus;

  number_of_hubs_with_feature = 0;
  busses = usb_get_busses();
  if (busses == NULL)
    {
      perror ("failed to access USB");
      return -1;
    }

  for (bus = busses; bus; bus = bus->next)
    {
      struct usb_device *dev;

      for (dev = bus->devices; dev; dev = dev->next)
	{
	  usb_dev_handle *uh;
	  int print = 0;

	  if (dev->descriptor.bDeviceClass != USB_CLASS_HUB)
	    continue;

	  if (listing
	      || (verbose
		  && ((atoi (bus->dirname) == busnum && dev->devnum == devnum)
		      || hub == number_of_hubs_with_feature)))
	    print = 1;

	  uh = usb_open (dev);

	  if (uh != NULL)
	    {
	      char buf[1024];
	      int len;
	      int nport;
	      struct usb_hub_descriptor *uhd = (struct usb_hub_descriptor *)buf;
	      if ((len = usb_control_msg (uh, USB_DIR_IN | USB_RT_HUB,
					  USB_REQ_GET_DESCRIPTOR,
					  USB_DT_HUB << 8, 0,
					  buf, sizeof (buf), CTRL_TIMEOUT))
		  > sizeof (struct usb_hub_descriptor))
		{
		  if (!(uhd->wHubCharacteristics[0] & HUB_CHAR_PORTIND)
		      && (uhd->wHubCharacteristics[0] & HUB_CHAR_LPSM) >= 2)
		    continue;

		  if (print)
		    printf ("Hub #%d at %s:%03d\n",
			    number_of_hubs_with_feature,
			    bus->dirname, dev->devnum);

		  switch ((uhd->wHubCharacteristics[0] & HUB_CHAR_LPSM))
		    {
		    case 0:
		      if (print)
			fprintf (stderr, " INFO: ganged switching.\n");
		      break;
		    case 1:
		      if (print)
			fprintf (stderr, " INFO: individual power switching.\n");
		      break;
		    case 2:
		    case 3:
		      if (print)
			fprintf (stderr, " WARN: No power switching.\n");
		      break;
		    }

		  if (print
		      && !(uhd->wHubCharacteristics[0] & HUB_CHAR_PORTIND))
		    fprintf (stderr, " WARN: Port indicators are NOT supported.\n");
		}
	      else
		{
		  perror ("Can't get hub descriptor");
		  usb_close (uh);
		  continue;
		}

	      nport = buf[2];
	      hubs[number_of_hubs_with_feature].busnum = atoi (bus->dirname);
	      hubs[number_of_hubs_with_feature].devnum = dev->devnum;
	      hubs[number_of_hubs_with_feature].dev = dev;
	      hubs[number_of_hubs_with_feature].indicator_support =
		(uhd->wHubCharacteristics[0] & HUB_CHAR_PORTIND)? 1 : 0;
	      hubs[number_of_hubs_with_feature].nport = nport;

	      number_of_hubs_with_feature++;

	      if (verbose)
		hub_port_status (uh, nport);

	      usb_close (uh);
	    }
	}
    }

  return number_of_hubs_with_feature;
}

int
get_hub (int busnum, int devnum)
{
  int i;

  for (i = 0; i < number_of_hubs_with_feature; i++)
    if (hubs[i].busnum == busnum && hubs[i].devnum == devnum)
      return i;

  return -1;
}

int
send_command(int hub, int request, int feature, int index) {

 usb_dev_handle *uh = NULL;

 if (hubs[hub].dev!=NULL) {

   uh = usb_open(hubs[hub].dev);

   if (uh == NULL) {
     return -1;
   } else {

     if (usb_control_msg(uh,USB_RT_PORT,request,feature,index,NULL,0,CTRL_TIMEOUT)) {
       return -1;
     }

   }

   usb_close(uh);

 } else {
  return -1;
 }

 return 0;

}
