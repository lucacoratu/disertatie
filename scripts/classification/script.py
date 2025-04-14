#The dataset was downloaded from https://github.com/Morzeux/HttpParamsDataset/blob/master/payload_full.csv

import csv
import requests

with open('payload_full.csv', newline='') as csvfile:
    csv_reader = csv.reader(csvfile, delimiter=',', quotechar='|')
    for row in csv_reader:
        if row[3] == '"norm"':
            param = row[0][1:-1]

            requests.get(f'http://127.0.0.1:8080/', params={'param': param})
