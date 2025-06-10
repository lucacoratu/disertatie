#Model save
import pickle
import time

import pandas as pd

# Modelling
from sklearn.preprocessing import StandardScaler
from sklearn.svm import SVC
from sklearn.metrics import accuracy_score, confusion_matrix, precision_score, recall_score, classification_report
from sklearn.model_selection import train_test_split

names = ['UrlLength', 'NumberParams', 'NumberSpecialChars', 'RatioSpecialChars', 'NumberRoundBrackets', 'NumberSquareBrackets', 'NumberCurlyBrackets', 'NumberApostrophes', 'NumberQuotationMarks', 'NumberDots', 'NumberSlash', 'NumberBackslash', 'NumberComma', 'NumberColon', 'NumberSemicolon', 'NumberMinus', 'NumberPlus','NumberLessGreater', 'DistanceDots', 'DistanceSlash', 'DistanceBackslash', 'DistanceComma', 'DistanceColon', 'DistanceSemicolon', 'DistanceMinus', 'DistancePlus']

print(len(names))

#Read the csv files
df_benign = pd.read_csv("../../agent/datasets/benign.csv", header=None, names=names)
df_lfi = pd.read_csv("../../agent/datasets/lfi.csv", header=None, names=names)
df_sqli = pd.read_csv("../../agent/datasets/sqli.csv", header=None, names=names)
df_xss = pd.read_csv("../../agent/datasets/xss.csv", header=None, names=names)
df_ssti = pd.read_csv("../../agent/datasets/ssti.csv", header=None, names=names)


import pickle

# Model import
from sklearn.svm import SVC

#Load the trained model from disk
with open('../../agent/detection/ai/models/svc.pkl', 'rb') as f:
    loaded_model = pickle.load(f)


#Benign entry
names = ['UrlLength', 'NumberParams', 'NumberSpecialChars', 'RatioSpecialChars', 'NumberRoundBrackets', 'NumberSquareBrackets', 'NumberCurlyBrackets', 'NumberApostrophes', 'NumberQuotationMarks', 'NumberDots', 'NumberSlash', 'NumberBackslash', 'NumberComma', 'NumberColon', 'NumberSemicolon', 'NumberMinus', 'NumberPlus','NumberLessGreater', 'DistanceDots', 'DistanceSlash', 'DistanceBackslash', 'DistanceComma', 'DistanceColon', 'DistanceSemicolon', 'DistanceMinus', 'DistancePlus']

num_tries = 1000
start = time.time()
for i in range(num_tries):
    new_df = pd.DataFrame([df_benign.iloc[i]])
    new_df.columns = names
    # #Reshape the features to the accepted data for the svc model
    # a = pd.DataFrame(np.array(features).reshape(1, -1))
    # #Add the names of the columns
    # a.columns = names

    # Make the prediction
    prediction = loaded_model.predict(new_df)

    #Return the prediction to the runner
    #print(prediction[0])
end = time.time()
print('Avg:',(end - start)/num_tries * 1000, 'miliseconds')