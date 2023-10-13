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
            print(str(slot_count))
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
            if slot_count == max_proposer_slots-1 and os.path.exists("att.json"):
                time.sleep(8)
                with open('att.json', 'r') as input_file:
                    json_objects = [json.loads(line) for line in input_file]
                output_json = json.dumps(json_objects, indent=4)
                print("111")
                with open('uatt.json', 'w') as output_file:
                    output_file.write(output_json)
            if slot_count%32 == 0:
                if os.path.exists("att.json"):
                    os.remove("att.json")
                if os.path.exists("uatt.json"):
                    os.remove("uatt.json")
                time.sleep(12)
        else:
            time.sleep(1)

if __name__ == "__main__":
    main()
