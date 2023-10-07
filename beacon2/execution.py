import os
import json
import time
from datetime import datetime, timedelta

def main():
    with open('genesistime.txt', 'r') as file:
        datetime_str = file.read().strip()
    t = datetime.strptime(datetime_str[:19], "%Y-%m-%d %H:%M:%S")
    if os.path.exists("uatt.json"):
        os.remove("uatt.json")
    max_proposer_slots = 0
    last_checked_slot_count = 0
    while True:
        current_time = datetime.now()
        time_difference = current_time - t
        slot_duration = timedelta(seconds=12)
        slot_count = time_difference // slot_duration
        if slot_count > last_checked_slot_count:
            print(slot_count)
            last_checked_slot_count = slot_count
            if os.path.exists("duties.json"):
                with open("duties.json", "r") as duties:
                    duties_data = json.load(duties)
                for duty in duties_data["duties"]:
                    if "proposer_slots" in duty:
                        proposer_slots = max(duty["proposer_slots"])
                        if max_proposer_slots is None or proposer_slots > max_proposer_slots:
                            max_proposer_slots = proposer_slots
                with open("proposer_slot.txt", "w") as file:
                    file.write(str(max_proposer_slots))
                os.remove("duties.json")
            if os.path.exists("att.json"):
                with open("att.json", "r") as att_file:
                    att_data = json.load(att_file)
                if att_data["data"]["slot"] < max_proposer_slots:
                    if os.path.exists("uatt.json"):
                        with open("uatt.json", "r") as uatt_file:
                            uatt_data = json.load(uatt_file)
                        if isinstance(uatt_data, list):
                            uatt_data.append(att_data)
                        else:
                            uatt_data = [uatt_data, att_data]
                        with open("uatt.json", "w") as uatt_file:
                            json.dump(uatt_data, uatt_file, indent=2)
                    else:
                        with open("uatt.json", "w") as uatt_file:
                            json.dump(att_data, uatt_file, indent=2)
                os.remove("att.json")
            else:
                time.sleep(1)
            if (slot_count+1)%8 == 0:
                print("epoch end")
                max_proposer_slots = 0
                if os.path.exists("uatt.json"):
                    time.sleep(3)
                    os.remove("uatt.json")


if __name__ == "__main__":
    main()
