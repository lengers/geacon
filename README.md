# Geacon

**Using Go to implement CobaltStrike's Beacon**

Forked from [the original geacon project](https://github.com/darkr4y/geacon) by the great **[@darkr4y](https://github.com/darkr4y)**.

Greatly improved by the work of [@DidierStevens](https://github.com/DidierStevens) and his tool [cs-parse-http-traffic.py](https://github.com/DidierStevens/Beta/blob/master/cs-parse-http-traffic.py).

----

*This project is for learning protocol analysis and reverse engineering only, if someone's rights have been violated, please contact me to remove the project, and the last DO NOT USE IT ILLEGALLY*



## How to play

1. Setup the teamserver and start a http lisenter, the teamserver will generate the file `.cobaltstrike.beacon_keys`.
2. Compile the BeaconTool with Jetbrains Idea, use command `java -jar BeaconTool.jar ` to convert java keystore to PEM format.
3. Run the `build.sh` file. This will: 
   1. replace the RSA key pair in the file `cmd/config/config.go` (the RSA private key is not required, I wrote it in the code just for the record)
   2. Compile the geacon whatever platform you want to run. Adjust the GO variables in the script beforehand (for example, use the command `export GOOS="darwin" && export GOARCH="amd64" && go build cmd/main.go` to compile an executable binary running on MacOS). 
4. Having fun ! PR and issue is welcome ;)
5. Geacon has just been tested on CobaltStrike 4.6.1 and currently only supports a c2profile if you write the encryption and decryption routines yourself.

Please note: The branch `main` contains the http geacon, the branch `tcp` contains an adjusted version of geacon that is compatible with Cobalt Strike as a tcp beacon.

## Screenshot

Get the Geacon's command execution results on Linux.
![login](https://github.com/darkr4y/geacon/raw/master/screenshots/sc.png)



## Protocol analysis

To be continued, I will update as soon as I have time ...

## Geacon constraints

1. Geacon is built in Go, and Go binaries are not small. The compiled geacon is larger than the buffer allowed by Cobalt Strike for the `spawn` commands, which means that these commands do not work if you compiled geacon using the default Go compiler. There are some solutions to this, and none of them are ideal:
   1. Use `-compiler gccgo -gccgoflags '-static-libgo'` to compile geacon. This will result in dynamically linked dependencies on libc in the version currently present in the build environment. That means you are relying on these dependencies, and if they are not present on your target machine, geacon will fail to run.
   2. You can compile geacon locally, upload it to a temp directory, and run it from there. This is not OPSEC safe, as it leaves files on disk.

## Todo

1. ~~Support CobaltStrike 4.x~~

2. Fix the OS icon issue in session table -> WONTDO
3. String encoding issue
