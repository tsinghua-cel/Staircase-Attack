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

* Start the execution layer of Eth 2.0.
```
./execution.sh
```

* Start the honest validators. 

First, reopen a new bash.

```
cd beacon1
./beacon.sh
```
Second, reopen a new bash again.

```
cd beacon1
./validator.sh
```

* Start the Byzantine validators. 

First, reopen a new bash.

```
cd beacon2
./beacon.sh
```
Second, reopen a new bash again.

```
cd beacon1
./validator.sh
```

* Reopen a new bash and start the monitor script.

```
cd beacon2
python3 execution.py
```

## **3 Experimental results**<a id="chapter-003"></a>

A typical epoch that has been attacked. The attestor role of honest validator in epoch 17 is shown here.

<img src=./figs/honest1.png width=100% />

The last Byzantine validator release its block in the 17th slot of epoch 18. 
<img src=./figs/byzantine1.png width=100% />


After the block is released, an organization happens.
<img src=./figs/honest3.png width=100% />

Then after two epochs, the honest validators in the first period of epoch 17 receives penalties.
<img src=./figs/honest2.png width=100% />


## 4 Acknowledgments<a id="chapter-004"></a>
Attacking source code is developed based on [Prysm](https://github.com/prysmaticlabs/prysm/tree/v4.0.3-patchFix), an popular implementation of Ethereum PoS.
