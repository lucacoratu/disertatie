#https://www.kaggle.com/datasets/syedsaqlainhussain/sql-injection-dataset?resource=download

import csv
import requests
import sys

csv.field_size_limit(sys.maxsize)

print("Sending sqli params from payload_full.csv")

with open('payload_full.csv', newline='') as csvfile:
    csv_reader = csv.reader(csvfile, delimiter=',', quotechar='|')
    for row in csv_reader:
        if row[2] == '"sqli"':
            param = row[0][1:-1]
            requests.get(f'http://127.0.0.1:8080/', params={'param': param})

print("Sending sqli params from SQLiV3.csv")

with open('SQLiV3.csv', newline='') as csvfile:
    csv_reader = csv.reader(csvfile, delimiter=',', quotechar='|')
    for row in csv_reader:
        if row[1] == '1':
            param = row[0][1:-1]
            requests.get("http://127.0.0.1:8080/", params={"param": param})