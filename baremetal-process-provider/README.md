# Provider BaremetalProvider

`baremetalprovider` is a [Crossplane](https://crossplane.io/) Provider
that manage the creation and deletion of a Linux Process. 

The `Process` resource manage the remote Linux process by downloading a file and executing it. 

The `EmqxProcess` resource is specific for the Emqx broker process: it download the tar source file, set the API keys for the access and run the start command. 

The required parameters for the resource creation are:


- `host` (string): the address of the Linux machine that hosts the Broker
- `username` (string): username of the Linux machine
- `password` (string): Password of the user
- `brokerAPIKey` (string): Apikey that needs to be installed in the broker to grant the access to the Topic Provider. The api Key has the format `service:key`

