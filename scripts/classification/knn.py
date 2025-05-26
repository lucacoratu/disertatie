import pickle

import pandas as pd

# Modelling
from sklearn.neighbors import KNeighborsClassifier
from sklearn.metrics import accuracy_score, confusion_matrix, precision_score, recall_score, classification_report
from sklearn.model_selection import train_test_split

names = ['UrlLength', 'NumberParams', 'NumberSpecialChars', 'RatioSpecialChars', 'NumberRoundBrackets', 'NumberSquareBrackets', 'NumberCurlyBrackets', 'NumberApostrophes', 'NumberQuotationMarks', 'NumberDots', 'NumberSlash', 'NumberBackslash', 'NumberComma', 'NumberColon', 'NumberSemicolon', 'NumberMinus', 'NumberPlus','NumberLessGreater', 'DistanceDots', 'DistanceSlash', 'DistanceBackslash', 'DistanceComma', 'DistanceColon', 'DistanceSemicolon', 'DistanceMinus', 'DistancePlus']

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

# scaler = StandardScaler()
# X_scaled = scaler.fit_transform(X)

#Prepare the train data and test data
X_train, X_test, y_train, y_test = train_test_split(
    X, y, test_size=0.2, random_state=42
)

knn = KNeighborsClassifier(n_neighbors=5) 
knn.fit(X_train, y_train)

#Predict
# Daca decomentez linia asta acuratetea este 76%
# X_test = pd.concat([X_test, df_test.drop('label', axis=1)], ignore_index=True)
# y_test = pd.concat([y_test, df_test['label']], ignore_index=True)

y_pred = knn.predict(X_test)

print("Accuracy:", accuracy_score(y_test, y_pred))
print("\nClassification Report:\n", classification_report(y_test, y_pred))

# model = KNeighborsClassifier(n_neighbors=5)
# model.fit(X, y)

# with open('../../agent/detection/ai/models/knn.pkl', 'wb') as f:
#     pickle.dump(model, f)