# Staircase Attack

<font size=4> This is the staircase attack using the newest Prysm implementation (V4.2.0). </font>

## Results

The results can also be found in PDF `results for 1616.pdf`.

### Feasibility of Warm-up Attack

<div>	
    <center>
    <img src="./figs/block-31.png"
         style=100%/>
    <br>
    <font size=3> Figure 1: Block 31 (root 0xbb) is proposed in slot 31.</font>
    </center>
</div>

<div>	
    <center>
    <img src="./figs/block-32.png"
         style=100%/>
    <br>
    <font size=3> Figure 2: Block 32 (root 0xcc) is withheld by a Byzantine validator and delayed for 8 seconds.</font>
    </center>
</div>

<div>	
    <center>
    <img src="./figs/warm-up-attack.png"
         style=100%/>
    <br>
    <font size=3> Figure 3: Attestations from honest validators in slot 32 to 35. As block 32 is delayed, the target of attestations in slot 32 is set as block 31 (the first INFO). After block 32 is released, the attestations from honest validators will have their target as block 32 (2nd to the 5th INFO).  </font>
    </center>
</div>

### Staircase Attack with Latest Version

<div>	
    <center>
    <img src="./figs/reorg-occurred.png"
         style=100%/>
    <br>
    <font size=3> Figure 4: Re-produced results using the Capella version. </font>
    </center>
</div>

<div>	
    <center>
    <img src="./figs/attestations.png"
         style=100%/>
    <br>
    <font size=3> Figure 5: Attestations from honest validators before and after reorganization. It can be seen that the fork is Capella.  </font>
    </center>
</div>

## **1 Run the Testnet Step by Step**<a id="chapter-002"></a>
We establish 1000 validators for testing. Among them, the honest validators run the code in folder beacon1 while the Byzantine validators run the code in folder beacon2.

* First, clone the repository.

* Next, start the testnet. The script will start one execution client geth, two beacon clients and two validator clients.
```shell
./start.sh
```
The testnet is start if you see this:
```shell
===========Start execution layer================
INFO [10-15|15:09:13.903] Maximum peer count                       ETH=50 LES=0 total=50
INFO [10-15|15:09:13.906] Smartcard socket not found, disabling    err="stat /run/pcscd/pcscd.comm: no such file or directory"
INFO [10-15|15:09:13.914] Set global gas cap                       cap=50,000,000
INFO [10-15|15:09:13.920] Allocated cache and file handles         database=/root/staircase/gethdata/geth/chaindata cache=16.00MiB handles=16
INFO [10-15|15:09:14.006] Opened ancient database                  database=/root/staircase/gethdata/geth/chaindata/ancient/chain readonly=false
INFO [10-15|15:09:14.006] Writing custom genesis block 
INFO [10-15|15:09:14.009] Persisted trie from memory database      nodes=33 size=4.81KiB time="193.033µs" gcnodes=0 gcsize=0.00B gctime=0s livenodes=1 livesize=0.00B
INFO [10-15|15:09:14.010] Freezer shutting down 
INFO [10-15|15:09:14.011] Successfully wrote genesis state         database=chaindata                               hash=904da6..76899a
INFO [10-15|15:09:14.011] Allocated cache and file handles         database=/root/staircase/gethdata/geth/lightchaindata cache=16.00MiB handles=16
INFO [10-15|15:09:14.113] Opened ancient database                  database=/root/staircase/gethdata/geth/lightchaindata/ancient/chain readonly=false
INFO [10-15|15:09:14.113] Writing custom genesis block 
INFO [10-15|15:09:14.116] Persisted trie from memory database      nodes=33 size=4.81KiB time="260.302µs" gcnodes=0 gcsize=0.00B gctime=0s livenodes=1 livesize=0.00B
INFO [10-15|15:09:14.117] Successfully wrote genesis state         database=lightchaindata                          hash=904da6..76899a
INFO [10-15|15:09:14.226] Maximum peer count                       ETH=50 LES=0 total=50
INFO [10-15|15:09:14.226] Smartcard socket not found, disabling    err="stat /run/pcscd/pcscd.comm: no such file or directory"
INFO [10-15|15:09:14.228] Set global gas cap                       cap=50,000,000
Your new account is locked with a password. Please give a password. Do not forget this password.
!! Unsupported terminal, password will be echoed.
Password: 
Repeat password: 
Address: {123463a4b065722e99115d6c222f267d9cabb524}
INFO[0000] No genesis time specified, defaulting to now()  prefix=genesis
INFO[0000] Specified a chain config file: config.yml     prefix=genesis
INFO[0003] Done writing genesis state to genesis.ssz     prefix=genesis
INFO[0003] Command completed                             prefix=genesis
===========Start honest validators==============
===========Start Byzantine validators===========
```

* Check if all clients are starting correctly.
```shell
ps -ef
```
If you see this, you can wait and enjoy it!!!
```shell
root     2873829       1  0 15:09 pts/3    00:00:04 ./geth --http --http.api eth,engine --datadir=gethdata --allow-insecure-unlock --unlock=0x123
root     2873875       1 12 15:09 pts/3    00:00:56 ./beacon1/beacon-chain --datadir=beacon1/beacondata --min-sync-peers=0 --genesis-state=genesi
root     2873876       1  7 15:09 pts/3    00:00:31 ./beacon1/validator --datadir=beacon1/validatordata --accept-terms-of-use --interop-num-valid
root     2873877       1  0 15:09 pts/3    00:00:00 python3 ./beacon1/balance_monitor.py
root     2873919       1  9 15:09 pts/3    00:00:42 ./beacon2/beacon-chain --datadir=beacon2/beacondata --min-sync-peers=1 --genesis-state=genesi
root     2873920       1  5 15:09 pts/3    00:00:22 ./beacon2/validator --datadir=beacon2/validatordata --accept-terms-of-use --interop-start-ind
root     2873970       1  0 15:09 pts/3    00:00:00 python3 ./beacon2/role_monitor.py
```


* Last, stop the testnet after running for some time. 
```
./stop.sh
```