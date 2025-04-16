#Model save
import pickle

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


df_benign['label'] = 'benign'
df_lfi['label'] = 'lfi'
df_sqli['label'] = 'sqli'
df_xss['label'] = 'xss'
df_ssti['label'] = 'ssti'

#Concatenate all the datasets
full_df = pd.concat([df_benign, df_lfi, df_sqli, df_xss, df_ssti] , ignore_index=True)

#Randomize the rows
full_df = full_df.sample(frac=1, random_state=42).reset_index(drop=True)

#Extract the features
X = full_df.drop('label', axis=1)
y = full_df['label']

#Prepare the train data and test data
X_train, X_test, y_train, y_test = train_test_split(
    X, y, test_size=0.2, random_state=42
)

svm = SVC(kernel='rbf', C=1, gamma='scale')  # kernel can be 'linear', 'poly', 'rbf', 'sigmoid'
svm.fit(X_train, y_train)


y_pred = svm.predict(X_test)

print("Accuracy:", accuracy_score(y_test, y_pred))
print("\nClassification Report:\n", classification_report(y_test, y_pred))

#Train the model with all dataset and save it to file

svm = SVC(kernel='rbf', C=1, gamma='scale')  # kernel can be 'linear', 'poly', 'rbf', 'sigmoid'
svm.fit(X, y)

with open('../../agent/detection/ai/models/svc.pkl', 'wb') as f:
    pickle.dump(svm, f)