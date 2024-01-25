# Staircase Attack

<font size=4> This is the staircase attack using the newest Prysm implementation [(V4.2.0)](https://github.com/prysmaticlabs/prysm/releases/tag/v4.2.0). </font>

## Results

<font size=4> The results can also be found in PDF, `results for 1616.pdf`. </font>

### Feasibility of Warm-up Attack

<div>	
    <img src="./figs/block-31.png"
         style=100%/>
    <br>
    <font size=3> Figure 1: Block 31 (root 0xbb) is proposed in slot 31.</font>
</div>

<div>
    <img src="./figs/block-32.png"
         style=100%/>
    <br>
    <font size=3> Figure 2: Block 32 (root 0xcc) is withheld by a Byzantine validator and delayed for 8 seconds.</font>
</div>

<div>	
    <img src="./figs/warm-up-attack.png"
         style=100%/>
    <br>
    <font size=3> Figure 3: Attestations from honest validators in slot 32 to 35. As block 32 is delayed, the target of attestations in slot 32 is set as block 31 (the first INFO). After block 32 is released, the attestations from honest validators will have their target as block 32 (2nd to the 5th INFO).  </font>
</div>

### Staircase Attack with Latest Version

<div>	
    <img src="./figs/reorg-occurred.png"
         style=100%/>
    <br>
    <font size=3> Figure 4: Re-produced results using the Capella version. </font>
</div>

<div>	
    <img src="./figs/attestations.png"
         style=100%/>
    <br>
    <font size=3> Figure 5: Attestations from honest validators before and after reorganization. It can be seen that the fork is Capella.  </font>
</div>

## Run the Testnet Step by Step
We establish 1000 validators for testing. Among them, the honest validators run the code in folder beacon1 while the Byzantine validators run the code in folder beacon2.

* First, clone the repository.

* Next, start the testnet. The script will start one execution client geth, two beacon clients and two validator clients.
```shell
./start.sh
```
* Stop the testnet after running for some time. 
```
./stop.sh
```