import subprocess
import json

command = "./validator.sh"

process = subprocess.Popen(command, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)


while True:
    output = process.stdout.readline()
    if output == '' and process.poll() is not None:
        break
    if "aggregation_bits:" in output:
        with open("uatt.json", "a") as json_file:
            json_file.write(json.dumps(output, indent=4))
            json_file.write('\n')
    if output:
        print(output.strip())
process.wait()
