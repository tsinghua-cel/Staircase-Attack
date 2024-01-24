# Staircase Attack

<font size=4> This is the staircase attack using the newest Prysm implementation (V4.2.0). </font>

### Feasibility of Warm-up Attack

<div>	
    <center>
    <img src="./figs/block-31.png"
         style=100%/>
    <br>
    <font size=4> Figure 1: Block 31 with root 0xbb. </font>
    </center>
</div>

<div>	
    <center>
    <img src="./figs/block-32.png"
         style=100%/>
    <br>
    <font size=4> Figure 2: Block 32 with root 0xcc. </font>
    </center>
</div>

<div>	
    <center>
    <img src="./figs/warm-up-attack.png"
         style=100%/>
    <br>
    <font size=4> Figure 3: Attestations from honest validators in slots 32 to 35. </font>
    </center>
</div>

### Staircase Attack with Latest Version

<div>	
    <center>
    <img src="./figs/reorg-occurred.png"
         style=100%/>
    <br>
    <font size=4> Figure 4: Reorganization. </font>
    </center>
</div>

<div>	
    <center>
    <img src="./figs/attestations.png"
         style=100%/>
    <br>
    <font size=4> Figure 5: Attestations from honest validators before and after reorganization. It can be seen that the fork is Capella.  </font>
    </center>
</div>