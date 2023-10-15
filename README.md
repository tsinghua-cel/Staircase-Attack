# Staircase-Attack

This is an attack implementation against incentive of Ethereum PoS.

## Table of Contents

1. [Description](#chapter-001)<br>
2. [Run the Testnet Step by Step](#chapter-002)<br>
3. [Experimental Results](#chapter-003)<br>
4. [Acknowledgments](#chapter-004)<br>


## **1 Description**<a id="chapter-001"></a>

**Built for research use**: This testnet is designed for research purposes and can conduct the Staircase attack on incentive of Ethereum PoS. It enables complete reproducibility of the attacks presented in the research paper and includes source code with historical data. 


## **2 Run the Testnet Step by Step**<a id="chapter-002"></a>
We establish 1000 validators for testing. Among them, the honest validators run the code in folder beacon1 while the Byzantine validators run the code in folder beacon2.

* First, clone a repository here:
```
git clone https://github.com/rauljordan/eth-pos-devnet && cd eth-pos-devnet
```

* Next, start the testnet. The script will start one execution client geth, two beacon client and two validator.
```
./start.sh
```


* Last, stop the testnet after running for some times. 
```
./stop.sh
```

## **3 Experimental results**<a id="chapter-003"></a>

After stop the testnet, you can check the file in out for the details of clienst running.

The stake of hoenst validator in each epoch is shown in ./beacon1/balance.txt.

We show a typical epoch that has been attacked. The attestor role of honest validator in epoch 17 is shown here.

<img src=./figs/honest1.png width=50% />

The last Byzantine validator release its block in the 17th slot of epoch 18. 
<img src=./figs/byzantine1.png width=50% />


After the block is released, an organization happens.
<img src=./figs/honest3.png width=50% />

Then after two epochs, the honest validators in the first period of epoch 17 receives penalties.
<img src=./figs/honest2.png width=50% />


## 4 Acknowledgments<a id="chapter-004"></a>
Attacking source code is developed based on [Prysm](https://github.com/prysmaticlabs/prysm/tree/v4.0.3-patchFix), an popular implementation of Ethereum PoS.
