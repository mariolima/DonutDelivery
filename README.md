## Description

This is a simple webserver that integrates with [Donut](https://github.com/TheWover/donut) using Binject's [go-donut](https://github.com/Binject/go-donut) and delivers shellcode given a particular (PE) .exe URL - sent as a GET parameter. 

Ideally used with your custom shellcode loader, you can make it arbitrarily load any [Donutable](https://github.com/Flangvik/SharpCollection) PE by remotely requesting its shellcode and merely providing the tool name.

## How to use

```bash
go install github.com/mariolima/donutdelivery@latest
donutdelivery  -l '127.0.0.1:8087' -path '/secretdonutpath' -up 'urlsecretparam' -sp 'argssecretparam'

curl http://127.0.0.1:8087/secretdonutpath?urlsecretparam=https://github.com/Flangvik/SharpCollection/raw/master/NetFramework_4.0_x64/Seatbelt.exe&argssecretparam=antivirus > loader.bin
```

## Help

```
Usage of donutdelivery:
  -l string
    	The listening address of the server (default "127.0.0.1:8087")
  -path string
    	HTTP path to be used for donut delivery (default "/donut")
  -ap string
    	HTTP parameter where the comma separated PE arguments will be sent (default Random)
  -secret string
    	Secret to be sent in the requests along with the PE url (Authorization header if param not defined). 
        No authentication if not defined.
  -sp string
    	HTTP paramter where the secret will be sent (instead of Authorization token header)
  -up string
    	HTTP parameter where the PE url will be sent (default Random)
```
