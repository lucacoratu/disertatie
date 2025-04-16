#https://www.kaggle.com/datasets/syedsaqlainhussain/cross-site-scripting-xss-dataset-for-deep-learning

import csv
import requests

with open('XSS_dataset.csv', newline='') as csvfile:
    csv_reader = csv.reader(csvfile, delimiter=',', quotechar='|')
    for row in csv_reader:
        if row[2] == '1':
            param = row[1][1:-1]
            requests.get("http://127.0.0.1:8080/", params={"param": param})