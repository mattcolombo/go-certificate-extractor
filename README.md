# go-certificate-extractor
A tool for extracting certificates and keys from a p12 store when OpenSSL is not available. Mainly useful for Windows systems that do not already have OpenSSL and where user does not have privileges to install it.

:warning: **NOTE:** all the below instructions are mainly aimed at Windows users. As such also all commands provided are intended to be used with Powershell. When using a different shell the syntax of some of the commands may change slightly. 

Also note that anyone using Linux can achieve the same things with OpenSSL (which comes installed with any reasonable distribution) in a much easier way using native OpenSSL commands. Even Windows users that are able to install OpenSSL would be in principle better served using that. This tool therefore is chiefly aimed at Windows users that are unable to obtain OpenSSL due to system restrictions of some kind.

## What the tool does

This tool will read a p12 formatted certificate store, extract the certificate and private key and create files containing:
* the certificate in PEM format
* the base64 encoded string containing the PEM certificate
* the private key in PEM format
* the base64 encoded string containing the PEM private key

## How to use

First of all download the `go-certificate-extractor_windows_amd64_vN_n.exe` file from the release page in Github (found under this main [Github page](https://github.com/mattcolombo/go-certificate-extractor)) and save it wherever it is best convenient on the disk. For the commands below to work out of the box, it's best to also rename the `.exe` file to simply `go-certificate-extractor.exe`. If not, please adjust the below commands to reflect the name of the file.

To run the tool navigate to the folder where you wish to store the output and run the below command
```(Powershell)
\path\to\go-certificate-extractor.exe <\path\to\p12file> <common name> <p12 store password>
```
As an example (in the case that both the executable and the p12 are in the same folder where we wish to have the output)
```(Powershell)
.\go-certificate-extractor.exe .\my-p12-file.p12 my-common-name.co.uk kjshdu62t34
```
The tool will generate 4 files in the current folder as described above:
* `<common name>-signed.crt` (certificate in PEM format)
* `<common name>-signed.crt-b64enc.txt` (base64 encoded string containing the PEM certificate)
* `<common name>-key.pem` (private key in PEM format)
* `<common name>-key.pem-b64enc.txt` (base64 encoded string containing the PEM private key)

## How to build 

Building the tool works the same way as building any other go application. In the specific Windows case the command below can be used. Notice this works for Powershell; other types of shell may require slightly different syntax.
```(Powershell)
$env:GOOS="windows"; $env:GOARCH="amd64"; go build .\go-certificate-extractor\go-certificate-extractor.go
```