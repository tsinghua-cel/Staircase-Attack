import json
import time
import os

def main():
      validator_balance={}

      for i in range(32):
            key = i
            value = [32000000000]
            validator_balance[key] = value


      while True:
            if os.path.exists("balance.json"):
                  with open('balance.json', 'r') as file:
                        new_data = json.load(file)

                  new_balances = new_data['balances_after_epoch_transition']
                  for i in range(32):
                        validator_balance[i].append(new_balances[i])
                  
                  with open("balance.txt", "w") as file:
                        for key, value in validator_balance.items():
                              file.write(f"{key}: {value}\n")

                  os.remove('balance.json')
            time.sleep(60)


if __name__ == "__main__":
    main()