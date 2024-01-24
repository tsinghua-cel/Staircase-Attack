# Staircase Attack

This is the staircase attack using the newest Prysm implementation (V4.2.0).

## Feasibility of Warm-up Attack

<div>	
    <center>
    <img src="./figs/block-31.png"
         style=100%/>
    <br>
    Figure 1: Block 31 with root 0xbb.
    </center>
</div>

<div>	
    <center>
    <img src="./figs/block-32.png"
         style=100%/>
    <br>
    Figure 2: Block 32 with root 0xcc.
    </center>
</div>

<div>	
    <center>
    <img src="./figs/warm-up-attack.png"
         style=100%/>
    <br>
    Figure 3: Attestations from honest validators in slots 32 to 35.
    </center>
</div>

## Staircase Attack with Latest Version

<div>	
    <center>
    <img src="./figs/reorg-occurred.png"
         style=100%/>
    <br>
    Figure 4: Reorganization.
    </center>
</div>

<div>	
    <center>
    <img src="./figs/attestations.png"
         style=100%/>
    <br>
    Figure 5: Attestations from honest validators before and after reorganization. It can be seen that the fork is Capella. 
    </center>
</div>