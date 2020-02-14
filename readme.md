# AZMigrate

## Description

AZMigrate is a tool to migrate construct agencies

## Build
To build AZMigrate, run the script in /build that corresponds to your desired target architecture. 
The built executable will be in the /dist folder

Note: Ensure you run the build script from within the build directory

## Run
### Prerequisites
The following environment variables need to be set:
- CONSTRUCT_ACCESS_KEY: The construct subsystem access key
### Command Line Arguments
#### Optional Args
- path - Specify an absolute path for the agencies.yml. If not used, AZMigrate will check the working directory for agencies.yml
- agency - The agency to migrate. Multiple agencies can be specified in a comma separated list. This will cause agencies.yml to be ignored.
- target - The target host id to migrate agencies to. If not specified, agencies will be migrated to the Azure-MT-Host
- enable - check Enabled for each agency, with no other changes
- disable - uncheck Enabled for each agency, with no other changed
### Examples
- Migrate agencies defined in agencies.yml to Azure-MT-Host
> azmigrate.exe

- Migrate agencies defined in agencies.yml to a specific host
> azmigrate.exe -target "d6e733d3-634c-4f75-8193-4b8b68c28292"

- Migrate specific agencies to Azure-MT-Host
> azmigrate.exe -agency CRC10,CRC9

Migrate specific agencies to a specific host 
> azmigrate.exe -agency CRC10,CRC9 -target "d6e733d3-634c-4f75-8193-4b8b68c28292"

Use a different .yml file
> azmigrate.exe -path C:\Temp\test.yml

Enable a set of agencies (no other change)
> azmigrate.exe -agency CRC10,CRC9 -enable

Disable a set of agencies (no other change)
> azmigrate.exe -agency CRC10,CRC9 -disable