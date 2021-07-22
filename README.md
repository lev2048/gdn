# GDN

gdn can turn process into daemon.

### Usage

    gdn: run command as daemon

    	<command>   run your command
    	list        show running commands
    	stop <proc>  stop a command by SIGTERM
    	log <proc>   view log of command

    	help        show help
    	version     show version

### Example

    $ gdn <command>

### Where are log files?

All log files are stored in /var/gdn based on proc name

## License

Licensed under The GPLv3 License