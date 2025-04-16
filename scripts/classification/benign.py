import csv
import requests
import sys

csv.field_size_limit(sys.maxsize)

print("Sending benign parameters from payload_full.csv")

with open('payload_full.csv', newline='') as csvfile:
    csv_reader = csv.reader(csvfile, delimiter=',', quotechar='|')
    for row in csv_reader:
        if row[3] == '"norm"':
            param = row[0][1:-1]

            requests.get(f'http://127.0.0.1:8080/', params={'param': param})

print("Sending benign parameters from SQLiV3.csv")

with open('SQLiV3.csv', newline='') as csvfile:
    csv_reader = csv.reader(csvfile, delimiter=',', quotechar='|')
    for row in csv_reader:
        if row[1] == '0':
            param = row[0][1:-1]
            requests.get(f'http://127.0.0.1:8080/', params={'param': param})

print("Sending benign traffic from XSS_dataset.csv")

with open('XSS_dataset.csv', newline='') as csvfile:
    csv_reader = csv.reader(csvfile, delimiter=',', quotechar='|')
    for row in csv_reader:
        if row[1] == '0':
            param = row[0][1:-1]
            requests.get(f'http://127.0.0.1:8080/', params={'param': param})